package websocket

import (
	domainLeaderboard "quiz-realtime/internal/domain/leaderboard"
)

type LeaderboardUpdateResponse struct {
	Type        string                    `json:"type"`
	SessionID   string                    `json:"session_id"`
	UserID      string                    `json:"user_id"`
	Score       int                       `json:"score"`
	Leaderboard []domainLeaderboard.Entry `json:"leaderboard"`
}
