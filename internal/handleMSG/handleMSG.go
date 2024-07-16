package handleMSG

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.From.ID, "")
	switch {
	case !handleDB.IsUserInDb(message.From.ID) && message.Text != "":
		handleDB.AddNewUser(message.From)
		err := handleDB.AddMsg(message.From.ID, message.Text)
		if err != nil {
			msg.Text = "Произошла ошибка при добавлении сообщения в базу данных"
		} else {
			msg.Text = "Сообщение успешно сохранено"
		}
		break

	case handleDB.IsUserInDb(message.From.ID) && message.Text != "":
		if handleDB.NumOfMsgSaved(message.From.ID) < 35 {
			saveMsg := fmt.Sprintf(message.From.FirstName + " \n@" + message.From.UserName + "\n" + message.Time().Format("02-01-2006 15:04:05") + "\n" + message.Text)
			err := handleDB.AddMsg(message.From.ID, saveMsg)
			if err != nil {
				msg.Text = "Произошла ошибка при добавлении сообщения в базу данных"
			} else {
				msg.Text = "Сообщение успешно сохранено"
			}
		} else {
			msg.Text = "Превышен лимит числа сообщений для сохранения, удалите старое чтобы освободить место"
		}
		break

	case message.Text != "":

	default:
		msg.Text = "Данный формат не поддерживается"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Не удалось ответить на команду")
		panic(err)
	}
}
