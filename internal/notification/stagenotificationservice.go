package notification

import (
	"context"
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StageNotificationService struct {
	ctx            context.Context
	db             *gorm.DB
	pump           NotificationPump[model.StageNotification]
	handlerFactory NotificationHandlerFactoryFunc[model.StageNotification]
	channels       map[uint]chan NotificationHandlerResult[model.StageNotification]
	log            *log.Entry
}

func NewStageNotificationService(db *gorm.DB, pump NotificationPump[model.StageNotification], handlerFactory StageNotificationHandlerFactoryFunc) *StageNotificationService {
	return &StageNotificationService{
		ctx:            context.Background(),
		db:             db,
		pump:           pump,
		handlerFactory: handlerFactory,
		log:            log.WithFields(log.Fields{}),
		channels:       make(map[uint]chan NotificationHandlerResult[model.StageNotification]),
	}
}

// stub out hosting.Service interface on StageNotificationService
func (s *StageNotificationService) Start() {
	s.pump.SetReceiver(s.receive)
	s.pump.Start()

	s.start()
}

func (s *StageNotificationService) Stop() {
	s.pump.Stop()
}

func (s *StageNotificationService) GetName() string {
	return "Stage Notification Service"
}

func (s *StageNotificationService) start() {
	for {
		// loop over s.results
		// if done, remove from map
		for _, channel := range s.channels {
			select {
			case result := <-channel:
				id := result.Notification.ID
				s.log.Infof("Handler (notification [%d])", id)

				if result.Error != nil {
					s.log.Errorf("Error handling stage notification %d: %s", id, result.Error)
				}
				if result.Done {
					result.Notification.Done()
					s.db.Save(result.Notification)
					delete(s.channels, id)
				}
			default:
				continue
			}
		}
	}
}

func (s *StageNotificationService) receive(notification *model.StageNotification) error {
	if s.isCurrentlyBeingHandled(notification.ID) {
		return fmt.Errorf("already handling notification [%d]", notification.ID)
	}

	handler, err := s.handlerFactory()
	if err != nil {
		return err
	}

	context := NewNotificationHandlerContext[model.StageNotification](s.ctx, notification)
	s.channels[notification.ID] = context.Channel()

	go handler.Handle(context)

	return nil
}

func (s *StageNotificationService) isCurrentlyBeingHandled(id uint) bool {
	_, ok := s.channels[id]
	return ok
}
