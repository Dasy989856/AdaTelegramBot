package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик команд.
func (b *BotTelegram) handlerCommand(msg *tgbotapi.Message) error {
	switch msg.Command() {
	case "start":
		if err := b.cmdStart(msg); err != nil {
			return err
		}
		return nil
	default:
		if err := b.handlerMessage(msg); err != nil {
			return err
		}
		// botMsg := tgbotapi.NewMessage(userId, `Неизвестная команда 🥲`)
		// if err := b.sendMessage(userId, botMsg); err != nil {
		// 	return fmt.Errorf("error send unknow command error: %w", err)
		// }
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

	// Отправка меню /start.
	if err := b.sendStartMenu(userId); err != nil {
		return err
	}

	// Очистка чата.
	if err := b.cleareAllChat(userId); err != nil {
		return err
	}

	return nil
}

// Отправка стартового меню.
func (b *BotTelegram) sendStartMenu(userId int64) error {
	// Установка шага пользователя.
	if err := b.db.SetStepUser(userId, "start"); err != nil {
		return err
	}

	text := "Возможности телеграмм бота Ада:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление событиями.", "ad_event"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Биржа рекламных интеграций.", "exchange"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистики.", "statistic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тех. поддержка.", "help"),
		),
	)

	// Создание/получение startMessage которое не удаляется.
	startMessageId, err := b.db.GetStartMessageId(userId)
	if err != nil {
		if err := updateStartMenu(b, userId, startMessageId, keyboard, text); err != nil {
			return err
		}
		return nil
	}

	// Изменение startMenu.
	if err := editStartMenu(b, userId, startMessageId, keyboard, text); err != nil {
		// Попытка создать новое старт меню.
		if err := updateStartMenu(b, userId, startMessageId, keyboard, text); err != nil {
			return err
		}
	}

	return nil
}

// Обновление startMenu.
func updateStartMenu(b *BotTelegram, userId int64, startMessageId int, keyboard tgbotapi.InlineKeyboardMarkup, text string) error {
	menuMsg := tgbotapi.NewMessage(userId, text)
	menuMsg.ReplyMarkup = keyboard

	// Создание нового startMessage.
	newStartMessage, err := b.bot.Send(menuMsg)
	if err != nil {
		return fmt.Errorf("error send startMenu: %w", err)
	}

	// Удаление если возможно старого startMessage.
	b.cleareMessage(userId, startMessageId)

	// Установка нового startMessage.
	if err := b.db.UpdateStartMessageId(userId, newStartMessage.MessageID); err != nil {
		return err
	}

	return nil
}

func editStartMenu(b *BotTelegram, userId int64, startMessageId int, keyboard tgbotapi.InlineKeyboardMarkup, text string) error {
	menuMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, startMessageId, text, keyboard)
	if _, err := b.bot.Send(menuMsg); err != nil {
		return fmt.Errorf("error edit startMenu: %w", err)
	}
	return nil
}
