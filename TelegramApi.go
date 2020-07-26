package main

import (
	"fmt"
	"net/http"
	"strconv"
)

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
