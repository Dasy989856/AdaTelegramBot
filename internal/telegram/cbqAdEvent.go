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
	messageID := cbq.Message.MessageID
	adEventType := strings.ToLower(cbqSteps[1])

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      adEventType,
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.partner")

	var botMsg string
	switch adEventType {
	case "sale":
		botMsg = "Отлично! Теперь отправь мне ссылку на покупателя. Пример: @buyer"
	case "buy":
		botMsg = "Отлично! Теперь отправь мне ссылку на продавца. Пример: @saler"
	case "barter":
		botMsg = "Отлично! Теперь отправь мне ссылку на партнера по бартеру. Пример: @barter"
	default:
		delete(b.cashAdEvents, userId)
		sendRestart(b, userId)
		return fmt.Errorf("unknow type adEvent. cbq: %v", cbqSteps)
	}

	_, err := b.bot.Send(tgbotapi.NewEditMessageText(userId, messageID, botMsg))
	if err != nil {
		return fmt.Errorf("error send botMsg from adEventSale: %w", err)
	}

	return nil
}
