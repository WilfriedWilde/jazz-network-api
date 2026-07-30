package repository

import "fmt"

func formatDistinct(distinct bool) string {
	if distinct {
		return "DISTINCT"
	}

	return ""
}

func formatJoin(original string, tablesToJoin []string) string {
	var s string

	if len(tablesToJoin) > 0 {
		for _, t := range tablesToJoin {
			s += fmt.Sprintf("JOIN %s ON %s.%s_id = %s.id", original, t, t)
		}
	}

	return s
}

func formatRowFilters(filters map[string]any) string {
	var s string

	if len(filters) > 0 {
		for key, value := range filters {
			s += fmt.Sprintf("WHERE %s = %v", key, value)
		}
	}

	return s
}
