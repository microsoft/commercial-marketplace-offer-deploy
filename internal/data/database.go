package data

import (
	"fmt"
	"path/filepath"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	Instance() *gorm.DB
}

type DatabaseOptions struct {
	UseInMemory bool
	Dsn         string
}

type database struct {
	db *gorm.DB
}

const (
	DataDatabasePath = "./"
	DatabaseName     = "modm"
	DatabaseFileName = DatabaseName + ".db"
	InMemoryDsn      = "file::memory:"
)

// The default db path for the database if nothing is set. Default value is DataDatabasePath
var defaultDatabasePath string = DataDatabasePath

func SetDefaultDatabasePath(path string) {
	defaultDatabasePath = path
}

// Db implements DbContext
func (ctx *database) Instance() *gorm.DB {
	return ctx.db
}

// Gets a new db using the provided Options. If options are nil, the default DSN is used
func NewDatabase(options *DatabaseOptions) Database {
	if options == nil {
		options = &DatabaseOptions{Dsn: getDefaultDsn()}
	}

	var db *gorm.DB
	if options.UseInMemory {
		db = newInMemoryDatabase()
	} else {
		db = newDatabase(options.Dsn)
	}

	if log.GetLevel() == log.DebugLevel {
		db = db.Debug()
	}

	return &database{db}
}

func newInMemoryDatabase() *gorm.DB {
	return newDatabase(InMemoryDsn)
}

// creates a new database, establishing an open connection with a new session
// path: the folder path to the database file
func newDatabase(dsn string) *gorm.DB {
	db, err := createInstance(dsn)

	if err != nil {
		log.Fatalf("Could not open database %s. Error: %v", dsn, err)
	}

	return db
}

func createInstance(dsn string, models ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("could not open and connect to database at %s: %w", dsn, err)
	}

	if err := db.AutoMigrate(&model.Deployment{}, &model.Stage{}, &model.EventHook{}, &model.EventHookAuditEntry{}, &model.InvokedOperation{}, &model.InvokedOperationAttribute{}, &model.StageNotification{}); err != nil {
		return nil, fmt.Errorf("could not migrate models %T: %w", models, err)
	}

	if tx := db.Exec("PRAGMA foreign_keys = ON", nil); tx.Error != nil {
		return nil, fmt.Errorf("unable to turn on foreign keys in sqlite db: %w", tx.Error)
	}
	return db, nil
}

func getDefaultDsn() string {
	return filepath.Join(defaultDatabasePath, DatabaseFileName)
}
