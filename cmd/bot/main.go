package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aaltgod/telegram-bot/internal/bot"
	"github.com/aaltgod/telegram-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	logger := logger.NewLogger()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	api := bot.NewApi(logger, fmt.Sprintf(
		"http://%s:%s", os.Getenv("STORAGE_HOST"),
		os.Getenv("HTTP_STORAGE_SERVICE_PORT"),
	))

	botApi, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logger.Fatal(err)
	}

	b := bot.NewBot(logger, botApi, api)

	logger.Infoln("Bot is running")
	if err := b.Start(); err != nil {
		logger.Fatal(err)
	}
}
