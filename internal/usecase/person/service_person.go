package usecase

import (
	"context"
	"fmt"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/internal/helper"
	"github.com/klassmann/cpfcnpj"
	"sync"
)

type PersonUseCase struct {
	repo PersonMYSQLRepo
}

func NewPersonUseCase(repo PersonMYSQLRepo) *PersonUseCase {
	return &PersonUseCase{repo: repo}
}

func (uc *PersonUseCase) Migrate(ctx context.Context, legacyPeople []entity.LegacyPerson) ([]entity.LegacyPerson, []error) {
	buckets := helper.SplitBetweenAvailableBuckets(legacyPeople, uc.repo.GetBucketsAvailable())

	ch := make(chan entity.ChanelLegacyPerson)
	var wg sync.WaitGroup

	for i, legacyPeopleOnBuckets := range buckets {
		wg.Add(1)
		go func(idxBucket int, legacyPeopleOnBuckets []entity.LegacyPerson, ch chan entity.ChanelLegacyPerson) {
			defer wg.Done()

			var people []entity.Person = nil
			var convertedLegacyPeople []entity.LegacyPerson = nil

			for _, legacyPerson := range legacyPeopleOnBuckets {
				person := convertLegacyPersonToPerson(legacyPerson)
				if person != nil {
					convertedLegacyPeople = append(convertedLegacyPeople, legacyPerson)
					people = append(people, *person)
				}
			}

			if people != nil {
				err := uc.repo.StorePeople(ctx, people)
				ch <- entity.ChanelLegacyPerson{
					LegacyPeople: convertedLegacyPeople,
					Err:          err,
				}
			}
		}(i, legacyPeopleOnBuckets, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var migratedLegacyPerson []entity.LegacyPerson = nil
	var bucketsErrors []error = nil

	for lec := range ch {
		if lec.Err != nil {
			bucketsErrors = append(bucketsErrors, lec.Err)
		} else {
			migratedLegacyPerson = append(migratedLegacyPerson, lec.LegacyPeople...)
		}
	}

	return migratedLegacyPerson, bucketsErrors
}

func convertLegacyPersonToPerson(legacyPerson entity.LegacyPerson) *entity.Person {

	cpf := fmt.Sprintf("%011.0f", legacyPerson.NI)
	if cpfcnpj.ValidateCPF(cpf) && legacyPerson.BirthDate != nil {
		return &entity.Person{
			NI:         cpf,
			Name:       legacyPerson.Name,
			PersonType: entity.PF,
			BirthDate:  legacyPerson.BirthDate,
		}
	} else {
		cnpj := fmt.Sprintf("%014.0f", legacyPerson.NI)
		if cpfcnpj.ValidateCNPJ(cpf) && legacyPerson.BirthDate == nil {
			return &entity.Person{
				NI:         cnpj,
				Name:       legacyPerson.Name,
				PersonType: entity.PJ,
			}
		}
		return nil
	}
}
