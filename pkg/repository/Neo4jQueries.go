package repository

const (
	CREATE_SERVICES      = `UNWIND $services AS service_props MERGE (s:Service {name: service_props.name})`
	CREATE_TRACES        = `UNWIND $traces AS trace_props MERGE (t:Trace {traceId: trace_props.traceId})`
	CREATE_SPANS         = `UNWIND $spans AS span_props MERGE (s:Span {spanId: span_props.spanId}) SET s = span_props`
	CREATE_LINK_SERVICES = `
		UNWIND $spans AS span_props
		MATCH (service:Service {name: span_props.serviceName})
		MATCH (span:Span {spanId: span_props.spanId})
		MERGE (service)-[:EMITTED]->(span)
	`
	CREATE_LINK_TRACE_TO_SPAN = `
		UNWIND $spans AS span_props
		MATCH (trace:Trace {traceId: span_props.traceId})
		MATCH (span:Span {spanId: span_props.spanId})
		MERGE (trace)-[:CONTAINS]->(span)
	`
	CREATE_LINK_PARENT_TO_CHILD = `
		UNWIND $spans AS span_props
		WITH span_props
		WHERE span_props.parentSpanId IS NOT NULL AND span_props.parentSpanId <> ""
		MATCH (parent:Span {spanId: span_props.parentSpanId})
		MATCH (child:Span {spanId: span_props.spanId})
		MERGE (parent)-[:PARENT_OF]->(child)
	`
)

/**
// 4. Create relationships in separate queries within the same transaction.
		// Link Services to Spans
		queryServiceToSpan := `
			UNWIND $spans AS span_props
			MATCH (service:Service {name: span_props.serviceName})
			MATCH (span:Span {spanId: span_props.spanId})
			MERGE (service)-[:EMITTED]->(span)`
		if _, err := tx.Run(ctx, queryServiceToSpan, map[string]interface{}{"spans": spanMaps}); err != nil {
			return nil, fmt.Errorf("error creating service->span relationship: %w", err)
		}

		// Link Traces to Spans
		queryTraceToSpan := `
			UNWIND $spans AS span_props
			MATCH (trace:Trace {traceId: span_props.traceId})
			MATCH (span:Span {spanId: span_props.spanId})
			MERGE (trace)-[:CONTAINS]->(span)`
		if _, err := tx.Run(ctx, queryTraceToSpan, map[string]interface{}{"spans": spanMaps}); err != nil {
			return nil, fmt.Errorf("error creating trace->span relationship: %w", err)
		}

		// Link parent Spans to child Spans
		queryParentToChild := `
			UNWIND $spans AS span_props
			WHERE span_props.parentSpanId IS NOT NULL AND span_props.parentSpanId <> ""
			MATCH (parent:Span {spanId: span_props.parentSpanId})
			MATCH (child:Span {spanId: span_props.spanId})
			MERGE (parent)-[:PARENT_OF]->(child)`
		if _, err := tx.Run(ctx, queryParentToChild, map[string]interface{}{"spans": spanMaps}); err != nil {
			return nil, fmt.Errorf("error creating parent->child relationship: %w", err)
		}

		return nil, nil
	}

*/
