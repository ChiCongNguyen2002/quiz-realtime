package postgres

import (
	"database/sql"

	"quiz-realtime/internal/domain/leaderboard"

	"github.com/google/uuid"
)

type ScoreRepository struct {
	DB *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{DB: db}
}

func (r *ScoreRepository) SaveScore(sessionID string, userID string, score int) error {
	_, err := r.DB.Exec(`
		INSERT INTO scores (id, session_id, user_id, score, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (session_id, user_id)
		DO UPDATE SET
			score = EXCLUDED.score,
			updated_at = NOW()
	`, uuid.NewString(), sessionID, userID, score)

	return err
}

func (r *ScoreRepository) GetTopScores(sessionID string, limit int) ([]leaderboard.Entry, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.DB.Query(`
		SELECT user_id, score
		FROM scores
		WHERE session_id = $1
		ORDER BY score DESC
		LIMIT $2
	`, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []leaderboard.Entry

	for rows.Next() {
		var e leaderboard.Entry
		if err := rows.Scan(&e.UserID, &e.Score); err != nil {
			return nil, err
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
