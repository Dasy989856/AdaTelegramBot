package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик сообщений.
func (b *BotTelegram) handlerMessage(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID
	fmt.Printf("Info %s: user=%s; MSG=%s;\n", time.Now().Format("2006-01-02 15:04:05.999"), msg.From.UserName, msg.Text)
	step, err := b.db.GetStepUser(userId)
	if err != nil {
		return err
	}

	// Сообщение обрабатываеются отталкиваясь от текущего шага пользователя.
	switch step {
	case "ad_event.create.partner":
		if err := adEventPartner(b, msg); err != nil {
			log.Println("error in adEventPartner: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.channel":
		if err := adEventChannel(b, msg); err != nil {
			log.Println("error in adEventChannel: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.price":
		if err := adEventPrice(b, msg); err != nil {
			log.Println("error in adEventPrice: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.date_posting":
		if err := adEventDatePosting(b, msg); err != nil {
			log.Println("error in adEventDatePosting: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.date_delete":
		if err := adEventDateDelete(b, msg); err != nil {
			log.Println("error in adEventDateDelete: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.partner":
		if err := adEventUpdatePartner(b, msg); err != nil {
			log.Println("error in adEventUpdatePartner: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.channel":
		if err := adEventUpdateChannel(b, msg); err != nil {
			log.Println("error in adEventUpdateChannel: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.price":
		if err := adEventUpdatePrice(b, msg); err != nil {
			log.Println("error in adEventUpdatePrice: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.date_posting":
		if err := adEventUpdateDatePosting(b, msg); err != nil {
			log.Println("error in adEventUpdateDatePosting: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.date_delete":
		if err := adEventUpdateDateDelete(b, msg); err != nil {
			log.Println("error in adEventUpdateDateDelete: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.arrival_of_subscribers":
		if err := adEventUpdateArrivalOfSubscribers(b, msg); err != nil {
			log.Println("error in adEventUpdateArrivalOfSubscribers: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	default:
		botMsg := tgbotapi.NewMessage(userId, "Не получается обработать сообщение... 😔")
		botMsg.ParseMode = tgbotapi.ModeHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
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
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную ссылку на пользователя, попробуйте снова.
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

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	adEvent.Partner = msg.Text
	b.db.SetStepUser(userId, "ad_event.create.channel")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на пользователя добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, "✍️ Теперь требуется отправить ссылку на рекламируемый Вами канал.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, "✍️ Теперь требуется отправить ссылку на канал, в котором выйдет Ваша реклама.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, "✍️ Теперь требуется отправить ссылку на канал, с которым будет взаимный пиар.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBarter:
		botMsg := tgbotapi.NewMessage(userId, "✍️ Теперь требуется отправить ссылку на канал/магазин партнера по бартеру.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	default:
		if err := b.sendRequestRestartMsg(userId); err != nil {
			return err
		}
		return fmt.Errorf("unknow type adEvent. typeEvent: %s", adEvent.Type)
	}

	return nil
}

func adEventChannel(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную ссылку на канал, попробуйте снова.
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

	adEvent, err := b.getAdEventCreatingCache(userId)
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
		botMsg := tgbotapi.NewMessage(userId, "✍️ Теперь требуется отправить стоимость рекламного поста.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, "✍️ Теперь требуется отправить стоимость рекламного поста.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить стоимость поста взаимного пиара.
		<b>Пример:</b> 0 (Если взаимный пиар был без доплаты)
		Можно указать <b>'-сумма'</b> если была доплата с Вашей стороны.
		Можно указать <b>'+сумма'</b> если доплатили Вам.`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBarter:
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить прибыль с бартера.
		<b>Пример:</b> 0 (Если считать прибыль не требуется)`)
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
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную стоимость, попробуйте снова.
		<b>Пример:</b> 1000`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
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

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	switch adEvent.Type {
	case models.TypeSale:
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить дату и время размещения рекламного поста.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить дату и время размещения рекламного поста.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить дату и время размещения поста взаимного пиара.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBarter:
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить дату и время размещения бартера.
		`+exampleDate)
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

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DatePosting = msg.Text

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время размещения рекламы добавлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	switch adEvent.Type {
	case models.TypeSale:
		b.db.SetStepUser(userId, "ad_event.create.date_delete")
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить дату и время удаления рекламного поста.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		// Отправка завершающего создания ad события сообщения.
		if err := adEventCreateLastMessage(b, userId, adEvent); err != nil {
			return err
		}
	case models.TypeMutual:
		b.db.SetStepUser(userId, "ad_event.create.date_delete")
		botMsg := tgbotapi.NewMessage(userId, `✍️ Теперь требуется отправить дату и время удаления поста взаимного пиара.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBarter:
		// Отправка завершающего создания ad события сообщения.
		if err := adEventCreateLastMessage(b, userId, adEvent); err != nil {
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

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DateDelete = msg.Text

	// Сравнение даты размещения и удаления.
	durationDatePosting, err := sdk.ParseUserDateToTime(adEvent.DatePosting)
	if err != nil {
		return fmt.Errorf("error parse durationDatePosting: %w", err)
	}

	durationDateDelete, err := sdk.ParseUserDateToTime(adEvent.DateDelete)
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
		botMsg := tgbotapi.NewMessage(userId, `🎉 <b>Дата и время удаления рекламы добавлены!</b>`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeBuy:
		botMsg := tgbotapi.NewMessage(userId, `🎉 <b>Дата и время удаления рекламы добавлены!</b>`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case models.TypeMutual:
		botMsg := tgbotapi.NewMessage(userId, `🎉 <b>Дата и время удаления поста взаимного пиара добавлены!</b>`)
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

	// Отправка завершающего создания ad события сообщения.
	if err := adEventCreateLastMessage(b, userId, adEvent); err != nil {
		return err
	}

	return nil
}

func adEventCreateLastMessage(b *BotTelegram, userId int64, adEvent *models.AdEvent) error {
	text := "<b>✍️ Вы хотите создать данное событие?</b>"
	text = text + createTextAdEventDescription(adEvent)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "ad_event.create.end"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", "start"),
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

func adEventUpdatePartner(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную ссылку на партнера, попробуйте снова.
		<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.Partner = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на партнера обновлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateChannel(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную ссылку на канал, попробуйте снова.
		<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.Channel = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на канал обновлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdatePrice(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxPrice.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную стоимость, попробуйте снова.
		<b>Пример:</b> 1000`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	price, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return err
	}

	adEvent.Price = price

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Стоимость обновлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateDatePosting(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DatePosting = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время размещения рекламы обновлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateDateDelete(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DateDelete = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время удаления рекламы обновлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateArrivalOfSubscribers(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxArrivalOfSubscribers.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректный приход подписчиков, попробуйте снова.
		<b>Пример:</b> 1000`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	arrivalOfSubscribers, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return err
	}
	adEvent.ArrivalOfSubscribers = arrivalOfSubscribers

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Приход подписчиков обновлен!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

