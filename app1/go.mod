module github.com/e-gov/opentelemetry-poc/app1

go 1.13

replace (
	go.opentelemetry.io/otel => ../../opentelemetry-go
	go.opentelemetry.io/otel/exporters/stdout => ../../opentelemetry-go/exporters/stdout
	go.opentelemetry.io/otel/sdk => ../../opentelemetry-go/sdk
)

require (
	go.opentelemetry.io/otel v0.10.0
	go.opentelemetry.io/otel/exporters/stdout v0.10.0
	go.opentelemetry.io/otel/sdk v0.10.0
)
