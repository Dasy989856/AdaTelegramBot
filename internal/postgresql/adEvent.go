package postgresql

import (
	"AdaTelegramBot/internal/models"
	"fmt"
)

// Создание дефолтного ad события.
// Exemple:
/*
INSERT INTO public.ad_events
(ready, user_id, "type", partner, chanel, price, date_posting, date_delete)
VALUES(true, 1, 'buy', 'https://t.me/nikname', 'https://t.me/chanelname', 1000, '2023-05-08 17:22:19', '2023-05-08 18:22:19')
RETURNING id;
*/
func (t *TelegramBotDB) AdEventCreation(event *models.AdEvent) (eventId int64, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}() 

	sql := fmt.Sprintf(`INSERT INTO public.%s (ready, user_id, "type", partner, channel, price, date_posting, date_delete)
	values (true, $1, $2, $3, $4, $5, $6, $7) RETURNING id;`, adEventsTable)
	if err := tx.QueryRow(sql, event.UserId, event.Type, event.Partner, event.Channel, event.Price, event.DatePosting, event.DateDelete).Scan(&eventId); err != nil {
		return 0, fmt.Errorf("error creation event: %w", err)
	}

	return eventId, nil
}

// Добавление информации о приходе подписчиков.
// Exemple:
/*
UPDATE public.ad_events
	SET arrival_of_subscribers=123
	WHERE id=1;
*/
func (t *TelegramBotDB) UpdateSubscribeInAdEvent(eventId, subscribers int64) (err error) {
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