package repo_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/internal/usecase/person/repo"
	"github.com/diogoalbuquerque/migration-customers/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_StorePeople_Success(t *testing.T) {
	people := []entity.Person{createPerson(), createCompany()}

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	mock.ExpectBegin()

	query := "INSERT INTO CLIENTE_LIBERADO \\(NR_IDENTIFICADOR, NM_CLIENTE, TP_CLIENTE, DT_NASCIMENTO\\) VALUES \\(\\?, \\?, \\?, \\?\\);"

	for i, insertableObject := range people {
		prepare := mock.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(insertableObject.NI, insertableObject.Name, insertableObject.PersonType, insertableObject.BirthDate).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	mock.ExpectCommit()

	err = repoPerson.StorePeople(context.TODO(), people)
	assert.NoError(t, err, "This result should not have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_StorePeople_BeginError(t *testing.T) {
	connection, _, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	err = repoPerson.StorePeople(context.TODO(), []entity.Person{})
	assert.Error(t, err, "This result should have errors.")
}

func Test_StorePeople_PrepareError(t *testing.T) {

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = repoPerson.StorePeople(context.TODO(), []entity.Person{createPerson()})
	assert.Error(t, err, "This result should have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_StorePeople_ExecError(t *testing.T) {

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	mock.ExpectBegin()

	query := "INSERT INTO CLIENTE_LIBERADO \\(NR_IDENTIFICADOR, NM_CLIENTE, TP_CLIENTE, DT_NASCIMENTO\\) VALUES \\(\\?, \\?, \\?, \\?\\);"
	mock.ExpectPrepare(query)

	mock.ExpectRollback()

	err = repoPerson.StorePeople(context.TODO(), []entity.Person{createCompany()})
	assert.Error(t, err, "This result should have errors.")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("An error occurred and the following expectations were not met: %s", err)
	}

	assert.NoError(t, err, "This result should not have errors.")
}

func Test_StorePeople_CommitError(t *testing.T) {
	people := []entity.Person{createPerson(), createCompany()}

	connection, mock, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	mock.ExpectBegin()

	query := "INSERT INTO CLIENTE_LIBERADO \\(NR_IDENTIFICADOR, NM_CLIENTE, TP_CLIENTE, DT_NASCIMENTO\\) VALUES \\(\\?, \\?, \\?, \\?\\);"

	for i, insertableObject := range people {
		prepare := mock.ExpectPrepare(query)
		prepare.ExpectExec().
			WithArgs(insertableObject.NI, insertableObject.Name, insertableObject.PersonType, insertableObject.BirthDate).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	mock.ExpectCommit().WillReturnError(errors.New("person_mysql - StorePeople - Commit - Mock Error"))

	err = repoPerson.StorePeople(context.TODO(), people)
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

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	nc := repoPerson.GetBucketsAvailable()

	assert.Equal(t, 2, nc, "This result should be the same.")

}

func Test_GetBucketsAvailable_Default(t *testing.T) {

	connection, _, err := sqlmock.New()
	assert.Nil(t, err, fmt.Sprintf("Error '%s' was not expected when opening a stub database connection", err))

	mysqlTest := createMYSQLTest(connection)
	repoPerson := repo.NewPersonMYSQLRepo(mysqlTest)

	defer repoPerson.Close()

	nc := repoPerson.GetBucketsAvailable()

	assert.Equal(t, 1, nc, "This result should be the same.")

}

func createMYSQLTest(db *sql.DB) *mysql.MYSQL {
	return &mysql.MYSQL{
		MySQLDatabase: "Portal",
		DB:            db,
	}
}

func createPerson() entity.Person {
	birthDate := time.Now()
	return entity.Person{NI: "88626159226", Name: "User One", PersonType: entity.PF, BirthDate: &birthDate}
}

func createCompany() entity.Person {
	return entity.Person{NI: "87511532000171", Name: "Company One", PersonType: entity.PJ}
}
