package main

import "fmt"

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
