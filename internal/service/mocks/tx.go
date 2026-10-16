package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type Tx struct {
	mock.Mock
}

func (m *Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *Tx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *Tx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *Tx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *Tx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(pgx.BatchResults)
}

func (m *Tx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

func (m *Tx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

func (m *Tx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *Tx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	calledArgs := m.Called(ctx, sql, args)
	if calledArgs.Get(0) == nil {
		return nil, calledArgs.Error(1)
	}
	return calledArgs.Get(0).(pgx.Rows), calledArgs.Error(1)
}

func (m *Tx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	calledArgs := m.Called(ctx, sql, args)
	if calledArgs.Get(0) == nil {
		return nil
	}
	return calledArgs.Get(0).(pgx.Row)
}

func (m *Tx) Conn() *pgx.Conn {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgx.Conn)
}
