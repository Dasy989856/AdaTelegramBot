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
