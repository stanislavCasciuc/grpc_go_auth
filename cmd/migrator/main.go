package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migratePath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migratePath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "table", "migrations", "migrations table name")
	flag.Parse()

	if storagePath == "" || migratePath == "" {
		panic("storage and migrate paths are required")
	}

	m, err := migrate.New(
		"file://"+migratePath, fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied")
}
