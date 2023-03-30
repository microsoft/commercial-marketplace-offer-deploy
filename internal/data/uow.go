package data

import (
	"gorm.io/gorm"
)

type UOW interface {
	RegisterNew(entity any)
	RegisterDirty(entity any)
	RegisterClean(entity any)
	RegisterDeleted(entity any)
	Commit() error
}

type UnitOfWork struct {
	database Database
	newEntities	[]any
	modifiedEntities	[]any
	cleanEntities	[]any
	deletedEntities	[]any
}

func NewUnitOfWork(database Database) UOW {
	return &UnitOfWork{
		database: database,
	}
}

func (uow *UnitOfWork) RegisterNew(entity any) {
	uow.newEntities = append(uow.newEntities, entity)
}

func (uow *UnitOfWork) RegisterDirty(entity any) {
	uow.modifiedEntities = append(uow.modifiedEntities, entity)
}

func (uow *UnitOfWork) RegisterClean(entity any) {	
	uow.cleanEntities = append(uow.cleanEntities, entity)
}

func (uow *UnitOfWork) RegisterDeleted(entity any) {
	uow.deletedEntities = append(uow.deletedEntities, entity)
}	

func (uow *UnitOfWork) Commit() error {
	db := uow.database.Instance()
	db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range uow.newEntities {
			if err := tx.Create(entity).Error; err != nil {
				return err
			}
		}
		for _, entity := range uow.modifiedEntities{
			if err := tx.Save(entity).Error; err != nil {
				return err
			}
		}
		for _, entity := range uow.deletedEntities{
			if err := tx.Delete(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
	uow.newEntities = nil
	uow.modifiedEntities = nil
	uow.cleanEntities = nil
	uow.deletedEntities = nil
	return nil
}