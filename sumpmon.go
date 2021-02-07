package sumpmon

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lineus/go-sqlitelogs"
	_ "github.com/mattn/go-sqlite3" // use sqlite3
)

// Logger - a thing
type Logger struct {
	db *sql.DB
}

// SaveLog - saves a log
func (logger Logger) SaveLog(action string, result string) (sql.Result, error) {
	var epoch = time.Now().Unix()
	stmt, err := logger.db.Prepare("INSERT INTO logs(epoch, action, result) values(?,?,?);")
	if err != nil {
		log.Fatal("PREPARE FAILED: ", err)
	}
	res, err := stmt.Exec(epoch, action, result)
	if err != nil {
		log.Fatal("INSERT FAILED: ", err)
	}

	return res, err
}

// Alive - has there been a log saved in the db in the last hour
func (logger Logger) Alive() bool {
	stmt, err := logger.db.Prepare("SELECT epoch FROM logs ORDER BY id DESC LIMIT 1;")
	if err != nil {
		log.Fatal("Prepare Failed: ", err)
	}

	var epoch int64
	err = stmt.QueryRow().Scan(&epoch)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		fmt.Println("error error")
		epoch = time.Now().Unix()
	}

	fmt.Printf("epoch: %+v", epoch)
	last := time.Unix(epoch, 0)
	now := time.Now()
	fmt.Println("last: ", last)
	fmt.Println("now: ", now)
	return last.After(now.Add(-60 * time.Minute))
}

// GetAllLogs - returns a slice of SqliteLogs, all of them in fact.
func (logger Logger) GetAllLogs() []sqlitelogs.SqliteLog {
	ret := make([]sqlitelogs.SqliteLog, 1)
	return ret
}

// GetLogsBetween - return all of the logs betwixt two times
func (logger Logger) GetLogsBetween(start time.Time, end time.Time) []sqlitelogs.SqliteLog {
	ret := make([]sqlitelogs.SqliteLog, 1)
	return ret
}

// Init - connect to the db and get your Logger instance
func Init(dsn string) (sqlitelogs.SqliteLogger, error) {
	var l Logger
	db, err := sql.Open("sqlite3", dsn)
	l.db = db

	_, err = l.db.Exec(`CREATE TABLE IF NOT EXISTS logs(
		id INTEGER PRIMARY KEY ASC,
		epoch INTEGER,
		action TEXT,
		result TEXT
	)`)

	if err != nil {
		log.Fatal("CREATE TABLE FAILED: ", err)
	}

	return l, err
}
