package utilities

import (
	"strconv"
	"strings"
)

// ValidateIPAddress validates an IP address
// An IP address is a string in the form "A.B.C.D" where the value of A, B, C, and D are integers between 0 and 255
func ValidateIPAddress(ip string) ([]int, bool) {
	// split the IP address by dots
	IPString := strings.Split(ip, ".")
	// if the length of the split is not 4, return false
	if len(IPString) != 4 {
		return nil, false
	}
	// for each part of the split, convert it to an integer
	ipParts := make([]int, 4)
	for i, part := range IPString {
		// if the conversion fails, return false
		var err error
		ipParts[i], err = strconv.Atoi(part)
		if err != nil {
			return nil, false

		}
		// if the integer is not between 0 and 255, return false
		if ipParts[i] < 0 || ipParts[i] > 255 {
			return nil, false

		}
	}
	return ipParts, true
}
