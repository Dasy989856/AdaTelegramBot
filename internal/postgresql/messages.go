package postgresql

import (
	"fmt"
)

// Добавляет ID сообщения пользователя.
func (t *TelegramBotDB) AddUserMessageId(userId int64, messageId int) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`INSERT INTO public.%s (id, user_id) values ($1, $2) ON CONFLICT DO NOTHING;`, messageIdsTable)
	if _, err := tx.Exec(sql, messageId, userId); err != nil {
		return fmt.Errorf("error insert messageId: %w", err)
	}

	return nil
}

// Удаление messageId пользователя.
func (t *TelegramBotDB) DeleteUserMessageId(messageId int) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`DELETE FROM public.%s WHERE id=$1;`, messageIdsTable)
	if _, err := tx.Exec(sql, messageId); err != nil {
		return fmt.Errorf("error delete messageId: %w", err)
	}

	return nil
}

// Возвращает список messageIds пользователя.
func (t *TelegramBotDB) GetUserMessageIds(userId int64) (messageIds []int, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	messageIds = make([]int, 0, 50)

	sql := fmt.Sprintf(`SELECT id FROM public.%s WHERE user_id=$1;`, messageIdsTable)
	rows, err := tx.Query(sql, userId)
	if err != nil {
		return nil, fmt.Errorf("error select messageIds: %w", err)
	}

	for rows.Next() {
		var messageId int
		if err := rows.Scan(&messageId); err != nil {
			return nil, fmt.Errorf("error scan messageId in GetUserMessageIds: %w", err)
		}
		messageIds = append(messageIds, messageId)
	}

	return messageIds, nil
}

// Возвращает startMessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
func (t *TelegramBotDB) GetStartMessageId(userId int64) (messageId int, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Получение startMessageId пользователя.
	sql := fmt.Sprintf(`SELECT start_message_id FROM public.%s WHERE id=$1;`, usersTable)
	if err := tx.QueryRow(sql, userId).Scan(&messageId); err != nil {
		return 0, fmt.Errorf("error select startMessageId: %w", err)
	}

	if messageId == 0 {
		return 0, fmt.Errorf("startMessageId quil 0")
	}

	return messageId, nil
}

// Обновление startMessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
func (t *TelegramBotDB) UpdateStartMessageId(userId int64, messageId int) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`UPDATE public.%s SET start_message_id=$1 WHERE id=$2;`, usersTable)
	if _, err := tx.Exec(sql, messageId, userId); err != nil {
		return fmt.Errorf("error update startMessageId: %w", err)
	}

	return nil
}

// Возвращает adMessageId. Это сообщение которое не удаляется, купленная в боте реклама.
func (t *TelegramBotDB) GetAdMessageId(userId int64) (messageId int, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Получение adMessageId пользователя.
	sql := fmt.Sprintf(`SELECT ad_message_id FROM public.%s WHERE id=$1;`, usersTable)
	if err := tx.QueryRow(sql, userId).Scan(&messageId); err != nil {
		return 0, fmt.Errorf("error select GetAdMessageId: %w", err)
	}

	if messageId == 0 {
		return 0, fmt.Errorf("adMessageId quil 0")
	}

	return messageId, nil
}

// Обновление adMessageId. Это сообщение которое не удаляется, купленная в боте реклама.
func (t *TelegramBotDB) UpdateAdMessageId(userId int64, messageId int) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`UPDATE public.%s SET ad_message_id=$1 WHERE id=$2;`, usersTable)
	if _, err := tx.Exec(sql, messageId, userId); err != nil {
		return fmt.Errorf("error update adMessageId: %w", err)
	}

	return nil
}