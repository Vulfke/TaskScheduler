package main

import "fmt"

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
