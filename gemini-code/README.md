# Custom OpenTelemetry Trace Receiver

This project contains a simple Go service that acts as a receiver for OpenTelemetry (OTLP) traces. It's designed to show how you can build a custom backend to receive and process telemetry data from an OpenTelemetry Collector.How It WorksInstrumented Application: Sends traces to the OpenTelemetry Collector.OpenTelemetry Collector: Receives the traces, batches them, and forwards them to two destinations:The console (via the logging exporter).Our custom Go service (via the otlp exporter).Go Service (main.go): Listens on gRPC port 4317, receives the trace data, and prints a summary of the spans to its console.PrerequisitesGo: Version 1.21 or newer.OpenTelemetry Collector: Download the "contrib" distribution from the official releases page.How to RunFollow these steps in separate terminal windows.1. Run the Go Trace ServiceThis is your custom backend. It will wait for the collector to send it data.# First, download the necessary Go modules
go mod tidy

## Then, run the service

go run main.go
You should see the output: OTLP trace receiver listening on gRPC port 43172. Run the OpenTelemetry CollectorThe collector will receive data from an application and forward it to your Go service.# Make sure 'otel-collector-config.yaml' is in the same directory

## The executable name may vary based on your OS (e.g., otelcol-contrib_windows_amd64.exe)

./otelcol-contrib --config otel-collector-config.yaml
3. Send a Sample TraceNow, you need an application to generate and send a trace. The easiest way to do this is to run a simple, pre-instrumented application. You can use the OpenTelemetry demo application or a simple script.Here is an example using curl to send a minimal OTLP/JSON trace to the collector's HTTP endpoint. Save the following as trace.json:{
  "resourceSpans": [
    {
      "resource": {
        "attributes": [
          { "key": "service.name", "value": { "stringValue": "curl-test-service" } }
        ]
      },
      "scopeSpans": [
        {
          "scope": { "name": "manual-test" },
          "spans": [
            {
              "traceId": "0102030405060708090a0b0c0d0e0f10",
              "spanId": "1112131415161718",
              "name": "manual-curl-span",
              "kind": "SPAN_KIND_INTERNAL",
              "startTimeUnixNano": "1700000000000000000",
              "endTimeUnixNano": "1700000001000000000",
              "attributes": [
                { "key": "http.method", "value": { "stringValue": "GET" } }
              ]
            }
          ]
        }
      ]
    }
  ]
}
Now send it using curl:curl -X POST -H "Content-Type: application/json" \
  -d @trace.json <http://localhost:4318/v1/traces>
Expected OutcomeIn the Collector's terminal, you will see the trace data printed by the logging exporter.In the Go service's terminal, you will see the output from our custom code, which should look something like this:Received traces!
Number of resource spans: 1
  Resource: service.name=curl-test-service
    Scope: manual-test
    Number of spans in this scope: 1
      - Span: manual-curl-span (TraceID: 102030405060708090a0b0c0d0e0f10, SpanID: 1112131415161718, ParentSpanID: )
        Attributes: http.method=GET
This confirms that your entire pipeline is working correctly! You can now build on this Go service to store, analyze, or forward the traces as needed.