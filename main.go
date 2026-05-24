package main

import (
	"database/sql"
	"discord-gh-webhooks-bot/db"
	"discord-gh-webhooks-bot/parser"
	"fmt"
	"log"
	"os"

	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	db := db.BotDB{}

	if err := db.Open("./bot.db"); err != nil {
		log.Panic(err)
	}
	defer db.Close()
	if err := db.InitDB(); err != nil {
		log.Panic(err)
	}
	err, route_type, channel_name_format := db.GetRoute(123, "litarin1/discord-gh-webhooks-bot")
	if err == sql.ErrNoRows {
		log.Println("litarin1/discord-gh-webhooks-bot does not exist")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Println("litarin1/discord-gh-webhooks-bot:", route_type, channel_name_format.String)
	}

	var bytes []byte
	bytes, err = os.ReadFile("routes.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	parsed, err := parser.Parse(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed.GetRoute("litarin1", ""))
}
