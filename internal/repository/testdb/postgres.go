package testdb

import (
	"log"
	"path/filepath"

	"github.com/Rasulikus/qaservice/internal/config"
	"github.com/Rasulikus/qaservice/internal/migrator"
	"github.com/Rasulikus/qaservice/internal/repository"
	"gorm.io/gorm"
)

var (
	truncateSQL = `
	TRUNCATE TABLE
		answers,
	    questions
	RESTART IDENTITY CASCADE;
	`
)

type TestDB struct {
	DB     *gorm.DB
	cfg    *config.DBConfig
	migDir string
}

func NewTestDB(cfg *config.DBConfig, migDir string) *TestDB {
	if cfg == nil {
		log.Print("config is nil use default config")
		cfg = &config.DBConfig{
			User: "admin",
			Pass: "mypassword",
			Host: "localhost",
			Port: "5432",
			Name: "qaservice_test",
		}
	}

	if migDir == "" {
		log.Print("migDir is empty use default migDir")
		migDir = filepath.Join("..", "..", "..", "migrations")
	}

	testDB, err := repository.NewDB(cfg, migDir)
	if err != nil {
		log.Fatal(err)
	}

	return &TestDB{
		DB:     testDB,
		cfg:    cfg,
		migDir: migDir,
	}
}

func (db *TestDB) RecreateTables() {
	mig, err := migrator.NewMigrator(db.cfg.PostgresURL(), db.migDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := mig.Up(); err != nil {
		log.Fatal(err)
	}
	db.CleanDB()
}

func (db *TestDB) CleanDB() {
	if err := db.DB.Exec(truncateSQL).Error; err != nil {
		log.Fatalf("clean db: %v", err)
	}
}

func (db *TestDB) Close() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("get sql db: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("close db: %v", err)
	}
}
