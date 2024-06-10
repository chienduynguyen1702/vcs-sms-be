package services

import "strconv"

func parseStringToInt(id string) (int, error) {
	return strconv.Atoi(id)
}

func parseStringToUint(id string) (uint, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return uint(intID), nil
}
