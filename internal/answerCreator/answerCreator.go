package answerCreator

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func StartTextSint(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	msg.Text = "Перешлите мне сообщения о которых вам нужно будет напомнить позже, не забудьте указать время. \n" +
		"Левая кнопка выбор времени получения уведомлений. \n" +
		"Правая кнопка показать все отложенные сообщения."
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выбор времени", "timeset"),
			tgbotapi.NewInlineKeyboardButtonData("Просмотреть сообщения", "checkAll"),
		),
	)
	return msg
}

func TimeMsgSint(msg tgbotapi.MessageConfig, userData *tgbotapi.User) (tgbotapi.MessageConfig, tgbotapi.InlineKeyboardMarkup) {
	user := handleDB.GetUserData(userData.ID)
	timesAvailable := make([]string, 24)
	for i := range timesAvailable {
		if i < 9 {
			timesAvailable[i] = "0" + strconv.Itoa(i+1)
		} else {
			timesAvailable[i] = strconv.Itoa(i + 1)
		}
		timesAvailable[i] += ":00"
		if user.Time[i] == true {
			timesAvailable[i] += "\U00002705"
		}
	}
	msg.Text = fmt.Sprintf("Выберете время когда хотите получать уведомление" +
		"\nИзменение большого числа времён может производиться с небольшой задержкой" +
		"\nВремя указано по МСК")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[0], "time1"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[1], "time2"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[2], "time3"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[3], "time4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[4], "time5"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[5], "time6"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[6], "time7"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[7], "time8"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[8], "time9"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[9], "time10"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[10], "time11"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[11], "time12"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[12], "time13"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[13], "time14"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[14], "time15"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[15], "time16"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[16], "time17"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[17], "time18"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[18], "time19"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[19], "time20"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[20], "time21"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[21], "time22"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[22], "time23"),
			tgbotapi.NewInlineKeyboardButtonData(timesAvailable[23], "time24"),
		),
	)
	msg.ReplyMarkup = keyboard
	return msg, keyboard
}

func UserSavedMsgs(TgId int64, msg tgbotapi.MessageConfig) {
	user := handleDB.GetUserData(TgId)
	for i, notification := range user.Msg {
		if notification != "" && notification != "empty" {
			msg.Text = notification
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Удалить напоминание", "delete_"+strconv.Itoa(i)),
				))
			msg.DisableWebPagePreview = true
			msg.ParseMode = "HTML"
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Не удалось ответить на команду")
				panic(err)
			}
		}

	}

}
