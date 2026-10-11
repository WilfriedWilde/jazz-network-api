package repository

import (
	"fmt"
	"strings"
	"time"
)

func buildIDQueryArguments(query string, ids []int64) (string, []any) {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	return fmt.Sprintf(query, strings.Join(placeholders, ", ")), args
}

func getStopDate(start string, i int) (string, error) {
	startDate, err := time.Parse(time.DateOnly, start)
	if err != nil {
		return "", err
	}

	date := startDate.AddDate(0, 0, i)

	return date.Format(time.DateOnly), nil
}
