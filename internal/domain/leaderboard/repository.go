package leaderboard

type Repository interface {
	UpdateScore(sessionID string, userID string, score int) error
	GetLeaderboard(sessionID string) ([]Entry, error)
}

type ScoreRepository interface {
	SaveScore(sessionID string, userID string, score int) error
	GetTopScores(sessionID string, limit int) ([]Entry, error)
}
