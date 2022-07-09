package internal

import (
	"context"
	"github.com/diogoalbuquerque/migration-customers/config"
	"github.com/diogoalbuquerque/migration-customers/internal/usecase"
	luc "github.com/diogoalbuquerque/migration-customers/internal/usecase/legacy"
	repoLegacy "github.com/diogoalbuquerque/migration-customers/internal/usecase/legacy/repo"
	puc "github.com/diogoalbuquerque/migration-customers/internal/usecase/person"
	"github.com/diogoalbuquerque/migration-customers/internal/usecase/person/repo"
	"github.com/diogoalbuquerque/migration-customers/pkg/db2"
	"github.com/diogoalbuquerque/migration-customers/pkg/logger"
	"github.com/diogoalbuquerque/migration-customers/pkg/mysql"
	"github.com/diogoalbuquerque/migration-customers/pkg/secret"
)

func Run(ctx context.Context, cfg config.Config, awsSecret secret.AwsSecret) error {

	l := logger.New(cfg.Log.Level)

	db2DB, err := db2.New(awsSecret, db2.MaxIdleConns(cfg.DB2.IdleConnMax), db2.MaxOpenConns(cfg.DB2.OpenConnMax), db2.MaxLifetime(cfg.DB2.LifeConnMax))

	if err != nil {
		l.Error(err)
		return err
	}

	defer db2DB.Close()

	lr := repoLegacy.NewLegacyPersonDB2Repo(db2DB)

	mysqlDB, err := mysql.New(awsSecret, mysql.MaxIdleConns(cfg.MYSQL.IdleConnMax), mysql.MaxOpenConns(cfg.MYSQL.OpenConnMax), mysql.MaxLifetime(cfg.MYSQL.LifeConnMax))

	if err != nil {
		l.Error(err)
		return err
	}

	defer mysqlDB.Close()

	pr := repo.NewPersonMYSQLRepo(mysqlDB)

	return usecase.NewMigrationUseCase(luc.NewLegacyUseCase(lr), puc.NewPersonUseCase(pr), l).Migration(ctx, cfg.Limit)
}
