package timeManage

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	"SMNotifyBot/internal/loger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"time"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func WaitForTime() {

	for {
		//проверка на то пора ли делать уведомление
		time.Sleep(10 * time.Second)
		if time.Now().Minute() == 0 {
			callForNotification(time.Now().Hour())
			time.Sleep(60 * time.Second)
		}
		//проверка на то не отвечает ли телеграмм на запросы
		status := checkTelegramAPI()
		if status != http.StatusOK {
			logEntry := fmt.Sprintf(
				"API Telegram не работает. Код ответа:, %d, ожидаю ответа сервера\n",
				status)
			loger.LogToFile(logEntry)
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
			msg.Text = "У вас включено напоминание, но нет сохраненных сообщений"
			bot.Send(msg)
		}
	}
	loger.LogToFile("Уведомление о напоминания произведено успешно")
}

func checkTelegramAPI() int {
	resp, err := http.Get("https://api.telegram.org/bot" + api.GetApiToken() + "/getMe")
	if err != nil {
		loger.LogToFile("Ошибка при запросе к API: " + fmt.Sprintf("%s", err))
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
