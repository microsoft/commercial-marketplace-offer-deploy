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

func NewUnitOfWork() *UnitOfWork {
	return &UnitOfWork{}
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
	return nil
}