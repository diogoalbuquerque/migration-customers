package db2

import "time"

type Option func(db2 *DB2)

func MaxOpenConns(size int) Option {
	return func(d *DB2) {
		d.maxOpenConns = size
	}
}

func MaxIdleConns(size int) Option {
	return func(d *DB2) {
		d.maxIdleConns = size
	}
}

func MaxLifetime(timeout int) Option {
	return func(d *DB2) {
		d.maxLifetime = time.Duration(timeout) * time.Second
	}
}
