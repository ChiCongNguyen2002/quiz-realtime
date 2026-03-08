package postgres

import (
	"quiz-realtime/internal/domain/session"

	"gorm.io/gorm"
)

type ParticipantRepository struct {
	DB *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{DB: db}
}

func (r *ParticipantRepository) AddParticipant(sessionID string, userID string) error {
	participant := &session.Participant{
		SessionID: sessionID,
		UserID:    userID,
	}
	return r.DB.Create(participant).Error
}
