package fake

import (
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// fake database that's in memory
type Database struct {
	db *gorm.DB
	t  *testing.T
}

func (d *Database) Instance() *gorm.DB {
	return d.db
}

func (d *Database) Setup(action func(db *gorm.DB)) *Database {
	require.NotNil(d.t, action)
	action(d.db)

	return d
}

func NewFakeDatabase(t *testing.T) *Database {
	return &Database{
		db: data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance(),
		t:  t,
	}
}
