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
		errMsg := tgbotapi.NewMessage(msg.Chat.ID, `Неизвестная команда 🥲`)
		if _, err := b.bot.Send(errMsg); err != nil {
			return err
		}
		return nil
	}
}

// Команда /start
func (b *BotTelegram) cmdStart(msg *tgbotapi.Message) error {
	if err := b.db.DefaultUserCreation(msg.Chat.ID, msg.Chat.UserName, msg.Chat.FirstName); err != nil {
		return err
	}

	if err := b.sendMenuMsg(msg.Chat.ID); err != nil {
		return err
	}

	return nil
}

// Отправка стартового меню.
func (b *BotTelegram) sendMenuMsg(chatID int64) error {
	menuMsg := tgbotapi.NewMessage(chatID, "Возможности телеграмм бота Ада:")
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
	if _, err := b.bot.Send(menuMsg); err != nil {
		return fmt.Errorf("error send menuMsg: %w", err)
	}

	if err := b.db.SetStepUser(chatID, "start"); err != nil {
		return err
	}

	return nil
}