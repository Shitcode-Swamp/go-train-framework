package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	errCodeDocumentValidationFailed = 121
	errCodeWriteConflict            = 112
)

func FromMongoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("%w: %w", store.ErrNotFound, err)
	}
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("%w: %w", store.ErrDuplicateKey, err)
	}

	var serverErr mongo.ServerError
	if errors.As(err, &serverErr) {
		switch {
		case serverErr.HasErrorCode(errCodeDocumentValidationFailed):
			return fmt.Errorf("%w: %w", store.ErrConstraintViolation, err)
		case serverErr.HasErrorCode(errCodeWriteConflict) || serverErr.HasErrorLabel("TransientTransactionError"):
			return fmt.Errorf("%w: %w", store.ErrConflict, err)
		}
	}

	if mongo.IsNetworkError(err) || mongo.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %w", store.ErrConnection, err)
	}

	return err
}
