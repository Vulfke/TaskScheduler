package main

import "fmt"

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

		for !task.Done() {
		}

		scheduler.api.SendMessage(Message{
			chatId: scheduler.api.ChatId(),
			text:   task.Message(),
		})

	}(task)
}
