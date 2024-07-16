package handleDB

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	cfg := mysql.Config{
		User:   "evgenu",
		Passwd: getDataForDB().Password,
		Net:    "tcp",
		Addr:   getDataForDB().IP,
		DBName: "SMproj",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	db.SetConnMaxLifetime(time.Minute)
	if err != nil {
		log.Fatal(err)
		fmt.Println(2)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		fmt.Println(3)
	}
}

func getDataForDB() Pass { //read password from file
	file, _ := os.Open("config/configDB.json")
	decoder := json.NewDecoder(file)
	configuration := Pass{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	return configuration
}

func IsUserInDb(TgID int64) bool {
	connectDB()
	defer db.Close()
	row := db.QueryRow("SELECT id FROM users WHERE tgId = ?", TgID)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Println("нет такого юзера: ", err)
		return false
	}
	if id != 0 {
		return true
	}
	return false
}

func IsUserHaveMsg(TgID int64) bool {
	connectDB()
	defer db.Close()
	row := db.QueryRow("SELECT MsgSaved FROM users WHERE tgId = ?", TgID)
	var msgCount int
	if err := row.Scan(&msgCount); err != nil {
		log.Println(err)
		return false
	}
	if msgCount > 0 {
		return true
	}
	return false
}
func GetUserMsgs(TgID int64) UserProfile {
	connectDB()
	defer db.Close()
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
		log.Println(err)

	}
	fmt.Println(user)
	return user
}
func AddNewUser(user *tgbotapi.User) {
	connectDB()
	defer db.Close()
	_, err := db.Exec("INSERT INTO users (tgId,nameFromTg,userName,timesChosen,msgSaved) VALUES (?,?,?,?,?)",
		user.ID, user.FirstName+user.LastName, user.UserName, 0, 0)
	if err != nil {
		log.Println("addNewUser: ", err)
	}
}

func NumOfMsgSaved(TgID int64) int {
	connectDB()
	defer db.Close()
	row := db.QueryRow("SELECT MsgSaved FROM users WHERE tgId = ?", TgID)
	var msgCount int
	if err := row.Scan(&msgCount); err != nil {
		log.Println(err)
	}
	return msgCount

}

func AddMsg(TgID int64, text string) error {
	user := GetUserMsgs(TgID)
	var empty int
	for i, val := range user.Msg {
		if val == "" {
			empty = i + 1
			break
		}
	}
	memorySell := "msg" + strconv.Itoa(empty)
	connectDB()
	defer db.Close()
	query := "UPDATE users SET " + memorySell + " = '" + text + "' WHERE tgId = " + "?"
	fmt.Println(query)
	res, err := db.Exec(query, TgID)
	fmt.Println(res)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
