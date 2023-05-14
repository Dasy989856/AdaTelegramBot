package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func cbqAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "<b>📓 Управление событиями:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать событие", "ad_event.create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Просмотреть события", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventMenu: %w", err)
	}

	return nil
}

func cbqAdEventCreate(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "<b>Выберите тип события:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продажа рекламы", "ad_event.create.sale"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Покупка рекламы", "ad_event.create.buy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Взаимный пиар", "ad_event.create.mutual"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Кастомное.", "ad_event.create.custom"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
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
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := `Теперь требуется отправить мне ссылку на покупателя.
	<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
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
		Type:      models.TypeBuy,
	}
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := `Теперь требуется отправить мне ссылку на продавца.
	<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
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
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := `Теперь требуется отправить мне ссылку на пратнера по взаимному пиару.
	<b>Пример:</b> @AdaTelegramBot или https://t.me/AdaTelegramBot`
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateMutual: %w", err)
	}

	return nil
}

func cbqAdEventCreateEnd(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	// Валидация события.
	if !fullDataAdEvent(adEvent) {
		botMsg := tgbotapi.NewMessage(userId, "Были введены не все данные, что бы повторить воспользуйтесь командой <b>/start</b>")
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
		return nil
	}

	// Сохранение события в бд.
	adEventId, err := b.db.AdEventCreation(adEvent)
	if err != nil {
		return err
	}

	// Отправка сообщения.
	text := fmt.Sprintf("<b>🎊 Отлично! Событие добавлено! Индификатор события: %d.</b>", adEventId)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateEnd: %w", err)
	}

	// Очистка кэша.
	delete(b.adEventCreatingCache, userId)
	return nil
}

func cbqAdEventView(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "Выберите тип событий:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все типы", "ad_event.view.any"),
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
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventView: %w", err)
	}

	return nil
}

func cbqAdEventViewAny(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeYesterday())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeToday())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeTomorrow())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData("Предыдущая неделя", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeLastWeek())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Текущая неделя", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeThisWeek())+";any;1"),
			// tgbotapi.NewInlineKeyboardButtonData("Следующая неделя", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeNextWeek())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeLastMonth())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeThisMonth())+";any;1"),
			// tgbotapi.NewInlineKeyboardButtonData("Следующий месяц", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeNextMonth())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData("Предыдущий год", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeLastYear())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeThisYear())+";any;1"),
			// tgbotapi.NewInlineKeyboardButtonData("Следующий год", "ad_event.view.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeNextYear())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAny: %w", err)
	}

	return nil
}

func cbqAdEventViewSelect(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID
	lenRow := viper.GetInt("ada_bot.len_dinamic_row")

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных.
	data, err := parseDataAdEventView(cbqData)
	if err != nil {
		return err
	}
	fmt.Println(data)

	// Проврека данных.
	if _, ok := b.adEventCreatingCache[userId]; !ok {
		// Получение данных из БД.
		adEvents, err := b.db.GetRangeAdEventsOfUser(userId, data.TypeAdEvent, data.StartDate, data.EndDate)
		if err != nil {
			return err
		}

		// Разбиение событий и сохранение в кэш.
		b.adEventsCache[userId] = sdk.ChunkSlice(adEvents, lenRow)
	}

	// Отображение событий.
	text, keyboard, err := createTextAndKeyboardForAdEventView(b, userId, data)
	if err != nil {
		return err
	}

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.DisableWebPagePreview = true
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAnyAll: %w", err)
	}

	return nil
}

func parseDataAdEventView(cbqData string) (data *models.CbqDataForCbqAdEventViewSelect, err error) {
	// ad_event.view.any.select?14.05.2023 00:00;14.05.2023 23:59;any;1
	dataSlice := strings.Split(cbqData, ";")
	if len(dataSlice) != 4 {
		return nil, fmt.Errorf("dataSlice incorrect. dataSlice: %v", dataSlice)
	}
	data = new(models.CbqDataForCbqAdEventViewSelect)

	data.StartDate, err = sdk.ParseUserDateToTime(dataSlice[0])
	if err != nil {
		return nil, err
	}

	data.EndDate, err = sdk.ParseUserDateToTime(dataSlice[1])
	if err != nil {
		return nil, err
	}

	data.TypeAdEvent = models.TypeAdEvent(dataSlice[2])
	if err != nil {
		return nil, err
	}

	pageForDisplay, err := strconv.Atoi(dataSlice[3])
	if err != nil {
		return nil, fmt.Errorf("error pasge PageForDisplay: %w", err)
	}
	data.PageForDisplay = pageForDisplay

	return data, nil
}

func createTextAndKeyboardForAdEventView(b *BotTelegram, userId int64, data *models.CbqDataForCbqAdEventViewSelect) (string, tgbotapi.InlineKeyboardMarkup, error) {
	lenRow := viper.GetInt("ada_bot.len_dinamic_row")

	adEvents, err := b.getAdEventsCache(userId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}

	if len(adEvents) == 0 {
		text := `<b>🗓 Нет событий.</b>`
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view.any"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
			),
		)
		
		return text, keyboard, nil
	}

	// Создание кнопок.
	text := fmt.Sprintf(`<b>🗓 Выбранные события. Страница %d/%d. </b>
	✔️ Выберите номер события на <b>кнопках ниже</b> для редактирования события.
	`, data.PageForDisplay, len(adEvents))

	bufButtonRows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	bufButtonRow := make([]tgbotapi.InlineKeyboardButton, 0, lenRow)
	for i, adEvent := range adEvents[data.PageForDisplay-1] {
		buttonId := fmt.Sprintf("%d", i+1)
		buttonData := fmt.Sprintf("adEventId%d", adEvent.Id)
		button := tgbotapi.NewInlineKeyboardButtonData(buttonId, buttonData)
		bufButtonRow = append(bufButtonRow, button)

		text = text + fmt.Sprintf("\n<b>    ✍️ Событие № %s</b>:", buttonId)
		text = text + createTextAdEventDescription(&adEvent)
	}
	bufButtonRows = append(bufButtonRows, bufButtonRow)

	if len(adEvents) > 1 {
		pageRow := createPageRowForViewAdEvent(data, len(adEvents))
		bufButtonRows = append(bufButtonRows, pageRow)
	}

	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view.any"),
	)
	bufButtonRows = append(bufButtonRows, backRow)

	startMenuRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
	)
	bufButtonRows = append(bufButtonRows, startMenuRow)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(bufButtonRows...)

	return text, keyboard, nil
}

func createPageRowForViewAdEvent(data *models.CbqDataForCbqAdEventViewSelect, maxPage int) []tgbotapi.InlineKeyboardButton {
	buffButton := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	if data.PageForDisplay-1 > 0 {
		textDataPreviousPage := fmt.Sprintf("ad_event.view.select?%s;%s;%d",
			sdk.ParseTimeToRangeDate(data.StartDate, data.EndDate), data.TypeAdEvent, data.PageForDisplay-1)
		buffButton = append(buffButton, tgbotapi.NewInlineKeyboardButtonData("<<", textDataPreviousPage))
	}

	if data.PageForDisplay+1 <= maxPage {
		textDataNextPage := fmt.Sprintf("ad_event.view.select?%s;%s;%d",
			sdk.ParseTimeToRangeDate(data.StartDate, data.EndDate), data.TypeAdEvent, data.PageForDisplay+1)
		buffButton = append(buffButton, tgbotapi.NewInlineKeyboardButtonData(">>", textDataNextPage))
	}

	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view.any"),
	)

	return tgbotapi.NewInlineKeyboardRow(buffButton...)
}
