package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

// NextDate возвращает ближайшую дату выполнения, строго большую now.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата задачи: %w", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("правило повторения не задано")
	}

	switch parts[0] {
	case "d":
		return nextByDays(now, date, parts)
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("некорректное ежегодное правило")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateLayout), nil
			}
		}
	case "w":
		return nextByWeekdays(now, date, parts)
	case "m":
		return nextByMonthDays(now, date, parts)
	default:
		return "", fmt.Errorf("неподдерживаемое правило повторения")
	}
}

func afterNow(date, now time.Time) bool {
	year, month, day := now.Date()
	return date.After(time.Date(year, month, day, 0, 0, 0, 0, date.Location()))
}

func nextByDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("некорректное правило дней")
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval < 1 || interval > 400 {
		return "", fmt.Errorf("некорректный интервал дней")
	}
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			return date.Format(dateLayout), nil
		}
	}
}

func nextByWeekdays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("некорректное правило недель")
	}
	weekdays, err := parseValues(parts[1], 1, 7)
	if err != nil {
		return "", fmt.Errorf("некорректный день недели: %w", err)
	}
	for {
		date = date.AddDate(0, 0, 1)
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		if weekdays[weekday] && afterNow(date, now) {
			return date.Format(dateLayout), nil
		}
	}
}

func nextByMonthDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("некорректное правило месяцев")
	}

	days := make(map[int]bool)
	for _, item := range strings.Split(parts[1], ",") {
		value, err := strconv.Atoi(item)
		if err != nil || value == 0 || value < -2 || value > 31 {
			return "", fmt.Errorf("некорректный день месяца")
		}
		days[value] = true
	}

	months := make(map[int]bool)
	if len(parts) == 3 {
		var err error
		months, err = parseValues(parts[2], 1, 12)
		if err != nil {
			return "", fmt.Errorf("некорректный месяц: %w", err)
		}
	} else {
		for month := 1; month <= 12; month++ {
			months[month] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) || !months[int(date.Month())] {
			continue
		}
		day := date.Day()
		lastDay := date.AddDate(0, 1, -day).Day()
		if days[day] || (days[-1] && day == lastDay) || (days[-2] && day == lastDay-1) {
			return date.Format(dateLayout), nil
		}
	}
}

func parseValues(value string, min, max int) (map[int]bool, error) {
	if value == "" {
		return nil, fmt.Errorf("пустое значение")
	}
	values := make(map[int]bool)
	for _, item := range strings.Split(value, ",") {
		parsed, err := strconv.Atoi(item)
		if err != nil || parsed < min || parsed > max {
			return nil, fmt.Errorf("значение вне допустимого диапазона")
		}
		values[parsed] = true
	}
	return values, nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if nowParam := r.FormValue("now"); nowParam != "" {
		parsedNow, err := time.Parse(dateLayout, nowParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	nextDate, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(nextDate))
}
