package handleUpdate

import (
	"SMNotifyBot/internal/handleButton"
	"SMNotifyBot/internal/handleMSG"
	"SMNotifyBot/internal/loger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update tgbotapi.Update) {
	switch {
	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton.HandleButton(update.CallbackQuery)
		break
	// Handle messages
	case update.Message != nil:
		handleMSG.HandleMessage(update.Message)
		break

	// If something goes wrong
	default:
		loger.LogToFile("Ошибка: При обработке update в handleUpdate произошла ошибка. Содержание update: " + fmt.Sprintln(update))
	}
}
