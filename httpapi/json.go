package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"jazz-network-api/model"
	"log"
	"mime"
	"net/http"
	"strings"
)

const maxBodyBytes = 1 << 20

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("error encoding json: %s", err.Error())
	}
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var zero T

	// Require a JSON Content-Type while allowing parameters such as charset=utf-8.
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.EqualFold(mediaType, "application/json") {
		return zero, &model.RequestError{
			Status:  http.StatusUnsupportedMediaType,
			Code:    "unsupported_media_type",
			Message: "Content-Type must be application/json",
		}
	}

	// Limit how much of the request body can be read to prevent oversized payloads.
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	// Reject fields that are not declared by T, catching misspelled or unsupported input.
	decoder.DisallowUnknownFields()

	// Decode into a pointer so that JSON null can be distinguished from a valid zero value.
	var input *T
	if err := decoder.Decode(&input); err != nil {
		// EOF on the first decode means that the body is empty or contains only whitespace.
		if errors.Is(err, io.EOF) {
			return zero, &model.RequestError{
				Status:  http.StatusBadRequest,
				Code:    "malformed_request",
				Message: "request body must not be empty",
			}
		}

		// MaxBytesReader reports an oversized body while the decoder is reading it.
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return zero, &model.RequestError{
				Status:  http.StatusRequestEntityTooLarge,
				Code:    "malformed_request",
				Message: "request body is too large",
			}
		}

		// UnmarshalTypeError identifies a known field, or the request root,
		// whose JSON value cannot be decoded into the expected Go type.
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			field := typeErr.Field
			if field == "" {
				field = "request body"
			}

			return zero, &model.RequestError{
				Status:  http.StatusBadRequest,
				Code:    "malformed_request",
				Message: field + " has an invalid JSON type",
			}
		}

		// DisallowUnknownFields reports unsupported properties with this stable
		// encoding/json prefix; expose only the quoted property name, not the full error.
		const unknownFieldPrefix = "json: unknown field "
		if strings.HasPrefix(err.Error(), unknownFieldPrefix) {
			field := strings.TrimPrefix(err.Error(), unknownFieldPrefix)
			return zero, &model.RequestError{
				Status:  http.StatusBadRequest,
				Code:    "malformed_request",
				Message: "unknown property " + field,
			}
		}

		// SyntaxError covers invalid JSON tokens, while ErrUnexpectedEOF covers
		// otherwise valid JSON that ends before the current value is complete.
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) {
			return zero, &model.RequestError{
				Status:  http.StatusBadRequest,
				Code:    "malformed_request",
				Message: "request body contains malformed JSON",
			}
		}

		// Keep unrecognized decoder failures generic so internal diagnostics are not exposed.
		return zero, &model.RequestError{
			Status:  http.StatusBadRequest,
			Code:    "malformed_request",
			Message: "request body contains invalid JSON",
		}
	}

	// Reject the JSON literal null, which otherwise decodes successfully into a nil pointer.
	if input == nil {
		return zero, &model.RequestError{
			Status:  http.StatusBadRequest,
			Code:    "malformed_request",
			Message: "request body must contain one JSON object",
		}
	}

	// A second decode must reach EOF; otherwise the body contains another value or trailing data.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return zero, &model.RequestError{
			Status:  http.StatusBadRequest,
			Code:    "malformed_request",
			Message: "request body must contain one JSON object",
		}
	}

	// The body contains exactly one non-null JSON value matching T.
	return *input, nil
}
