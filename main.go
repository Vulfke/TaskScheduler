package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	loginConfigPath := "config/login.json"
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
	var loginConfig LoginInfo

	err = json.Unmarshal(loginConfigBytes, &loginConfig)
	if err != nil {
		fmt.Println(err)
	}

	exampleLoginConfig := map[string]string{
		"token":  loginConfig.Token,
		"chatId": loginConfig.ChatID,
	}

	taskScheduler := NewTaskScheduler(exampleLoginConfig)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	commandProcessor := NewCommandProcessor(taskScheduler)

	fmt.Println("What can I do for you?")

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
