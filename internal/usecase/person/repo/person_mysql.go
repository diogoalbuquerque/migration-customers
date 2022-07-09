package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/pkg/mysql"
	"time"
)

const (
	_defaultMinAvailableBuckets = 1
)

type PersonMYSQLRepo struct {
	*mysql.MYSQL
}

func NewPersonMYSQLRepo(mysql *mysql.MYSQL) *PersonMYSQLRepo {
	return &PersonMYSQLRepo{mysql}
}

func (r *PersonMYSQLRepo) StorePeople(ctx context.Context, people []entity.Person) error {

	query := "INSERT INTO CLIENTE_LIBERADO (NR_IDENTIFICADOR, NM_CLIENTE, TP_CLIENTE, DT_NASCIMENTO) VALUES (?, ?, ?, ?);"

	transaction, err := r.DB.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("person_mysql - StorePeople - Begin: %w", err)
	}

	defer transaction.Rollback()

	for _, person := range people {
		prepareContext, transactionError := transaction.PrepareContext(ctx, query)

		if transactionError != nil {
			return fmt.Errorf("person_mysql - StorePeople - [%s] - PrepareContext: %w", person.NI, transactionError)
		}

		_, transactionError = prepareContext.ExecContext(ctx, person.NI, person.Name, person.PersonType, getValueTimeOrNullSqlValue(person.BirthDate))

		if transactionError != nil {
			return fmt.Errorf("person_mysql - StorePeople - [%s] - ExecContext: %w", person.NI, transactionError)
		}

	}

	if err = transaction.Commit(); err != nil {
		return fmt.Errorf("person_mysql - StorePeople - Commit: %w", err)
	}

	return nil
}

func (r *PersonMYSQLRepo) GetBucketsAvailable() int {

	dbStats := r.DB.Stats()
	availableBuckets := (dbStats.MaxOpenConnections - dbStats.OpenConnections) / 2
	if availableBuckets < _defaultMinAvailableBuckets {
		return _defaultMinAvailableBuckets
	}
	return availableBuckets
}

func getValueTimeOrNullSqlValue(value *time.Time) interface{} {
	if value == nil {
		return sql.NullTime{}
	}
	return value
}
