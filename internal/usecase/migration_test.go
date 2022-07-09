package usecase_test

import (
	"context"
	"errors"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/internal/usecase"
	"github.com/diogoalbuquerque/migration-customers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	_mockLimit    = 10
	_mockLogLevel = "error"
)

type test struct {
	name string
	mock func()
	err  interface{}
}

var successResult = []entity.LegacyPerson{createLegacyPerson(), createLegacyCompany()}
var errorLoadDatasource = errors.New("legacy_person_db2 - GetLegacyPeople - Mock Error")
var errorStorePeople = []error{errors.New("person_mysql - StorePeople - [48815782000143] - Mock Error")}
var errorMigrate = errors.New("run - Migrate")
var errorUpdateLegacyPeople = []error{errors.New("legacy_person_db2 - UpdateLegacyPeople - [48815782000143] - Mock Error")}
var errorReconciliation = errors.New("run - Reconciliation")

func serviceMigration(t *testing.T) (*usecase.MigrationUseCase, *MockServiceLegacy, *MockServicePerson) {
	t.Helper()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	serviceLegacy := NewMockServiceLegacy(mockCtl)
	servicePerson := NewMockServicePerson(mockCtl)

	migrationUseCase := usecase.NewMigrationUseCase(serviceLegacy, servicePerson, logger.New(_mockLogLevel))

	return migrationUseCase, serviceLegacy, servicePerson
}

func Test_Migration(t *testing.T) {
	t.Parallel()
	ctx, _ := xray.BeginSegment(context.TODO(), "CONTEXT_TEST")
	migration, serviceLegacy, servicePerson := serviceMigration(t)

	tests := []test{
		{
			name: "datasource error",
			mock: func() {
				serviceLegacy.EXPECT().LoadDatasource(ctx, _mockLimit).Return(nil, errorLoadDatasource)
			},
			err: errorLoadDatasource,
		},
		{
			name: "datasource empty",
			mock: func() {
				serviceLegacy.EXPECT().LoadDatasource(ctx, _mockLimit).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "migrate error",
			mock: func() {
				gomock.InOrder(
					serviceLegacy.EXPECT().LoadDatasource(ctx, _mockLimit).Return(successResult, nil),
					servicePerson.EXPECT().Migrate(ctx, successResult).Return(nil, errorStorePeople))

			},
			err: errorMigrate,
		},
		{
			name: "reconciliation error",
			mock: func() {
				gomock.InOrder(
					serviceLegacy.EXPECT().LoadDatasource(ctx, _mockLimit).Return(successResult, nil),
					servicePerson.EXPECT().Migrate(ctx, successResult).Return(successResult, errorStorePeople),
					serviceLegacy.EXPECT().Reconciliation(ctx, successResult).Return(nil, errorUpdateLegacyPeople))

			},
			err: errorReconciliation,
		},
		{
			name: "success",
			mock: func() {
				gomock.InOrder(
					serviceLegacy.EXPECT().LoadDatasource(ctx, _mockLimit).Return(successResult, nil),
					servicePerson.EXPECT().Migrate(ctx, successResult).Return(successResult, nil),
					serviceLegacy.EXPECT().Reconciliation(ctx, successResult).Return(successResult, nil))

			},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			tc.mock()

			err := migration.Migration(ctx, _mockLimit)

			require.Equal(t, tc.err, err)

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
