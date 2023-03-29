package data

type UOW interface {
	RegisterNew(entity any)
	RegisterDirty(entity any)
	RegisterClean(entity any)
	RegisterDeleted(entity any)
	Commit() error
}

type UnitOfWork struct {
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
	return nil
}