package notification

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
)

type StageNotificationService struct {
	ctx            context.Context
	pump           *StageNotificationPump
	handlerFactory StageNotificationHandlerFactoryFunc
	results        map[uint]chan stageNotificationHandlerResult
	log            *log.Entry
}

func NewStageNotificationService(pump *StageNotificationPump, handlerFactory StageNotificationHandlerFactoryFunc) *StageNotificationService {
	return &StageNotificationService{
		ctx:            context.Background(),
		pump:           pump,
		handlerFactory: handlerFactory,
		log:            log.WithFields(log.Fields{}),
		results:        make(map[uint]chan stageNotificationHandlerResult),
	}
}

// stub out hosting.Service interface on StageNotificationService
func (s *StageNotificationService) Start() {
	s.pump.SetReceiver(s.receive)
	s.pump.Start()

	go s.start()
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
		for _, done := range s.results {
			select {
			case result := <-done:
				s.log.Infof("Stage notification [%d] handler completed", result.Id)

				if result.Error != nil {
					s.log.Errorf("Error handling stage notification %d: %s", result.Id, result.Error)
				}
				delete(s.results, result.Id)
				return
			default:
				continue
			}
		}
	}
}

func (s *StageNotificationService) receive(notification *model.StageNotification) error {
	handler, err := s.handlerFactory()
	if err != nil {
		return err
	}

	id := notification.ID
	done := make(chan stageNotificationHandlerResult, 2)

	s.results[id] = done

	context := &stageNotificationHandlerContext{
		ctx:          s.ctx,
		notification: notification,
		done:         done,
	}
	go handler.Handle(context)

	return nil
}
