package main

import (
	"database/sql"
	"fmt"
	"sync" // experiments with mutex - failed
	_ "github.com/mattn/go-sqlite3"
)


var database *sql.DB
var mutex sync.Mutex
func InitDb() {
	var err error
	database, err = sql.Open("sqlite3", "./idStorage.db")
	// defer database.Close() - segfault when uncomm... idk wtf
	errorPrint("Database open file error:", err)
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS msg2usr (msgId INTEGER PRIMARY KEY, usrId INTEGER)")
	errorPrint("Database prepare, if file db not exist error:", err)
	_, err = statement.Exec()
	errorPrint("Database prepare Exec:", err)
}

func SaveToDb(msgId int, usrId int64) { // Problem HERE
	// var mutex sync.Mutex
	 mutex.Lock()
	
	statement, err := database.Prepare("INSERT INTO msg2usr (msgId, usrId) VALUES (?, ?)")
	errorPrint("SaveToDb prepare:", err)
	_, err = statement.Exec(msgId, usrId)
	errorPrint("SaveToDb Exec:", err)
	
	mutex.Unlock()
}

func SearchInDb(targetId int) int64 {
	mutex.Lock()
	rows, err := database.Query("SELECT msgId, usrId FROM msg2usr")
	errorPrint("SearchInDb() Query:", err)

	var msgId int
	var usrId int64

	for rows.Next() {
		err = rows.Scan(&msgId, &usrId)
		errorPrint("Loop, rows.Scan", err)
		if msgId == targetId {
			mutex.Unlock()
			return usrId
		}
	}
	mutex.Unlock()
	return ownerID
}

func errorPrint(comment string, err error) {
	if err != nil {
		fmt.Println(comment, err)
	}
}
