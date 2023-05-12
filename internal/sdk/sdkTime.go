package sdk

import (
	"fmt"
	"time"
)

// Парсинг даты в time.Time
func ParseDateToTime(timeString string) (*time.Time, error) {
	// layout := "2006-01-02T15:04:00+03:00"
	layout := "02.01.2006 15:04"

	t, err := time.Parse(layout, timeString)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	return &t, nil
}

// Парсинг time.Time в дату.
func ParseTimeToDate(time *time.Time) string {
	return time.Format("02.01.2006 15:04")
}

// Парсинг time.Time в диапозон времени.
func ParseTimeToRangeDate(timeStart, timeEnd *time.Time) (rangeDate string) {
	return timeStart.Format("02.01.2006 15:04") + ";" + timeEnd.Format("02.01.2006 15:04")
}
// Возвращает метки времени, начало и конец вчерашнего дня.
func GetTimeRangeYesterday() (start, end *time.Time) {
	// Получаем текущую метку времени
	now := time.Now()

	// Вычитаем 1 день, чтобы получить метку времени вчерашнего дня
	yesterday := now.AddDate(0, 0, -1)

	// Получаем метку времени начала вчерашнего дня
	startOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

	// Получаем метку времени конца вчерашнего дня
	endOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), yesterday.Location())

   return &startOfYesterday, &endOfYesterday
}

// Возвращает метки времени, начало и конец текущего дня.
func GetTimeRangeToday() (start, end *time.Time) {
	// Получение текущей метки времени.
	now := time.Now()

	// Получение метки времени начала текущего дня.
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Получение метки времени конца текущего дня.
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)

	return &startOfDay, &endOfDay
}

// Возвращает метки времени, начало и конец завтрешнего дня.
func GetTimeRangeTomorrow() (start, end *time.Time) {
	// Получаем текущую метку времени
    now := time.Now()

    // Добавляем 1 день, чтобы получить метку времени завтрашнего дня
    tomorrow := now.AddDate(0, 0, 1)

    // Получаем метку времени начала завтрашнего дня
    startOfTomorrow := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

    // Получаем метку времени конца завтрашнего дня
    endOfTomorrow := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), tomorrow.Location())

	return &startOfTomorrow, &endOfTomorrow
}

// Возвращает метки времени, начало и конец предыдущей недели.
func GetTimeRangeLastWeek() (start, end *time.Time) {
	// Получаем текущую метку времени
    now := time.Now()
    // Определяем первый день недели и часовой пояс
    firstDay := time.Monday

    // Вычисляем метку времени начала текущей недели
    startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-firstDay), 0, 0, 0, 0, now.Location())

    // Вычисляем метку времени начала предыдущей недели
    startOfLastWeek := startOfWeek.AddDate(0, 0, -7)

    // Вычисляем метку времени конца предыдущей недели
    endOfLastWeek := startOfWeek.Add(-time.Nanosecond)

	return &startOfLastWeek, &endOfLastWeek
}

// Возвращает метки времени, начало и конец текущей недели.
func GetTimeRangeThisWeek() (start, end *time.Time) {
	// Получаем текущую метку времени и часовой пояс
    now := time.Now()
    // Определяем первый день недели
    firstDay := time.Monday

    // Вычисляем метку времени начала текущей недели
    startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-firstDay), 0, 0, 0, 0, now.Location())

    // Вычисляем метку времени конца текущей недели
    endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	return &startOfWeek, &endOfWeek
}

// Возвращает метки времени, начало и конец следующей недели.
func GetTimeRangeNextWeek() (start, end *time.Time) {
	// Получаем текущую метку времени и часовой пояс
    now := time.Now()
    // Определяем первый день недели
    firstDay := time.Monday

    // Вычисляем метку времени начала следующей недели
    startOfNextWeek := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-firstDay)+7, 0, 0, 0, 0, now.Location())

    // Вычисляем метку времени конца следующей недели
    endOfNextWeek := startOfNextWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	return &startOfNextWeek, &endOfNextWeek
}

// Возвращает метки времени, начало и конец предыдущего месяца.
func GetTimeRangeLastMonth() (start, end *time.Time) {
	// Получение текущей метки времени.
	now := time.Now()

	// Получение метки времени первого дня текущего месяца.
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени первого дня предыдущего месяца.
	firstOfLastMonth := firstOfMonth.AddDate(0, -1, 0)

	// Получение метки времени начала предыдущего месяца.
	startOfLastMonth := firstOfLastMonth.AddDate(0, 0, 0)

	// Получение метки времени конца предыдущего месяца.
	endOfLastMonth := firstOfMonth.Add(-time.Nanosecond)

	return &startOfLastMonth, &endOfLastMonth
}

// Возвращает метки времени, начало и конец текущего месяца.
func GetTimeRangeThisMonth() (start, end *time.Time) {
	// Получение текущей метки времени
	now := time.Now()

	// Получение первого дня текущего месяца.
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени начала текущего месяца.
	startOfMonth := firstOfMonth.AddDate(0, 0, 0)

	// Получение метки времени конца текущего месяца.
	endOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return &startOfMonth, &endOfMonth
}

// Возвращает метки времени, начало и конец следующего месяца.
func GetTimeRangeNextMonth() (start, end *time.Time) {
	// Получение текущей метки времени.
	now := time.Now()

	// Получение метки времени первого дня текущего месяца.
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени первого дня следующего месяца.
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	// Получение метки времени начала следующего месяца.
	startOfNextMonth := firstOfNextMonth.AddDate(0, 0, 0)

	// Получение метки времени конца следующего месяца.
	endOfNextMonth := firstOfNextMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return &startOfNextMonth, &endOfNextMonth
}

// Возвращает метки времени, начало и конец предыдущего года.
func GetTimeRangeLastYear() (start, end *time.Time) {
	// Получение текущей метки времени
	now := time.Now()

	// Получение метки времени начала предыдущего года
	startOfLastYear := time.Date(now.Year()-1, time.January, 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени конца предыдущего года
	endOfLastYear := time.Date(now.Year()-1, time.December, 31, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	return &startOfLastYear, &endOfLastYear
}

// Возвращает метки времени, начало и конец текущего года.
func GetTimeRangeThisYear() (start, end *time.Time) {
	// Получение текущей метки времени
	now := time.Now()

	// Получение метки времени начала текущего года
	startOfYear := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени конца текущего года
	endOfYear := time.Date(now.Year(), time.December, 31, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	return &startOfYear, &endOfYear
}

// Возвращает метки времени, начало и конец следующего года.
func GetTimeRangeNextYear() (start, end *time.Time) {
	// Получение текущей метки времени
	now := time.Now()

	// Получение метки времени начала следующего года
	startOfNextYear := time.Date(now.Year()+1, time.January, 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени конца следующего года
	endOfNextYear := time.Date(now.Year()+1, time.December, 31, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	return &startOfNextYear, &endOfNextYear
}

// Возвращает метки времени, начало и конец указанного месяца текущего года.
func GetTimeRangeSelectedMonthThisYes(month time.Month) (start, end *time.Time) {
	// Получение текущей метки времени
	now := time.Now()

	// Получение метки времени начала указанного месяца текущего года
	startOfMonth := time.Date(now.Year(), month, 1, 0, 0, 0, 0, now.Location())

	// Получение метки времени конца указанного месяца текущего года
	endOfMonth := time.Date(now.Year(), month+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	return &startOfMonth, &endOfMonth
}