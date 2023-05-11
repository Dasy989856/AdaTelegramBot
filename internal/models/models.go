package models

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

// Ошибки.
var (
	ErrUserNotFound = fmt.Errorf("user not found")
	// Example: "22.08.2022 16:30"
	RegxAdEventDate = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01]).(0[1-9]|1[0-2]).(\d{4}) ([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	// Example: "https://t.me/nikname", "@nikname"
	RegxUrlType1 = regexp.MustCompile(`^https:\/\/t\.me\/[a-zA-Z0-9_]+$`)
	// Example: "https://t.me/nikname", "@nikname"
	RegxUrlType2 = regexp.MustCompile(`^@[a-zA-Z0-9_]+$`)
	// Example: "1000"
	RegxPrice = regexp.MustCompile(`[0-9]+`)
	// Example: "1"
	RegxID = regexp.MustCompile(`[0-9]+`)
)

// Типы событий.
var (
	TypeAny    = "any"
	TypeSale   = "sale"
	TypeBuy    = "buy"
	TypeMutual = "mutual"
)

// Пользователь при регистрации.
type User struct {
	Id        int64  `json:"id"`                        // Chat_ID
	CreatedAt string `json:"createdAt" db:"created_at"` // Дата создания.
	Name      string `json:"name" db:"name"`            // Имя пользователя.
	UserURL   string `json:"userUrl" db:"user_url"`     // Ссылка пользователя.
	Step      string `json:"stap" db:"stap"`            // Шаг пользвателя (на каком шаге находится пользователь)
	Login     string `json:"login" db:"login"`
	Password  string `json:"password" db:"password"`
}

// Ad событие.
type AdEvent struct {
	Id                   int64  `json:"id" db:"id"`
	CreatedAt            string `json:"createdAt" db:"created_at"`                        // Дата создания события.
	Ready                bool   `json:"ready" db:"ready"`                                 // Состояние события (Временно не используется)
	UserId               int64  `json:"userId" db:"user_id"`                              // Id пользователя.
	Type                 string `json:"type" db:"type"`                                   // Тип события. (sale, buy ...)
	Partner              string `json:"partner" db:"partner"`                             // Ссылка партнера. (Продавец / Покупатель)
	Channel              string `json:"channel" db:"channel"`                             // Ссылка на канал. (Продавец / Покупатель)
	Price                int64  `json:"price" db:"price"`                                 // Цена.
	DatePosting          string `json:"datePosting" db:"date_posting"`                    // Дата постинга.
	DateDelete           string `json:"dateDelete" db:"date_delete"`                      // Дата удаления поста.
	ArrivalOfSubscribers int64  `json:"arrivalOfSubscribers" db:"arrival_of_subscribers"` // Приход подписчиков.
}

// Если ad событе полностью заполенно - возвращается true. Иначе false.
func (ae *AdEvent) AllData() bool {
	if ae.UserId == 0 {
		log.Println("not found ae.UserId event")
		return false
	}

	if ae.Type == "" {
		log.Println("not found ae.Type event")
		return false
	}

	if ae.CreatedAt == "" {
		log.Println("not found ae.CreatedAt event")
		return false
	}

	if ae.DatePosting == "" {
		log.Println("not found ae.DatePosting event")
		return false
	}

	if ae.DateDelete == "" {
		log.Println("not found ae.DateDelete event")
		return false
	}

	if ae.Partner == "" {
		log.Println("not found ae.Partner event")
		return false
	}

	if ae.Channel == "" {
		log.Println("not found ae.Channel event")
		return false
	}

	return true
}

// БД для телеграмм бота.
type TelegramBotDB interface {
	// Закрытие БД.
	Close() error

	// User

	// Получение данных пользователя.
	GetUserData(userId int64) (user *User, err error)
	// Создание пользователя.
	DefaultUserCreation(chatId int64, userUrl, firstName string) error

	// AdEvent

	// Получение ad события.
	GetAdEvent(eventId int64) (*AdEvent, error)
	// Получение всех ad событий пользователя запрашиваемого типа.
	GetAdEventsOfUser(userId int64, typeAdEvent string) ([]AdEvent, error)
	// Создание ad события.
	AdEventCreation(event *AdEvent) (int64, error)
	// Удаление ad события.
	AdEventDelete(eventId int64) error
	// Проверка доступа пользователя к ad событию.
	// CheckUserAccessToAdEvent(userId, eventId int64) (bool, error)
	// Обновление информации о приходе подписчиков.
	UpdateSubscribesInAdEvent(eventId, subscribers int64) error
	// Установка шага пользователя.
	SetStepUser(userId int64, step string) error
	// Получение текущего шага пользователя.
	GetStepUser(userId int64) (step string, err error)
	// Подучение id незавершенного ad события.
	GetUnfinishedAdEventId(userId int64) (eventId int64, err error)

	// Messages

	// Добавление messageId пользователя.
	AddUsermessageId(userId int64, messageId int) error
	// Удаление messageId пользователя.
	DeleteUsermessageId(messageId int) error
	// Возвращает список messageIds пользователя.
	GetUsermessageIds(userId int64) ([]int, error)
	// Возвращает startmessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
	GetStartmessageId(userId int64) (messageId int, err error)
	// Обновление startmessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
	UpdateStartmessageId(userId int64, messageId int) (err error)
	// Возвращает admessageId. Это сообщение которое не удаляется, купленная в боте реклама.
	GetAdmessageId(userId int64) (messageId int, err error)
	// Обновление AdmessageId. Это сообщение которое не удаляется, купленная в боте реклама.
	UpdateAdmessageId(userId int64, messageId int) (err error)
}

// Парсинг даты в time.Time
func ParseDateToTime(timeString string) (*time.Time, error) {
	layout := "2006-01-02T15:04:00+03:00"

	t1, err := time.Parse(layout, timeString)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	return &t1, nil
}

// Парсинг time.Time в дату.
func ParseTimeToDate(time *time.Time) string {
	return time.Format("02.01.2006 15:04")
}
