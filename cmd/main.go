package main

import (
	"abs-by-ammka-bot/internal/postgresql"
	"abs-by-ammka-bot/internal/telegram"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Инициализация конфигурации
	if err := initConfig(); err != nil {
		log.Panic("LOGGER: ", err)
		return
	}

	// Подключение к БД.
	db, err := postgresql.NewDB()
	if err != nil {
		log.Panic("LOGGER: ", err)
		return
	}

	// Инициализация телеграмм бота.
	tgBot, err := telegram.NewBotTelegram(postgresql.NewTelegramBotDB(db))
	if err != nil {
		log.Panic("LOGGER: ", err)
		return
	}

	// Запуск бота.
	if err := tgBot.StartBotUpdater(); err != nil {
		log.Panic("LOGGER: ", err)
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
