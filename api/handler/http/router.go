package http

import (
	"net/http"

	ws "quiz-realtime/internal/infrastructure/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *QuizHandler, hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/sessions")
	api.POST("", handler.CreateSession)
	api.POST("/:session_id/join", handler.JoinSession)
	api.POST("/:session_id/submit", handler.SubmitAnswer)
	api.GET("/:session_id/leaderboard", handler.GetLeaderboard)

	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}
