package utilities

import (
	"strconv"
	"time"
)

func ParseStringToInt(id string) (int, error) {
	return strconv.Atoi(id)
}

func ParseStringToUint(id string) (uint, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return uint(intID), nil
}

func ParseUintToString(id uint) string {
	return strconv.Itoa(int(id))
}

// parse %40 to @
func ParseEmail(email string) string {
	return email
}

// format YYYY-MM-DD
func ParseDateToString(date time.Time) string {
	return date.Format("2006-01-02")
}

func ParseStringToDate(dateString string) time.Time {
	date, _ := time.Parse("2006-01-02", dateString)
	return date
}

func ParsePageAndLimit(page, limit string) (int, int) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	return pageInt, limitInt
}
