package handleUpdate

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleButton"
	"SMNotifyBot/internal/handleMSG"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

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
		fmt.Println("Что-то пошло не так при использовании бота")
	}
}
