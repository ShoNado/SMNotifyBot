package timeManage

import (
	"SMNotifyBot/internal/answerCreator"
	api "SMNotifyBot/internal/handleAPI"
	"SMNotifyBot/internal/handleDB"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"time"
)

var bot, _ = tgbotapi.NewBotAPI(api.GetApiToken())

func WaitForTime() {
	// Имя файла для записи логов
	var LogFileName = "logs.txt"
	// Создаем файл и проверяем наличие
	file, err := os.OpenFile(LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close()
	for {
		time.Sleep(10 * time.Second)
		if time.Now().Minute() == 0 {
			callForNotification(time.Now().Hour())
			time.Sleep(60 * time.Second)
		}
		status := checkTelegramAPI()
		currentTime := time.Now().Format("15:04:05 01-02-2006")
		if status != http.StatusOK {
			logEntry := fmt.Sprintf("%s, API Telegram не работает. Код ответа:, %d, ожидаю ответа сервера\n", currentTime, status)
			if _, err := file.WriteString(logEntry); err != nil {
				fmt.Printf("Ошибка при записи в файл: %vn", err)
			}
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
	log.Println("Уведомление о напоминания произведено успешно")
}

func checkTelegramAPI() int {
	resp, err := http.Get("https://api.telegram.org/bot" + api.GetApiToken() + "/getMe")
	if err != nil {
		fmt.Println("Ошибка при запросе к API: ", err)
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
