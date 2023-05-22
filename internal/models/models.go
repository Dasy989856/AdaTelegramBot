package models

import (
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Пользователь при регистрации.
type User struct {
	Id             int64  `json:"id"`                                   // Chat_ID
	CreatedAt      string `json:"createdAt" db:"created_at"`            // Дата создания.
	Name           string `json:"name" db:"name"`                       // Имя пользователя.
	UserURL        string `json:"userUrl" db:"user_url"`                // Ссылка пользователя.
	Step           string `json:"stap" db:"stap"`                       // Шаг пользвателя (на каком шаге находится пользователь)
	Login          string `json:"login" db:"login"`                     // Логин пользователя.
	PasswordHash   string `json:"password" db:"password"`               // Хэш пользовательского пароля.
	StartMessageId int    `json:"startMessageId" db:"start_message_id"` // Индефикатор сообщения startMessage.
	AdMessageId    int    `json:"adMessageId" db:"ad_message_id"`       // Индефикатор сообщения adMessage.
	InfoMessageId  int    `json:"infoMessageId" db:"info_message_id"`   // Индефикатор сообщения infoMessage.
	DateLastAlert  string `json:"DateLastAlert" db:"date_last_alert"`   // Даты последнего оповещения пользователя.
}

// Ошибки.
var (
	ErrUserNotFound = fmt.Errorf("user not found")
	// Example: "22.08.2022 16:30"
	RegxAdEventDate = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01]).(0[1-9]|1[0-2]).(\d{4}) ([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	// Example: "https://t.me/nikname", "https://www.instagram.com/nikname.store/"
	RegxUrlType1 = regexp.MustCompile(`^https://[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+(/[a-zA-Z0-9-]*)*`)
	// Example: "@nikname"
	RegxUrlType2 = regexp.MustCompile(`^@[a-zA-Z0-9_]+$`)
	// Example: 1000
	RegxPrice = regexp.MustCompile(`[0-9]+`)
	// Example: 1000
	RegxArrivalOfSubscribers = regexp.MustCompile(`[0-9]+`)
	// Example: 1
	RegxId = regexp.MustCompile(`[0-9]+`)
)

// Типы CallbackQuery.

type CbqStatic tgbotapi.CallbackQuery  // CallbackQuery без CbqData
type CbqDinamic tgbotapi.CallbackQuery // CallbackQuery с CbqData
type CbqPath []string                  // Путь
type CbqData []byte

var CbqSep string = "?" // Разделитель между cbqPath{[]string} и cbqData{json}

// Тип AdEvent.
type TypeAdEvent string

var (
	TypeAny    TypeAdEvent = "any"
	TypeSale   TypeAdEvent = "sale"
	TypeBuy    TypeAdEvent = "buy"
	TypeMutual TypeAdEvent = "mutual"
	TypeBarter TypeAdEvent = "barter"
)

// Ad событие.
type AdEvent struct {
	Id                   int64       `json:"id" db:"id"`
	CreatedAt            string      `json:"createdAt" db:"created_at"`                        // Дата создания события.
	Ready                bool        `json:"ready" db:"ready"`                                 // Состояние события (Временно не используется)
	UserId               int64       `json:"userId" db:"user_id"`                              // Id пользователя.
	Type                 TypeAdEvent `json:"type" db:"type"`                                   // Тип события. (sale, buy ...)
	Partner              string      `json:"partner" db:"partner"`                             // Ссылка партнера. (Продавец / Покупатель)
	Channel              string      `json:"channel" db:"channel"`                             // Ссылка на канал. (Продавец / Покупатель)
	Price                int64       `json:"price" db:"price"`                                 // Цена.
	DatePosting          string      `json:"datePosting" db:"date_posting"`                    // Дата размещения. "02.01.06 15:04"
	DateDelete           string      `json:"dateDelete" db:"date_delete"`                      // Дата удаления поста. "02.01.06 15:04"
	ArrivalOfSubscribers int64       `json:"arrivalOfSubscribers" db:"arrival_of_subscribers"` // Приход подписчиков.
}

// Данные для создания статистики.
type DataForStatistics struct {
	CountAdEventSale   int64 // Кол-во проданных реклам.
	CountAdEventBuy    int64 // Кол-во купленных реклам.
	CountAdEventMutaul int64 // Кол-во взаимных пиаров.
	CountAdEventBarter int64 // Кол-во бартеров.
	Profit             int64 // Прибыль.
	Losses             int64 // Убытки.
}

// Сессия пользователя.
type Session struct {
	DomainPath string // Наименование основной цепочки.
	Step       int64  // Шаг в цепочке.
	StateMsg   string // Состояние ожидающих данных в Msg.
	Cache      map[string]interface{}
}

// БД для телеграмм бота.
type TelegramBotDB interface {
	// Закрытие БД.
	Close() error

	// Получение данных пользователя.
	GetUserData(userId int64) (user *User, err error)
	// Создание пользователя.
	DefaultUserCreation(chatId int64, userUrl, firstName string) error
	// Получение последней даты оповещения.
	GetTimeLastAlert(userId int64) (timeLastAlert time.Time, err error)
	// Обновление последней даты оповещения.
	UpdateTimeLastAlert(userId int64, timeLastAlert time.Time) error

	// Получение ad события.
	GetAdEvent(adEventId int64) (*AdEvent, error)
	// Получение всех ad событий в указаном диапазоне времени.
	GetRangeAdEvents(typeAdEvent TypeAdEvent, startDate, endDate time.Time) ([]AdEvent, error)
	// Получение всех ad событий пользователя запрашиваемого типа.
	GetAdEventsOfUser(userId int64, typeAdEvent TypeAdEvent) ([]AdEvent, error)
	// Получение всех ad событий пользователя запрашиваемого типа в указаном диапазоне времени.
	GetRangeAdEventsOfUser(userId int64, typeAdEvent TypeAdEvent, startDate, endDate time.Time) ([]AdEvent, error)
	// Создание ad события.
	AdEventCreation(adEvent *AdEvent) (int64, error)
	// Удаление ad события.
	AdEventDelete(eventId int64) error
	// Обновление информации о приходе подписчиков.
	AdEventUpdate(adEvent *AdEvent) error
	// Установка шага пользователя.
	SetStepUser(userId int64, step string) error
	// Получение текущего шага пользователя.
	GetStepUser(userId int64) (step string, err error)
	// Подучение id незавершенного ad события.
	GetUnfinishedAdEventId(userId int64) (eventId int64, err error)

	// Добавление messageId пользователя.
	AddUserMessageId(userId int64, messageId int) error
	// Удаление messageId пользователя.
	DeleteUsermessageId(messageId int) error
	// Возвращает список messageIds пользователя.
	GetUserMessageIds(userId int64) ([]int, error)
	// Возвращает startmessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
	GetStartMessageId(userId int64) (messageId int, err error)
	// Обновление startmessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
	UpdateStartMessageId(userId int64, messageId int) (err error)
	// Возвращает admessageId. Это сообщение которое не удаляется, купленная в боте реклама.
	GetAdMessageId(userId int64) (messageId int, err error)
	// Обновление AdmessageId. Это сообщение которое не удаляется, купленная в боте реклама.
	UpdateAdMessageId(userId int64, messageId int) (err error)

	// Получение данных пользователя для статистики.
	GetRangeDataForStatistics(userId int64, typeAdEvent TypeAdEvent, startDate, endDate time.Time) (data *DataForStatistics, err error)
}

type CbqDataForCbqAdEventViewSelect struct {
	StartDate      time.Time   // Начальная дата событий.
	EndDate        time.Time   // Конечная дата событий.
	TypeAdEvent    TypeAdEvent // Тип событий.
	PageForDisplay int         // Страница для отображения.
}
