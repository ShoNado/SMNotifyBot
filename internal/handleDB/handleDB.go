package handleDB

import (
	"database/sql"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"time"
)

var db *sql.DB

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

func connectDB() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Second * 3)
	db.SetMaxOpenConns(1000)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println(pingErr)
	}
}

func CreateDB() {
	connectDB()

	filePath := "config/userDB.sql"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}
	query := string(data)

	_, err = db.Exec(query) //
	if err != nil {
		panic(err)
	}
}

func getDataForDB() Pass { //read password from file
	file, _ := os.Open("config/configDB.json")
	decoder := json.NewDecoder(file)
	configuration := Pass{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println(err)
	}
	return configuration
}

func IsUserInDb(TgID int64) bool {
	connectDB()
	row := db.QueryRow("SELECT id FROM users WHERE tgId = ?", TgID)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Println("Новый юзер: ", err, " , инициирую добавление в ДБ")
		return false
	}
	if id != 0 {
		return true
	}
	return false
}

func GetUserData(TgID int64) UserProfile {
	connectDB()
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
	connectDB()
	_, err := db.Exec("INSERT INTO users (tgId,nameFromTg,userName,timesChosen,msgSaved) VALUES (?,?,?,?,?)",
		user.ID, user.FirstName+user.LastName, user.UserName, 0, 0)
	if err != nil {
		log.Println("addNewUser: ", err)
	}
}

func NumOfMsgSaved(TgID int64) int {
	connectDB()
	row := db.QueryRow("SELECT MsgSaved FROM users WHERE tgId = ?", TgID)
	var msgCount int
	if err := row.Scan(&msgCount); err != nil {
		log.Println(err)
	}
	return msgCount
}

func AddMsg(TgID int64, text string) error {
	user := GetUserData(TgID)
	var empty int
	for i, val := range user.Msg {
		if val == "" || val == "empty" {
			empty = i + 1
			break
		}
	}
	memorySell := "msg" + strconv.Itoa(empty)
	connectDB()
	query := "UPDATE users SET " + memorySell + " = '" + text + "' WHERE tgId = " + "?"
	_, err1 := db.Exec(query, TgID)
	num := NumOfMsgSaved(TgID) + 1
	_, err2 := db.Exec("UPDATE users SET msgSaved  = ? WHERE tgId = ?", num, TgID)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	if err2 != nil {
		log.Println(err2)
		return err2
	}
	return nil
}

func DeleteMSG(TgID int64, deleteID int) bool {
	connectDB()
	memorySell := "msg" + strconv.Itoa(deleteID+1)
	query := "SELECT " + memorySell + " FROM users WHERE tgId = " + "?"
	row := db.QueryRow(query, TgID)
	var msgData string
	if err := row.Scan(&msgData); err != nil {
		log.Println(err)
	}
	if msgData == "" || msgData == "empty" {
		return false
	} else {
		query := "UPDATE users SET " + memorySell + " = 'empty' WHERE tgId = " + "?"
		_, err := db.Exec(query, TgID)
		if err != nil {
			log.Println(err)
			return false
		}
		num := NumOfMsgSaved(TgID) - 1
		_, err = db.Exec("UPDATE users SET msgSaved  = ? WHERE tgId = ?", num, TgID)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
}

func NumTimes(TgID int64) int {
	connectDB()
	row := db.QueryRow("SELECT timesChosen FROM users WHERE tgId = ?", TgID)
	var timeCount int
	if err := row.Scan(&timeCount); err != nil {
		log.Println(err)
	}
	return timeCount
}

func Addtime(TgID int64, time string) bool {
	connectDB()
	query := "UPDATE users SET " + time + " = 1 WHERE tgId = " + "?"
	_, err := db.Exec(query, TgID)
	if err != nil {
		log.Println(err)
		return false
	}
	num := NumTimes(TgID) + 1
	_, err = db.Exec("UPDATE users SET timesChosen  = ? WHERE tgId = ?", num, TgID)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func IsTimeTaken(TgID int64, time string) bool {
	connectDB()
	query := "SELECT " + time + " FROM users WHERE tgId = ?"
	row := db.QueryRow(query, TgID)
	var timeIs bool
	if err := row.Scan(&timeIs); err != nil {
		log.Println(err)
	}
	return timeIs
}

func DeleteTime(TgID int64, time string) bool {
	connectDB()
	query := "UPDATE users SET " + time + " = 0 WHERE tgId = " + "?"
	_, err := db.Exec(query, TgID)
	if err != nil {
		log.Println(err)
		return false
	}
	num := NumTimes(TgID) - 1
	_, err = db.Exec("UPDATE users SET timesChosen  = ? WHERE tgId = ?", num, TgID)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetIdsByHour(hour int) []int64 {
	connectDB()
	var usersID []int64
	query := ""
	if hour > 0 && hour < 24 {
		query = "SELECT tgId FROM users WHERE time" + strconv.Itoa(hour) + " = ?"
	} else {
		query = "SELECT tgId FROM users WHERE time24 = ?"
	}
	tgIDData, _ := db.Query(query, true)
	for tgIDData.Next() {
		var tgID int64
		if err := tgIDData.Scan(&tgID); err != nil {
			return nil
		}
		usersID = append(usersID, tgID)
	}
	return usersID
}
