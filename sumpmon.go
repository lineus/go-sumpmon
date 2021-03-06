package sumpmon

import (
	"database/sql"
	"errors"
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
		epoch = time.Now().Unix()
	} else if err != nil {
		log.Fatal("Query Failed: ", err)
	}

	last := time.Unix(epoch, 0)
	now := time.Now()
	return last.After(now.Add(-60 * time.Minute))
}

// GetAllLogs - returns a slice of SqliteLogs, all of them in fact.
func (logger Logger) GetAllLogs() ([]sqlitelogs.SqliteLog, error) {
	stmt, err := logger.db.Prepare("SELECT * FROM logs;")
	if err != nil {
		return nil, err
	}

	cursor, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	return cursorToSlice(cursor)
}

// GetLogsBetween - return all of the logs betwixt two times
func (logger Logger) GetLogsBetween(start time.Time, end time.Time) ([]sqlitelogs.SqliteLog, error) {
	stmt, err := logger.db.Prepare("SELECT * FROM logs WHERE epoch BETWEEN ? AND ?;")
	if err != nil {
		return nil, err
	}

	cursor, err := stmt.Query(start.Unix(), end.Unix())
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	return cursorToSlice(cursor)
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

func cursorToSlice(cursor *sql.Rows) ([]sqlitelogs.SqliteLog, error) {
	ret := make([]sqlitelogs.SqliteLog, 0)
	var err error
	for cursor.Next() {
		var i int
		var e int64
		var a string
		var r string

		err = cursor.Scan(&i, &e, &a, &r)
		if err != nil {
			return ret, err
		}
		ret = append(ret, sqlitelogs.SqliteLog{
			ID:     i,
			Epoch:  e,
			Action: a,
			Result: r,
		})
	}

	return ret, err
}
