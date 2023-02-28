package main

import (
	telegrambot "github.com/arturo-source/tramalicantebot/telegram-bot"
)

func main() {
	err := telegrambot.Run()
	if err != nil {
		panic(err)
	}
}
