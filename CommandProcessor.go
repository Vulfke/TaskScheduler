package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// So I need to create some kind of command processor, and maybe some task parser however

// We are talking about command line interface.
// So command processor, I suppose there should be some commands
// addtask jsonfile returns taskid
// stoptask taskid returns true/false
// addtasks jsonfile returns array of taskids.
// well that's all I suppose.
// Also there should be some interface for telegram
// But to be honest That quite far away.

// JsonParser Probably that is some kind of unmarshaler or anything.
// Also I need to get map of strings. But I could use real dtos.

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
