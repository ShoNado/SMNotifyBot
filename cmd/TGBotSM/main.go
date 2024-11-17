package main

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	handleMSG "SMNotifyBot/internal/handleUpdate"
	"SMNotifyBot/internal/timeManage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(api.GetApiToken())
	if err != nil {
		log.Println(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	handleDB.CreateDB()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(updates)

	go timeManage.WaitForTime()

}

func receiveUpdates(updates tgbotapi.UpdatesChannel) {
	for {
		update := <-updates
		handleMSG.HandleUpdate(update)
	}
}
