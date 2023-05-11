package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {

	sec := d / time.Second
	h := sec / 3600
	m := sec/60 - h*60
	s := sec - h*3600 - m*60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func ParseClock(t time.Time) (time.Time, error) {
	h, m, s := t.Clock()
	parsed, err := time.Parse("15:04:05", fmt.Sprintf("%02d:%02d:%02d", h, m, s))
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func TodayRemainder(t time.Time) time.Duration {
	tomorrow := t.Add(time.Hour * 24)
	layout := "2006-01-02"
	tomorrow, _ = time.Parse(layout, tomorrow.Format(layout))
	return tomorrow.Sub(t)
}
