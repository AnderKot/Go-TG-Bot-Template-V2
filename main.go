package main

import (
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/ini.v1"
	"gorm.io/driver/sqlite"
	"os"
)

func main() {
	config, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	runArgs := config.Section("runArgs")

	botAPI, err := tgBotAPI.NewBotAPI(runArgs.Key("BotKey").String())
	if err != nil {
		panic(err)
	}

	database := InitBase(sqlite.Open("Database.db"))
	defer database.ExportTemplates()

	bot := NewBotService(botAPI, database)

	defer bot.final()
	go bot.Start()
}

//postgres.Open("host="+dbHostName+" port="+dbPort+" user="+dbLogin+" password="+dbPass+" dbname=service sslmode=disable")

/*
	dbHostName, exists := os.LookupEnv("DB_HOST_NAME")
	if !exists {
		if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
		}
		dbHostName, exists = os.LookupEnv("DB_HOST_NAME")
		if !exists {
			panic("DB_HOST_NAME")
		}
	}

	dbLogin, exists := os.LookupEnv("DB_USERS_USER")
	if !exists {
		panic("DB_USERS_USER")
	}

	dbPort, exists := os.LookupEnv("DB_USERS_PORT")
	if !exists {
		panic("DB_USERS_PORT")
	}

	dbPass, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		panic("DB_PASSWORD")
	}
*/
