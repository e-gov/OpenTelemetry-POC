package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/label"
)

var (
	võti1      = label.Key("POC/võti1")
	võti2      = label.Key("POC/võti2")
	lemonsKey  = label.Key("POC/lemons")
	anotherKey = label.Key("POC/another")
)

func main() {
	// InstallNewPipeline loob uue mõõtetulemuste eksportimise "torustiku".
	pusher, err := stdout.InstallNewPipeline([]stdout.Option{
		stdout.WithQuantiles([]float64{0.5, 0.9, 0.99}), // seadistussuvand
		stdout.WithPrettyPrint(),                        // seadistussuvand
	}, nil)
	if err != nil {
		log.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}
	defer pusher.Stop()

	// Loob trasseerija.
	tracer := global.Tracer("POC/app2")
	// Loob mõõtja.
	meter := global.Meter("POC/app2")

	// Rida märgendeid.
	commonLabels := []label.KeyValue{
		lemonsKey.Int(10),
		label.String("A", "1"),
		label.String("B", "2"), label.String("C", "3"),
	}

	oneMetricCB := func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, commonLabels...)
	}
	// Avaldis tagastab instrumendi tüübist netric.Float64ValueObserver. Tagastatav
	// väärtus ei paku huvi, sest tagasikutsef-n (callback) on vaatlejale (Observer)
	// edastatud. Näidu võtab vaatleja. Instrumendi nimi on POC.instrument-1.
	// ValueObserver tähendab tõmbepõhimõttel mõõduvõtmist.
	_ = metric.Must(meter).NewFloat64ValueObserver(
		"POC.instrument-1",
		oneMetricCB, // tagasikutsef-n
		metric.WithDescription("A ValueObserver set to 1.0"), // seadistussuvand
	)

	// Avaldis tagastab instrumendi tüübist metric.Float64ValueRecorder. Instrumendi
	// nimi on POC.instrument-2.
	valuerecorderTwo := metric.Must(meter).NewFloat64ValueRecorder("POC.instrument-2")

	// Background annab tühja konteksti.
	ctx := context.Background()

	// Konteksti täiendamine võti-väärtus paaridega.
	ctx = correlation.NewContext(ctx,
		võti1.String("väärtus1"),
		võti2.String("väärtus2"),
	)

	// Selle mõttest ei saa aru. Instrumendile lisatakse märgendid, moodustub uus
	// instrument... Mõte on vist selles, et hakkab registreerima väärtusi märgendite
	// iga kombinatsiooni kohta. Vt kasutust allpool, alamvahemikus.
	valuerecorder := valuerecorderTwo.Bind(commonLabels...)
	defer valuerecorder.Unbind()

	// Koheselt täidetav f-navaldis (IIFE).
	err = func(ctx context.Context) error {
		// Loob uue vahemiku (span). trace.Span on liidesetüüp.
		var span trace.Span

		// Alustab vahemiku trasseerimist.
		ctx, span = tracer.Start(ctx, "operation")
		defer span.End()

		// Lisab sündmuse.
		span.AddEvent(ctx, "Nice operation!", label.Int("Sündmuse märgend", 100))
		span.SetAttributes(anotherKey.String("yes"))

		meter.RecordBatch(
			// Note: call-site variables added as context Entries:
			correlation.NewContext(ctx, anotherKey.String("xyz")),
			commonLabels,
			valuerecorderTwo.Measurement(2.0),
		)

		// Teine IIFE.
		return func(ctx context.Context) error {
			// Uus vahemik. Järgmised kolm laused moodustavad idioomi.
			var span trace.Span
			ctx, span = tracer.Start(ctx, "Sub operation...")
			defer span.End()

			// Lisab sündmused
			span.SetAttributes(lemonsKey.String("five"))
			span.AddEvent(ctx, "Sub span event")

			valuerecorder.Record(ctx, 1.3)

			return nil
		}(ctx)

	}(ctx)

	if err != nil {
		panic(err)
	}
}
