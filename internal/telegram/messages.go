package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик сообщений.
func (b *BotTelegram) handlerMessage(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID
	step, err := b.db.GetStepUser(userId)
	if err != nil {
		return err
	}

	// Сообщение обрабатываеются отталкиваясь от текущего шага пользователя.
	switch step {
	case "ad_event.create.partner":
		if err := adEventPartner(b, msg); err != nil {
			return err
		}
	case "ad_event.create.chanel":
		if err := adEventChanel(b, msg); err != nil {
			return err
		}
	case "ad_event.create.price":
		if err := adEventPrice(b, msg); err != nil {
			return err
		}
	case "ad_event.create.date_posting":
		if err := adEventDatePosting(b, msg); err != nil {
			return err
		}
	case "ad_event.create.date_delete":
		if err := adEventDateDelete(b, msg); err != nil {
			return err
		}
	default:
		botMsg := tgbotapi.NewMessage(userId, "Не получается обработать сообщение... 😔")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}

	}

	return nil
}

func adEventPartner(b *BotTelegram, msg *tgbotapi.Message) error {
	// Example: https://t.me/nikname ; @nikname
	regxType1 := regexp.MustCompile(`^https:\/\/t\.me\/[a-zA-Z0-9_]+$`)
	regxType2 := regexp.MustCompile(`^@[a-zA-Z0-9_]+$`)
	userId := msg.Chat.ID

	if !regxType1.MatchString(msg.Text) && !regxType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, "Вы ввели некорректную ссылку на пользователя, попробуйте снова. Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot")
		if err := b.sendMessage(userId, botMsg); err != nil {
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
	b.db.SetStepUser(msg.Chat.ID, "ad_event.create.chanel")

	switch adEvent.Type {
	case "sale":
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправьте мне ссылку на рекламируемый Вами канал.")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case "buy":
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправьте мне ссылку на канал, в котором выйдет Ваша реклама.")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case "mutal":
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправьте мне ссылку на канал, в котором выйдет Ваша реклама.")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := sendRequestRestartMsg(b, userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	return nil
}

func adEventChanel(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную ссылку на канал, попробуйте снова. Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if models.RegxUrlType1.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.Channel = msg.Text

	if adEvent.Type == "mutual" {
		b.db.SetStepUser(msg.Chat.ID, "ad_event.create.date_posting")
	} else {
		b.db.SetStepUser(msg.Chat.ID, "ad_event.create.price")
	}

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправьте мне стоимость.")
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventPrice(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxPrice.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную стоимость, попробуйте снова. Пример: 1000")
		if err := b.sendMessage(userId, botMsg); err != nil {
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
	b.db.SetStepUser(msg.Chat.ID, "ad_event.create.date_posting")

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправьте дату размещения рекламы. Формат `22.08.2022 16:30`")
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventDatePosting(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную дату, попробуйте снова. Пример: 22.08.2022 16:30")
		if err := b.sendMessage(userId, botMsg); err != nil {
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
	b.db.SetStepUser(msg.Chat.ID, "ad_event.create.date_delete")

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправьте дату удаления рекламы. Формат `22.08.2022 16:30`")
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventDateDelete(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную дату, попробуйте снова. Пример: 22.08.2022 16:30")
		if err := b.sendMessage(userId, botMsg); err != nil {
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

	// Сравнение даты постинга и удаления.
	durationDatePosting, err := models.ParseDateToTime(adEvent.DatePosting)
	if err != nil {
		return fmt.Errorf("error parse durationDatePosting: %w", err)
	}

	durationDateDelete, err := models.ParseDateToTime(adEvent.DateDelete)
	if err != nil {
		return fmt.Errorf("error parse durationDateDelete: %w", err)
	}

	if durationDateDelete.Sub(*durationDatePosting) <= 0 {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели дату удаления поста меньше даты размещения поста, попробуйте снова.")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Дата удаления рекламы добавлена успешно!")
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	// Показать событие.
	{
		botMsgText := createAdEnentDescription(adEvent)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да.", "ad_event.create.end"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отменить.", "start"),
			),
		)
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, botMsgText)
		botMsg.ReplyMarkup = keyboard
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	}

	return nil
}