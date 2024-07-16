package handleUpdate

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleButton"
	"SMNotifyBot/internal/handleCommand"
	"SMNotifyBot/internal/handleMSG"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		fmt.Printf("%v %v (%v): %v\n", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.ID, update.Message.Text)
	} else if update.CallbackQuery != nil {
		fmt.Printf("%v %v (%v): кнопка %v\n", update.CallbackQuery.From.UserName, update.CallbackQuery.From.FirstName, update.CallbackQuery.From.ID, update.CallbackQuery.Data)
	}
	switch {
	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton.HandleButton(update.CallbackQuery)
		break
	// Handle messages
	case update.Message.IsCommand():
		handleCommand.HandleCommand(update.Message)
		break
	// Handle commands
	case update.Message != nil:
		handleMSG.HandleMessage(update.Message)
		break

	// If something goes wrong
	default:
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Что-то пошло не так, попробуйте снова")
		bot.Send(msg)
	}
}
