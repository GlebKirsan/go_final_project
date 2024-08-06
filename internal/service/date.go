package service

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const YYYYMMDD = "20060102"
const HOURS_PER_DAY = 24

type DateService struct {
}

func NewDateService() *DateService {
	return &DateService{}
}

func convertSliceToInt(s []string) ([]int, error) {
	result := []int{}
	for _, element_s := range s {
		element, err := strconv.Atoi(element_s)
		if err != nil {
			return []int{}, err
		}
		result = append(result, element)
	}
	return result, nil
}

func nextDateForDay(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("repeat for days has wrong format")
	}

	digit, err := strconv.Atoi(rules[1])
	if err != nil {
		return "", errors.New("cannot parse day number")
	}
	if digit > 400 {
		return "", errors.New("days number too large")
	}

	for {
		date = date.AddDate(0, 0, digit)
		if After(date, now) {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func nextDateForYear(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) > 1 {
		return "", errors.New("repeat for years has wrong format")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if After(date, now) {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func nextDateForWeek(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("repeat for weeks has wrong format")
	}

	weekdays_s, err := convertSliceToInt(strings.Split(rules[1], ","))
	if err != nil {
		return "", err
	}
	weekdays := map[time.Weekday]bool{}
	for _, day := range weekdays_s {
		if day < 1 || day > 7 {
			return "", errors.New("repeat for weeks has wrong format")
		}
		weekdays[time.Weekday(day%7)] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if After(date, now) && weekdays[date.Weekday()] {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func nextDateForMonth(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) < 2 || len(rules) > 3 {
		return "", errors.New("repeat for weeks has wrong format")
	}

	days_s, err := convertSliceToInt(strings.Split(rules[1], ","))
	if err != nil {
		return "", err
	}

	days := map[int]bool{}
	for _, day := range days_s {
		if day < -2 || day > 31 {
			return "", errors.New("month repeat format error: day out of bounds")
		}
		days[day] = true
	}

	months := map[time.Month]bool{}
	if len(rules) == 3 {
		months_s, err := convertSliceToInt(strings.Split(rules[2], ","))
		if err != nil {
			return "", nil
		}

		for _, month := range months_s {
			if month < 1 || month > 12 {
				return "", errors.New("month repeat format error: month out of number")
			}
			months[time.Month(month)] = true
		}
	} else {
		for i := range 12 {
			months[time.Month(i+1)] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if !months[date.Month()] {
			date = date.AddDate(0, 1, 0)
			date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
			continue
		}

		lastMonthDay := date.AddDate(0, 1, -date.Day()+1)
		toMonthEnd := date.Sub(lastMonthDay).Hours() / HOURS_PER_DAY
		if !days[date.Day()] && !days[int(toMonthEnd)] {
			continue
		}

		if After(date, now) {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func (service *DateService) NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	parsed, err := time.Parse(YYYYMMDD, date)
	if err != nil {
		return "", err
	}

	rules := strings.Split(repeat, " ")
	if rules[0] == "d" {
		return nextDateForDay(now, parsed, rules)
	} else if rules[0] == "y" {
		return nextDateForYear(now, parsed, rules)
	} else if rules[0] == "w" {
		return nextDateForWeek(now, parsed, rules)
	} else if rules[0] == "m" {
		return nextDateForMonth(now, parsed, rules)
	}
	return "", errors.New("wrong repeat format")
}

func Before(d1 string, d2 string) bool {
	return d1 < d2
}

func (service *DateService) Before(d1 string, d2 string) bool {
	return Before(d1, d2)
}

func After(d1 time.Time, d2 time.Time) bool {
	return d1.Truncate(24 * time.Hour).After(d2.Truncate(24 * time.Hour))
}

func Parse(date string) (time.Time, error) {
	return time.Parse(YYYYMMDD, date)
}
