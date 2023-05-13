package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"strconv"
	"time"

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
	case "ad_event.detele":
		if err := adEventDelete(b, msg); err != nil {
			return err
		}

	default:
		botMsg := tgbotapi.NewMessage(userId, "Не получается обработать сообщение... 😔")
		botMsg.ParseMode = tgbotapi.ModeHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("В главное меню.", "start"),
			),
		)
		botMsg.ReplyMarkup = keyboard

		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}

	}

	return nil
}

func adEventPartner(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы ввели некорректную ссылку на пользователя, попробуйте снова.
		<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if models.RegxUrlType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.Partner = msg.Text
	b.db.SetStepUser(userId, "ad_event.create.chanel")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на пользователя добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, "Теперь требуется отправить мне ссылку на рекламируемый Вами канал.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, "Теперь требуется отправить мне ссылку на канал, в котором выйдет Ваша реклама.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, "Теперь требуется отправить мне ссылку на канал, с которым будет взаимный пиар.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := b.sendRequestRestartMsg(userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	return nil
}

func adEventChanel(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы ввели некорректную ссылку на канал, попробуйте снова.
		<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if models.RegxUrlType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	adEvent.Channel = msg.Text
	b.db.SetStepUser(userId, "ad_event.create.price")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на канал добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, "Теперь требуется отправить стоимость рекламного поста.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, "Теперь требуется отправить стоимость рекламного поста.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить стоимость поста взаимного пиара.
		<b>Пример:</b> 0 (Если взаимный пиар был без доплаты)
		Можно указать <b>'-сумма'</b> если была доплата с Вашей стороны.
		Можно указать <b>'+сумма'</b> если доплатили Вам.`)
		// botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату размещения поста взаимного пиара.
		// <b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := b.sendRequestRestartMsg(userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	return nil
}

func adEventPrice(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxPrice.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы ввели некорректную стоимость, попробуйте снова.
		<b>Пример:</b> 1000`)
		botMsg.ParseMode = tgbotapi.ModeHTML
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
	b.db.SetStepUser(userId, "ad_event.create.date_posting")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Цена добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату и время размещения рекламного поста.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату и время размещения рекламного поста.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату и время размещения поста взаимного пиара.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := b.sendRequestRestartMsg(userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	return nil
}

func adEventDatePosting(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы ввели некорректную дату, попробуйте снова.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
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
	b.db.SetStepUser(userId, "ad_event.create.date_delete")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время размещения рекламы добавлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату и время удаления рекламного поста.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату и время удаления рекламного поста.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `Теперь требуется отправить дату и время удаления поста взаимного пиара.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := b.sendRequestRestartMsg(userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	return nil
}

func adEventDateDelete(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы ввели некорректную дату, попробуйте снова.
		<b>Пример:</b> `+sdk.ParseTimeToDate(time.Now()))
		botMsg.ParseMode = tgbotapi.ModeHTML
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

	// Сравнение даты размещения и удаления.
	durationDatePosting, err := sdk.ParseDateToTime(adEvent.DatePosting)
	if err != nil {
		return fmt.Errorf("error parse durationDatePosting: %w", err)
	}

	durationDateDelete, err := sdk.ParseDateToTime(adEvent.DateDelete)
	if err != nil {
		return fmt.Errorf("error parse durationDateDelete: %w", err)
	}

	if durationDateDelete.Sub(durationDatePosting) <= 0 {
		botMsg := tgbotapi.NewMessage(userId, "Вы ввели дату удаления поста меньше даты размещения поста, попробуйте снова.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Ответ.
	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, `Дата и время удаления рекламного поста добавлены успешно!`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, `Дата и время удаления рекламного поста добавлены успешно!`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `Дата и время удаления поста взаимного пиара добавлены успешно!`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := b.sendRequestRestartMsg(userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent")
	}

	if err := adEventCreateLastMessage(b, userId, adEvent); err != nil {
		return err
	}

	return nil
}

func adEventCreateLastMessage(b *BotTelegram, userId int64, adEvent *models.AdEvent) error {
	text := "<b>    ✍️ Вы хотите создать данное событие?</b>"
	text = text + createAdEventDescriptionText(adEvent)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да.", "ad_event.create.end"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить.", "start"),
		),
	)
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.ReplyMarkup = keyboard
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}
	return nil
}

// TODO no work
func adEventDelete(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxID.MatchString(msg.Text) {
		botMsgText := `
		Вы ввели некорректный ID. Попробуйте снова.
		<b>Пример:</b> 1.`
		botMsg := tgbotapi.NewMessage(userId, botMsgText)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEventId, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return fmt.Errorf("error passing adEventId: %w", err)
	}

	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}

	if err := adEventDeleteLastMessage(b, userId, adEvent); err != nil {
		return err
	}

	return nil
}

func adEventDeleteLastMessage(b *BotTelegram, userId int64, adEvent *models.AdEvent) error {
	text := "<b>    ✍️ Вы хотите создать данное событие?</b>"
	text = text + createAdEventDescriptionText(adEvent)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да.", "ad_event.create.end"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить.", "start"),
		),
	)
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.ReplyMarkup = keyboard
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}
	return nil
}
