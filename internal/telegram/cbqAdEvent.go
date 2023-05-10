package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cbqAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "Управление событиями:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать событие.", "ad_event.create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Просмотреть события.", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "start"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventMenu: %w", err)
	}

	return nil
}

func cbqAdEventCreate(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "Выберите тип события:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продажа рекламы.", "ad_event.create.sale"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Покупка рекламы.", "ad_event.create.buy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Взаимный пиар.", "ad_event.create.mutual"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Кастомное.", "ad_event.create.castom"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "ad_event"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreate: %w", err)
	}

	return nil
}

func cbqAdEventCreateSale(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      models.TypeSale,
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := `
	Теперь требуется отправить мне ссылку на покупателя.
	Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot`

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageId, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateSale: %w", err)
	}

	return nil
}

func cbqAdEventCreateBuy(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      "buy",
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := `
	Теперь требуется отправить мне ссылку на продавца.
	Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot`

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageId, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateBuy: %w", err)
	}

	return nil
}

func cbqAdEventCreateMutual(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Ready:     true,
		Type:      models.TypeMutual,
	}
	b.cashAdEvents[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	botMsg := `
	Теперь требуется отправить мне ссылку на пратнера по взаимному пиару.
	Пример: @AdaTelegramBot или https://t.me/AdaTelegramBot`

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageText(userId, messageId, botMsg)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateMutal: %w", err)
	}

	return nil
}

func cbqAdEventCreateEnd(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID

	adEvent, err := getAdEventFromCash(b, userId)
	if err != nil {
		return err
	}

	// Валидация события.
	if !adEvent.AllData() {
		botMsg := tgbotapi.NewMessage(userId, "Были введены не все данные, попробуйте снова.")
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		if err := b.cmdStart(cbq.Message); err != nil {
			return err
		}
		return nil
	}

	// Сохранение события в бд.
	adEventId, err := b.db.AdEventCreation(adEvent)
	if err != nil {
		return err
	}

	botMsgString := fmt.Sprintf("Отлично! Событие добавлено! Индификатор события: %d.", adEventId)
	botMsg := tgbotapi.NewMessage(userId, botMsgString)
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	botMsg = tgbotapi.NewMessage(userId, "Для возврата в главное меню воспользуйтесь командой /start.")
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func cbqAdEventView(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "Выберите тип событий:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все.", "ad_event.view.all"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Проданная реклама.", "ad_event.view.sale"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Купленная реклама.", "ad_event.view.buy"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Взаимный пиар.", "ad_event.view.mutual"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Кастомное.", "ad_event.create.castom"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "ad_event"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventView: %w", err)
	}

	return nil
}

func cbqAdEventViewAll(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.messageId

	text := "Выберите фильтр для всех событий:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сегодня.", "ad_event.view.all.today"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Текущая неделя.", "ad_event.view.all.this_week"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Следующая неделя.", "ad_event.view.all.next_week"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Прошлая неделя.", "ad_event.view.all.last_week"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Кастомное.", "ad_event.create.castom"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "ad_event.view"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventView: %w", err)
	}

	return nil
}

func cbqAdEventViewAllToday(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение списка событий.

	text := "Выбранные события:"
	var keyboard tgbotapi.InlineKeyboardMarkup

	// Получение списка элементов из базы данных или API
	elements := getElements()

	// Создание списка кнопок
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, element := range elements {
		button := tgbotapi.NewInlineKeyboardButtonData(element.Name, element.ID)
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{button})
	}

	// Создание клавиатуры
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(keyboard)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please choose:")
	msg.ReplyMarkup = replyMarkup
	bot.Send(msg)
}

func button(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	query := update.CallbackQuery
	data := query.Data

	// Обработка нажатия на кнопку
	element := getElementByID(data)
	text := fmt.Sprintf("Вы выбрали элемент %s", element.Name)
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	bot.Send(msg)
}
