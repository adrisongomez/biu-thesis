package models

type TraceNode struct {
	TraceID string `json:"traceId"`
}

func (t *TraceNode) Raw() map[string]any {
	return map[string]any{"traceId": t.TraceID}
}

func NewTraceNode(traceId string) TraceNode {
	return TraceNode{TraceID: traceId}
}
