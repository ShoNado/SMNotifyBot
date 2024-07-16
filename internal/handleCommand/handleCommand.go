package handleCommand

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.From.ID, "")
	switch {
	case message.Command() == "start":
		msg = answerCreator.StartTextSint(msg)
		if !handleDB.IsUserInDb(message.From.ID) {
			handleDB.AddNewUser(message.From)
		}
	case message.Command() == "info":
		msg.Text = "Данный бот сохраняет сообщения которые вы ему пересылаете, а затем уведомляет о них, в заданное время. " +
			"Для навигации используйте встроенные кнопки под сообщениями."
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться в начало", "start"),
			),
		)

	case message.Command() == "settime":
		msg = answerCreator.TimeMsgSint(msg)

	case message.Command() == "checkall":

		if handleDB.IsUserInDb(message.From.ID) || handleDB.IsUserHaveMsg(message.From.ID) {
			answerCreator.UserSavedMsgs(message.From.ID, msg)
			return
		} else {
			msg.Text = "В данный момент у вас нет сохраненных сообщений"
		}

	default:
		msg.Text = "Неизвестная команда"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Не удалось ответить на команду")
		panic(err)
	}

}
