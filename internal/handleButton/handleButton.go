package handleButton

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleButton(query *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(query.From.ID, "")
	callback := tgbotapi.NewCallback(query.ID, "") //answer for call with empty string
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	}
	switch {
	case query.Data == "start":
		msg = answerCreator.StartTextSint(msg)
		if !handleDB.IsUserInDb(query.From.ID) {
			handleDB.AddNewUser(query.From)
		}
	case query.Data == "checkAll":
		if handleDB.IsUserInDb(query.From.ID) || handleDB.IsUserHaveMsg(query.From.ID) {
			answerCreator.UserSavedMsgs(query.From.ID, msg)
			return
		} else {
			msg.Text = "В данный момент у вас нет сохраненных сообщений"
		}

	case query.Data == "timeset":
		msg = answerCreator.TimeMsgSint(msg)

	default:
		msg.Text = "где ты эту кнопку нашел??"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Не удалось ответить на команду")
		panic(err)
	}

}
