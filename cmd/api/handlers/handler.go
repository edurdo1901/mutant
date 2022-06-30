package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"prueba.com/internal/mutant"
)

type Mutant interface {
	IsMutant(context.Context, []string) (bool, error)
	Stats(context.Context) (mutant.StatsResponse, error)
}

type Handler struct {
	mutant Mutant
}

func New(mutant Mutant) *Handler {
	return &Handler{mutant}
}

// API start handlers.
func (h *Handler) API(engine *gin.Engine) {
	engine.Use(errorHandler)
	engine.POST("/mutant", ismutant(h.mutant))
	engine.GET("/stats", stats(h.mutant))
	engine.GET("/", home)
}

// handler home
func home(c *gin.Context) {
	c.JSON(http.StatusOK, "Prueba mercado libre")
}
