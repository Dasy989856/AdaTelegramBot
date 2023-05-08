package telegram

import (
	"fmt"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик сообщений.
func (b *BotTelegram) handlerMessage(msg *tgbotapi.Message) error {
	step, err := b.db.GetStepUser(msg.Chat.ID)
	if err != nil {
		return err
	}

	// Сообщение обрабатываеются отталкиваясь от текущего шага пользователя.
	switch step {
	case "ad_event.partner":
		if err := adEventPartner(b, msg); err != nil {
			return err
		}
	case "ad_event.chanel":
		if err := adEventChanel(b, msg); err != nil {
			return err
		}
	case "ad_event.price":
		if err := adEventPrice(b, msg); err != nil {
			return err
		}
	case "ad_event.date_posting":
		if err := adEventDatePosting(b, msg); err != nil {
			return err
		}
	case "ad_event.date_delete":
		if err := adEventDateDelete(b, msg); err != nil {
			return err
		}
	default:
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Не получается обработать сообщение... 😔")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
	}

	return nil
}

// step: ad_event.partner
func adEventPartner(b *BotTelegram, msg *tgbotapi.Message) error {
	// Example: https://t.me/nikname ; @nikname
	regxType1 := regexp.MustCompile(`https:\/\/t\.me\/[A-Za-z0-9]+`)
	regxType2 := regexp.MustCompile(`@[A-Za-z0-9]+`)
	userId := msg.Chat.ID

	if !regxType1.MatchString(msg.Text) && !regxType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, "Вы ввели некорректную ссылку на пользователя, попробуйте снова.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if regxType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.Partner = msg.Text
	b.db.SetStepUser(msg.Chat.ID, "ad_event.chanel")

	switch adEvent.Type {
	case "sale":
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь мне ссылку на рекламируемый канал.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
	case "buy":
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь мне ссылку на канал, в котором выйдет твоя реклама.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
	default:
		if err := sendRestart(b, userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	return nil
}

// step: ad_event.chanel
func adEventChanel(b *BotTelegram, msg *tgbotapi.Message) error {
	// Example: https://t.me/nikname ; @nikname
	regxType1 := regexp.MustCompile(`https:\/\/t\.me\/[A-Za-z0-9]+`)
	regxType2 := regexp.MustCompile(`@[A-Za-z0-9]+`)
	userId := msg.Chat.ID

	if !regxType1.MatchString(msg.Text) && !regxType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную ссылку на канал, попробуйте снова.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if regxType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.Channel = msg.Text
	b.db.SetStepUser(msg.Chat.ID, "ad_event.price")

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь мне стоимость.")
	if _, err := b.bot.Send(botMsg); err != nil {
		return err
	}

	return nil
}

// step: ad_event.price
func adEventPrice(b *BotTelegram, msg *tgbotapi.Message) error {
	regxPrice := regexp.MustCompile(`[0-9]+`)
	userId := msg.Chat.ID

	if !regxPrice.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную стоимость, попробуйте снова.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	price, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return err
	}

	adEvent.Price = price
	b.db.SetStepUser(msg.Chat.ID, "ad_event.date_posting")

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь дату размещения рекламы. Формат `2022-08-22 16:30`")
	if _, err := b.bot.Send(botMsg); err != nil {
		return err
	}

	return nil
}

// step: ad_event.date_posting
func adEventDatePosting(b *BotTelegram, msg *tgbotapi.Message) error {
	// Example: "2022-08-22 16:30"
	regxDate := regexp.MustCompile(`^(\d{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):[0-5][0-9]$`)
	userId := msg.Chat.ID

	if !regxDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную дату, попробуйте снова.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.DatePosting = msg.Text
	b.db.SetStepUser(msg.Chat.ID, "ad_event.date_delete")

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь дату удаления рекламы. Формат `2022-08-22 16:30`")
	if _, err := b.bot.Send(botMsg); err != nil {
		return err
	}

	return nil
}

// step: ad_event.date_delete
func adEventDateDelete(b *BotTelegram, msg *tgbotapi.Message) error {
	// Example: "2022-08-22 16:30"
	regxDate := regexp.MustCompile(`^(\d{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):[0-5][0-9]$`)
	userId := msg.Chat.ID

	if !regxDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную дату, попробуйте снова.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.DateDelete = msg.Text
	b.db.SetStepUser(msg.Chat.ID, "start")

	// Сохранение события в бд.
	if !adEvent.AllData() {
		return fmt.Errorf("adEvent have not full data")
	}

	adEventId, err := b.db.AdEventCreation(adEvent)
	if err != nil {
		return err
	}

	botMsgString := fmt.Sprintf("Отлично! Событие добавлено! ID события: %d.", adEventId)
	botMsg := tgbotapi.NewMessage(msg.Chat.ID, botMsgString)
	if _, err := b.bot.Send(botMsg); err != nil {
		return err
	}

	return nil
}
