package main

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	handleMSG "SMNotifyBot/internal/handleUpdate"
	"SMNotifyBot/internal/loger"
	"SMNotifyBot/internal/timeManage"
	"bufio"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(api.GetApiToken())
	if err != nil {
		log.Println(err)
	}
	bot.Debug = false
	loger.LogToFile("Authorized on account " + bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	handleDB.CreateTabelForData()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(updates, ctx)

	go timeManage.WaitForTime()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			cancel()
			return
		}
	}

	// Ожидаем завершения
	<-done
	loger.LogToFile("Работа завершена")
}

func receiveUpdates(updates tgbotapi.UpdatesChannel, ctx context.Context) {
	for {
		for {
			select {
			case <-ctx.Done():
				loger.LogToFile("Получена команда на завершение. Остановка работы.")
				return
			default:
				update := <-updates
				handleMSG.HandleUpdate(update)
			}
		}

	}
}
