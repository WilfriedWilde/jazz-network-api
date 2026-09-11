package model

import "fmt"

type ErrBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrFields struct {
	Code   string       `json:"code"`
	Errors []FieldError `json:"errors"`
}

type RequestError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *RequestError) Error() string {
	return e.Message
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Status int
	Code   string
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %d field error(s)", len(e.Fields))
}
