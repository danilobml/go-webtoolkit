package gowebtoolkit

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

type Tools struct {
	MaxJsonSize        int
	AllowUnknownFields bool
}

// GenerateRandomString returns a random string of the given size using the characters of randomStringSource
func (t *Tools) GenerateRandomString(size int) string {
	s, r := make([]rune, size), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}

// JSON Operations
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576
	if t.MaxJsonSize != 0 {
		maxBytes = t.MaxJsonSize
	}
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	if !t.AllowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only one JSON")
	}

	return nil
}

func (t *Tools) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload = jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	t.WriteJSON(w, statusCode, payload)

	return nil
}
