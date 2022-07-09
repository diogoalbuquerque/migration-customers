package usecase_test

import (
	"context"
	"errors"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	usecase "github.com/diogoalbuquerque/migration-customers/internal/usecase/person"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type test struct {
	name string
	mock func()
	res  interface{}
	err  interface{}
}

var successResult = []entity.LegacyPerson{createLegacyPerson(), createLegacyCompany(), createInvalidLegacyPerson()}
var errorStorePeople = errors.New("person_mysql - StorePeople - Mock Error")

func servicePerson(t *testing.T) (*usecase.PersonUseCase, *MockPersonMYSQLRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockPersonMYSQLRepo(mockCtl)

	personUseCase := usecase.NewPersonUseCase(repo)

	return personUseCase, repo
}

func Test_Migrate(t *testing.T) {
	t.Parallel()

	servicePerson, repo := servicePerson(t)

	tests := []test{
		{
			name: "success",
			mock: func() {
				gomock.InOrder(
					repo.EXPECT().GetBucketsAvailable().Return(2),
					repo.EXPECT().StorePeople(context.TODO(), gomock.Any()).Return(nil))
			},
			res: successResult[0:2],
			err: []error(nil),
		},
		{
			name: "error",
			mock: func() {
				gomock.InOrder(
					repo.EXPECT().GetBucketsAvailable().Return(2),
					repo.EXPECT().StorePeople(context.TODO(), gomock.Any()).Return(errorStorePeople))
			},
			res: []entity.LegacyPerson(nil),
			err: []error{errorStorePeople},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			tc.mock()

			res, errs := servicePerson.Migrate(context.TODO(), successResult)

			require.ElementsMatch(t, tc.res, res)
			require.Equal(t, tc.err, errs)

		})
	}
}
func createLegacyPerson() entity.LegacyPerson {
	birthDate := time.Date(2004, time.April, 6, 0, 0, 0, 0, time.UTC)
	return entity.LegacyPerson{NI: 88626159226, Name: "User One", BirthDate: &birthDate}
}

func createLegacyCompany() entity.LegacyPerson {
	return entity.LegacyPerson{NI: 87511532000171, Name: "Company One"}
}

func createInvalidLegacyPerson() entity.LegacyPerson {
	birthDate := time.Date(2004, time.April, 6, 0, 0, 0, 0, time.UTC)
	return entity.LegacyPerson{NI: 99988877766, Name: "User NI Invalid", BirthDate: &birthDate}
}
