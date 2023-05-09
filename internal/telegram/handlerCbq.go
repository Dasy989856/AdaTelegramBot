package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик CallbackQuery.
func (b *BotTelegram) handlerCbq(update *tgbotapi.Update) error {
	fmt.Println(update.CallbackQuery.Message.Text)
	fmt.Println(update.CallbackQuery.Data)
	cbqSteps := strings.Split(update.CallbackQuery.Data, ".")
	if len(cbqSteps) <= 0 {
		return fmt.Errorf("error len cbqSteps")
	}

	switch cbqSteps[0] {
	case "start":
		if err := b.cmdStart(update.CallbackQuery.Message); err != nil {
			return err
		}
	case "ad_event":
		if err := cbqHandlerAdEvent(b, update.CallbackQuery); err != nil {
			return err
		}
	}

	return nil
}
