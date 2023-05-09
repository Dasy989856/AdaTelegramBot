package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик CallbackQuery ad_event.
func cbqHandlerAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	switch cbq.Data {
	case "ad_event":
		if err := cbqAdEvent(b, cbq); err != nil {
			return err
		}
	case "ad_event.create":
		if err := cbqAdEventCreate(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.sale":
		if err := cbqAdEventCreateSale(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.buy":
		if err := cbqAdEventCreateBuy(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.mutual":
		if err := cbqAdEventCreateMutual(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.custom":
		if err := cbqAdEventCreateCustom(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.end":
		if err := cbqAdEventCreateEnd(b, cbq); err != nil {
			return err
		}
	}

	return nil
}

// Главное меню CallbackQuery ad_event.
func cbqAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	keyboard, text := menuAdEvent()

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageID, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventMenu: %w", err)
	}

	return nil
}

func cbqAdEventCreate(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	keyboard, text := menuAdEventCreate()

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageID, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreate: %w", err)
	}

	return nil
}

func cbqAdEventCreateSale(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      "sale",
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := `
	Отправьте мне ссылку на покупателя.
	Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot`

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageID, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateSale: %w", err)
	}

	return nil
}

func cbqAdEventCreateBuy(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      "buy",
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := `
	Отправьте мне ссылку на продавца.
	Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot`

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageID, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateBuy: %w", err)
	}

	return nil
}

func cbqAdEventCreateMutual(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      "mutual",
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := `
	Отправьте мне ссылку на продавца.
	Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot`

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageID, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateMutal: %w", err)
	}

	return nil
}

// TODO Заглушка.
func cbqAdEventCreateCustom(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      "custom",
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := "Отлично! Теперь отправьте мне ссылку на продавца. Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot"

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageID, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateCustom: %w", err)
	}

	return nil
}


func cbqAdEventCreateEnd(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID

	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	// Валидация события.
	if !adEvent.AllData() {
		botMsg := tgbotapi.NewMessage(userId, "Были введены не все данные, попробуйте снова.")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		if err := b.cmdStart(cbq.Message); err != nil {
			return err
		}
		return nil
	}

	// Сохранение события в бд.
	adEventId, err := b.db.AdEventCreation(adEvent)
	if err != nil {
		return err
	}

	botMsgString := fmt.Sprintf("Отлично! Событие добавлено! Индификатор события: %d.", adEventId)
	botMsg := tgbotapi.NewMessage(userId, botMsgString)
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

// Создание описания ad события.
func createAdEnentDescription(a *models.AdEvent) (descriptionAdEvent string) {
	switch a.Type {
	case "sale":
		descriptionAdEvent = fmt.Sprintf(`
		Ваше событие:
		- Покупатель: %s,
		- Канал покупателя: %s,
		- Цена продажи: %d, 
		- Дата постинга рекламы: %s,
		- Дата удаления рекламы: %s
		
		Сохранить событие?`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case "buy":
		descriptionAdEvent = fmt.Sprintf(`
		Ваше событие:
		- Продавец: %s,
		- Канал продавца: %s,
		- Цена продажи: %d, 
		- Дата постинга рекламы: %s,
		- Дата удаления рекламы: %s

		Сохранить событие?`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case "mutal":
		descriptionAdEvent = fmt.Sprintf(`
		Ваше событие:
		- Партнер по ВП: %s,
		- Канал партнера по ВП: %s,
		- Дата постинга рекламы: %s,
		- Дата удаления рекламы: %s

		Сохранить событие?`, a.Partner, a.Channel, a.DatePosting, a.DateDelete)
	}

	return descriptionAdEvent
}