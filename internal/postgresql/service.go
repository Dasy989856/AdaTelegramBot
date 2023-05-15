package postgresql

import (
	"AdaTelegramBot/internal/models"
	"database/sql"
	"fmt"
	"math"
	"time"

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
	query := fmt.Sprintf(`SELECT (id, created_at, name, user_url, step, login, password)
	FROM public.%s WHERE id=$1;`, usersTable)
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
		Id:           chatId,
		Name:         firstName,
		UserURL:      "@" + userUrl,
		Step:         "start",
		Login:        userUrl,
		PasswordHash: "123",
	}

	// Создание default пользователя.
	sql := fmt.Sprintf(`INSERT INTO public.%s (id, name, user_url, step, login, password)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING;`, usersTable)
	if _, err := tx.Exec(sql, dU.Id, dU.Name, dU.UserURL, dU.Step, dU.Login, dU.PasswordHash); err != nil {
		return fmt.Errorf("error create default user. err: %w", err)
	}

	return nil
}

// Установка шага пользователя.
func (t *TelegramBotDB) SetStepUser(userId int64, step string) (err error) {
	tx := t.db.MustBegin()
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
func (t *TelegramBotDB) GetStepUser(userId int64) (step string, err error) {
	tx := t.db.MustBegin()
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

// Получение данных пользователя для статистики.
func (t *TelegramBotDB) GetRangeDataForStatistics(userId int64, typeAdEvent models.TypeAdEvent, startDate, endDate time.Time) (data *models.DataForStatistics, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Данные для создание статистик.
	var d models.DataForStatistics

	listAdEvents, err := t.GetRangeAdEventsOfUser(userId, typeAdEvent, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, adEvent := range listAdEvents {
		switch adEvent.Type {
		case models.TypeSale:
			d.CountAdEventSale++
			d.Profit += adEvent.Price
		case models.TypeBuy:
			d.CountAdEventBuy++
			d.Losses += adEvent.Price
		case models.TypeMutual:
			d.CountAdEventMutaul++
			if adEvent.Price > 0 {
				d.Profit += adEvent.Price
			} else {
				d.Losses += int64(math.Abs(float64(adEvent.Price)))
			}
		case models.TypeBarter:
			d.CountAdEventBarter++
			if adEvent.Price > 0 {
				d.Profit += adEvent.Price
			} else {
				d.Losses += int64(math.Abs(float64(adEvent.Price)))
			}
		}
	}

	return &d, nil
}

// Получение временной метки последнего предупреждения.
func (t *TelegramBotDB) GetTimeLastAlert(userId int64) (timeLastAlert time.Time, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var dateLastAlert string
	query := fmt.Sprintf(`SELECT date_last_alert FROM public.%s WHERE id=$1;`, usersTable)
	if err := tx.QueryRow(query, userId).Scan(&dateLastAlert); err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, models.ErrUserNotFound
		}
		return time.Time{}, fmt.Errorf("error scan user data. err: %w", err)
	}

	timeLastAlert, err = parseDateDataBaseToTime(dateLastAlert)
	if err != nil {
		return time.Time{}, fmt.Errorf("error pase date to time: %w", err)
	}
	return timeLastAlert, nil
}

// Обновление временной метки последнего предупреждения.
func (t *TelegramBotDB) UpdateTimeLastAlert(userId int64, timeLastAlert time.Time) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	dateLastAlert, err := parseTimeToDateDataBase(timeLastAlert)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`UPDATE public.%s SET date_last_alert=$1
	WHERE id=$2;`, usersTable)
	if _, err := tx.Exec(sql, dateLastAlert, userId); err != nil {
		return fmt.Errorf("error update date_last_alert user %d: %w", userId, err)
	}

	return nil
}

// Парсинг даты из БД в time.Time
func parseDateDataBaseToTime(timeString string) (time.Time, error) {
	layout := "2006-01-02T15:04:05Z"
	var t time.Time

	defaultTimeZoneInDataBase, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}, fmt.Errorf("error create defaultTimeZoneInDataBase: %w", err)
	}
	t = t.In(defaultTimeZoneInDataBase)

	t, err = time.ParseInLocation(layout, timeString, defaultTimeZoneInDataBase)
	if err != nil {
		return t, fmt.Errorf("error parsing date: %w", err)
	}

	return t, nil
}

// Парсинг time.Time в дату из БД
func parseTimeToDateDataBase(t time.Time) (string, error) {
	defaultTimeZoneInDataBase, err := time.LoadLocation("UTC")
	if err != nil {
		return "", fmt.Errorf("error create defaultTimeZoneInDataBase: %w", err)
	}
	t = t.In(defaultTimeZoneInDataBase)

	return t.Format("2006-01-02T15:04:05Z"), nil
}
