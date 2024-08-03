package date

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const YYYYMMDD = "20060102"
const HOURS_PER_DAY = 24

func ConvertSliceToInt(s []string) ([]int, error) {
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
		if date.After(now) {
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
		if date.After(now) {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func nextDateForWeek(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("repeat for weeks has wrong format")
	}

	weekdays_s, err := ConvertSliceToInt(strings.Split(rules[1], ","))
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
		if date.After(now) && weekdays[date.Weekday()] {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func nextDateForMonth(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) < 2 || len(rules) > 3 {
		return "", errors.New("repeat for weeks has wrong format")
	}

	days_s, err := ConvertSliceToInt(strings.Split(rules[1], ","))
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
		months_s, err := ConvertSliceToInt(strings.Split(rules[2], ","))
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

		if date.After(now) {
			return date.Format(YYYYMMDD), nil
		}
	}
}

func NextDate(now time.Time, date time.Time, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	rules := strings.Split(repeat, " ")
	if rules[0] == "d" {
		return nextDateForDay(now, date, rules)
	} else if rules[0] == "y" {
		return nextDateForYear(now, date, rules)
	} else if rules[0] == "w" {
		return nextDateForWeek(now, date, rules)
	} else if rules[0] == "m" {
		return nextDateForMonth(now, date, rules)
	}
	return "", errors.New("wrong repeat format")
}
