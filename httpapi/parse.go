package httpapi

import (
	"net/http"
	"strconv"
)

func parsePathID(pathID string) (int64, error) {
	id, err := strconv.ParseInt(pathID, 10, 64)
	if err != nil {
		return 0, &requestError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: "id must be a positive integer",
		}
	}

	return id, nil
}
