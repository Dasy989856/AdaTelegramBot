package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик CallbackQuery ad_event.
func cbqHandlerAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery, cbqSteps []string) error {
	// Вывод меню.
	if len(cbqSteps) == 1 {
		if err := cbqAdEventMenu(b, cbq); err != nil {
			return err
		}
		return nil
	}

	// Обработчик типов ad_event.
	if len(cbqSteps) >= 2 {
		if err := cbqHandlerAdEventType(b, cbq, cbqSteps); err != nil {
			return err
		}
	}

	return nil
}

// Главное меню CallbackQuery ad_event.
func cbqAdEventMenu(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	textMsg := "Выберите тип события:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продажа рекламы.", "ad_event.sale"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Покупка рекламы.", "ad_event.buy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Взаимный пиар.", "ad_event.mutual"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageID, textMsg, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventMenu: %w", err)
	}

	return nil
}

// Обработчик CallbackQuery ad_event.type.
func cbqHandlerAdEventType(b *BotTelegram, cbq *tgbotapi.CallbackQuery, cbqSteps []string) error {
	userId := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID
	adEventType := strings.ToLower(cbqSteps[1])

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("02.01.2006 15:04:05.999"),
		Ready:     true,
		Type:      adEventType,
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.partner")

	var botMsg string
	switch adEventType {
	case "sale":
		botMsg = "Отлично! Теперь отправьтете мне ссылку на покупателя. Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot"
	case "buy":
		botMsg = "Отлично! Теперь отправьтете мне ссылку на продавца. Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot"
	case "mutual":
		botMsg = "Отлично! Теперь отправьтете мне ссылку на партнера по взаимному пиару. Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot"
	default:
		delete(b.cashAdEvents, userId)
		sendRestart(b, userId)
		return fmt.Errorf("unknow type adEvent. cbq: %v", cbqSteps)
	}

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageID, botMsg)); err != nil {
		return fmt.Errorf("error send botMsg from adEventSale: %w", err)
	}

	return nil
}
