package models

import (
	"encoding/hex"

	common "go.opentelemetry.io/proto/otlp/common/v1"
	trace "go.opentelemetry.io/proto/otlp/trace/v1"
)

type SpanNode struct {
	TraceID           string            `json:"traceId"`
	SpanID            string            `json:"spanId"`
	ParentSpanID      string            `json:"parentSpanId,omitempty"`
	Name              string            `json:"name"`
	Kind              string            `json:"kind"`
	ServiceName       string            `json:"serviceName"`
	StartTimeUnixNano uint64            `json:"startTimeUnixNano"`
	EndTimeUnixNano   uint64            `json:"endTimeUnixNano"`
	DurationNanos     uint64            `json:"durationNanos"`
	Attributes        map[string]string `json:"attributes"`
}

func (s *SpanNode) Raw() map[string]any {
	resp := map[string]any{
		"id":                s.SpanID,
		"traceId":           s.TraceID,
		"spanId":            s.SpanID,
		"parentSpanId":      s.ParentSpanID,
		"name":              s.Name,
		"kind":              s.Kind,
		"serviceName":       s.ServiceName,
		"startTimeUnixNano": s.StartTimeUnixNano,
		"endTimeUnixNano":   s.EndTimeUnixNano,
		"durationNanos":     s.DurationNanos,
	}

	for k, v := range s.Attributes {
		if _, ok := resp[k]; ok {
			continue
		}
		resp[k] = v
	}

	return resp
}

func NewSpanNodeFromV1Span(span *trace.Span, serviceName string) SpanNode {
	return SpanNode{
		TraceID:           hex.EncodeToString(span.GetTraceId()),
		SpanID:            hex.EncodeToString(span.GetSpanId()),
		ParentSpanID:      hex.EncodeToString(span.GetParentSpanId()),
		Name:              span.GetName(),
		Kind:              span.GetKind().String(),
		ServiceName:       serviceName,
		StartTimeUnixNano: span.GetStartTimeUnixNano(),
		EndTimeUnixNano:   span.GetEndTimeUnixNano(),
		DurationNanos:     span.GetEndTimeUnixNano() - span.GetStartTimeUnixNano(),
		Attributes:        attributesToMap(span.GetAttributes()),
	}
}

func attributesToMap(attrs []*common.KeyValue) map[string]string {
	m := make(map[string]string)
	for _, attr := range attrs {
		m[attr.GetKey()] = attr.GetValue().GetStringValue()
	}
	return m
}
