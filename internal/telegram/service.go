package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"log"
	"strings"

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
	bot.Debug = false

	return &BotTelegram{bot: bot, db: db, cashAdEvents: make(map[int64]*models.AdEvent)}, nil
}

// Инициализация канала событий.
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
			if err := b.db.AddUsermessageId(update.Message.Chat.ID,
				update.Message.MessageID); err != nil {
				return err
			}

			if err := b.handlerCommand(update.Message); err != nil {
				log.Println(err)
			}
			continue
		}

		// Обработка сообщений.
		if update.Message != nil {
			if err := b.db.AddUsermessageId(update.Message.Chat.ID,
				update.Message.MessageID); err != nil {
				return err
			}

			if err := b.handlerMessage(update.Message); err != nil {
				log.Println(err)
			}
			continue
		}

		// Обработка CallbackQuery.
		if update.CallbackQuery != nil {
			if err := b.db.AddUsermessageId(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID); err != nil {
				return err
			}

			if err := b.handlerCbq(update.CallbackQuery); err != nil {
				log.Println(err)
			}
			continue
		}
	}

	return fmt.Errorf("updates chanel closed")
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

// Получение хэша ad события.
func getAdEventFromCash(b *BotTelegram, userId int64) (*models.AdEvent, error) {
	adEvent, ok := b.cashAdEvents[userId]
	if ok {
		return adEvent, nil
	}

	if err := sendRequestRestartMsg(b, userId); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("adEvent cache not found")
}

// Отправка в чат сообщения о повторной попытке.
func sendRequestRestartMsg(b *BotTelegram, userId int64) error {
	b.db.SetStepUser(userId, "start")
	botMsg := tgbotapi.NewMessage(userId, "К сожалению что то пошло не так. Выберите действие из меню /start повторно. 🥲")
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error send message in sendRestartMessage: %w", err)
	}
	return nil
}

// Очистка сообщения.
func (b *BotTelegram) cleareMessage(userId int64, messageId int) error {
	if err := b.db.DeleteUsermessageId(messageId); err != nil {
		return err
	}

	deleteMsg := tgbotapi.NewDeleteMessage(userId, messageId)
	if _, err := b.bot.Send(deleteMsg); err != nil {
		return fmt.Errorf("error cleare messageId%d: %w", messageId, err)
	}
	return nil
}

// Очистка чата.
func (b *BotTelegram) cleareAllChat(userId int64) error {
	startmessageId, err := b.db.GetStartmessageId(userId)
	if err != nil {
		return err
	}

	messageIds, err := b.db.GetUsermessageIds(userId)
	if err != nil {
		return err
	}

	// Удаление всех сообщений кроме startMessage.
	for _, messageId := range messageIds {
		if startmessageId == messageId {
			continue
		}
		b.cleareMessage(userId, messageId)
	}

	return nil
}

// Отправка сообщения пользователю.
func (b *BotTelegram) sendMessage(userId int64, c tgbotapi.Chattable) error {
	botMsg, err := b.bot.Send(c)
	if err != nil {
		return err
	}

	if err := b.db.AddUsermessageId(userId, botMsg.MessageID); err != nil {
		return err
	}

	return nil
}

// Изменение сообщения.
func editMessage(b *BotTelegram, userId int64, startmessageId int, keyboard tgbotapi.InlineKeyboardMarkup, text string) error {
	menuMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, startmessageId, text, keyboard)
	if _, err := b.bot.Send(menuMsg); err != nil {
		return fmt.Errorf("error edit startMenu: %w", err)
	}
	return nil
}

// Создание текст-описания ad события.
func createAdEventDescription(a *models.AdEvent) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		Ваше событие:
		- Покупатель: %s,
		- Канал покупателя: %s,
		- Цена продажи: %d, 
		- Дата постинга рекламы: %s,
		- Дата удаления рекламы: %s
		
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case "buy":
		descriptionAdEvent = fmt.Sprintf(`
		Ваше событие:
		- Продавец: %s,
		- Канал продавца: %s,
		- Цена продажи: %d, 
		- Дата постинга рекламы: %s,
		- Дата удаления рекламы: %s
		
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case "mutal":
		descriptionAdEvent = fmt.Sprintf(`
		Ваше событие:
		- Партнер по ВП: %s,
		- Канал партнера по ВП: %s,
		- Дата постинга рекламы: %s,
		- Дата удаления рекламы: %s
		
		`, a.Partner, a.Channel, a.DatePosting, a.DateDelete)
	}

	return descriptionAdEvent
}

// TODO not used functions vvv

// Проверка cbq на динамические данные. Возвращает данные и идификатор успешности.
func cbqGetData(cbq *tgbotapi.CallbackQuery) (data string, ok bool) {
	cbqPart := strings.Split(cbq.Data, ":")
	//Dinamic type
	if len(cbqPart) == 2 {
		return cbqPart[1], true
	}
	return "", false
}
