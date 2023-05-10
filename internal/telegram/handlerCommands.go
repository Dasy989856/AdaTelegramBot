package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
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

	// Отправка рекламы.
	if err := b.sendAdMessage(userId); err != nil {
		return err
	}

	// Отправка меню /start.
	if err := b.sendStartMessage(userId); err != nil {
		return err
	}

	// Очистка чата.
	if err := b.cleareAllChat(userId); err != nil {
		return err
	}

	return nil
}

// Отправка startMessage.
func (b *BotTelegram) sendStartMessage(userId int64) error {
	// Установка шага пользователя.
	if err := b.db.SetStepUser(userId, "start"); err != nil {
		return err
	}

	text := `📓 Возможности телеграмм бота Ада:`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление событиями.", "ad_event"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Биржа рекламных интеграций.", "exchange"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика.", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тех. поддержка.", "help"),
		),
	)

	// Создание/получение startMessage.
	startmessageId, err := b.db.GetStartmessageId(userId)
	if err != nil {
		if err := updateStartMessage(b, userId, startmessageId, keyboard, text); err != nil {
			return err
		}
		return nil
	}

	// Изменение startMenu.
	if err := editMessageReplyMarkup(b, userId, startmessageId, keyboard, text); err != nil {
		// Попытка создать новое старт меню.
		if err := updateStartMessage(b, userId, startmessageId, keyboard, text); err != nil {
			return err
		}
	}

	return nil
}

// Обновление startMessage.
func updateStartMessage(b *BotTelegram, userId int64, startmessageId int, keyboard tgbotapi.InlineKeyboardMarkup, text string) error {
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ReplyMarkup = keyboard

	// Создание нового startMessage.
	newStartMessage, err := b.bot.Send(botMsg)
	if err != nil {
		return fmt.Errorf("error send new startMessage: %w", err)
	}

	// Удаление если возможно старого startMessage.
	b.cleareMessage(userId, startmessageId)

	// Установка нового startMessage.
	if err := b.db.UpdateStartmessageId(userId, newStartMessage.MessageID); err != nil {
		return err
	}

	return nil
}

// Отправка adMessage.
func (b *BotTelegram) sendAdMessage(userId int64) error {
	text := `📓 Реклама: Присоединяйся к прекрасному каналу @ammka22`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Спасибо!", "ad"),
		),
	)

	// Создание/получение adMessage.
	admessageId, err := b.db.GetAdmessageId(userId)
	if err != nil {
		if err := updateAdMessage(b, userId, admessageId, keyboard, text); err != nil {
			return err
		}
		return nil
	}

	// Изменение adMessage.
	if err := editMessageReplyMarkup(b, userId, admessageId, keyboard, text); err != nil {
		// Попытка создать новое старт меню.
		if err := updateAdMessage(b, userId, admessageId, keyboard, text); err != nil {
			return err
		}
	}

	return nil
}

// Обновление adMessage.
func updateAdMessage(b *BotTelegram, userId int64, admessageId int, keyboard tgbotapi.InlineKeyboardMarkup, text string) error {
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ReplyMarkup = keyboard

	if viper.GetBool("ada_bot.adMessage") {
		// Создание нового adMessage.
		newAdMessage, err := b.bot.Send(botMsg)
		if err != nil {
			return fmt.Errorf("error send new adMessage: %w", err)
		}

		// Установка нового adMessage.
		if err := b.db.UpdateAdmessageId(userId, newAdMessage.MessageID); err != nil {
			return err
		}
	}

	// Удаление если возможно старого adMessage.
	b.cleareMessage(userId, admessageId)

	return nil
}
