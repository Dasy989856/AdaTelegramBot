package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
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
func createAdEventDescriptionText(a *models.AdEvent) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Покупатель:</b> %s
		- <b>Канал покупателя:</b> %s
		- <b>Цена продажи:</b> %d
		- <b>Дата постинга рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Продавец:</b> %s
		- <b>Канал продавца:</b> %s
		- <b>Цена покупки:</b> %d
		- <b>Дата постинга рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Партнер по взаимному пиару:</b> %s
		- <b>Канал партнера по взаимному пиару:</b> %s
		- <b>Цена взаимного пиара:</b> %d
		- <b>Дата постинга рекламы:</b> %s
		- <b>Дата удаления рекламы:</b> %s
		`, a.Partner, a.Channel, a.Price, a.DatePosting, a.DateDelete)
	}

	return descriptionAdEvent
}