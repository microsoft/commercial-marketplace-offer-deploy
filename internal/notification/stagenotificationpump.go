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
	receive       NotificationPumpReceiveFunc[model.StageNotification]
	isRunning     bool
	stopChannel   chan bool
	sleepDuration time.Duration
}

func (p *StageNotificationPump) SetReceiver(receive NotificationPumpReceiveFunc[model.StageNotification]) {
	p.receive = receive
}

func (p *StageNotificationPump) Start() {
	if p.receive == nil {
		log.Error("pump has not receiver set")
		return
	}

	if p.isRunning {
		return
	}

	time.Sleep(p.sleepDuration)
	p.isRunning = true

	log.Tracef("Starting StageNotificationPump with sleep duration %v", p.sleepDuration)

	go func() {
		for {
			select {
			case <-p.stopChannel:
				p.isRunning = false
				return
			default:
				// do we want to read all unsent notifications at once to reduce db calls?
				notification, ok := p.read()
				if ok {
					err := p.receive(notification)
					if err != nil {
						log.Error(err)
					}
				}
				time.Sleep(p.sleepDuration)
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
	tx := p.db.Where("is_done = ?", false).First(record)

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
