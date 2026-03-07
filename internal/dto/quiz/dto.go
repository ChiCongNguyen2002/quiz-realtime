package quiz

import (
	domainLeaderboard "quiz-realtime/internal/domain/leaderboard"
)

type SubmitAnswerRequest struct {
	UserID  string `json:"user_id"`
	Answers []struct {
		QuestionID string `json:"question_id"`
		Answer     string `json:"answer"`
	} `json:"answers"`
}

type SubmitAnswerResponse struct {
	SessionID   string                    `json:"session_id"`
	UserID      string                    `json:"user_id"`
	Score       int                       `json:"score"`
	Leaderboard []domainLeaderboard.Entry `json:"leaderboard"`
}

type CreateSessionRequest struct {
	QuizID string `json:"quiz_id"`
}

type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	QuizID    string `json:"quiz_id"`
}

type JoinSessionRequest struct {
	UserID string `json:"user_id"`
}

type JoinSessionResponse struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}

type GetLeaderboardResponse struct {
	SessionID   string                    `json:"session_id"`
	Leaderboard []domainLeaderboard.Entry `json:"leaderboard"`
}
