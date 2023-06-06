package notification

import (
	"time"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StageNotificationPump struct {
	db      *gorm.DB
	Receive func(notification *model.StageNotification) error
	isRunning bool
	stopChannel chan bool
}

func (p *StageNotificationPump) Start() {
	if p.isRunning {
		return
	}

	p.isRunning = true

	go func() {
		for {
			select {
			case <-p.stopChannel:
				p.isRunning = false
				return
			default:
				notification, ok := p.read()
				if !ok {
					time.Sleep(30 * time.Second)
					continue
				}

				err := p.Receive(notification)
				if err != nil {				
					log.Error(err)
					continue
				}
			}
		}
	}()
}

func (p *StageNotificationPump) Stop() {
	if p.isRunning {
		p.stopChannel <- true
	}
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
