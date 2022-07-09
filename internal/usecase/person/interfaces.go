package usecase

import (
	"context"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
)

type (
	PersonMYSQLRepo interface {
		StorePeople(ctx context.Context, people []entity.Person) error
		GetBucketsAvailable() int
	}
)
