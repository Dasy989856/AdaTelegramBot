package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "Статистика:"
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

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventMenu: %w", err)
	}

	return nil
}

func cbqStatisticsBrief(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID
	// Получение данных и бд.

	// Сборка сообщения.
	text := "Выберите период:"
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

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreate: %w", err)
	}

	return nil
}

func cbqStatisticsBriefSelect(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных из БД.
	startDate, endDate := sdk.GetTimeRangeToday()
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
		return fmt.Errorf("error edit msg in cbqAdEventMenu: %w", err)
	}

	return nil
}
