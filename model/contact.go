package model

type Contact struct {
	ID          int64  `json:"id"`
	Method      string `json:"method"`
	Status      string `json:"status"`
	LastAttempt string `json:"last_attempt"`
	Notes       string `json:"notes"`
}
