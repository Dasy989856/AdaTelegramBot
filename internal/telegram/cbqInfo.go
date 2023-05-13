package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cbqInfo(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := `<b>ℹ️ Информация о проекте ℹ️</b>
	Данный бот был создан для помощи администраторам телеграмм каналов или их менеджеров.

	<b>🌟 Что умеет бот на текущий момент:</b>
	🔸Создавать различные события:
		- Покупка рекламы
		- Продажа рекламы
		- Взаимный пиар
	🔸Напоминать о предстоящих событиях
	🔸Отслеживать финансовые показатели
	
	<b>🌟 Функционал который появится:</b>
	🔸Отправка писем с отчетностю на почту
	🔸Создание рекламной биржи 

	<b>🌟 Что такое рекламная биржа?</b>
	<b>Рекламная биржа</b> - платформа, которая позволяет продавать и покупать рекламные интеграции. На ней можно:
	🔸 Размещать объявления о продаже или покупки рекламных интеграций
	🔸 Размещать и искать объявления о взаимном продвижении
	🔸 Заключать безопасные сделки
	🔸 Общаться с партнерами
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqHelpInfo: %w", err)
	}

	return nil
}
