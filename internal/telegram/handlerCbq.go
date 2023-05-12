package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parseCbq(cbq *tgbotapi.CallbackQuery) (path []string, data string, err error) {
	cbqDataSlice := strings.Split(cbq.Data, "?")
	if len(cbqDataSlice) < 1 {
		return nil, "", fmt.Errorf("len cbq incorrect. cbq: %s ", cbq.Data)
	}

	cbqPathSlice := strings.Split(cbqDataSlice[0], ".")
	if len(cbqPathSlice) < 1 {
		return nil, "", fmt.Errorf("len cbq path incorrect. cbq: %s ", cbq.Data)
	}

	switch len(cbqDataSlice) {
	case 1:
		return cbqPathSlice, "", nil
	case 2:
		return cbqPathSlice, cbqDataSlice[1], nil
	default:
		return nil, "", fmt.Errorf("len cbq incorrect. cbq: %s", cbq.Data)
	}
}

func (b *BotTelegram) handlerCbq(cbq *tgbotapi.CallbackQuery) error {
	fmt.Println("CBQ: " + cbq.Data)

	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch path[0] {
	case "start":
		if err := b.cmdStart(cbq.Message); err != nil {
			return err
		}
	case "ad_event":
		if err := handlerCbqAdEvent(b, cbq); err != nil {
			return err
		}
	case "statistics":
		if err := handlerCbqStatistics(b, cbq); err != nil {
			return err
		}
	case "help":
		fmt.Println("help NO WORK")
		// if err := handlerCbqAdEvent(b, cbq); err != nil {
		// 	return err
		// }
	}

	return nil
}

func handlerCbqAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch strings.Join(path, ".") {
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
	case "ad_event.view.any":
		if err := cbqAdEventViewAny(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.any.all":
		if err := cbqAdEventViewAnyAll(b, cbq); err != nil {
			return err
		}
		// case "ad_event.edit.":
		// 	if err := cbqAdEventEdit(b, cbq); err != nil {
		// 		return err
		// 	}
	}
	return nil
}

func handlerCbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch strings.Join(path, ".") {
	case "statistics":
		if err := cbqStatistics(b, cbq); err != nil {
			return err
		}
	case "statistics.brief":
		if err := cbqStatisticsBrief(b, cbq); err != nil {
			return err
		}
	case "statistics.brief.select":
		if err := cbqStatisticsBriefSelect(b, cbq); err != nil {
			return err
		}
	}
	return nil
}

// TODO no work
func handlerCbqHelp(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	switch cbq.Data {
	case "help.feature":
		// if err := cbqAdEvent(b, cbq); err != nil {
		// 	return err
		// }
	}
	return nil
}
