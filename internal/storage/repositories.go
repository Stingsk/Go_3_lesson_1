package storage

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	TrainingUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("training '%s' not found", e.TrainingUUID)
}

type Repository interface {
	AddTraining(ctx context.Context) error

	GetTraining(ctx context.Context) error

	UpdateTraining(
		ctx context.Context,
		trainingUUID string,
		updateFn func(ctx context.Context) error,
	) error
}
