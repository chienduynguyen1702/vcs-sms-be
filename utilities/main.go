package utilities

import "strconv"

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