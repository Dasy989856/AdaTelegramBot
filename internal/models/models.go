package models

import "fmt"

// Ошибки.
var (
	ErrUserNotFound = fmt.Errorf("user not found")
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

// Событие.
type Event struct {
	Id                   int64  `json:"id" db:"id"`
	UserId               int64  `json:"userId" db:"user_id"`
	Type                 string `json:"type" db:"name"`                                   // Тип события. (sale, buy)
	CreatedAt            string `json:"createdAt" db:"created_at"`                        // Дата создания события.
	PostingDate          string `json:"postingDate" db:"posting_date"`                    // Дата постинга.
	PartnerURL           string `json:"partnerName" db:"partner_url"`                     // Имя партнера. (Продавец / Покупатель)
	Price                int    `json:"price" db:"price"`                                 // Цена.
	ArrivalOfSubscribers int    `json:"arrivalOfSubscribers" db:"arrival_of_subscribers"` // Приход подписчиков.
}

// БД для телеграмм бота.
type TelegramBotDB interface {
	// Получение данных пользователя.
	GetUserData(userId int64) (user *User, err error)
	// Создание пользователя.
	DefaultUserCreation(chatId int64, userUrl, firstName string) error
	// Добавление события.
	AdEventCreation(event *Event) (int64, error)
	// Удаление события.
	AdEventDelete(eventId int64) error
	// Установка шага пользователя.
	SetStepUser(userId int64, step string) error
	// Получение текущего шага пользователя.
	GetStepUser(userId int64) (step string, err error)
	// Закрытие БД.
	Close() error
}
