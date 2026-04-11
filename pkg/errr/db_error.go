package errr

import (
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

// ParseDBError converts database errors into application errors.
func ParseDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return New(http.StatusConflict, "resource already exists")
		case pgerrcode.ForeignKeyViolation:
			return New(http.StatusBadRequest, "referenced resource does not exist")
		case pgerrcode.NotNullViolation:
			return New(http.StatusBadRequest, "required field is missing")
		case pgerrcode.CheckViolation:
			return New(http.StatusBadRequest, "field value violates check constraint")
		case pgerrcode.InvalidTextRepresentation:
			return New(http.StatusBadRequest, "invalid input format")
		default:
			return New(http.StatusInternalServerError, "database error")
		}
	}

	return err
}
