package telemetry

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	logssdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Some common variables
	backgroundCtx := context.Background()
	resource := newResource(backgroundCtx)

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up logger.
	loggerProvider, err := newLoggerProvider(backgroundCtx, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)

	// Wire up logger to OpenTelemetry
	otelLogger := slog.New(otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{}))
	slog.SetDefault(otelLogger)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(backgroundCtx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(backgroundCtx, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	return
}

// Creates a new OenTelemetry logger provider in Go.
// NOTE: No official logger provider exists in OpenTelemetry,
// so using a temporary implemementation.
//
// ctx context.Context, resource *resource.Resource
// *logssdk.LoggerProvider, error
func newLoggerProvider(ctx context.Context, resource *resource.Resource) (*logssdk.LoggerProvider, error) {
	logExporter, _ := otlplogs.NewExporter(ctx)
	loggerProvider := logssdk.NewLoggerProvider(
		logssdk.WithBatcher(logExporter),
		logssdk.WithResource(resource),
	)

	return loggerProvider, nil
}

// Create a new TextMapPropagator.
//
// No parameters.
// Returns a TextMapPropagator.
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// Creates a new OpenTelemetry resource that represents this service.
//
// ctx context.Context
// *resource.Resource
func newResource(ctx context.Context) *resource.Resource {
	hostName, _ := os.Hostname()
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName("todo-api-go"),
			semconv.ServiceVersion("v1.0.0"),
			semconv.HostName(hostName),
		),
	)

	return res
}

// Create a new trace.TracerProvider and an error.
//
// No parameters.
// Returns a *trace.TracerProvider and an error.
func newTraceProvider(ctx context.Context, resource *resource.Resource) (*trace.TracerProvider, error) {
	stdoutTraceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	otlpTraceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(resource),

		trace.WithBatcher(stdoutTraceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),

		trace.WithBatcher(otlpTraceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
	)

	return traceProvider, nil
}

// Create a new MeterProvider.
//
// No parameters.
// Returns a *metric.MeterProvider and an error.
func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	stdoutMetricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	otlpMetricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(stdoutMetricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),

		metric.WithReader(metric.NewPeriodicReader(otlpMetricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)

	return meterProvider, nil
}
