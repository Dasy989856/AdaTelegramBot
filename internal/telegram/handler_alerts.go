package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"log"
	"math"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Оповещение о предстоящих событиях.
func (b *BotTelegram) handlerAlerts() (err error) {
	var cashAdEvents []models.AdEvent
	for {
		time.Sleep(15 * time.Second)

		timeStart, _ := sdk.GetTimeRangeToday()
		_, timeEnd := sdk.GetTimeRangeTomorrow()
		cashAdEvents, err = b.db.GetRangeAdEvents(models.TypeAny, timeStart, timeEnd)
		if err != nil {
			return fmt.Errorf("handlerAlerts: error GetRangeAdEvents: %w", err)
		}

		for _, aE := range cashAdEvents {
			// Проврека последнего оповещения.
			timeLastAlert, err := b.db.GetTimeLastAlert(aE.UserId)
			if err != nil {
				return err
			}

			// Оповещение не чаще чем раз в 1 минуту.
			if int64(math.Abs(time.Since(timeLastAlert).Minutes())) > 1 {
				if err := aletrPosting(b, &aE); err != nil {
					return err
				}
				if err := aletrDelete(b, &aE); err != nil {
					return err
				}
			}
		}
	}
}

// Оповещение о размещении рекламы.
func aletrPosting(b *BotTelegram, aE *models.AdEvent) error {
	timeDateStart, err := sdk.ParseUserDateToTime(aE.DateStart)
	if err != nil {
		return fmt.Errorf("aletrPosting: error ParseUserDateToTime: %w", err)
	}

	// Событие прошло.
	if time.Since(timeDateStart).Minutes() > 0 {
		return nil
	}

	// Сохранение подписчиков на момент выхода рекламы.
	if int64(math.Abs(time.Since(timeDateStart).Minutes())) == 0 {
		currentSub, err := getCurrentSubscriptionFromTelegramChannel(aE.Channel)
		if err != nil {
			return fmt.Errorf("aletrPosting: error getCurrentSubscriptionFromTelegramChannel: %w", err)
        }

		if err := b.db.UpdatePartnerChannelSubscribersInStart(aE.Id, currentSub); err != nil {
			return fmt.Errorf("aletrPosting: error UpdatePartnerChannelSubscribersInStart: %w", err)
		}

		// TODO сохранить кол-во подписчиков канала пользователя.

		if err := b.db.UpdateTimeLastAlert(aE.UserId, time.Now()); err != nil {
			return fmt.Errorf("aletrPosting: error UpdateTimeLastAlert: %w", err)
		}
	}

	timeLeft := int64(math.Abs(time.Since(timeDateStart).Minutes()))
	if checkTimeAlert(aE.UserId, timeLeft) {
		text := createTextAlertForAdEventPosting(aE, timeLeft)
		botMsg := tgbotapi.NewMessage(aE.UserId, text)
		botMsg.ParseMode = tgbotapi.ModeHTML
		botMsg.DisableWebPagePreview = true
		if err := b.sendAlertMessage(aE.UserId, botMsg); err != nil {
			return fmt.Errorf("aletrPosting: error sendAlertMessage: %w", err)
		}
	}

	return nil
}

// Оповещение о удалении рекламы.
func aletrDelete(b *BotTelegram, aE *models.AdEvent) error {
	timeDateEnd, err := sdk.ParseUserDateToTime(aE.DateEnd)
	if err != nil {
		log.Println(fmt.Errorf("aletrDelete: error ParseUserDateToTime: %w", err))
		return err
	}

	// Событие прошло.
	if time.Since(timeDateEnd).Minutes() > 0 {
		return nil
	}

	// Сохранение подписчиков на момент завершения рекламы.
	if int64(math.Abs(time.Since(timeDateEnd).Minutes())) == 0 {
		currentSub, err := getCurrentSubscriptionFromTelegramChannel(aE.Channel)
		if err != nil {
			return fmt.Errorf("aletrDelete: error getCurrentSubscriptionFromTelegramChannel: %w", err)
        }

		if err := b.db.UpdatePartnerChannelSubscribersInEnd(aE.Id, currentSub); err != nil {
			return fmt.Errorf("aletrDelete: error UpdatePartnerChannelSubscribersInEnd: %w", err)
		}

		// TODO сохранить кол-во подписчиков канала пользователя.

		if err := b.db.UpdateTimeLastAlert(aE.UserId, time.Now()); err != nil {
			return fmt.Errorf("aletrDelete: error UpdateTimeLastAlert: %w", err)
		}
	}

	timeLeft := int64(math.Abs(time.Since(timeDateEnd).Minutes()))
	// Удаления  отображаются только за 1 час.
	if timeLeft > 60 {
		return nil
	}

	if checkTimeAlert(aE.UserId, timeLeft) {
		text := createTextAlertForAdEventDelete(aE, timeLeft)
		botMsg := tgbotapi.NewMessage(aE.UserId, text)
		botMsg.ParseMode = tgbotapi.ModeHTML
		botMsg.DisableWebPagePreview = true
		if err := b.sendAlertMessage(aE.UserId, botMsg); err != nil {
			return fmt.Errorf("error edit msg in aletrDelete: %w", err)
		}
	}

	return nil
}

// Проверка доступа к оповещениям
func checkTimeAlert(userId, timeLeft int64) bool {
	// var timeAlerts []int64
	// TODO Смотрим на какое время установил предупреждения полульзователь.
	timeAlerts := []int64{1440, 60, 30, 10, 5}

	for _, timeAlert := range timeAlerts {
		if timeAlert == timeLeft {
			return true
		}
	}

	return false
}
