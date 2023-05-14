package telegram

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"fmt"
	"time"
)

func createStaticsBriefText(d *models.DataForStatistics) string {
	return fmt.Sprintf(`
	<b>      📈 Статистика</b>
<b>Количество проданных реклам:</b> %d
<b>Количество купленных реклам:</b> %d
<b>Количество взаимных пиаров:</b> %d
<b>Прибыль:</b> %d рублей
<b>Траты:</b> %d рублей
<b>Чистая прибыль:</b> %d рублей
`, d.CountAdEventSale, d.CountAdEventBuy, d.CountAdEventMutaul, d.Profit, d.Losses, d.Profit-d.Losses)
}

// Создание текст-описания ad события.
func createTextAdEventDescription(a *models.AdEvent) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Покупатель:</b> %s
		- <b>Канал покупателя:</b> %s
		- <b>Цена продажи:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Продавец:</b> %s
		- <b>Канал продавца:</b> %s
		- <b>Цена покупки:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Партнер по взаимному пиару:</b> %s
		- <b>Канал партнера по взаимному пиару:</b> %s
		- <b>Цена взаимного пиара:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	}

	return descriptionAdEvent
}

// Создание текста оповещения для размещения рекламы.
func createTextAlertForAdEventPosting(a *models.AdEvent, minutesLeftAlert int64) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		Через %s Вы должны разместить рекламу. Подробнее:
		- <b>Покупатель:</b> %s
		- <b>Канал покупателя:</b> %s
		- <b>Цена продажи:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, getTextTime(minutesLeftAlert), a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		Через %s Ваша реклама будет размещена. Подробнее:
		- <b>Продавец:</b> %s
		- <b>Канал продавца:</b> %s
		- <b>Цена покупки:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, getTextTime(minutesLeftAlert), a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		Через %s у Вас начнется взаимный пиар. Подробнее:
		- <b>Партнер по взаимному пиару:</b> %s
		- <b>Канал партнера по взаимному пиару:</b> %s
		- <b>Цена взаимного пиара:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, getTextTime(minutesLeftAlert), a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	}

	return descriptionAdEvent
}

// Создание текста оповещения для удаления рекламы.
func createTextAlertForAdEventDelete(a *models.AdEvent, minutesLeftAlert int64) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		Через %s Вы должны удалить рекламу. Подробнее:
		- <b>Покупатель:</b> %s
		- <b>Канал покупателя:</b> %s
		- <b>Цена продажи:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, getTextTime(minutesLeftAlert), a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		Через %s Ваша реклама будет удалена. Подробнее:
		- <b>Продавец:</b> %s
		- <b>Канал продавца:</b> %s
		- <b>Цена покупки:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, getTextTime(minutesLeftAlert), a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		Через %s у Вас закончится взаимный пиар. Подробнее:
		- <b>Партнер по взаимному пиару:</b> %s
		- <b>Канал партнера по взаимному пиару:</b> %s
		- <b>Цена взаимного пиара:</b> %d
		- <b>Дата размещения рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, getTextTime(minutesLeftAlert), a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	}

	return descriptionAdEvent
}

// Получение правильного текста в зависимости от времени.
func getTextTime(minutes int64) string {
	var textTime string
	if minutes/60 < 1 {
		// Минуты
		if minutes == 1 {
			textTime = fmt.Sprintf("<b>%d</b> минута", minutes)
		} else if minutes >= 2 && minutes <= 4 {
			textTime = fmt.Sprintf("<b>%d</b> минуты", minutes)
		} else {
			textTime = fmt.Sprintf("<b>%d</b> минут", minutes)
		}
	} else {
		// Часы
		hours := minutes / 60
		switch {
		case hours == 1 || hours == 21:
			textTime = fmt.Sprintf("<b>%d</b> час", hours)
		case hours >= 2 && hours <= 4 || hours >= 22 && hours <= 24:
			textTime = fmt.Sprintf("<b>%d</b> часа", hours)
		default:
			textTime = fmt.Sprintf("<b>%d</b> часов", hours)
		}
	}

	return textTime
}

// Возвращает пример даты.
func getTextExampleDate() (string, error) {
	date, err := sdk.ParseTimeToUserDate(time.Now())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`
	В данный момент бот использует только время по МСК 'UTC+3'.
	<b>Пример:</b> %s `, date), nil
}
