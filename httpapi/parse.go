package httpapi

import (
	"fmt"
	"jazz-network-api/model"
	"net/http"
	"net/url"
	"strconv"
)

var (
	deleteKeysQuery = map[string]struct{}{
		"force": {},
	}
)

func parsePathID(pathID string) (int64, error) {
	id, err := strconv.ParseInt(pathID, 10, 64)
	if err != nil {
		return 0, &model.RequestError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: "id must be a positive integer",
		}
	}

	return id, nil
}

func parseStrictQuery(r *http.Request, allowedKeys map[string]struct{}) (url.Values, error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("malformed query: %w", err)
	}

	for k := range query {
		if _, found := allowedKeys[k]; !found {
			return nil, fmt.Errorf("invalid query key: %s", k)
		}
	}

	return query, nil
}

func parseBool(value, key string) (bool, error) {
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("%s should be a boolean", key)
	}
}

func parseDelete(q []string) (bool, error) {
	if len(q) > 1 {
		return false, fmt.Errorf("only one 'force' query parameter allowed")
	}

	var force bool
	var err error

	if len(q) > 0 {
		force, err = parseBool(q[0], "force")
		if err != nil {
			return false, err
		}
	}

	return force, nil
}
