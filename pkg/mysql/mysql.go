package mysql

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/diogoalbuquerque/migration-customers/pkg/secret"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	_defaultMaxOpenConns  = 10
	_defaultMaxIdleConns  = 2
	_defaultMaxLifetime   = time.Second * 10
	_defaultMySQLDatabase = "Default"
)

type MYSQL struct {
	maxOpenConns  int
	maxIdleConns  int
	maxLifetime   time.Duration
	MySQLDatabase string
	DB            *sql.DB
}

func New(awsSecret secret.AwsSecret, opts ...Option) (*MYSQL, error) {
	mysql := &MYSQL{
		maxOpenConns:  _defaultMaxOpenConns,
		maxIdleConns:  _defaultMaxIdleConns,
		maxLifetime:   _defaultMaxLifetime,
		MySQLDatabase: _defaultMySQLDatabase,
	}

	for _, opt := range opts {
		opt(mysql)
	}

	con := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", awsSecret.MySQLUsername, awsSecret.MySQLPassword, awsSecret.MySQLHost,
		awsSecret.MySQLPort, awsSecret.MySQLDatabase)
	db, err := xray.SQLContext(awsSecret.MySQLEngine, con)
	if err != nil {
		return nil, fmt.Errorf("mysql - New - Open: %w", err)
	}

	db.SetMaxOpenConns(mysql.maxOpenConns)
	db.SetMaxIdleConns(mysql.maxIdleConns)
	db.SetConnMaxLifetime(mysql.maxLifetime)

	mysql.MySQLDatabase = awsSecret.MySQLDatabase

	mysql.DB = db
	return mysql, nil

}

func (d *MYSQL) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
