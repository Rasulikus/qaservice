package migrator

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	DB  *sql.DB
	Dir string
}

func NewMigrator(dsn, dir string) (*Migrator, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}

	return &Migrator{
		DB:  db,
		Dir: dir,
	}, nil
}

func (m *Migrator) Up() error {
	if err := goose.Up(m.DB, m.Dir); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	log.Print("migrate up ok")
	return nil
}

func (m *Migrator) Down() error {
	if err := goose.Down(m.DB, m.Dir); err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}
	log.Print("migrate down ok")
	return nil
}

func (m *Migrator) Reset() error {
	if err := goose.Reset(m.DB, m.Dir); err != nil {
		return fmt.Errorf("migrate reset: %w", err)
	}
	log.Print("migrate reset ok")
	return nil
}

func (m *Migrator) Close() error {
	return m.DB.Close()
}
