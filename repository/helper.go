package repository

import "strings"

func getRowSlice(s string) []string {
	split := strings.Split(s, ";")
	rowSlice := make([]string, len(split))

	for i := range split {
		rowSlice[i] = strings.TrimSpace(strings.ToLower(split[i]))
	}

	return rowSlice
}
