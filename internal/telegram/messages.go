package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик сообщений.
func (b *BotTelegram) handlerMessage(msg *tgbotapi.Message) error {
	step, err := b.db.GetStepUser(msg.Chat.ID)
	if err != nil {
		return err
	}

	// Сообщение обрабатываеются отталкиваясь от текущего шага пользователя.
	switch step {
	case "test":
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Тестовое сообщение.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
	default:
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Не получается обработать сообщение... 😔")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
	}

	return nil
}
