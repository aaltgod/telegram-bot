package main

import (
	"fmt"
	"os"

	"github.com/aaltgod/telegram-bot/internal/api"
	myhttp "github.com/aaltgod/telegram-bot/internal/api/delivery/http"
	"github.com/aaltgod/telegram-bot/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger := logger.NewLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalln(err)
	}

	handler := myhttp.NewHandler(logger, fmt.Sprintf(
		"http://%s:%s", os.Getenv("STORAGE_HOST"),
		os.Getenv("HTTP_STORAGE_SERVICE_PORT"),
	))

	api := api.NewServer(logger, handler)
	if err := api.Start(); err != nil {
		logger.Fatal(err)
	}
}
