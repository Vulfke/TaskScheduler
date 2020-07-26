package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

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
	logger  *log.Logger
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

	// Should be moved to another section
	file, err := os.Create(fmt.Sprintf("logs/%s-%d.txt", pname, pid))
	if err != nil {
		log.Println("cannot create file for task", err)
	}

	logout := bufio.NewWriter(file)
	logger := log.New(logout, "ProcessEndedTask:", 0)

	return ProcessEndedTask{
		pid,
		pname,
		message,
		logger,
	}
}

func (t ProcessEndedTask) Done() bool {

	windows := true
	linux := false

	if windows {

		out, err := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(t.pid)).Output()
		t.logger.Println(string(out), "  ", fmt.Sprintf("\" PID eq %d\"", t.pid))
		// return true
		if err != nil {
			t.logger.Println("Windows, exec.Process() Error \r\n", err)
		}

		if strings.Contains(string(out), "No tasks are running") {
			return true
		}
		return false
	}

	if linux {

		process, err := os.FindProcess(t.pid)
		if err != nil {
			t.logger.Println("FindProcessError\r\n", err)
			return true
		}

		err = process.Signal(syscall.Signal(0))

		t.logger.Println("Process.Signal Error\r\n", err)

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

// UnknownTask
// ***********************************************************************************************

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
