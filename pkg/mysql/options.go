package mysql

import "time"

type Option func(mysql *MYSQL)

func MaxOpenConns(size int) Option {
	return func(m *MYSQL) {
		m.maxOpenConns = size
	}
}

func MaxIdleConns(size int) Option {
	return func(m *MYSQL) {
		m.maxIdleConns = size
	}
}

func MaxLifetime(timeout int) Option {
	return func(m *MYSQL) {
		m.maxLifetime = time.Duration(timeout) * time.Second
	}
}
