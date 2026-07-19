package apperr

import "strconv"

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ParseQueryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return n
}
