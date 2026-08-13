package pgsql

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/contract"
)

const (
	errCodeUniqueViolation      = "23505"
	errCodeForeignKeyViolation  = "23503"
	errCodeCheckViolation       = "23514"
	errCodeNotNullViolation     = "23502"
	errCodeSerializationFailure = "40001"
	errCodeDeadlockDetected     = "40P01"

	errClassConnectionException = "08"
)

func TranslateError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %v", contract.ErrNotFound, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch {
		case pgErr.Code == errCodeUniqueViolation:
			return fmt.Errorf("%w: %v", contract.ErrDuplicateKey, err)
		case pgErr.Code == errCodeForeignKeyViolation || pgErr.Code == errCodeCheckViolation || pgErr.Code == errCodeNotNullViolation:
			return fmt.Errorf("%w: %v", contract.ErrConstraintViolation, err)
		case pgErr.Code == errCodeSerializationFailure || pgErr.Code == errCodeDeadlockDetected:
			return fmt.Errorf("%w: %v", contract.ErrConflict, err)
		case strings.HasPrefix(pgErr.Code, errClassConnectionException):
			return fmt.Errorf("%w: %v", contract.ErrConnection, err)
		}
		return err
	}

	var netErr net.Error
	if errors.As(err, &netErr) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", contract.ErrConnection, err)
	}

	return err
}
