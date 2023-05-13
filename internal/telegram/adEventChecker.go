package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"log"
	"time"
)

// Проврека дат событий и оповещение об их исполнении.
func (b *BotTelegram) adEventChecker() (err error) {
	var cashAdEvents []models.AdEvent
	for {
		time.Sleep(15 * time.Second)

		timeStart, timeEnd := sdk.GetTimeRangeToday()
		fmt.Println(timeStart, timeEnd)
		cashAdEvents, err = b.db.GetRangeAdEvents(models.TypeAny, timeStart, timeEnd)
		if err != nil {
			log.Println(fmt.Errorf("error get AdEvents from DB: %w", err))
			return err
		}

		for _, aE := range cashAdEvents {
			// Даты удаления.
			timeDateDelete, err := sdk.ParseUserDateToTime(aE.DateDelete)
			if err != nil {
				log.Println(fmt.Errorf("error parsing date in adEventChecker: %w", err))
				return err
			}

			fmt.Println(time.Since(timeDateDelete).Minutes())
			if time.Since(timeDateDelete).Minutes() < 15 {
				fmt.Println(aE)
				fmt.Println("<15")
			} else {
				fmt.Println(aE)
				fmt.Println(">15")
			}
		}
	}
}
