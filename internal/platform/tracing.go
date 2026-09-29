package platform

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// shutdownTracing — отправляет накопленные спаны и закрывает соединение. Вызывается в main при выходе.
var shutdownTracing = func(context.Context) error { return nil }

// InitTracing настраивает глобальный OpenTelemetry: откуда трейсы (service.name), куда их слать,
// какую долю сохранять и как передавать контекст между сервисами (заголовок traceparent).
// Если TRACING_OTLP_ENDPOINT не задан — ничего не делает: otel.Tracer(...) вернёт no-op
// и весь код трейсинга работает вхолостую, без накладных расходов.
func InitTracing(ctx context.Context, cfg *config.Config, service string, log *slog.Logger) error {
	// W3C Trace Context — стандартный формат заголовка traceparent.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if cfg.Tracing.Endpoint == "" {
		log.Info("tracing disabled: TRACING_OTLP_ENDPOINT is empty")
		return nil
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Tracing.Endpoint),
		otlptracegrpc.WithInsecure(), // локально без TLS
	)
	if err != nil {
		return fmt.Errorf("create otlp exporter: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		// Спаны копятся и уходят пачками в фоне — запросы не ждут отправки.
		sdktrace.WithBatcher(exporter),
		// service.name — по нему трейсы делятся по сервисам в Tempo (как job в Prometheus и service в Loki).
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
		// ParentBased: если вызывающий сервис уже решил «сохраняем трейс» — продолжаем его,
		// иначе решаем сами с вероятностью SampleRatio.
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Tracing.SampleRatio))),
	)
	otel.SetTracerProvider(provider)
	shutdownTracing = provider.Shutdown

	log.Info("tracing enabled", "endpoint", cfg.Tracing.Endpoint, "sample_ratio", cfg.Tracing.SampleRatio)
	return nil
}

// ShutdownTracing дожидается отправки последних спанов. Без него то, что накопилось
// в батче перед остановкой процесса, потерялось бы.
func ShutdownTracing(ctx context.Context) {
	if err := shutdownTracing(ctx); err != nil {
		slog.Error("shutdown tracing", "error", err)
	}
}
