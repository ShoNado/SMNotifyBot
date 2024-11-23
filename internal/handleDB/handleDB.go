package handleDB

import (
	"SMNotifyBot/internal/loger"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

//var db *sql.DB

type Pass struct {
	Password string
	IP       string
}
type UserProfile struct {
	Id          int
	TgID        int64
	NameFromTg  string
	UserName    string
	TimesChosen int
	MsgSaved    int
	Time        [24]bool
	Msg         [35]string
}

func connectDB() *sql.DB {

	db, err := sql.Open("sqlite3", "config/userDB.db")
	if err != nil {
		loger.LogToFile("Критическая ошибка: Не удалось подключиться к базе данных")
		panic(err)
	}

	return db
}

func disconnectDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		loger.LogToFile("Критическая ошибка: Не удалось отключиться от базы данных")
		panic(err)
	}
}

func CreateTabelForData() {
	db := connectDB()
	defer disconnectDB(db)
	filePath := "config/userDB.sql"
	data, err := os.ReadFile(filePath)
	if err != nil {
		loger.LogToFile("Ошибка: Не удалось считать файл конфигурации " + fmt.Sprintln(err))
	}

	query := string(data)

	_, err = db.Exec(query) //
	if err != nil {
		loger.LogToFile("Ошибка: Не удалось создать таблицу для хранения данных " + fmt.Sprintln(err))
		panic(err)
	}
}

func IsUserInDb(TgID int64) bool {
	db := connectDB()
	defer disconnectDB(db)
	row := db.QueryRow("SELECT id FROM users WHERE tgId = ?", TgID)
	var id int
	if err := row.Scan(&id); err != nil {
		if fmt.Sprint(err) == "no such table: users " {
			loger.LogToFile(" Ошибка: Пропала таблица для данных, создаю заново")
			CreateTabelForData()
		}
		loger.LogToFile("Новый юзер, так как: " + fmt.Sprint(err) + " , инициирую добавление в ДБ")
		return false
	}
	if id != 0 {
		return true
	}
	return false
}

func GetUserData(TgID int64) UserProfile {
	db := connectDB()
	defer disconnectDB(db)
	row := db.QueryRow("SELECT * FROM users WHERE tgId = ?", TgID)
	var user UserProfile
	if err := row.Scan(&user.Id, &user.TgID, &user.NameFromTg, &user.UserName, &user.TimesChosen, &user.MsgSaved,
		&user.Time[0], &user.Time[1], &user.Time[2], &user.Time[3], &user.Time[4], &user.Time[5], &user.Time[6],
		&user.Time[7], &user.Time[8], &user.Time[9], &user.Time[10], &user.Time[11], &user.Time[12], &user.Time[13],
		&user.Time[14], &user.Time[15], &user.Time[16], &user.Time[17], &user.Time[18], &user.Time[19], &user.Time[20],
		&user.Time[21], &user.Time[22], &user.Time[23], &user.Msg[0], &user.Msg[1], &user.Msg[2], &user.Msg[3], &user.Msg[4],
		&user.Msg[5], &user.Msg[6], &user.Msg[7], &user.Msg[8], &user.Msg[9], &user.Msg[10], &user.Msg[11], &user.Msg[12],
		&user.Msg[13], &user.Msg[14], &user.Msg[15], &user.Msg[16], &user.Msg[17], &user.Msg[18], &user.Msg[19], &user.Msg[20],
		&user.Msg[21], &user.Msg[22], &user.Msg[23], &user.Msg[24], &user.Msg[25], &user.Msg[26], &user.Msg[27], &user.Msg[28],
		&user.Msg[29], &user.Msg[30], &user.Msg[31], &user.Msg[32], &user.Msg[33], &user.Msg[34]); err != nil {
	}
	return user
}

func AddNewUser(user *tgbotapi.User) {
	db := connectDB()
	defer disconnectDB(db)
	//goland:noinspection ALL
	_, err := db.Exec("INSERT INTO users (tgId, nameFromTg, userName, timesChosen, msgSaved) VALUES (?,?,?,?,?)",
		user.ID, user.FirstName+user.LastName, user.UserName, 0, 0)
	if err != nil {
		loger.LogToFile("addNewUser: " + fmt.Sprintln(err))
	}
}

func NumOfMsgSaved(TgID int64) int {
	db := connectDB()
	defer disconnectDB(db)
	row := db.QueryRow("SELECT MsgSaved FROM users WHERE tgId = ?", TgID)
	var msgCount int
	if err := row.Scan(&msgCount); err != nil {
		loger.LogToFile(err)
	}
	return msgCount
}

func AddMsg(TgID int64, text string) error {
	db := connectDB()
	defer disconnectDB(db)

	user := GetUserData(TgID)
	var empty int
	for i, val := range user.Msg {
		if val == "" || val == "empty" {
			empty = i + 1
			break
		}
	}
	memorySell := "msg" + strconv.Itoa(empty)
	query := "UPDATE" + " users SET " + memorySell + "  = ? WHERE tgId = ?"
	_, err1 := db.Exec(query, text, TgID)
	num := NumOfMsgSaved(TgID) + 1
	_, err2 := db.Exec("UPDATE users SET msgSaved  = ? WHERE tgId = ?", num, TgID)
	if err1 != nil {
		loger.LogToFile(err1)
		return err1
	}
	if err2 != nil {
		loger.LogToFile(err2)
		return err2
	}
	return nil
}

func DeleteMSG(TgID int64, deleteID int) bool {
	db := connectDB()
	defer disconnectDB(db)
	memorySell := "msg" + strconv.Itoa(deleteID+1)
	query := "SELECT " + memorySell + " FROM users WHERE tgId = " + "?"
	row := db.QueryRow(query, TgID)
	var msgData string
	if err := row.Scan(&msgData); err != nil {
		loger.LogToFile(err)
	}
	if msgData == "" || msgData == "empty" {
		return false
	} else {
		query := "UPDATE " + " users SET  " + memorySell + "  = 'empty' WHERE tgId = ?"
		_, err := db.Exec(query, TgID)
		if err != nil {
			loger.LogToFile(err)
			return false
		}
		num := NumOfMsgSaved(TgID) - 1
		_, err = db.Exec("UPDATE users SET msgSaved  = ? WHERE tgId = ?", num, TgID)
		if err != nil {
			loger.LogToFile(err)
			return false
		}
		return true
	}
}

func NumTimes(TgID int64) int {
	db := connectDB()
	defer disconnectDB(db)
	row := db.QueryRow("SELECT timesChosen FROM users WHERE tgId = ?", TgID)
	var timeCount int
	if err := row.Scan(&timeCount); err != nil {
		loger.LogToFile(err)
	}
	return timeCount
}

func Addtime(TgID int64, time string) bool {
	db := connectDB()
	defer disconnectDB(db)
	query := "UPDATE" + " users SET " + time + " = 1 WHERE tgId = ?"
	_, err := db.Exec(query, TgID)
	if err != nil {
		loger.LogToFile(err)
		return false
	}
	num := NumTimes(TgID) + 1
	_, err = db.Exec("UPDATE users SET timesChosen  = ? WHERE tgId = ?", num, TgID)
	if err != nil {
		loger.LogToFile(err)
		return false
	}
	return true
}

func IsTimeTaken(TgID int64, time string) bool {
	db := connectDB()
	defer disconnectDB(db)
	query := "SELECT " + time + " FROM users WHERE tgId = ?"
	row := db.QueryRow(query, TgID)
	var timeIs bool
	if err := row.Scan(&timeIs); err != nil {
		loger.LogToFile(err)
	}
	return timeIs
}

func DeleteTime(TgID int64, time string) bool {
	db := connectDB()
	defer disconnectDB(db)
	query := "UPDATE " + "users SET " + time + "  = 0 WHERE tgId = ?"
	_, err := db.Exec(query, TgID)
	if err != nil {
		loger.LogToFile(err)
		return false
	}
	num := NumTimes(TgID) - 1
	_, err = db.Exec("UPDATE users SET timesChosen  = ? WHERE tgId = ?", num, TgID)
	if err != nil {
		loger.LogToFile(err)
		return false
	}
	return true
}

func GetIdsByHour(hour int) []int64 {
	db := connectDB()
	defer disconnectDB(db)
	var usersID []int64
	queryTime := ""
	if hour > 0 && hour < 24 {
		queryTime = "time" + strconv.Itoa(hour)
	} else {
		queryTime = "time24"
	}
	query := "SELECT tgId " + "FROM users WHERE " + queryTime + " = ?"
	tgIDData, _ := db.Query(query, 1)
	for tgIDData.Next() {
		var tgID int64
		if err := tgIDData.Scan(&tgID); err != nil {
			loger.LogToFile(err)
			return nil
		}
		usersID = append(usersID, tgID)
	}
	return usersID
}
