package usecase_test

import (
	"context"
	"errors"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	usecase "github.com/diogoalbuquerque/migration-customers/internal/usecase/legacy"
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

const (
	_mockLimit = 10
)

var successResult = []entity.LegacyPerson{createLegacyPerson(), createLegacyCompany()}
var errorGetLegacyPeople = errors.New("legacy_person_db2 - GetLegacyPeople - Mock Error")
var errorUpdateLegacyPeople = errors.New("legacy_person_db2 - UpdateLegacyPeople - Mock Error")

func serviceLegacy(t *testing.T) (*usecase.LegacyUseCase, *MockLegacyPersonDB2Repo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockLegacyPersonDB2Repo(mockCtl)

	legacyUseCase := usecase.NewLegacyUseCase(repo)

	return legacyUseCase, repo
}

func Test_LoadDatasource(t *testing.T) {
	t.Parallel()

	serviceLegacy, repo := serviceLegacy(t)

	tests := []test{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().GetLegacyPeople(context.TODO(), _mockLimit).Return(successResult, nil)
			},
			res: successResult,
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repo.EXPECT().GetLegacyPeople(context.TODO(), _mockLimit).Return(nil, errorGetLegacyPeople)
			},
			res: []entity.LegacyPerson(nil),
			err: errorGetLegacyPeople,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := serviceLegacy.LoadDatasource(context.TODO(), _mockLimit)

			require.Equal(t, tc.res, res)
			require.Equal(t, tc.err, err)
		})
	}

}

func Test_Reconciliation(t *testing.T) {
	t.Parallel()

	serviceLegacy, repo := serviceLegacy(t)

	tests := []test{
		{
			name: "success",
			mock: func() {
				gomock.InOrder(
					repo.EXPECT().GetBucketsAvailable().Return(2),
					repo.EXPECT().UpdateLegacyPeople(context.TODO(), []entity.LegacyPerson{createLegacyPerson()}).Return(nil))
			},
			res: []entity.LegacyPerson{createLegacyPerson()},
			err: []error(nil),
		},
		{
			name: "error",
			mock: func() {
				gomock.InOrder(
					repo.EXPECT().GetBucketsAvailable().Return(2),
					repo.EXPECT().UpdateLegacyPeople(context.TODO(), []entity.LegacyPerson{createLegacyPerson()}).Return(errorUpdateLegacyPeople))
			},
			res: []entity.LegacyPerson(nil),
			err: []error{errorUpdateLegacyPeople},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, errs := serviceLegacy.Reconciliation(context.TODO(), []entity.LegacyPerson{createLegacyPerson()})

			require.Equal(t, tc.res, res)
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
