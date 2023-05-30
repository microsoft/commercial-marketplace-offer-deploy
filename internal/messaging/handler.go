package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/labstack/gommon/log"
)

//region message handler context

type MessageHandlerContext struct {
	context         context.Context
	ReceivedMessage *azservicebus.ReceivedMessage
}

func (c *MessageHandlerContext) Context() context.Context {
	return c.context
}

//endregion message handler context

// message handler

const MessageHandlerMethodName = "Handle"

type ServiceBusMessageHandler interface {
	Handle(ctx context.Context, message *azservicebus.ReceivedMessage) error
}

// inner handler we can use to dispatch to our message handlers
type serviceBusMessageHandler struct {
	handler      any
	handleMethod reflect.Value
	messageType  reflect.Type
}

// Handle implements MessageHandler for service bus usage
func (h *serviceBusMessageHandler) Handle(ctx context.Context, message *azservicebus.ReceivedMessage) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic handling message type: %s", h.messageType.Name())
		}
	}()

	typedMessage := reflect.New(h.messageType)

	if message != nil {
		err := json.Unmarshal(message.Body, typedMessage.Interface())
		if err != nil {
			return fmt.Errorf("message unmarshal failure: %w", err)
		}
	}

	context := MessageHandlerContext{
		context:         ctx,
		ReceivedMessage: message,
	}

	result := h.handleMethod.Call([]reflect.Value{
		typedMessage,
		reflect.ValueOf(context),
	})

	// should be a return of error
	errorRef := result[0].Interface()

	if errorRef != nil {
		return errorRef.(error)
	}

	return nil
}

// Creates a new service bus message handler from any handler struct by wrapping it
// required implementation of the handler: { Handle(message *T, context MessageHandlerContext) error }
func NewServiceMessageHandler(handler any) (ServiceBusMessageHandler, error) {
	err := validateHandler(handler)
	if err != nil {
		return nil, err
	}
	return &serviceBusMessageHandler{
		handler:      handler,
		handleMethod: getHandleMethod(handler),
		messageType:  getMessageType(handler),
	}, nil
}

func getHandleMethod(handler any) reflect.Value {
	return reflect.ValueOf(handler).MethodByName(MessageHandlerMethodName)
}

func getMessageType(handler any) reflect.Type {
	handlerType := reflect.TypeOf(handler)
	method, _ := handlerType.MethodByName(MessageHandlerMethodName)
	messageParameterIndex := 1

	parameterType := method.Type.In(messageParameterIndex).Elem()
	return parameterType
}

func validateHandler(handler any) error {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Ptr {
		return fmt.Errorf("handler must be a pointer to a struct")
	}
	if handlerType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("handler must be a struct")
	}

	method, ok := handlerType.MethodByName(MessageHandlerMethodName)

	if !ok {
		return fmt.Errorf("handler does not implement method with the name %s", MessageHandlerMethodName)
	}

	err := validateHandleFunc(method)
	if err != nil {
		return err
	}

	return nil
}

func validateHandleFunc(method reflect.Method) error {
	if method.Type.Kind() != reflect.Func {
		return fmt.Errorf("handleFunc must be a function")
	}

	//confirm signature

	funcParameterCount := method.Type.NumIn() - 1 //take one off for the struct ref in the method signature
	requiredParameterCount := 2

	if funcParameterCount != requiredParameterCount {
		return fmt.Errorf("handleFunc must have exactly 2 parameters, not %d", funcParameterCount)
	}

	messageParameterType := method.Type.In(1)
	if messageParameterType.Kind() == reflect.Ptr {
		messageParameterType = messageParameterType.Elem()
	}

	if messageParameterType.Kind() != reflect.Struct {
		return fmt.Errorf("handleFunc first parameter must be a pointer to a struct")
	}

	contextParameterType := method.Type.In(2)

	if contextParameterType.Kind() == reflect.Ptr {
		contextParameterType = contextParameterType.Elem()
	}

	if contextParameterType.Name() != reflect.TypeOf(MessageHandlerContext{}).Name() {
		return fmt.Errorf("handleFunc second parameter must be a pointer to a MessageHandlerContext")
	}

	return nil
}
