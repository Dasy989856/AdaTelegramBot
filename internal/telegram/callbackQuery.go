package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик CallbackQuery.
func (b *BotTelegram) handlerCallbackQuery(update *tgbotapi.Update) error {
	cbqSteps := strings.Split(update.CallbackQuery.Data, ".")
	if len(cbqSteps) < 2 {
		return fmt.Errorf("error len cbqSlice")
	}

	switch cbqSteps[0] {
	case "ad_event":
		if err := handlerAdEvent(b, update, cbqSteps); err != nil {
			return err
		}

	}

	return nil
}

// Обработчик CallbackQuery ad_event.
func handlerAdEvent(b *BotTelegram, update *tgbotapi.Update, cbqSteps []string) error {
	// Создание кэша ad события.
	b.cashAdEvents[update.Message.Chat.ID] = models.AdEvent{
		UserId:    update.Message.Chat.ID,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),

	}

	switch cbqSteps[1] {
	case "sale":
		if err := adEventSale(b, update.CallbackQuery); err != nil {
			return err
		}
	case "buy":
		if err := adEventBuy(b, update.CallbackQuery); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknow cbq[1] step. cbq: %v", cbqSteps)
	}

	return nil
}

func adEventSale(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	chatID := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	botMsg := "Отправьте мне URL покупателя."
	_, err := b.bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, botMsg))
	if err != nil {
		return fmt.Errorf("error send  botMsg from adEventSale: %w", err)
	}

	b.db.SetStepUser(chatID, "ad_event.partner")

	return nil
}

func adEventBuy(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	chatID := cbq.Message.Chat.ID
	messageID := cbq.Message.MessageID

	readyMsg := "Отлично, покупка рекламы добавлена!"
	_, err := b.bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, readyMsg))
	if err != nil {
		return fmt.Errorf("error send  readyMsg from addEventSale: %w", err)
	}

	return nil
}
