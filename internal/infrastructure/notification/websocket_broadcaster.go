package notification

import (
	"encoding/json"

	domainLeaderboard "quiz-realtime/internal/domain/leaderboard"
	dto "quiz-realtime/internal/dto/quiz"
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

func (b *WebsocketBroadcaster) BroadcastLeaderboardUpdated(resp dto.SubmitAnswerResponse) error {
	if b.Hub == nil {
		return nil
	}

	payload := struct {
		Type        string                    `json:"type"`
		SessionID   string                    `json:"session_id"`
		Leaderboard []domainLeaderboard.Entry `json:"leaderboard"`
	}{
		Type:        "leaderboard_update",
		SessionID:   resp.SessionID,
		Leaderboard: resp.Leaderboard,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	b.Hub.Broadcast <- data
	return nil
}
