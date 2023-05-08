package postgresql

import (
	"AdaTelegramBot/internal/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	usersTable      = "users"       // Пользователи.
	adEventsTable   = "ad_events"   // AD события.
	messageIdsTable = "message_ids" // ID сообщений пользователей.
)

type Config struct {
	Host     string
	Port     string
	NameDB   string
	Username string
	Password string
	ModeSSL  string
}

func NewDB() (*sqlx.DB, error) {
	cfg := Config{
		Host:     viper.GetString("postgre_sql.host"),
		Port:     viper.GetString("postgre_sql.port"),
		NameDB:   viper.GetString("postgre_sql.name_db"),
		Username: viper.GetString("postgre_sql.user_name"),
		Password: viper.GetString("postgre_sql.password"),
		ModeSSL:  viper.GetString("postgre_sql.ssl_mode"),
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password= %s host=%s port=%s dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.NameDB, cfg.ModeSSL))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type TelegramBotDB struct {
	db *sqlx.DB
}

func NewTelegramBotDB(db *sqlx.DB) *TelegramBotDB {
	return &TelegramBotDB{
		db: db,
	}
}

// Закрытие БД.
func (t *TelegramBotDB) Close() error {
	return t.db.Close()
}

// Получение данных пользователя.
func (t *TelegramBotDB) GetUserData(userId int64) (user *models.User, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	u := new(models.User)
	query := fmt.Sprintf(`SELECT (id, created_at, name, user_url, step, login, password) FROM public.%s WHERE id=$1;`, usersTable)
	if err := tx.QueryRow(query, userId).Scan(u); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("error scan user data. err: %w", err)
	}

	return u, nil
}

// Создание default пользователя.
func (t *TelegramBotDB) DefaultUserCreation(chatId int64, userUrl, firstName string) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Default User
	dU := models.User{
		Id:       chatId,
		Name:     firstName,
		UserURL:  "@" + userUrl,
		Step:     "start",
		Login:    userUrl,
		Password: "123",
	}

	// Создание default пользователя.
	sql := fmt.Sprintf(`INSERT INTO public.%s (id, name, user_url, step, login, password)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING;`, usersTable)
	if _, err := tx.Exec(sql, dU.Id, dU.Name, dU.UserURL, dU.Step, dU.Login, dU.Password); err != nil {
		return fmt.Errorf("error create default user. err: %w", err)
	}

	return nil
}

// Установка шага пользователя.
func (p *TelegramBotDB) SetStepUser(userId int64, step string) (err error) {
	tx := p.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`UPDATE public.%s SET step=$1 WHERE id=$2;`, usersTable)
	if _, err := tx.Exec(sql, step, userId); err != nil {
		return fmt.Errorf("error update step user: %w", err)
	}

	return nil
}

// Поулчение шага пользователя.
func (p *TelegramBotDB) GetStepUser(userId int64) (step string, err error) {
	tx := p.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`SELECT step FROM public.%s WHERE id=$1;`, usersTable)
	if err := tx.QueryRow(sql, userId).Scan(&step); err != nil {
		return "", err
	}

	return step, nil
}

// Добавляет ID сообщения пользователя.
func (p *TelegramBotDB) AddUserMessageId(userId int64, messageId int) (err error) {
	tx := p.db.MustBegin()
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
func (p *TelegramBotDB) DeleteUserMessageId(messageId int) (err error) {
	tx := p.db.MustBegin()
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
func (p *TelegramBotDB) GetUserMessageIds(userId int64) (messageIds []int, err error) {
	tx := p.db.MustBegin()
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
