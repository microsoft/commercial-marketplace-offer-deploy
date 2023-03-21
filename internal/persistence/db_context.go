package persistence

import "gorm.io/gorm"

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

// Db implements DbContext
func (ctx *database) Instance() *gorm.DB {
	return ctx.db
}

// Gets a new db context using the provided Options. If options are nil, the default DSN is used
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

	return &database{db}
}
