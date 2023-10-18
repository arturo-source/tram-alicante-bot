package main

import (
	"net/http"
	"os"

	telegrambot "github.com/arturo-source/tramalicantebot/telegram-bot"
)

func main() {
	// because fl0 deploy fails if you are not listening
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	err := telegrambot.Run()
	if err != nil {
		panic(err)
	}
}
