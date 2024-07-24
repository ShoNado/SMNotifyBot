package timeManage

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func WaitForTime() {
	for {
		time.Sleep(time.Minute)
		if time.Now().Minute() == 0 {
			callForNotification(time.Now().Hour())
		}
	}
}

func callForNotification(hour int) {
	userList := handleDB.GetIdsByHour(hour)
	for _, id := range userList {
		msg := tgbotapi.NewMessage(id, "")
		if handleDB.NumOfMsgSaved(id) > 0 {
			answerCreator.UserSavedMsgs(id, msg)
		} else {
			msg.Text = "В данный момент у вас нет сохраненных сообщений"
			bot.Send(msg)
		}
	}
}
