package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Handler func(ctx context.Context) error

type Query struct {
	Name     string
	QueryRaw string
}

type Client interface {
	DB() DB
	Close() error
}

type TxManager interface {
	ReadCommitted(ctx context.Context, fn Handler) error
}

type Transactor interface {
	BeginTx(ctx context.Context, txOpts pgx.TxOptions) (pgx.Tx, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

type NamedExecer interface {
	ScanOneCtx(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllCtx(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecer interface {
	ExecCtx(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryCtx(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowCtx(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
