package notification

import (
	"context"
)

type NotificationPumpReceiveFunc[T any] func(notification *T) error

type NotificationPump[T any] interface {
	Start()
	Stop()
	SetReceiver(receiver NotificationPumpReceiveFunc[T])
}

//region notification handler

// factory function for creating a notification handler
type NotificationHandlerFactoryFunc[T any] func() (NotificationHandler[T], error)

// notification handler function
type NotificationHandlerFunc[T any] func(context *NotificationHandlerContext[T])

// notification handler
type NotificationHandler[T any] interface {
	Handle(context *NotificationHandlerContext[T])
}

//endregion notification handler

//region context

type NotificationHandlerResult[T any] struct {
	Notification *T
	Error        error
}

type NotificationHandlerContext[T any] struct {
	ctx          context.Context
	done         chan NotificationHandlerResult[T]
	Notification *T
}

func NewNotificationHandlerContext[T any](ctx context.Context, notification *T) *NotificationHandlerContext[T] {
	return &NotificationHandlerContext[T]{
		ctx:          ctx,
		done:         make(chan NotificationHandlerResult[T], 1),
		Notification: notification,
	}
}

func (c *NotificationHandlerContext[T]) Context() context.Context {
	return c.ctx
}

func (c *NotificationHandlerContext[T]) Done(result NotificationHandlerResult[T]) {
	c.done <- result
}

func (c *NotificationHandlerContext[T]) Channel() chan NotificationHandlerResult[T] {
	return c.done
}

//endregion context
