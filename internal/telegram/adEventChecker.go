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
func (b *BotTelegram) adEventChecker() (err error) {
	var cashAdEvents []models.AdEvent
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("START CHECK")

		_, timeEnd := sdk.GetTimeRangeTomorrow()
		timeStart := time.Now()
		cashAdEvents, err = b.db.GetRangeAdEvents(models.TypeAny, timeStart, timeEnd)
		if err != nil {
			log.Println(fmt.Errorf("error get AdEvents from DB: %w", err))
			return err
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
	timeDatePosting, err := sdk.ParseUserDateToTime(aE.DatePosting)
	if err != nil {
		log.Println(fmt.Errorf("error parsing date in aletrPosting: %w", err))
		return err
	}

	timeLeftAlert := int64(math.Abs(time.Since(timeDatePosting).Minutes()))
	fmt.Println("POSTING: ", timeLeftAlert)

	if checkTimeAlert(aE.UserId, timeLeftAlert) {
		text := createAlertTextForAdEventPosting(aE, timeLeftAlert)
		botMsg := tgbotapi.NewMessage(aE.UserId, text)
		botMsg.ParseMode = tgbotapi.ModeHTML
		botMsg.DisableWebPagePreview = true
		if err := b.sendAlertMessage(aE.UserId, botMsg); err != nil {
			return fmt.Errorf("error edit msg in aletrPosting: %w", err)
		}
	}

	return nil
}

// Оповещение о удалении рекламы.
func aletrDelete(b *BotTelegram, aE *models.AdEvent) error {
	timeDateDelete, err := sdk.ParseUserDateToTime(aE.DateDelete)
	if err != nil {
		log.Println(fmt.Errorf("error parsing date in aletrDelete: %w", err))
		return err
	}

	timeLeftAlert := int64(math.Abs(time.Since(timeDateDelete).Minutes()))
	fmt.Println("DELETE: ", timeLeftAlert)

	// Удаления  отображаются только за 1 час.
	if timeLeftAlert > 60 {
		return nil
	}

	if checkTimeAlert(aE.UserId, timeLeftAlert) {
		text := createAlertTextForAdEventDelete(aE, timeLeftAlert)
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
func checkTimeAlert(userId, timeLeftAlert int64) bool {
	// var timeAlerts []int64
	// TODO Смотрим на какое время установил предупреждения полульзователь.
	timeAlerts := []int64{1440, 60, 30, 15, 5}

	for _, timeAlert := range timeAlerts {
		if timeAlert == timeLeftAlert {
			return true
		}
	}

	return false
}
