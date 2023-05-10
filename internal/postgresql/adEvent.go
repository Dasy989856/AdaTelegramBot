package postgresql

import (
	"AdaTelegramBot/internal/models"
	"fmt"
)

// Получение ad события.
func (t *TelegramBotDB) GetAdEvent(eventId int64) (adEvent *models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var aE models.AdEvent
	sql := fmt.Sprintf(`SELECT (id, created_at, user_id, "type", partner, channel, price, date_posting, date_delete, arrival_of_subscribers)
	FROM public.%s WHERE id=$1`, adEventsTable)
	if err := tx.QueryRow(sql, eventId).Scan(&aE); err != nil {
		return nil, fmt.Errorf("error creation event: %w", err)
	}
	
	fmt.Println(aE)

		// Изменение формата времени.
	// timeDatePosting, err := models.ParseDateToTime(event.DatePosting)
	// if err != nil {
	// 	return 0, err
	// }
	// event.DatePosting = timeDatePosting.Format("2006-01-02 15:04:05.999")

	// timeDateDelete, err := models.ParseDateToTime(event.DateDelete)
	// if err != nil {
	// 	return 0, err
	// }
	// event.DateDelete = timeDateDelete.Format("2006-01-02 15:04:05.999")

	return &aE, nil
}

// Создание ad события.
func (t *TelegramBotDB) AdEventCreation(event *models.AdEvent) (eventId int64, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Изменение формата времени.
	timeDatePosting, err := models.ParseDateToTime(event.DatePosting)
	if err != nil {
		return 0, err
	}
	event.DatePosting = timeDatePosting.Format("2006-01-02 15:04:05.999")
	timeDateDelete, err := models.ParseDateToTime(event.DateDelete)
	if err != nil {
		return 0, err
	}
	event.DateDelete = timeDateDelete.Format("2006-01-02 15:04:05.999")

	sql := fmt.Sprintf(`INSERT INTO public.%s (ready, user_id, "type", partner, channel, price, date_posting, date_delete)
	values (true, $1, $2, $3, $4, $5, $6, $7) RETURNING id;`, adEventsTable)
	if err := tx.QueryRow(sql, event.UserId, event.Type, event.Partner, event.Channel, event.Price,
		event.DatePosting, event.DateDelete).Scan(&eventId); err != nil {
		return 0, fmt.Errorf("error creation event: %w", err)
	}

	return eventId, nil
}

// Добавление информации о приходе подписчиков.
func (t *TelegramBotDB) UpdateSubscribesInAdEvent(eventId, subscribers int64) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`UPDATE public.%s SET arrival_of_subscribers=$1
	WHERE id=$2;`, adEventsTable)
	if _, err := tx.Exec(sql, subscribers, eventId); err != nil {
		return fmt.Errorf("error update arrival_of_subscribers eventId%d: %w", eventId, err)
	}

	return nil
}

// Удаление ad события.
func (t *TelegramBotDB) AdEventDelete(eventId int64) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`DELETE FROM public.%s WHERE id=$1;`, adEventsTable)
	if _, err := tx.Exec(sql, eventId); err != nil {
		return fmt.Errorf("error delete event: %w", err)
	}

	return nil
}

// Возвращает id незавершенного события.
func (t *TelegramBotDB) GetUnfinishedAdEventId(userId int64) (adEventId int64, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// SELECT id FROM public.ad_events WHERE ready=false and user_id=959606248;
	sql := fmt.Sprintf(`SELECT id FROM public.%s WHERE ready=false AND user_id=$1;`, adEventsTable)
	if err := tx.QueryRow(sql, userId).Scan(&adEventId); err != nil {
		return 0, fmt.Errorf("error scan id event: %w", err)
	}

	return adEventId, nil
}
