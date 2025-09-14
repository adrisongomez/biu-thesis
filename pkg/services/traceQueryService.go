package services

import (
	"net/http"

	"github.com/adrisongomez/thesis/pkg/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TraceQueryService struct {
	lg   *zap.SugaredLogger
	repo repository.TraceQueryRepository
}

func (t *TraceQueryService) GetTraceHandler(c echo.Context) error {
	traceID := c.Param("traceId")
	t.lg.Infow("GetTraceHandler got called with", "traceID", traceID)

	if traceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "traceId is required"})
	}

	spans, err := t.repo.GetTraceByID(c.Request().Context(), traceID)

	if err != nil {
		t.lg.Errorw("Error getting trace by ID", "erro", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve trace"})
	}
	if len(spans) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "trace not found"})
	}
	return c.JSON(http.StatusOK, spans)
}

func (t *TraceQueryService) GetTracesList(c echo.Context) error {
	t.lg.Info("GetTracesList got called")
	traces, err := t.repo.GetTraces(c.Request().Context())
	if err != nil {
		t.lg.Errorf("Error getting traces: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve traces"})
	}
	return c.JSON(http.StatusOK, traces)
}

func NewTraceQueryService(repo repository.TraceQueryRepository) TraceQueryService {
	return TraceQueryService{
		lg:   zap.S(),
		repo: repo,
	}
}
