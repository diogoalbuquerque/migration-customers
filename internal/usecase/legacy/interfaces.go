package usecase

import (
	"context"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
)

type (
	LegacyPersonDB2Repo interface {
		GetLegacyPeople(ctx context.Context, limit int) ([]entity.LegacyPerson, error)
		UpdateLegacyPeople(ctx context.Context, legacyPeople []entity.LegacyPerson) error
		GetBucketsAvailable() int
	}
)
