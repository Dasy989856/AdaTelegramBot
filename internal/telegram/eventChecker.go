package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"log"
	"time"
)

// Проврека дат событие и оповещение об их исполнении.
func (b *BotTelegram) adEventChecker() (err error) {
	var cashAdEvents []models.AdEvent
	for {
		time.Sleep(10 * time.Second)
		// time.Sleep(time.Minute)

		timeStart, timeEnd := sdk.GetTimeRangeToday()
		fmt.Println(timeStart, timeEnd)
		cashAdEvents, err = b.db.GetRangeAdEvents(models.TypeAny, timeStart, timeEnd)
		if err != nil {
			log.Println(fmt.Errorf("error get AdEvents from DB: %w", err))
			return err
		}
		fmt.Println("Get EVENT OK")

		for _, aE := range cashAdEvents {
			timeDateDelete, err := sdk.ParseDateToTime(aE.DateDelete)
			if err != nil {
				log.Println(fmt.Errorf("error parsing date AdEvent: %w", err))
			return err
			}
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
