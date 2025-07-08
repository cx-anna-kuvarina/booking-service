package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	ReadPool  *pgxpool.Pool
	WritePool *pgxpool.Pool
}

func NewDB(ctx context.Context, readDSN, writeDSN string) (*DB, error) {
	readPool, err := pgxpool.New(ctx, readDSN)
	if err != nil {
		return nil, err
	}

	writePool, err := pgxpool.New(ctx, writeDSN)
	if err != nil {
		readPool.Close()
		return nil, err
	}

	return &DB{
		ReadPool:  readPool,
		WritePool: writePool,
	}, nil
}

func (d *DB) Close() {
	d.ReadPool.Close()
	d.WritePool.Close()
}
