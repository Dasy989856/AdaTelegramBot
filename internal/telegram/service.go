package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// Структура телеграмм бота.
type BotTelegram struct {
	bot          *tgbotapi.BotAPI
	db           models.TelegramBotDB
	cashAdEvents map[int64]*models.AdEvent // Хэш-таблица ad событий.
}

// Создание телеграмм бота.
func NewBotTelegram(db models.TelegramBotDB) (*BotTelegram, error) {
	token := viper.GetString("token.telegram")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &BotTelegram{bot: bot, db: db, cashAdEvents: make(map[int64]*models.AdEvent)}, nil
}

// Запуск апдейтера.
func (b *BotTelegram) StartBotUpdater() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates := b.InitUpdatesChanel()
	if err := b.handlerUpdates(updates); err != nil {
		return err
	}
	return nil
}

func (b *BotTelegram) InitUpdatesChanel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	return b.bot.GetUpdatesChan(u)
}

// Обработчики сообщений.
func (b *BotTelegram) handlerUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		// Обработка команд.
		if update.Message != nil && update.Message.IsCommand() {
			if err := b.handlerCommand(update.Message); err != nil {
				log.Println(err)
			}
			continue
		}

		// Обработка сообщений.
		if update.Message != nil {
			if err := b.handlerMessage(update.Message); err != nil {
				log.Println(err)
			}
			continue
		}

		// Обработка CallbackQuery.
		if update.CallbackQuery != nil {
			if err := b.handlerCallbackQuery(&update); err != nil {
				log.Println(err)
			}
			continue
		}
	}

	return fmt.Errorf("updates chanel closed")
}

// Получение хэша ad события.
func getAdEventFromCash(b *BotTelegram, userId int64) (*models.AdEvent, error) {
	adEvent, ok := b.cashAdEvents[userId]
	if ok {
		return adEvent, nil
	}

	if err := sendRestart(b, userId); err != nil {
		return nil, err
	}
	
	return nil, fmt.Errorf("adEvent cache not found")
}

// Отправка в чат сообщения о повторной попытке.
func sendRestart(b *BotTelegram, userId int64) error {
	b.db.SetStepUser(userId, "start")
	botMsg := tgbotapi.NewMessage(userId, "К сожалению что то пошло не так. Выберите действие из меню /start повторно. 🥲")
	if _, err := b.bot.Send(botMsg); err != nil {
		return fmt.Errorf("error send message in sendRestartMessage: %w", err)
	}
	return nil
}

// TODO Очистка чата. Пока что не работает.
// func (b *BotTelegram) cleareAllChat(chatID int64) error {
// 	deleteMsg := tgbotapi.NewDeleteMessage(chatID, 0)
// 	if _, err := b.bot.Send(deleteMsg); err != nil {
// 		return fmt.Errorf("error cleare all chat: %w", err)
// 	}
// 	return nil
// }
