package persistence

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sigs.k8s.io/kustomize/kyaml/pathutil"
)

const (
	DefaultDatabasePath = "/data/cmod/"
	DatabaseName        = "commercial-marketplace-offer-deploy"
	DatabaseFileName    = DatabaseName + ".db"
	InMemoryDsn         = "file::memory:?cache=shared"
)

func newInMemoryDatabase() *gorm.DB {
	return newDatabase(InMemoryDsn)
}

// creates a new database, establishing an open connection with a new session
// path: the folder path to the database file
func newDatabase(dsn string) *gorm.DB {
	db, err := createConnection(dsn)

	if err != nil {
		log.Fatalf("Could not open database %s. Error: %v", dsn, err)
	}

	return db
}

func createConnection(dsn string, models ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("could not open and connect to database at %s: %w", dsn, err)
	}

	if err := db.AutoMigrate(models); err != nil {
		return nil, fmt.Errorf("could not migrate models %T: %w", models, err)
	}

	if tx := db.Exec("PRAGMA foreign_keys = ON", nil); tx.Error != nil {
		return nil, fmt.Errorf("unable to turn on foreign keys in sqlite db: %w", tx.Error)
	}
	return db, nil
}

func getDefaultDsn() string {
	return pathutil.Join(DefaultDatabasePath, DatabaseFileName)
}
