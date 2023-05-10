package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик CallbackQuery.
func (b *BotTelegram) handlerCbq(cbq *tgbotapi.CallbackQuery) error {
	fmt.Println(cbq.Message.Text)
	fmt.Println(cbq.Data)

	// Определение типа cbq.
	cbqPart := strings.Split(cbq.Data, ":")
	if len(cbqPart) < 1 || len(cbqPart) > 2 {
		return fmt.Errorf("error len cbqPart in static, cbqData: %s", cbq.Data)
	}

	// Статический cbq.
	if len(cbqPart) == 1 {
		if err := handlerCbqStatic(b, cbq); err != nil {
			return err
		}
	}

	// Динамический cbq.
	// if len(cbqPart) == 2 {
	//  TODO
	// }

	return nil
}

func handlerCbqStatic(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	cbqSteps := strings.Split(cbq.Data, ".")
	if len(cbqSteps) < 1 {
		return fmt.Errorf("error len cbqSteps in static, cbqData: %s", cbq.Data)
	}

	switch cbqSteps[0] {
	case "start":
		if err := b.cmdStart(cbq.Message); err != nil {
			return err
		}
	case "ad_event":
		if err := cbqHandlerAdEvent(b, cbq); err != nil {
			return err
		}
	}

	return nil
}