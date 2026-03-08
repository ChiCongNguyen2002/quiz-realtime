package postgres

import (
	"quiz-realtime/internal/domain/session"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) GetByID(id string) (*session.Session, error) {
	var sess session.Session
	err := r.DB.Where("id = ?", id).First(&sess).Error
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (r *SessionRepository) Create(quizID string) (*session.Session, error) {
	sess := &session.Session{
		ID:     uuid.NewString(),
		QuizID: quizID,
	}
	err := r.DB.Create(sess).Error
	return sess, err
}
