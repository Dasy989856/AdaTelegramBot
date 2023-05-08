package telegram

import (
	"abs-by-ammka-bot/internal/models"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

/*
		Полезные команды которые понадобятся.
Пример 1 (Добавление ссылок под строкой ввода сообщения):
	// numericKeyboard := tgbotapi.NewReplyKeyboard(
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("Добавить продажу рекламы."),
	// 		tgbotapi.NewKeyboardButton("Добавить покупку рекламы."),
	// 	))
	// newMsg.ReplyMarkup = numericKeyboard
*/

type BotTelegram struct {
	bot *tgbotapi.BotAPI
	db  models.TelegramBotDB
}

func NewBotTelegram(db models.TelegramBotDB) (*BotTelegram, error) {
	token := viper.GetString("token.telegram")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &BotTelegram{bot: bot, db: db}, nil
}

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
		// Команды.
		if update.Message != nil && update.Message.IsCommand() {
			if err := b.handlerCommand(update.Message); err != nil {
				return err
			}
			continue
		}

		// Сообщения.
		if update.Message != nil {
			if err := b.handlerMessage(update.Message); err != nil {
				return err
			}
			continue
		}

		// CallbackQuery.
		if update.CallbackQuery != nil {
			if err := b.handlerCallbackQuery(&update); err != nil {
				return err
			}
			continue
		}
	}

	return fmt.Errorf("handlerUpdates closed")
}

// Очистка чата.
func (b *BotTelegram) cleareAllChat(chatID int64) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, 0)
	if _, err := b.bot.Send(deleteMsg); err != nil {
		return fmt.Errorf("error cleare all chat: %w", err)
	}
	return nil
}