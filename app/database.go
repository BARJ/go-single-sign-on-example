package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const driver = "postgres"
const schema = "public"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Database struct {
	config DBConfig
}

func NewDatabase(config DBConfig) Database {
	return Database{config: config}
}

func (d Database) Connect() (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"sslmode=disable host=%s port=%d user=%s password=%s dbname=%s search_path=%s",
		d.config.Host,
		d.config.Port,
		d.config.User,
		d.config.Password,
		d.config.DBName,
		schema,
	)

	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func (d Database) MigrateUp(db *sql.DB) error {
	var driver database.Driver

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", d.config.DBName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to perform database migration: %v", err)
	}

	return nil
}
