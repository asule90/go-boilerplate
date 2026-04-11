package errr

import (
	"errors"
	"strings"
)

var (
	ErrNoRows       = errors.New("record not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal server error")
)

// GetLastNErrorMessage splits err.Error() by "::" and returns the last N segments joined by ": ".
func GetLastNErrorMessage(err error, n int) string {
	if err == nil {
		return ""
	}
	parts := strings.Split(err.Error(), "::")
	if n >= len(parts) {
		return err.Error()
	}
	return strings.Join(parts[len(parts)-n:], ": ")
}
