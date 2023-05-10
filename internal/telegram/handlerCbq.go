package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotTelegram) handlerCbq(cbq *tgbotapi.CallbackQuery) error {
	cbqSteps := strings.Split(cbq.Data, ".")
	if len(cbqSteps) < 1 {
		return fmt.Errorf("error len cbqSteps. cbqData: %s", cbq.Data)
	}

	fmt.Println("CBQ: " + cbq.Data)
	switch cbqSteps[0] {
	case "start":
		if err := b.cmdStart(cbq.Message); err != nil {
			return err
		}
	case "ad_event":
		if err := handlerCbqAdEvent(b, cbq); err != nil {
			return err
		}
	case "statistics":
		fmt.Println("statistics NO WORK")
		// if err := handlerCbqAdEvent(b, cbq); err != nil {
		// 	return err
		// }
	case "help":
		fmt.Println("help NO WORK")
		// if err := handlerCbqAdEvent(b, cbq); err != nil {
		// 	return err
		// }
	}

	return nil
}

func handlerCbqAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
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
	case "ad_event.create.end":
		if err := cbqAdEventCreateEnd(b, cbq); err != nil {
			return err
		}
	case "ad_event.view":
		if err := cbqAdEventView(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.all":
		if err := cbqAdEventViewAll(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.all.today":
		if err := cbqAdEventViewAllToday(b, cbq); err != nil {
			return err
		}
	}
	return nil
}

// TODO no work
func handlerCbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	switch cbq.Data {
	case "statistics.brief":
		// TODO no work
		if err := cbqAdEventView(b, cbq); err != nil {
			return err
		}
	case "statistics.full":
		// TODO no work
		if err := cbqAdEventViewAll(b, cbq); err != nil {
			return err
		}
	}
	return nil
}

// TODO no work
func handlerCbqHelp(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	switch cbq.Data {
	case "help.feature":
		if err := cbqAdEvent(b, cbq); err != nil {
			return err
		}
	}
	return nil
}
