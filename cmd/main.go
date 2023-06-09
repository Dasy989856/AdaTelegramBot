package main

import (
	"AdaTelegramBot/internal/postgresql"
	"AdaTelegramBot/internal/telegram"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Инициализация конфигурации
	if err := initConfig(); err != nil {
		log.Panic("main: error initConfig: ", err)
		return
	}

	// Подключение к БД.
	db := postgresql.NewDB()

	// Инициализация телеграмм бота.
	tgBot, err := telegram.NewBotTelegram(postgresql.NewTelegramBotDB(db))
	if err != nil {
		log.Panic("main: error telegram.NewBotTelegram: ", err)
		return
	}

	// Запуск бота.
	if err := tgBot.StartBotUpdater(); err != nil {
		log.Panic("main: error tgBot.StartBotUpdater: ", err)
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
