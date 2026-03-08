package postgres

import (
	"quiz-realtime/internal/domain/leaderboard"

	"gorm.io/gorm"
)

type ScoreRepository struct {
	DB *gorm.DB
}

func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{DB: db}
}

func (r *ScoreRepository) SaveScore(sessionID string, userID string, score int) error {
	var existing leaderboard.Score
	err := r.DB.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.DB.Create(&leaderboard.Score{
			SessionID: sessionID,
			UserID:    userID,
			Score:     score,
		}).Error
	} else if err != nil {
		return err
	}

	return r.DB.Model(&existing).Update("score", score).Error
}

func (r *ScoreRepository) GetTopScores(sessionID string, limit int) ([]leaderboard.Entry, error) {
	if limit <= 0 {
		limit = 10
	}

	var scores []leaderboard.Score
	err := r.DB.Where("session_id = ?", sessionID).
		Order("score DESC").
		Limit(limit).
		Find(&scores).Error

	if err != nil {
		return nil, err
	}

	entries := make([]leaderboard.Entry, 0, len(scores))
	for _, s := range scores {
		entries = append(entries, leaderboard.Entry{
			UserID: s.UserID,
			Score:  s.Score,
		})
	}

	return entries, nil
}
