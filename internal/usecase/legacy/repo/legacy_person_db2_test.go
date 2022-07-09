package repo_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/internal/usecase/legacy/repo"
	"github.com/diogoalbuquerque/migration-customers/pkg/db2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	_defaultLimit = 10
)

func Test_GetLegacyPeople_Success(t *testing.T) {
	legacyPerson := createLegacyPerson()
	legacyCompany := createLegacyCompany()

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	query := fmt.Sprintf("SELECT SLCS.NR_DOC, SLCS.TX_NOME, SLCS.DT_NASC FROM LIBERAR_CAD_SITE_COMPRD SLCS WHERE SLCS.IN_ATIVO = 1 AND SLCS.TS_PROC_PORTAL IS NULL ORDER BY TS_PROC FETCH FIRST %d ROWS ONLY FOR FETCH ONLY WITH UR;", _defaultLimit)

	rows := mock.NewRows(
		[]string{"NR_DOC", "TX_NOME", "DT_NASC"}).
		AddRow(legacyPerson.NI, legacyPerson.Name, legacyPerson.BirthDate).
		AddRow(legacyCompany.NI, legacyCompany.Name, legacyCompany.BirthDate)

	mock.ExpectQuery(query).WillReturnRows(rows)

	legacyPeople, err := repoLegacy.GetLegacyPeople(context.TODO(), _defaultLimit)

	assert.NoError(t, err, "This result should not have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")

	assert.NotEmpty(t, legacyPeople, "This result should not be empty.")

	assert.Equal(t, legacyPerson.NI, legacyPeople[0].NI, "The resul should be the same.")
	assert.Equal(t, legacyPerson.Name, legacyPeople[0].Name, "The resul should be the same.")
	assert.Equal(t, legacyPerson.BirthDate, legacyPeople[0].BirthDate, "The resul should be the same.")

	assert.Equal(t, legacyCompany.NI, legacyPeople[1].NI, "The resul should be the same.")
	assert.Equal(t, legacyCompany.Name, legacyPeople[1].Name, "The resul should be the same.")
	assert.Equal(t, legacyCompany.BirthDate, legacyPeople[1].BirthDate, "The resul should be the same.")
}

func Test_GetLegacyPeople_QueryContextError(t *testing.T) {

	connection, _, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	_, err = repoLegacy.GetLegacyPeople(context.TODO(), 10)

	assert.Error(t, err, "This result should have errors.")

}

func Test_GetLegacyPeople_RowsScanError(t *testing.T) {

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	query := fmt.Sprintf("SELECT SLCS.NR_DOC, SLCS.TX_NOME, SLCS.DT_NASC FROM LIBERAR_CAD_SITE_COMPRD SLCS WHERE SLCS.IN_ATIVO = 1 AND SLCS.TS_PROC_PORTAL IS NULL ORDER BY TS_PROC FETCH FIRST %d ROWS ONLY FOR FETCH ONLY WITH UR;", _defaultLimit)

	rows := mock.NewRows(
		[]string{"NR_DOC", "TX_NOME", "DT_NASC"}).AddRow(nil, nil, nil)

	mock.ExpectQuery(query).WillReturnRows(rows)

	_, err = repoLegacy.GetLegacyPeople(context.TODO(), _defaultLimit)

	assert.Error(t, err, "This result should have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")

}

func Test_UpdateLegacyPeople_Success(t *testing.T) {
	legacyPerson := createLegacyPerson()
	legacyCompany := createLegacyCompany()
	legacyPeople := []entity.LegacyPerson{legacyPerson, legacyCompany}

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	mock.ExpectBegin()

	query := "UPDATE LIBERAR_CAD_SITE_COMPRD SET TS_PROC_PORTAL = CURRENT_TIMESTAMP WHERE NR_DOC = \\?;"

	for i, insertableObject := range legacyPeople {
		prepare := mock.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(insertableObject.NI).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	mock.ExpectCommit()

	err = repoLegacy.UpdateLegacyPeople(context.TODO(), legacyPeople)
	assert.NoError(t, err, "This result should not have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_UpdateLegacyPeople_BeginError(t *testing.T) {
	connection, _, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	err = repoLegacy.UpdateLegacyPeople(context.TODO(), []entity.LegacyPerson{})
	assert.Error(t, err, "This result should have errors.")
}

func Test_UpdateLegacyPeople_PrepareError(t *testing.T) {

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = repoLegacy.UpdateLegacyPeople(context.TODO(), []entity.LegacyPerson{createLegacyPerson()})
	assert.Error(t, err, "This result should have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_UpdateLegacyPeople_ExecError(t *testing.T) {

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	mock.ExpectBegin()

	query := "UPDATE LIBERAR_CAD_SITE_COMPRD SET TS_PROC_PORTAL = CURRENT_TIMESTAMP WHERE NR_DOC = \\?;"
	mock.ExpectPrepare(query)

	mock.ExpectRollback()

	err = repoLegacy.UpdateLegacyPeople(context.TODO(), []entity.LegacyPerson{createLegacyCompany()})
	assert.Error(t, err, "This result should have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_UpdateLegacyPeople_CommitError(t *testing.T) {
	legacyPeople := []entity.LegacyPerson{createLegacyPerson(), createLegacyCompany()}

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	mock.ExpectBegin()

	query := "UPDATE LIBERAR_CAD_SITE_COMPRD SET TS_PROC_PORTAL = CURRENT_TIMESTAMP WHERE NR_DOC = \\?;"

	for i, insertableObject := range legacyPeople {
		prepare := mock.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(insertableObject.NI).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	mock.ExpectCommit().WillReturnError(errors.New("legacy_person_db2 - UpdateLegacyPeople - Commit - Mock Error"))

	err = repoLegacy.UpdateLegacyPeople(context.TODO(), legacyPeople)
	assert.Error(t, err, "This result should have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_GetBucketsAvailable_Success(t *testing.T) {

	connection, _, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	connection.SetMaxOpenConns(5)

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	nc := repoLegacy.GetBucketsAvailable()

	assert.Equal(t, 2, nc, "This result should be the same.")

}

func Test_GetBucketsAvailable_Default(t *testing.T) {

	connection, _, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	db2Test := createDB2Test(connection)
	repoLegacy := repo.NewLegacyPersonDB2Repo(db2Test)

	defer repoLegacy.Close()

	nc := repoLegacy.GetBucketsAvailable()

	assert.Equal(t, 1, nc, "This result should be the same.")

}

func createDB2Test(db *sql.DB) *db2.DB2 {
	return &db2.DB2{
		DB: db,
	}
}

func createLegacyPerson() entity.LegacyPerson {
	birthDate := time.Now()
	return entity.LegacyPerson{NI: 88626159226, Name: "User One", BirthDate: &birthDate}
}

func createLegacyCompany() entity.LegacyPerson {
	return entity.LegacyPerson{NI: 87511532000171, Name: "Company One"}
}
