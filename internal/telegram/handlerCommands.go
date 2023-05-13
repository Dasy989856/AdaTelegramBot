package telegram

import (
	"fmt"
	"log"

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
		// botMsg.ParseMode = tgbotapi.ModeHTML
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
	if err := b.db.DefaultUserCreation(userId, msg.Chat.UserName, msg.Chat.FirstName); err != nil {
		return err
	}

	// Отправка рекламы.
	if viper.GetBool("ada_bot.ad_message") {
		if err := b.sendAdMessage(userId); err != nil {
			return err
		}
	} else {
		if err := b.db.UpdateAdMessageId(userId, 0); err != nil {
			return err
		}
	}

	// TODO Отправка информации.
	// if viper.GetBool("ada_bot.info_message") {
	// 	if err := b.sendAdMessage(userId); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	if err := b.db.UpdateAdMessageId(userId, 0); err != nil {
	// 		return err
	// 	}
	// }

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

	// Создание botMsg startMessage.
	text := `📓 <b>Возможности телеграмм бота:</b>`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление событиями", "ad_event"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Биржа рекламных интеграций.", "exchange"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Информация", "info"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тех. поддержка", "help"),
		),
	)
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.ReplyMarkup = keyboard

	// Отправка botMsg startMessage.
	newStartMessage, err := b.bot.Send(botMsg)
	if err != nil {
		return fmt.Errorf("error send new startMessage: %w", err)
	}

	// Сохранение startMessageId.
	if err := b.db.AddUserMessageId(userId, newStartMessage.MessageID); err != nil {
		return err
	}

	// Удаление если возможно старого startMessage.
	startMessageId, err := b.db.GetStartMessageId(userId)
	if err != nil {
		log.Println("b.db.GetStartmessageId startMenu error: ", err)
	}
	b.cleareMessage(userId, startMessageId)

	// Установка нового startMessage.
	if err := b.db.UpdateStartMessageId(userId, newStartMessage.MessageID); err != nil {
		return err
	}

	return nil
}

// Отправка adMessage.
func (b *BotTelegram) sendAdMessage(userId int64) error {
	// Создание botMsg adMessage.
	text := `📓 <b>💵 РЕКЛАМА </b>`
	// keyboard := tgbotapi.NewInlineKeyboardMarkup(
	// 	tgbotapi.NewInlineKeyboardRow(
	// 		tgbotapi.NewInlineKeyboardButtonData("Управление событиями.", "ad_event"),
	// 	),
	// )
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	// botMsg.ReplyMarkup = keyboard

	// Отправка botMsg adMessage.
	newAdMessage, err := b.bot.Send(botMsg)
	if err != nil {
		return fmt.Errorf("error send new adMessage: %w", err)
	}

	// Сохранение adMessageId.
	if err := b.db.AddUserMessageId(userId, newAdMessage.MessageID); err != nil {
		return err
	}

	// Удаление если возможно старого startMessage.
	adMessageId, err := b.db.GetAdMessageId(userId)
	if err != nil {
		log.Println("b.db.GetStartmessageId startMenu error: ", err)
	}
	b.cleareMessage(userId, adMessageId)

	// Установка нового adMessage.
	if err := b.db.UpdateAdMessageId(userId, newAdMessage.MessageID); err != nil {
		return err
	}

	return nil
}
