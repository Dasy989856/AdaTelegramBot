package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Start

func menuStart() (keyboard tgbotapi.InlineKeyboardMarkup, text string) {
	text = "Возможности телеграмм бота Ада:"
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление событиями.", "ad_event"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Биржа рекламных интеграций.", "exchange"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика.", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тех. поддержка.", "help"),
		),
	)
	return
}

// AdEvent

func menuAdEvent() (keyboard tgbotapi.InlineKeyboardMarkup, text string) {
	text = "Управление событиями:"
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать событие.", "ad_event.create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить событие.", "ad_event.delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Просмотреть события.", "statistic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад.", "start"),
		),
	)
	return
}

func menuAdEventCreate() (keyboard tgbotapi.InlineKeyboardMarkup, text string) {
	text = "Выберите тип события:"
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
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
	return
}

// Exchange

// Statistics

// Help
