package http

import (
	nethttp "net/http"

	dto "quiz-realtime/internal/dto/quiz"
	appQuiz "quiz-realtime/internal/service/quiz"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	Service *appQuiz.Service
}

func (h *QuizHandler) CreateSession(c *gin.Context) {
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Service.CreateSession(req)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, resp)
}

func (h *QuizHandler) JoinSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req dto.JoinSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Service.JoinSession(sessionID, req)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, resp)
}

func (h *QuizHandler) SubmitAnswer(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req dto.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Service.SubmitAnswer(sessionID, req)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, resp)
}

func (h *QuizHandler) GetLeaderboard(c *gin.Context) {
	sessionID := c.Param("session_id")

	resp, err := h.Service.GetLeaderboard(sessionID)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, resp)
}
