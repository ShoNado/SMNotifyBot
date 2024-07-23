package handleButton

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func HandleButton(query *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(query.From.ID, "")
	callback := tgbotapi.NewCallback(query.ID, "") //answer for call with empty string

	if !handleDB.IsUserInDb(query.From.ID) {
		handleDB.AddNewUser(query.From)
	}
	switch {
	case query.Data == "start":
		msg = answerCreator.StartTextSint(msg)
	case query.Data == "checkAll":
		if handleDB.NumOfMsgSaved(query.From.ID) > 0 {
			answerCreator.UserSavedMsgs(query.From.ID, msg)
			return
		} else {
			msg.Text = "В данный момент у вас нет сохраненных сообщений"
		}

	case query.Data == "timeset":
		msg, _ = answerCreator.TimeMsgSint(msg, query.From)

	case strings.HasPrefix(query.Data, "delete_"):
		deleteId, err := strconv.Atoi(query.Data[7:])
		if err != nil {
			msg.Text = "Что-то пошло не так при удалении"
		}
		if handleDB.DeleteMSG(query.From.ID, deleteId) {
			msg.Text = "Сообщение успешно удалено"
			callback = tgbotapi.NewCallback(query.ID, "Сообщение успешно удалено")
			deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
			if _, err := bot.Request(deleteMsg); err != nil {
				log.Println(err)
			}
		} else {
			msg.Text = "Что-то пошло не так при удалении"
		}

	case strings.HasPrefix(query.Data, "time"):
		timeId := query.Data
		if !handleDB.IsTimeTaken(query.From.ID, timeId) && handleDB.Addtime(query.From.ID, timeId) {
			_, keyboard := answerCreator.TimeMsgSint(msg, query.From)
			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				query.Message.Chat.ID,
				query.Message.MessageID,
				keyboard)
			if _, err := bot.Send(editMsg); err != nil {
				log.Printf("Не удалось отредактровать сообщение")
			}
			if _, err := bot.Request(callback); err != nil {
				log.Println("не удалось ответить на колбек кнопки")
			}
			return
		} else if handleDB.IsTimeTaken(query.From.ID, timeId) && handleDB.DeleteTime(query.From.ID, timeId) {
			_, keyboard := answerCreator.TimeMsgSint(msg, query.From)
			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				query.Message.Chat.ID,
				query.Message.MessageID,
				keyboard)
			if _, err := bot.Send(editMsg); err != nil {
				log.Printf("Не удалось отредактровать сообщение")
			}
			if _, err := bot.Request(callback); err != nil {
				log.Println("не удалось ответить на колбек кнопки")
			}
			return
		} else {
			msg.Text = "Что-то пошло не так при изменении времени"
		}

	default:
		msg.Text = "где ты эту кнопку нашел??"
	}
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Не удалось ответить на команду")
	}

}
