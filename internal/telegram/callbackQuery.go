package telegram

import (
	"fmt"
	"strings"

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
		if err := handlerAdEvent(b, update.CallbackQuery, cbqSteps); err != nil {
			return err
		}
	}

	return nil
}
