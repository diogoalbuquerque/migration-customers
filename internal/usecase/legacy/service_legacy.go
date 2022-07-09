package usecase

import (
	"context"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/internal/helper"
	"sync"
)

type LegacyUseCase struct {
	repo LegacyPersonDB2Repo
}

func NewLegacyUseCase(repo LegacyPersonDB2Repo) *LegacyUseCase {
	return &LegacyUseCase{repo: repo}
}

func (uc *LegacyUseCase) LoadDatasource(ctx context.Context, limit int) ([]entity.LegacyPerson, error) {
	return uc.repo.GetLegacyPeople(ctx, limit)
}

func (uc *LegacyUseCase) Reconciliation(ctx context.Context, legacyPeople []entity.LegacyPerson) ([]entity.LegacyPerson, []error) {
	buckets := helper.SplitBetweenAvailableBuckets(legacyPeople, uc.repo.GetBucketsAvailable())

	ch := make(chan entity.ChanelLegacyPerson)
	var wg sync.WaitGroup

	for i, legacyPeopleOnBuckets := range buckets {
		wg.Add(1)
		go func(idxBucket int, legacyPeopleOnBuckets []entity.LegacyPerson, ch chan entity.ChanelLegacyPerson) {

			defer wg.Done()

			err := uc.repo.UpdateLegacyPeople(ctx, legacyPeopleOnBuckets)

			ch <- entity.ChanelLegacyPerson{
				LegacyPeople: legacyPeopleOnBuckets,
				Err:          err,
			}

		}(i, legacyPeopleOnBuckets, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var reconciledPeople []entity.LegacyPerson = nil
	var bucketsErrors []error = nil

	for lec := range ch {
		if lec.Err != nil {
			bucketsErrors = append(bucketsErrors, lec.Err)
		} else {
			reconciledPeople = append(reconciledPeople, lec.LegacyPeople...)
		}
	}

	return reconciledPeople, bucketsErrors
}
