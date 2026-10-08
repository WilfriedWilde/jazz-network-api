package repository

import (
	"fmt"
	"strings"
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
