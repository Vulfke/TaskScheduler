package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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

// logininfo:token:1119993912:AAGVQxcJZ-4zc2qpKqQn-Nwrlkk51Yys3UA
// webhook: https://api.telegram.org/bot1119993912:AAGVQxcJZ-4zc2qpKqQn-Nwrlkk51Yys3UA/getUpdates

type Message struct {
	chatId int
	text   string
}

func (msg Message) ToJson() string {

	return fmt.Sprintf("?chat_id=%d&text=%s", msg.chatId, msg.text)
}

type LoginInfo struct {
	chatId int32
	token  string
}

func NewLoginInfo(config GeneralTaskConfig) LoginInfo {
	return LoginInfo{
		token: config["token"],
	}
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
	fmt.Print(res)
}

type Task interface {
	Done() bool
}

type ExtendedTask interface {
	Task
	Result() GeneralTaskConfig
}

type ProcessEndedTask struct {
	pid   ProcessId
	pname ProcessName
}

func NewProcessEndedTask(config GeneralTaskConfig) ProcessEndedTask {
	var pid ProcessId = 0
	pname := "unknown"

	if val, ok := config.Get("pid"); ok {
		tpid, _ := strconv.Atoi(val)
		pid = tpid
	}

	if val, ok := config.Get("pname"); ok {
		pname = val
	}

	return ProcessEndedTask{
		pid,
		pname,
	}
}

func (t ProcessEndedTask) Done() bool {

	process, err := os.FindProcess(t.pid)
	if err != nil {
		// here I need reference on callback or maybe i don't
		return false
	}

	err = process.Signal(syscall.Signal(0))

	if err == nil {
		return true
	}

	return false
}

type UnknownTask struct{}

func (task UnknownTask) Done() bool {
	return true
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
	cnt := 0
	mod := 1000
	fmt.Println("task.Done() ", task.Done())
	for !task.Done() {
		if cnt%mod == 0 {
			fmt.Print(fmt.Sprintf("it is %d iteration", cnt))
		}
	}

	scheduler.api.SendMessage(Message{
		chatId: scheduler.api.ChatId(),
		text:   "Build is ended",
	})
}

func main() {
	exampleLoginConfig := map[string]string{
		"token":  "1119993912:AAGVQxcJZ-4zc2qpKqQn-Nwrlkk51Yys3UA",
		"chatId": "-429622161",
	}

	exampleTaskConfig := map[string]string{
		"type": "SomehtingElseJustForCheck",
	}

	taskScheduler := NewTaskScheduler(exampleLoginConfig)

	taskScheduler.StartTask(exampleTaskConfig)

}
