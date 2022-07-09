package db2

import (
	"database/sql"
	"fmt"
	"github.com/diogoalbuquerque/migration-customers/pkg/secret"
	_ "github.com/ibmdb/go_ibm_db"
	"time"
)

const (
	_defaultMaxOpenConns = 10
	_defaultMaxIdleConns = 2
	_defaultMaxLifetime  = time.Second * 10
)

type DB2 struct {
	maxOpenConns int
	maxIdleConns int
	maxLifetime  time.Duration
	DB           *sql.DB
}

func New(awsSecret secret.AwsSecret, opts ...Option) (*DB2, error) {
	db2 := &DB2{
		maxOpenConns: _defaultMaxOpenConns,
		maxIdleConns: _defaultMaxIdleConns,
		maxLifetime:  _defaultMaxLifetime,
	}

	for _, opt := range opts {
		opt(db2)
	}

	con := fmt.Sprintf("HOSTNAME=%v;DATABASE=%v;PORT=%v;UID=%v;PWD=%v",
		awsSecret.Db2Host, awsSecret.Db2Database, awsSecret.Db2Port, awsSecret.Db2Username, awsSecret.Db2Password)
	db, err := sql.Open("go_ibm_db", con)
	if err != nil {
		return nil, fmt.Errorf("db2 - New - Open: %w", err)
	}

	db.SetMaxOpenConns(db2.maxOpenConns)
	db.SetMaxIdleConns(db2.maxIdleConns)
	db.SetConnMaxLifetime(db2.maxLifetime)

	db2.DB = db
	return db2, nil

}

func (d *DB2) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
