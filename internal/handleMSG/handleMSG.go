package handleMSG

import (
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleCommand"
	"SMNotifyBot/internal/handleDB"
	"SMNotifyBot/internal/loger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		handleCommand.HandleCommand(message)
		return
	}
	msg := tgbotapi.NewMessage(message.From.ID, "")
	if !handleDB.IsUserInDb(message.From.ID) {
		handleDB.AddNewUser(message.From)
	}
	switch {
	case handleDB.IsUserInDb(message.From.ID) && message.Text != "":
		if handleDB.NumOfMsgSaved(message.From.ID) < 35 {
			var saveMsg string
			if message.ForwardFrom != nil {
				if message.ForwardFrom.UserName != "" {
					saveMsg = fmt.Sprintf(message.ForwardFrom.FirstName+" "+message.ForwardFrom.LastName+
						" \n<a href=\"t.me/%v\">Перейти в чат</a>\n"+"\n\n<blockquote>"+message.Text+"</blockquote>", message.ForwardFrom.UserName)
				} else {
					saveMsg = fmt.Sprintf(message.ForwardFrom.FirstName+" "+message.ForwardFrom.LastName+
						" \n<a href=\"tg://user?id=%v\">Открыть профиль</a>"+"\n"+"\n\n<blockquote>"+message.Text+"</blockquote>", message.ForwardFrom.ID)
				}
			} else {
				saveMsg = fmt.Sprintf(message.From.FirstName + "\n" + "\n\n<blockquote>" + message.Text + "</blockquote>")
			}
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

	default:
		msg.Text = "Данный формат не поддерживается"
	}

	if _, err := bot.Send(msg); err != nil {
		loger.LogToFile("Ошибка: Не удалось ответить на команду в handleMSG" + fmt.Sprintln(err))
	}
}
