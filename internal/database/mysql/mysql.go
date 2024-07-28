package mysql

import (
	"context"
	"fmt"
	"log"
	"movie-booking-app/users-service/internal/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type mysqlDatabase struct {
	Db *sqlx.DB
}

func NewMySQLDatabase(dbString string) (database.Database, error) {
	fmt.Println("dbString", dbString)
	db, err := sqlx.Open("mysql", dbString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(5)
	db.SetMaxOpenConns(50)

	return &mysqlDatabase{
		Db: db,
	}, nil
}

func (m *mysqlDatabase) GetDB() *sqlx.DB {
	return m.Db
}

func (m *mysqlDatabase) WrapWithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := m.Db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from panic:", r)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("transaction rollback failed: %v", rbErr)
			} else {
				log.Println("transaction rolled back")
			}
			err = fmt.Errorf("recovered from panic: %v", r)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("transaction rollback failed: %v", rbErr)
			} else {
				log.Println("transaction rolled back")
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("failed to commit transaction: ", cmErr)
				err = cmErr
			}
		}
	}()

	err = fn(tx)
	return err
}
