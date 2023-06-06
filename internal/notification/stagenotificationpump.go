package notification

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

type StageNotificationPump struct {
	db      *gorm.DB
	Receive func(notification *model.StageNotification) error
}

func (p *StageNotificationPump) Start() {
	// TODO: implement loop using read() and if true then p.receive()
}

func (p *StageNotificationPump) read() (*model.StageNotification, bool) {
	//write gorm query to read from database where Done is false
	record := &model.StageNotification{}
	tx := p.db.Where("done = ?", false).First(record)

	if tx.RowsAffected == 0 {
		return nil, false
	}
	return record, true
}
