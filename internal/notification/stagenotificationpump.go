package notification

import (
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const SleepDurationPumpDefault = 30 * time.Second

type StageNotificationPump struct {
	db            *gorm.DB
	receive       ReceiveNotificationFunc
	isRunning     bool
	stopChannel   chan bool
	sleepDuration time.Duration
}

type ReceiveNotificationFunc func(notification *model.StageNotification) error

func (p *StageNotificationPump) SetReceiver(receive ReceiveNotificationFunc) {
	p.receive = receive
}

func (p *StageNotificationPump) Start() {
	if p.receive == nil {
		log.Error("ReceiveNotificationFunc not set")
		return
	}

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
				// do we want to read all unsent notifications at once to reduce db calls?
				notification, ok := p.read()
				if !ok {
					time.Sleep(p.sleepDuration)
					continue
				}

				err := p.receive(notification)
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

func NewStageNotificationPump(db *gorm.DB, sleepDuration time.Duration) *StageNotificationPump {
	return &StageNotificationPump{
		db:            db,
		sleepDuration: sleepDuration,
		stopChannel:   make(chan bool),
	}
}
