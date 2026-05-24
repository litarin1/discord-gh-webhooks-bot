package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Channel = iota
	Category
	Webhook
)

var tables = map[string]string{
	"bot_routes": `
	CREATE TABLE IF NOT EXISTS bot_routes(
		server_id INTEGER NOT NULL,
		repo VARCHAR(140) NOT NULL,
		type INTEGER NOT NULL, -- 0=channel, 1=category, 2=webhook
		channel_name_format VARCHAR(100) NULL,
		PRIMARY KEY (server_id, repo),

		CONSTRAINT check_channel_format CHECK (type != 2 OR channel_name_format IS NOT NULL)
	);
	`,
}

/*
Usage:
```go
db := BotDB{}
if err := db.Open(); err != nil {
    panic(err)
}
defer db.Close()
if err := db.InitDB(); err != nil{
	panic(err)
}
// now do your stuff
```
*/

type BotDB struct {
	db *sql.DB
}

// func (db *BotDB) tableExists(table string) (exists bool, err error) {
// 	query := `SELECT EXISTS(
// 		SELECT 1 FROM sqlite_master
// 		WHERE type='table' AND name=?
// 	)`
// 	err = db.db.QueryRow(query, table).Scan(&exists)
// 	return
// }

func (db *BotDB) Open(dataSourceName string) (err error) {
	slog.Info("db Open()", "dataSourceName", dataSourceName)
	if db.db != nil {
		db.db.Close()
	}
	db.db, err = sql.Open("sqlite3", dataSourceName)
	return
}

func (db *BotDB) Close() error {
	slog.Info("db Close()")
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *BotDB) InitDB() error {
	// google.com/ai told me to do this for better compatibility
	_, err := db.db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	for table, query := range tables {
		slog.Info("Table", "name", table)
		_, err = db.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// server_id is int64 because sqlite cant store uint64
func (db *BotDB) GetRoute(server_id int64, repo string) (err error, route_type byte, channel_name_format sql.NullString) {
	res := db.db.QueryRow("SELECT type,channel_name_format FROM bot_routes WHERE server_id=? AND repo=?", server_id, repo)
	err = res.Scan(&route_type, &channel_name_format)
	return
}
