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
		return false
	}

	if ae.Type == "" {
		return false
	}

	if ae.CreatedAt == "" {
		return false
	}

	if ae.DatePosting == "" {
		return false
	}

	if ae.DateDelete == "" {
		return false
	}

	if ae.Partner == "" {
		return false
	}

	if ae.Channel == "" {
		return false
	}

	if ae.Type != "barter" && ae.Price == 0 {
		return false
	}

	return true
}

// БД для телеграмм бота.
type TelegramBotDB interface {
	// Получение данных пользователя.
	GetUserData(userId int64) (user *User, err error)
	// Создание пользователя.
	DefaultUserCreation(chatId int64, userUrl, firstName string) error
	// Создание дефолтного ad события.
	AdEventCreation(event *AdEvent) (int64, error)
	// Удаление события.
	AdEventDelete(eventId int64) error
	// Обновление информации о приходе подписчиков.
	UpdateSubscribesInAdEvent(eventId, subscribers int64) error
	// Установка шага пользователя.
	SetStepUser(userId int64, step string) error
	// Получение текущего шага пользователя.
	GetStepUser(userId int64) (step string, err error)
	// Подучение id незавершенного ad события.
	GetUnfinishedAdEventId(userId int64) (eventId int64, err error)
	// Добавление messageId пользователя.
	AddUserMessageId(userId int64, messageId int) error
	// Удаление messageId пользователя.
	DeleteUserMessageId(messageId int) error
	// Возвращает список messageIds пользователя.
	GetUserMessageIds(userId int64) ([]int, error)
	// Закрытие БД.
	Close() error
}
