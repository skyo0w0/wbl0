package storage

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var migrations = []func(tx *sqlx.Tx) error{
	initial,
}
var maxVersion = len(migrations)

func migrate(db *sqlx.DB) error {
	for v := 1; v <= maxVersion; v++ {
		err := migrateVersion(v, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateVersion(v int, db *sqlx.DB) error {
	var err error
	var tx *sqlx.Tx
	migrationFunc := migrations[v-1]

	if tx, err = db.BeginTxx(context.TODO(), nil); err != nil {
		log.Printf("migration[%d] failed to start transaction: %s\n", v, err.Error())
		return err
	}

	if err = migrationFunc(tx); err != nil {
		log.Printf("migration[%d] failed to migrate: %s\n", v, err.Error())
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("migration[%d] failed to commit changes: %s\n", v, err.Error())
		return err
	}

	return nil
}

func initial(tx *sqlx.Tx) error {
	query := `
	create table if not exists orders (
		order_uid text primary key,
		data text not null
	);
	`
	_, err := tx.ExecContext(context.TODO(), query)
	return err
}
