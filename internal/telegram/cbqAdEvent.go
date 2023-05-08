package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик CallbackQuery ad_event.
func handlerAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery, cbqSteps []string) error {
	userId := cbq.Message.Chat.ID
	adEventType := strings.ToLower(cbqSteps[1])

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      adEventType,
	}
	b.cashAdEvents[userId] = &adEvent

	switch adEventType {
	case "sale":
		if err := adEventSale(b, cbq); err != nil {
			return err
		}
	case "buy":
		if err := adEventBuy(b, cbq); err != nil {
			return err
		}
	case "barter":
		if err := adEventBarter(b, cbq); err != nil {
			return err
		}
	default:
		delete(b.cashAdEvents, userId)
		return fmt.Errorf("unknow cbq[1] step. cbq: %v", cbqSteps)
	}

	return nil
}

func adEventSale(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	chatID := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	b.db.SetStepUser(chatID, "ad_event.partner")

	botMsg := "Отлично! Теперь отправь мне ссылку на покупателя."
	_, err := b.bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, botMsg))
	if err != nil {
		return fmt.Errorf("error send botMsg from adEventSale: %w", err)
	}

	return nil
}

func adEventBuy(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	chatID := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	b.db.SetStepUser(chatID, "ad_event.partner")

	readyMsg := "Отлично! Теперь отправь мне ссылку на продавца."
	_, err := b.bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, readyMsg))
	if err != nil {
		return fmt.Errorf("error send readyMsg from addEventSale: %w", err)
	}

	return nil
}

func adEventBarter(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	chatID := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	b.db.SetStepUser(chatID, "ad_event.partner")

	readyMsg := "Отлично! Теперь отправь мне ссылку на партнера по бартеру."
	_, err := b.bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, readyMsg))
	if err != nil {
		return fmt.Errorf("error send readyMsg from addEventSale: %w", err)
	}

	return nil
}
