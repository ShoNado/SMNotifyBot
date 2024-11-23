package handleAPI

import (
	"SMNotifyBot/internal/loger"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TelegramBotToken string
}

func GetApiToken() string { //read token from file
	file, _ := os.Open("config/config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		loger.LogToFile("Не удалось подключиться к боту " + fmt.Sprintln(err))
	}
	return configuration.TelegramBotToken
}
