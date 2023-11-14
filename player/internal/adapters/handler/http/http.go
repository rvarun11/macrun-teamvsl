package http

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc services.PlayerService
}

func NewHTTPHandler(PlayerService services.PlayerService) *HTTPHandler {
	return &HTTPHandler{
		svc: PlayerService,
	}
}

func (h *HTTPHandler) RegisterPlayer(ctx *gin.Context) {
	var req playerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	_, err := h.svc.Register(req.toAggregate())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to register, something occured",
		})
		return
	}
	logger.Info("player registered")
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "New player created successfully",
	})
}

func (h *HTTPHandler) GetPlayer(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	p, err := h.svc.Get(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	res := fromAggregate(p)
	ctx.JSON(http.StatusOK, gin.H{
		"player": res,
	})
}

func (h *HTTPHandler) UpdatePlayer(ctx *gin.Context) {
	var req playerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	p, err := h.svc.Update(req.toAggregate())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	res := fromAggregate(p)
	ctx.JSON(http.StatusNoContent, gin.H{
		// 204 returns no body so this can be changed to 200 if body is needed
		"player":  res,
		"message": "Player updated successfully",
	})
}

func (h *HTTPHandler) ListPlayers(ctx *gin.Context) {
	players, err := h.svc.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, players)
}