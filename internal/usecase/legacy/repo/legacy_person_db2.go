package repo

import (
	"context"
	"fmt"
	"github.com/diogoalbuquerque/migration-customers/internal/entity"
	"github.com/diogoalbuquerque/migration-customers/pkg/db2"
)

const (
	_defaultMinAvailableBuckets = 1
)

type LegacyPersonDB2Repo struct {
	*db2.DB2
}

func NewLegacyPersonDB2Repo(db2 *db2.DB2) *LegacyPersonDB2Repo {
	return &LegacyPersonDB2Repo{db2}
}

func (r *LegacyPersonDB2Repo) GetLegacyPeople(ctx context.Context, limit int) ([]entity.LegacyPerson, error) {
	query := fmt.Sprintf("SELECT "+
		"    SLCS.NR_DOC, "+
		"    SLCS.TX_NOME, "+
		"    SLCS.DT_NASC "+
		"FROM LIBERAR_CAD_SITE_COMPRD SLCS "+
		"WHERE "+
		"     SLCS.IN_ATIVO = 1 "+
		"AND SLCS.TS_PROC_PORTAL IS NULL "+
		"ORDER BY TS_PROC FETCH FIRST %d ROWS ONLY FOR FETCH ONLY WITH UR;", limit)

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("legay_person_db2 - GetLegacyPeople - Query: %w", err)
	}
	defer rows.Close()

	var legacyPerson entity.LegacyPerson
	var legacyPeople []entity.LegacyPerson

	for rows.Next() {
		err = rows.Scan(&legacyPerson.NI, &legacyPerson.Name, &legacyPerson.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("legay_person_db2 - GetLegacyPeople - Scan: %w", err)
		}

		legacyPeople = append(legacyPeople, legacyPerson)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("legay_person_db2 - GetLegacyPeople - Rows: %w", err)
	}

	return legacyPeople, err
}

func (r *LegacyPersonDB2Repo) UpdateLegacyPeople(ctx context.Context, legacyPeople []entity.LegacyPerson) error {

	query := "UPDATE LIBERAR_CAD_SITE_COMPRD SET TS_PROC_PORTAL = CURRENT_TIMESTAMP WHERE NR_DOC = ?;"

	transaction, err := r.DB.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("legay_person_db2 - UpdateLegacyPeople - Begin: %w", err)
	}

	defer transaction.Rollback()

	for _, legacyPerson := range legacyPeople {
		prepareContext, transactionError := transaction.PrepareContext(ctx, query)

		if transactionError != nil {
			return fmt.Errorf("legay_person_db2 - UpdateLegacyPeople - [%f] - PrepareContext: %w", legacyPerson.NI, transactionError)
		}

		_, transactionError = prepareContext.ExecContext(ctx, legacyPerson.NI)

		if transactionError != nil {
			return fmt.Errorf("legay_person_db2 - UpdateLegacyPeople - [%f] - ExecContext: %w", legacyPerson.NI, transactionError)
		}

	}

	if err = transaction.Commit(); err != nil {
		return fmt.Errorf("legay_person_db2 - UpdateLegacyPeople - Commit: %w", err)
	}

	return nil
}

func (r *LegacyPersonDB2Repo) GetBucketsAvailable() int {

	dbStats := r.DB.Stats()
	availableBuckets := (dbStats.MaxOpenConnections - dbStats.OpenConnections) / 2
	if availableBuckets < _defaultMinAvailableBuckets {
		return _defaultMinAvailableBuckets
	}
	return availableBuckets
}
