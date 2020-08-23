// mõõtmiskatse-1 on rakendus, mis valmistab lihtsa OpenTelemetry mõõteriista
// (instrument), teeb sellega sarja mõõtmisi ja väljastab mõõtmistulemused
// standardväljundisse.
package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

func main() {
	// Valmista mõõteriist.
	meter1 := global.Meter("Mõõteriist 1")
	// Mõõdetav suurus.
	r, err := meter1.NewInt64Counter("pöördumiste arv")
	if err != nil {
		log.Fatal("Mõõdiku loomine ebaõnnestus: ", err.Error())
	}

	// Valmista ette eksporter.
	var stdoutOpt []stdout.Option
	var pushOpt []push.Option

	stdoutOpt = append(stdoutOpt, stdout.WithPrettyPrint())
	i, _ := time.ParseDuration("5s")
	pushOpt = append(pushOpt, push.WithPeriod(i))

	// "Controller organizes a periodic push of metric data."
	// Vaikimisi saadab iga 10 s järgi.
	pipeline, err := stdout.InstallNewPipeline(stdoutOpt, pushOpt)
	if err != nil {
		log.Fatal("Eksporteri loomine ebaõnnestus: ", err)
	}
	defer pipeline.Stop()

	// Tee mõõtmine.
	ctx := context.Background()

	r.Add(ctx, 1, label.String("A", "1"), label.String("B", "2"))
	r.Add(ctx, 1, label.String("A", "1"), label.String("B", "2"))

}
