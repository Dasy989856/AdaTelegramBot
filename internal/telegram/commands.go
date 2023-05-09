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
	if err := b.db.DefaultUserCreation(msg.Chat.ID, msg.Chat.UserName, msg.Chat.FirstName); err != nil {
		return err
	}

	if err := b.cleareAllChat(msg.Chat.ID); err != nil {
		return err
	}

	if err := b.sendMenuMsg(msg.Chat.ID); err != nil {
		return err
	}

	return nil
}

// Отправка стартового меню.
func (b *BotTelegram) sendMenuMsg(userId int64) error {
	menuMsg := tgbotapi.NewMessage(userId, "Возможности телеграмм бота Ада:")
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
	menuMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, menuMsg); err != nil {
		return fmt.Errorf("error send start menu: %w", err)
	}

	if err := b.db.SetStepUser(userId, "start"); err != nil {
		return err
	}

	return nil
}
