package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func InitDb() {
	var err error
	database, err = sql.Open("sqlite3", "./idStorage.db")
	errorPrint("Database open file error:", err)

	database.SetConnMaxLifetime(0)
	database.SetMaxOpenConns(100)
	database.SetMaxIdleConns(100)

	_, err = database.Exec("CREATE TABLE IF NOT EXISTS msg2usr (msgId INTEGER PRIMARY KEY, usrId INTEGER)")
	errorPrint("database.Exec:", err)
}

func SaveToDb(msgId int, usrId int64) {
	statement, err := database.Prepare("INSERT INTO msg2usr (msgId, usrId) VALUES (?, ?)")
	errorPrint("SaveToDb prepare:", err)
	_, err = statement.Exec(msgId, usrId)
	errorPrint("SaveToDb Exec:", err)
	statement.Close()
}

func SearchInDb(targetId int) (int64, error) {

	rows, err := database.Query("SELECT msgId, usrId FROM msg2usr")
	errorPrint("SearchInDb() Query:", err)

	var msgId int
	var usrId int64

	for rows.Next() {
		err = rows.Scan(&msgId, &usrId)
		errorPrint("Loop, rows.Scan", err)
		if msgId == targetId {
			rows.Close()
			return usrId, nil
		}
	}
	rows.Close()
	return usrId, errors.New("Autor not found")
}

func errorPrint(comment string, err error) {
	if err != nil {
		fmt.Println(comment, err)
	}
}
