package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Rasulikus/qaservice/internal/config"
	"github.com/Rasulikus/qaservice/internal/migrator"
	"github.com/Rasulikus/qaservice/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.DBConfig, migDir string) (*gorm.DB, error) {
	EnsureDatabase(cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.PostgresURL(),
	}))
	if err != nil {
		log.Fatal(err)
	}

	mig, err := migrator.NewMigrator(cfg.PostgresURL(), migDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := mig.Up(); err != nil {
		log.Fatal(err)
	}

	return db, err
}

func EnsureDatabase(cfg *config.DBConfig) {
	adminCfg := *cfg
	adminCfg.Name = "postgres"

	adminDB, err := sql.Open("pgx", adminCfg.PostgresURL())
	if err != nil {
		log.Fatalf("open admin db: %v", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		log.Fatalf("ping admin db: %v", err)
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);`
	if err := adminDB.QueryRow(query, cfg.Name).Scan(&exists); err != nil {
		log.Fatalf("check database exists: %v", err)
	}

	if exists {
		return
	}

	createSQL := fmt.Sprintf(`CREATE DATABASE %s`, cfg.Name)
	if _, err := adminDB.Exec(createSQL); err != nil {
		log.Fatalf("create database %s: %v", cfg.Name, err)
	}
}
func TranslateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrNotFound
	}
	return err
}
