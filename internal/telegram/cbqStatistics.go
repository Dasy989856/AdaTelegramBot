package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Страница 1
	keyboard1 = tgbotapi.NewInlineKeyboardMarkup(
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Вчера", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeThisYear())),
		// 	tgbotapi.NewInlineKeyboardButtonData("Сегодня", "statistics"),
		// 	tgbotapi.NewInlineKeyboardButtonData("Год", "statistics"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("День", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeThisYear())),
			tgbotapi.NewInlineKeyboardButtonData("Месяц", "statistics"),
			tgbotapi.NewInlineKeyboardButtonData("Год", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ввести вручную", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	keyboard2 = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вреча", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("2", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("3", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("4", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("5", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("6", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("7", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("8", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("9", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("10", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("11", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("12", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("13", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("14", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "start"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	keyboard3 = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текщая неделя", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1.05 - 7.05", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("8.05 - 14.05", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("15.05 - 21.05", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("22.05 - 28.05", "statistics.brief"),
			tgbotapi.NewInlineKeyboardButtonData("29.05 - 31.05", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "start"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	
)

func cbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "📈 <b>Статистика:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Краткая статистика", "statistics.brief"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Полная статистика", "statistics.full"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqStatistics: %w", err)
	}

	return nil
}

func cbqStatisticsBrief(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeYesterday())),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeToday())),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeTomorrow())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущая неделя", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeThisWeek())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeThisMonth())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "statistics.brief.select?"+sdk.ParseTimesToRangeDate(sdk.GetTimeRangeThisYear())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "statistics.brief"),
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

func cbqStatisticsBriefSelect(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	_, data, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	dataSlice := strings.Split(data, ";")
	if len(dataSlice) != 2 {
		return fmt.Errorf("dataSlice incorrect. dataSlice: %v", dataSlice)
	}

	startDate, err := sdk.ParseUserDateToTime(dataSlice[0])
	if err != nil {
		return err
	}
	endDate, err := sdk.ParseUserDateToTime(dataSlice[1])
	if err != nil {
		return err
	}

	// Получение данных из БД.
	d, err := b.db.GetRangeDataForStatistics(userId, models.TypeAny, startDate, endDate)
	if err != nil {
		return err
	}

	// Создание краткой статистики.
	text := createStaticsBriefText(d)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqStatisticsBriefSelect: %w", err)
	}

	return nil
}
