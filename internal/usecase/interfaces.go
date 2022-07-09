package usecase

import (
	"context"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
)

type (
	ServiceLegacy interface {
		LoadDatasource(ctx context.Context, limit int) ([]entity.LegacyPerson, error)
		Reconciliation(ctx context.Context, legacyPeople []entity.LegacyPerson) ([]entity.LegacyPerson, []error)
	}

	ServicePerson interface {
		Migrate(ctx context.Context, legacyPeople []entity.LegacyPerson) ([]entity.LegacyPerson, []error)
	}
)
