package notification

import (
	"encoding/json"
	"log"

	"quiz-realtime/internal/constants"
	quizDTO "quiz-realtime/internal/dto/quiz"
	wsDTO "quiz-realtime/internal/dto/websocket"
	ws "quiz-realtime/internal/infrastructure/websocket"
)

type WebsocketBroadcaster struct {
	Hub *ws.Hub
}

func NewWebsocketBroadcaster(hub *ws.Hub) *WebsocketBroadcaster {
	return &WebsocketBroadcaster{
		Hub: hub,
	}
}

func (b *WebsocketBroadcaster) BroadcastLeaderboardUpdated(resp quizDTO.SubmitAnswerResponse) error {
	if b.Hub == nil {
		return nil
	}

	payload := wsDTO.LeaderboardUpdateResponse{
		Type:        constants.WebSocketEventLeaderboardUpdate,
		SessionID:   resp.SessionID,
		UserID:      resp.UserID,
		Score:       resp.Score,
		Leaderboard: resp.Leaderboard,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal leaderboard payload: %v", err)
		return err
	}

	b.Hub.BroadcastToSession(resp.SessionID, data)
	return nil
}
