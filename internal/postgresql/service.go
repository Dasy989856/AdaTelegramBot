package postgresql

import (
	"abs-by-ammka-bot/internal/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	usersTable  = "users"
	eventsTable = "ad_events"
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

// Создание ad события.
func (t *TelegramBotDB) AdEventCreation(event *models.Event) (eventId int64, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`INSERT INTO public.%s (user_id, type, created_at, posting_date, partner_url,
		price, arrival_of_subscribers) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, eventsTable)
	if err := tx.QueryRow(sql, event.UserId, event.Type, event.CreatedAt, event.PostingDate, event.PartnerURL,
		event.Price, event.ArrivalOfSubscribers).Scan(&eventId); err != nil {
		return 0, fmt.Errorf("error creation event: %w", err)
	}

	return eventId, nil
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

	sql := fmt.Sprintf(`DELETE FROM public.%s WHERE id=$1;`, eventsTable)
	if _, err := tx.Exec(sql, eventId); err != nil {
		return fmt.Errorf("error delete event: %w", err)
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

	sql := fmt.Sprintf(`SELECT step FROM public.%s WHERE user_id=$1;`, usersTable)
	if err := tx.QueryRow(sql, userId).Scan(&step); err != nil {
		return "", err
	}

	return step, nil
}
