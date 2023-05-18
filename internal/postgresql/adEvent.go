package postgresql

import (
	"AdaTelegramBot/internal/models"
	"AdaTelegramBot/internal/sdk"
	"time"

	"database/sql"
	"fmt"
)

func (t *TelegramBotDB) GetAdEvent(adEventId int64) (adEvent *models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE id=$1;`, adEventsTable)

	var aE models.AdEvent
	var datePostingFromDB, dateDeleteFromDB string
	if err := tx.QueryRow(query, adEventId).Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
		&datePostingFromDB, &dateDeleteFromDB, &aE.ArrivalOfSubscribers); err != nil {
		return nil, fmt.Errorf("error scan AdEvent in GetAdEvent: %w", err)
	}
	aE.DatePosting, aE.DateDelete, err = editDateFromDataBaseToUserDate(datePostingFromDB, dateDeleteFromDB)
	if err != nil {
		return nil, err
	}

	return &aE, nil
}

func (t *TelegramBotDB) GetRangeAdEvents(typeAdEvent models.TypeAdEvent, startTime, endTime time.Time) (listAdEvent []models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// SELECT id, created_at, user_id, "type", partner, channel, price,
	// date_posting, date_delete, arrival_of_subscribers
	// FROM public.ad_events WHERE "type"='sale' AND ((date_posting BETWEEN '2022-05-13 00:00:00 +0300' AND '2024-05-13 00:00:00 +0300')
	// OR (date_delete BETWEEN '2022-05-13 00:00:00 +0300' AND '2024-05-13 00:00:00 +0300')) ORDER BY date_posting ASC;

	listAdEvent = make([]models.AdEvent, 0, 50)
	startDate, err := parseTimeToDateDataBase(startTime)
	if err != nil {
		return nil, err
	}
	endDate, err := parseTimeToDateDataBase(endTime)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if typeAdEvent == models.TypeAny {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE ((date_posting BETWEEN $1 AND $2)
		OR (date_delete BETWEEN $2 AND $1)) ORDER BY date_posting ASC;`, adEventsTable)

		rows, err = tx.Query(query, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	} else {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE "type"=$3 AND ((date_posting BETWEEN $1 AND $2)
		OR (date_delete BETWEEN $1 AND $2)) ORDER BY date_posting ASC;`, adEventsTable)

		rows, err = tx.Query(query, startDate, endDate, typeAdEvent)
		if err != nil {
			return nil, fmt.Errorf("error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	}

	for rows.Next() {
		var aE models.AdEvent
		var datePostingFromDB, dateDeleteFromDB string
		if err := rows.Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
			&datePostingFromDB, &dateDeleteFromDB, &aE.ArrivalOfSubscribers); err != nil {
			return nil, fmt.Errorf("error scan AdEvent in GetRangeAdEvents: %w", err)
		}

		aE.DatePosting, aE.DateDelete, err = editDateFromDataBaseToUserDate(datePostingFromDB, dateDeleteFromDB)
		listAdEvent = append(listAdEvent, aE)
	}

	return listAdEvent, nil
}

func (t *TelegramBotDB) GetAdEventsOfUser(userId int64, typeAdEvent models.TypeAdEvent) (listAdEvent []models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	listAdEvent = make([]models.AdEvent, 0, 50)
	var rows *sql.Rows
	if typeAdEvent == models.TypeAny {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE user_id=$1 ORDER BY date_posting ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId)
		if err != nil {
			return nil, fmt.Errorf("error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	} else {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE user_id=$1 AND "type"=$2 ORDER BY date_posting ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId, typeAdEvent)
		if err != nil {
			return nil, fmt.Errorf("error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	}

	for rows.Next() {
		var aE models.AdEvent
		var datePostingFromDB, dateDeleteFromDB string
		if err := rows.Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
			&datePostingFromDB, &dateDeleteFromDB, &aE.ArrivalOfSubscribers); err != nil {
			return nil, fmt.Errorf("error scan AdEvent in GetAdEventsOfUser: %w", err)
		}

		aE.DatePosting, aE.DateDelete, err = editDateFromDataBaseToUserDate(datePostingFromDB, dateDeleteFromDB)
		if err != nil {
			return nil, err
		}
		listAdEvent = append(listAdEvent, aE)
	}

	return listAdEvent, nil
}

func (t *TelegramBotDB) GetRangeAdEventsOfUser(userId int64, typeAdEvent models.TypeAdEvent, startTime, endTime time.Time) (listAdEvent []models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	//  SELECT id, created_at, user_id, "type", partner, channel, price,
	// 	date_posting, date_delete, arrival_of_subscribers
	// 	FROM public.ad_events WHERE user_id=959606248 AND ((date_posting BETWEEN '2021-01-01' AND '2023-12-31')
	// 	OR (date_delete BETWEEN '2021-01-01' AND '2023-12-31')) ORDER BY date_posting ASC

	listAdEvent = make([]models.AdEvent, 0, 50)
	startDate, err := parseTimeToDateDataBase(startTime)
	if err != nil {
		return nil, err
	}
	endDate, err := parseTimeToDateDataBase(endTime)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if typeAdEvent == models.TypeAny {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE user_id=$1 AND ((date_posting BETWEEN $2 AND $3)
		OR (date_delete BETWEEN $2 AND $3)) ORDER BY date_posting ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	} else {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price,
		date_posting, date_delete, arrival_of_subscribers
		FROM public.%s WHERE user_id=$1 AND "type"=$4 AND ((date_posting BETWEEN $2 AND $3)
		OR (date_delete BETWEEN $2 AND $3)) ORDER BY date_posting ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId, startDate, endDate, typeAdEvent)
		if err != nil {
			return nil, fmt.Errorf("error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	}

	for rows.Next() {
		var aE models.AdEvent
		var datePostingFromDB, dateDeleteFromDB string
		if err := rows.Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
			&datePostingFromDB, &dateDeleteFromDB, &aE.ArrivalOfSubscribers); err != nil {
			return nil, fmt.Errorf("error scan AdEvent in GetRangeAdEventsOfUser: %w", err)
		}

		aE.DatePosting, aE.DateDelete, err = editDateFromDataBaseToUserDate(datePostingFromDB, dateDeleteFromDB)

		listAdEvent = append(listAdEvent, aE)
	}

	return listAdEvent, nil
}

func (t *TelegramBotDB) AdEventCreation(event *models.AdEvent) (eventId int64, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if event.DateDelete == "" {
		event.DateDelete = "02.01.2006 15:04"
	}

	// Изменение формата времени.
	timeDatePosting, err := sdk.ParseUserDateToTime(event.DatePosting)
	if err != nil {
		return 0, err
	}
	event.DatePosting, err = parseTimeToDateDataBase(timeDatePosting)
	if err != nil {
		return 0, err
	}

	timeDateDelete, err := sdk.ParseUserDateToTime(event.DateDelete)
	if err != nil {
		return 0, err
	}
	event.DateDelete, err = parseTimeToDateDataBase(timeDateDelete)
	if err != nil {
		return 0, err
	}

	sql := fmt.Sprintf(`INSERT INTO public.%s (ready, user_id, "type", partner, channel, price, date_posting, date_delete)
	values (true, $1, $2, $3, $4, $5, $6, $7) RETURNING id;`, adEventsTable)
	if err := tx.QueryRow(sql, event.UserId, event.Type, event.Partner, event.Channel, event.Price,
		event.DatePosting, event.DateDelete).Scan(&eventId); err != nil {
		return 0, fmt.Errorf("error creation event: %w", err)
	}

	return eventId, nil
}

// Обновление данные ad события.
func (t *TelegramBotDB) AdEventUpdate(aE *models.AdEvent) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if aE.Type != "" {
		query := fmt.Sprintf(`UPDATE public.%s SET "type"=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Type, aE.Id); err != nil {
			return fmt.Errorf("error update type. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.Partner != "" {
		query := fmt.Sprintf(`UPDATE public.%s SET partner=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Partner, aE.Id); err != nil {
			return fmt.Errorf("error update partner. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.Channel != "" {
		query := fmt.Sprintf(`UPDATE public.%s SET channel=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Channel, aE.Id); err != nil {
			return fmt.Errorf("error update channel. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.Price != 0 {
		query := fmt.Sprintf(`UPDATE public.%s SET price=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Price, aE.Id); err != nil {
			return fmt.Errorf("error update price. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.DatePosting != "" {
		timeDatePosting, err := sdk.ParseUserDateToTime(aE.DatePosting)
		if err != nil {
			return err
		}
		aE.DatePosting, err = parseTimeToDateDataBase(timeDatePosting)
		if err != nil {
			return err
		}

		query := fmt.Sprintf(`UPDATE public.%s SET date_posting=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.DatePosting, aE.Id); err != nil {
			return fmt.Errorf("error update date_posting. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.DateDelete != "" {
		timeDateDelete, err := sdk.ParseUserDateToTime(aE.DateDelete)
		if err != nil {
			return err
		}
		aE.DateDelete, err = parseTimeToDateDataBase(timeDateDelete)
		if err != nil {
			return err
		}

		query := fmt.Sprintf(`UPDATE public.%s SET date_delete=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.DateDelete, aE.Id); err != nil {
			return fmt.Errorf("error update date_delete. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.ArrivalOfSubscribers != 0 {
		query := fmt.Sprintf(`UPDATE public.%s SET arrival_of_subscribers=$1
		WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.ArrivalOfSubscribers, aE.Id); err != nil {
			return fmt.Errorf("error update arrival_of_subscribers. eventId%d: %w", aE.Id, err)
		}
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

func editDateFromDataBaseToUserDate(datePostingFromDB, dateDeleteFromDB string) (userDatePosting, userDateDelete string, err error) {
	// User Date Posting
	timeDatePostingFromDB, err := parseDateDataBaseToTime(datePostingFromDB)
	if err != nil {
		return "", "", err
	}
	userDatePosting, err = sdk.ParseTimeToUserDate(timeDatePostingFromDB)
	if err != nil {
		return "", "", err
	}

	// User Date Delete
	timeDateDeleteFromDB, err := parseDateDataBaseToTime(dateDeleteFromDB)
	if err != nil {
		return "", "", err
	}
	userDateDelete, err = sdk.ParseTimeToUserDate(timeDateDeleteFromDB)
	if err != nil {
		return "", "", err
	}

	return userDatePosting, userDateDelete, nil
}
