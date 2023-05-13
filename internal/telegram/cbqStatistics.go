package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "📈 <b>Статистика:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Краткая статистика.", "statistics.brief"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Полная статистика.", "statistics.full"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "start"),
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
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeYesterday())),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeToday())),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeTomorrow())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущая неделя", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeLastWeek())),
			tgbotapi.NewInlineKeyboardButtonData("Текущая неделя", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeThisWeek())),
			tgbotapi.NewInlineKeyboardButtonData("Следующая неделя", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeNextWeek())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeLastMonth())),
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeThisMonth())),
			tgbotapi.NewInlineKeyboardButtonData("Следующий месяц", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeNextMonth())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий год", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeLastYear())),
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeThisYear())),
			tgbotapi.NewInlineKeyboardButtonData("Следующий год", "statistics.brief.select?"+sdk.ParseTimeToRangeDate(sdk.GetTimeRangeNextYear())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "statistics"),
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
	fmt.Println("DATA", data)

	dataSlice := strings.Split(data, ";")
	if len(dataSlice) != 2 {
		return fmt.Errorf("dataSlice incorrect. dataSlice: %v", dataSlice)
	}

	startDate, err := sdk.ParseDateToTime(dataSlice[0])
	if err != nil {
		return err
	}
	endDate, err := sdk.ParseDateToTime(dataSlice[1])
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
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "statistics.brief"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqStatisticsBriefSelect: %w", err)
	}

	return nil
}
