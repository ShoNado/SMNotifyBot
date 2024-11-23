package loger

import (
	"fmt"
	"os"
	"time"
)

func LogToFile(logText any) {
	// Имя файла для записи логов
	var LogFileName = "logs.txt"
	// Создаем файл и проверяем наличие
	file, err := os.OpenFile(LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close()

	var logEntry string
	timeStamp := time.Now().Format("15:04:05 01-02-2006")
	switch logText.(type) {
	case string:
		logEntry = fmt.Sprintf("%s %s \n", timeStamp, logText)
	case error:
		logEntry = fmt.Sprintf("%s Ошибка: Runtime err: %s \n", timeStamp, logText)
	default:
		logEntry = fmt.Sprintf("%s Ошибка: Не удалось установить тип ошибки \n", timeStamp)
	}

	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Printf("Ошибка при записи в файл: %vn", err)
	}
}
