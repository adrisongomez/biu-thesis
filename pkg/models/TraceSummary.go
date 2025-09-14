package models

type TraceSummary struct {
	TraceID      string `json:"traceId"`
	RootSpanName string `json:"rootSpanName"`
	ServiceName  string `json:"serviceName"`
	StartTime    uint64 `json:"startTime"` // Unix nano
	Duration     uint64 `json:"duration"`  // Nanos
	SpanCount    int64  `json:"spanCount"`
}
