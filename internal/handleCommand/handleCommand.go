package handleCommand

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	"SMNotifyBot/internal/loger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.From.ID, "")
	if !handleDB.IsUserInDb(message.From.ID) {
		handleDB.AddNewUser(message.From)
	}
	switch {
	case message.Command() == "start":
		msg = answerCreator.StartTextSint(msg)
	case message.Command() == "info":
		msg.Text = "Данный бот сохраняет сообщения которые вы ему пересылаете, а затем уведомляет о них, в заданное время. " +
			"Для навигации используйте встроенные кнопки под сообщениями."
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться в начало", "start"),
			),
		)

	case message.Command() == "settime":
		msg, _ = answerCreator.TimeMsgSint(msg, message.From)

	case message.Command() == "checkall":
		if handleDB.NumOfMsgSaved(message.From.ID) > 0 {
			answerCreator.UserSavedMsgs(message.From.ID, msg)
			return
		} else {
			msg.Text = "В данный момент у вас нет сохраненных сообщений"
		}

	default:
		msg.Text = "Неизвестная команда"
	}

	if _, err := bot.Send(msg); err != nil {
		loger.LogToFile("Ошибка: Не удалось ответить на команду в handleCommand" + fmt.Sprintln(err))
	}

}
