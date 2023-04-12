package log

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type LogMessage struct {
	JSONPayload string
}

type LogPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *LogMessage) error
}

type logPublisher struct {
	logMode string
	sender  string
}

func NewLogPublisher(sender string, logMode string) LogPublisher {
	bootStrapOTel()
	publisher := &logPublisher{sender: sender, logMode: logMode}

	return publisher
}

func (p *logPublisher) Publish(message *LogMessage) error {
	log.Printf("recieved logged mesage: %s", message)

	// TODO: write to open telemetry which will write to app insights

	return nil
}

// Open Telemetry

var tracer trace.Tracer

func bootStrapOTel() {

	ctx := context.Background()

	exp, err := newExporter()
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer("MODM_Tracer")

	ctx, span := tracer.Start(ctx, "MODM_tracer_start")
	defer span.End()

	span.AddEvent("Hello world event!")

}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("MODM"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func newExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
