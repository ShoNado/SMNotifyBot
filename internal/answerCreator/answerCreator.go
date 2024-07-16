package answerCreator

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

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

func TimeMsgSint(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	msg.Text = "Выберете время когда хотите получать уведомление"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("01:00", "time01"),
			tgbotapi.NewInlineKeyboardButtonData("02:00", "time02"),
			tgbotapi.NewInlineKeyboardButtonData("03:00", "time03"),
			tgbotapi.NewInlineKeyboardButtonData("04:00", "time04"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("05:00", "time05"),
			tgbotapi.NewInlineKeyboardButtonData("06:00", "time06"),
			tgbotapi.NewInlineKeyboardButtonData("07:00", "time07"),
			tgbotapi.NewInlineKeyboardButtonData("08:00", "time08"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("09:00", "time09"),
			tgbotapi.NewInlineKeyboardButtonData("10:00", "time10"),
			tgbotapi.NewInlineKeyboardButtonData("11:00", "time11"),
			tgbotapi.NewInlineKeyboardButtonData("12:00", "time12"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("13:00", "time13"),
			tgbotapi.NewInlineKeyboardButtonData("14:00", "time14"),
			tgbotapi.NewInlineKeyboardButtonData("15:00", "time15"),
			tgbotapi.NewInlineKeyboardButtonData("16:00", "time16"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("17:00", "time17"),
			tgbotapi.NewInlineKeyboardButtonData("18:00", "time18"),
			tgbotapi.NewInlineKeyboardButtonData("19:00", "time19"),
			tgbotapi.NewInlineKeyboardButtonData("20:00", "time20"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("21:00", "time21"),
			tgbotapi.NewInlineKeyboardButtonData("22:00", "time22"),
			tgbotapi.NewInlineKeyboardButtonData("23:00", "time23"),
			tgbotapi.NewInlineKeyboardButtonData("00:00", "time24"),
		),
	)
	return msg
}

func UserSavedMsgs(TgId int64, msg tgbotapi.MessageConfig) {
	var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())
	user := handleDB.GetUserMsgs(TgId)
	for i, notification := range user.Msg {
		if notification != "" {
			msg.Text = notification
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Удалить напоминание", "delete_"+strconv.Itoa(i)),
				))
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Не удалось ответить на команду")
				panic(err)
			}
		}

	}

}
