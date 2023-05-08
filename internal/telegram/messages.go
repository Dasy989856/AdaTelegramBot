package telegram

import (
	"fmt"
	"log"
	"regexp"

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
	case "ad_event.partner":
		if err := adEventPartner(b, msg); err != nil {
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

// step: ad_event.partner
func adEventPartner(b *BotTelegram, msg *tgbotapi.Message) error {
	// Example: https://t.me/nikname ; @nikname
	regxType1 := regexp.MustCompile(`https:\/\/t\.me\/[A-Za-z0-9]+`)
	regxType2 := regexp.MustCompile(`@[A-Za-z0-9]+/gm`)

	if !regxType1.MatchString(msg.Text) || !regxType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Вы ввели некорректную ссылку на пользователя, попробуйте снова.")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
	}

	// Приведение в единный тип.
	if regxType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	fmt.Println(msg.Text)

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, ok := b.cashAdEvents[msg.Chat.ID]
	if ok {
		adEvent.Partner = msg.Text
		b.db.SetStepUser(msg.Chat.ID, "ad_event.chanel")

		switch adEvent.Type {
		case "sale":
			botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь мне ссылку рекламируемый канал.")
			if _, err := b.bot.Send(botMsg); err != nil {
				return err
			}
		case "buy":
			botMsg := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Теперь отправь мне ссылку на канал, в котором выйдет твоя реклама.")
			if _, err := b.bot.Send(botMsg); err != nil {
				return err
			}
		}
		
	} else {
		botMsg := tgbotapi.NewMessage(msg.Chat.ID, "К сожалению процесс добавления придется начать повторно. 🥲")
		if _, err := b.bot.Send(botMsg); err != nil {
			return err
		}
		log.Println("error get cashAdEvents userId ", msg.Chat.ID)
		b.db.SetStepUser(msg.Chat.ID, "start")
	}

	return nil
}
