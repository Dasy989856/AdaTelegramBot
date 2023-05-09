package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик команд.
func (b *BotTelegram) handlerCommand(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	switch msg.Command() {
	case "start":
		if err := b.cmdStart(msg); err != nil {
			return err
		}
		return nil
	default:
		botMsg := tgbotapi.NewMessage(userId, `Неизвестная команда 🥲`)
		if err := b.sendMessage(userId, botMsg); err != nil {
			return fmt.Errorf("error send unknow command error: %w", err)
		}
		return nil
	}
}

// Команда /start
func (b *BotTelegram) cmdStart(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	// Регистрация пользователя.
	if err := b.db.DefaultUserCreation(msg.Chat.ID, msg.Chat.UserName, msg.Chat.FirstName); err != nil {
		return err
	}

	// Создание стартового сообщение которое не удаляется если его нет.
	startMessageId, err := b.db.GetStartMessageId(userId)
	if err != nil {
		return err
	}

	// Отправка меню /start.
	if err := b.sendStartMenu(userId, startMessageId); err != nil {
		return err
	}

	// Очистка чата.
	if err := b.cleareAllChat(msg.Chat.ID); err != nil {
		return err
	}

	return nil
}

// Отправка стартового меню.
func (b *BotTelegram) sendStartMenu(userId int64, startMessageId int) error {
	// Установка шага пользователя.
	if err := b.db.SetStepUser(userId, "start"); err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление событиями.", "ad_event"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика.", "statistic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продажа рекламы.", "exchange.sale"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Покупка рекламы.", "exchange.buy"),
		),
	)

	// Создание startMessage если его нет.
	if startMessageId == 0 {
		menuMsg := tgbotapi.NewMessage(userId, "Возможности телеграмм бота Ада:")
		menuMsg.ReplyMarkup = keyboard

		startMessage, err := b.bot.Send(menuMsg)
		if err != nil {
			return fmt.Errorf("error send start menu: %w", err)
		}

		if err := b.db.UpdateStartMessageId(userId, startMessage.MessageID); err != nil {
			return err
		}
	} else {
		menuMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, startMessageId, "Возможности телеграмм бота Ада:", keyboard)

		if _, err := b.bot.Send(menuMsg); err != nil {
			return fmt.Errorf("error send start menu: %w", err)
		}
	}

	return nil
}
