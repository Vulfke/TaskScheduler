package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type ProcessId = int
type ProcessName = string
type ChatId = string

type GeneralConfig interface {
	Get(field string) (string, bool)
	Set(field, value string)
}

type GeneralTaskConfig map[string]string

func (config GeneralTaskConfig) Get(field string) (string, bool) {
	fmt.Println(field, config)
	val, ok := config[field]
	return val, ok
}

func (config GeneralTaskConfig) Set(field, value string) {
	config[field] = value
}

type TaskSchedulerConfig = GeneralTaskConfig

type Message struct {
	chatId int
	text   string
}

func (msg Message) ToJson() string {

	return fmt.Sprintf("?chat_id=%d&text=%s", msg.chatId, msg.text)
}

type TelegramApi interface {
	Login(config GeneralConfig)
	SendMessage(msg Message)
	ChatId() int
}

type TelegramApiImpl struct {
	Client   http.Client
	EndPoint string
	Token    string
	chatId   int
}

func NewTelegramApiImpl(config GeneralConfig) TelegramApiImpl {
	client := http.Client{}
	webhook := "https://api.telegram.org/bot"
	token, ok := config.Get("token")
	if !ok {
		token = ""
	}
	chatIdStr, ok := config.Get("chatId")
	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		chatId = 0
	}
	fmt.Println(chatId, token)
	return TelegramApiImpl{
		Client:   client,
		EndPoint: webhook,
		Token:    token,
		chatId:   chatId,
	}
}

func (api TelegramApiImpl) ChatId() int {
	return api.chatId
}

func (api TelegramApiImpl) Login(config GeneralConfig) {
	api.Token, _ = config.Get("token")
	api.EndPoint = "https://api.telegram.org/bot"
}

func (api TelegramApiImpl) SendMessage(msg Message) {
	fmt.Println("Send Message")
	apiEndPoint := api.EndPoint + api.Token + "/sendMessage" + msg.ToJson()
	fmt.Println(apiEndPoint)
	res, err := api.Client.Post(apiEndPoint, "application/json", nil)
	if err != nil {
		fmt.Println("Lol there is error with net")
		return
	}
	defer res.Body.Close()
	fmt.Println(res)
}

type Task interface {
	Done() bool
	Message() string
}

type ExtendedTask interface {
	Task
	Result() GeneralTaskConfig
}

type ProcessEndedTask struct {
	pid     ProcessId
	pname   ProcessName
	message string
}

func NewProcessEndedTask(config GeneralTaskConfig) ProcessEndedTask {
	var pid ProcessId = 0
	pname := "unknown"
	message := "build is ended"

	if val, ok := config.Get("pid"); ok {
		tpid, _ := strconv.Atoi(val)
		pid = tpid
	}

	if val, ok := config.Get("pname"); ok {
		pname = val
	}

	if val, ok := config.Get("message"); ok {
		message = val
	}

	return ProcessEndedTask{
		pid,
		pname,
		message,
	}
}

func (t ProcessEndedTask) Done() bool {

	windows := true
	linux := false

	if windows {

		out, err := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(t.pid)).Output()
		fmt.Println(string(out), "  ", fmt.Sprintf("\" PID eq %d\"", t.pid))
		// return true
		if err != nil {
			log.Println("Windows, exec.Process() Error \r\n", err)
		}

		if strings.Contains(string(out), "No tasks are running") {
			return true
		}
		return false
	}

	if linux {

		process, err := os.FindProcess(t.pid)
		if err != nil {
			fmt.Println("FindProcessError\r\n", err)
			return true
		}

		err = process.Signal(syscall.Signal(0))

		fmt.Println("Process.Signal Error\r\n", err)

		if err != nil {
			return true
		}

		return false
	}

	return true
}

func (t ProcessEndedTask) Message() string {
	return t.message
}

type UnknownTask struct{}

func (task UnknownTask) Done() bool {
	return true
}

func (task UnknownTask) Message() string {
	return "There wasn't any real task actually"
}

func NewUnknownTask() UnknownTask {
	return UnknownTask{}
}

type TaskFabric struct {
}

func (fabric TaskFabric) Create(config GeneralTaskConfig) Task {

	taskType, ok := config.Get("type")

	if !ok {
		fmt.Println("There is no type in task config")
	}

	switch taskType {
	case "ProcessEndedTask":
		return NewProcessEndedTask(config)
	default:
		return NewUnknownTask()
	}
}

type TaskScheduler struct {
	api         TelegramApi
	mTaskList   []Task
	mTaskFabric TaskFabric
}

/**
* Ok I don't really like the idea of several config file for anything.
	I really hate config files, just because I need to make my life harder
	Just because of nothing literaly nothing that is matter.
	I dont' really understand why I sould create so agile application just for my personal use
	So I will close this with hardcoded thing and after that I will create configs.
	Just Because I am done with it.
*
*/

func NewTaskScheduler(config TaskSchedulerConfig) TaskScheduler {

	return TaskScheduler{
		NewTelegramApiImpl(config),
		[]Task{},
		TaskFabric{},
	}
}

func (scheduler TaskScheduler) StartTask(config GeneralTaskConfig) {
	task := scheduler.mTaskFabric.Create(config)
	fmt.Println(task, task.Done())
	go func(task Task) {
		cnt := 0
		mod := 1000
		fmt.Println("task.Done() ", task.Done())
		for !task.Done() {
			if cnt%mod == 0 {
				fmt.Print(fmt.Sprintf("it is %d iteration", cnt))
			}
			cnt += 1
		}

		scheduler.api.SendMessage(Message{
			chatId: scheduler.api.ChatId(),
			text:   task.Message(),
		})
	}(task)
}

/*
So I need to create some kind of command processor, and maybe some task parser however

We are talking about command line interface.
So command processor, I suppose there should be some commands
addtask jsonfile returns taskid
stoptask taskid returns true/false
addtasks jsonfile returns array of taskids.
well that's all I suppose.
Also there should be some interface for telegram
But to be honest That quite far away.

JsonParser Probably that is some kind of unmarshaler or anything. Also I need to get map of strings. But I could use real dtos.
*/

type Command struct {
}

type LoginInfo struct {
	Token  string `json: token`
	ChatID string `json: chatId`
}

type TaskDescription struct {
	TaskType     string    `json: taskType`
	Timeout      int       `json: timeout`
	AfterMessage string    `json: message`
	Command      Command   `json: command`
	ProcessID    ProcessId `json: processId`
}

type CommandProcessor struct {
	taskScheduler TaskScheduler
}

func (cp CommandProcessor) Process(sentense []byte) {
	row := strings.Trim(string(sentense), " \r\n")
	tokens := strings.Split(row, " ")
	command := tokens[0]

	switch command {
	case "addTask":
		cp.processAddTask(tokens)
	case "finishTask":
		cp.processFinishTask(tokens)
	case "exit":
		cp.exit()
	}
}

func (cp CommandProcessor) processAddTask(tokens []string) {
	taskConfigFilePath := tokens[1]

	fmt.Println("tokens:", tokens)

	wd, _ := os.Getwd()
	fmt.Println(taskConfigFilePath, wd)

	taskConfigBytes, err := ioutil.ReadFile(taskConfigFilePath)

	if err != nil {
		fmt.Println("CommandProcessor", err)
	}

	var config map[string]string

	err = json.Unmarshal(taskConfigBytes, &config)

	if err != nil {
		fmt.Println("CommandProcessor", err)
	}

	cp.taskScheduler.StartTask(config)

}

func (cp CommandProcessor) processFinishTask(tokens []string) {

}

func (cp CommandProcessor) exit() {

}

func NewCommandProcessor(taskScheduler TaskScheduler) CommandProcessor {
	return CommandProcessor{
		taskScheduler: taskScheduler,
	}
}

func main() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	loginConfigPath := "D:\\Documents\\GoProjects\\task-sheduler\\config\\login.json"
	loginConfigFile, err := os.Open(loginConfigPath)
	defer loginConfigFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	loginConfigBytes := make([]byte, 1024)
	n, errRead := loginConfigFile.Read(loginConfigBytes)
	fmt.Println(n)
	if errRead != nil {
		fmt.Println(errRead)
	}
	loginConfigBytes = loginConfigBytes[0:n]
	fmt.Println("Bytes: ", loginConfigBytes)
	var loginConfig LoginInfo

	err = json.Unmarshal(loginConfigBytes, &loginConfig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(loginConfig)

	exampleLoginConfig := map[string]string{
		"token":  loginConfig.Token,
		"chatId": loginConfig.ChatID,
	}

	taskScheduler := NewTaskScheduler(exampleLoginConfig)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	commandProcessor := NewCommandProcessor(taskScheduler)

	for {
		buffer := make([]byte, 1024)
		commandLineLength, err := reader.Read(buffer)
		buffer = buffer[:commandLineLength]
		if err != nil {
			fmt.Println(err)
			continue
		}

		commandProcessor.Process(buffer)
		writer.Write([]byte("Command is being processed"))
	}
}

// мама мыла раму конфидераты, мать их, задолбали уже всех американских матерей
