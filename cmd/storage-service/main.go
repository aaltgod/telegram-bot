package main

import (
	"fmt"
	"os"

	server "github.com/aaltgod/telegram-bot/internal/storage-service"
	myhttp "github.com/aaltgod/telegram-bot/internal/storage-service/delivery/http"
	"github.com/aaltgod/telegram-bot/internal/storage-service/repository/postgresql"
	"github.com/aaltgod/telegram-bot/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := logger.NewLogger()

	viper.SetConfigFile("newadmin.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Fatal("config file not found")
		} else {
			logger.Fatal("read config error")
		}
	}

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalln(err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("STORAGE_HOST"), os.Getenv("STORAGE_USER"), os.Getenv("STORAGE_PASSWORD"),
		os.Getenv("STORAGE_DB_NAME"), os.Getenv("STORAGE_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(dsn, err)
	}

	storage := postgresql.NewStorage(logger, db)
	if err := storage.Migrate(); err != nil {
		logger.Fatal(err)
	}

	id := viper.GetInt64("admin_id")
	name := viper.GetString("admin_name")

	if err := storage.InitAdmin(name, id); err != nil {
		logger.Fatal(err)
	}

	handler := myhttp.NewHandler(logger, storage)
	server := server.NewServer(logger, handler, storage)

	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}
}
