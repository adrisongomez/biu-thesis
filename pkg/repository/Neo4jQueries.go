package repository

const (
	CREATE_SERVICES      = `UNWIND $services AS service_props MERGE (s:Service {name: service_props.name})`
	CREATE_TRACES        = `UNWIND $traces AS trace_props MERGE (t:Trace {traceId: trace_props.traceId})`
	CREATE_SPANS         = `UNWIND $spans AS span_props MERGE (s:Span {spanId: span_props.spanId}) SET s = span_props`
	CREATE_RELATIONSHIPS = `
		// Link Services to Spans
		UNWIND $spans AS span_props
		MATCH (service:Service {name: span_props.serviceName})
		MATCH (span:Span {spanId: span_props.spanId})
		MERGE (service)-[:EMITTED]->(span);

		// Link Traces to Spans
		UNWIND $spans AS span_props
		MATCH (trace:Trace {traceId: span_props.traceId})
		MATCH (span:Span {spanId: span_props.spanId})
		MERGE (trace)-[:CONTAINS]->(span);

		// Link parent Spans to child Spans
		UNWIND $spans AS span_props
		WHERE span_props.parentSpanId IS NOT NULL AND span_props.parentSpanId <> ""
		MATCH (parent:Span {spanId: span_props.parentSpanId})
		MATCH (child:Span {spanId: span_props.spanId})
		MERGE (parent)-[:PARENT_OF]->(child);
	`
)
