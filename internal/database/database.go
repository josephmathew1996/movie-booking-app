package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetDB() *sqlx.DB
	WrapWithTransaction(context.Context, func(tx *sqlx.Tx) error) error
}
