package postgres

import (
	"database/sql"

	"quiz-realtime/internal/domain/session"

	"github.com/google/uuid"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) Create(quizID string) (*session.Session, error) {
	id := uuid.NewString()
	_, err := r.DB.Exec(`
		INSERT INTO quiz_sessions (id, quiz_id, started_at)
		VALUES ($1, $2, NOW())
	`, id, quizID)
	if err != nil {
		return nil, err
	}

	return &session.Session{
		ID:     id,
		QuizID: quizID,
	}, nil
}

func (r *SessionRepository) GetByID(id string) (*session.Session, error) {
	row, err := r.DB.Query(`
		SELECT id, quiz_id
		FROM quiz_sessions
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, sql.ErrNoRows
	}

	var s session.Session
	if err := row.Scan(&s.ID, &s.QuizID); err != nil {
		return nil, err
	}

	return &s, nil
}

type ParticipantRepository struct {
	DB *sql.DB
}

func NewParticipantRepository(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{DB: db}
}

func (r *ParticipantRepository) AddParticipant(sessionID string, userID string) error {
	_, err := r.DB.Exec(`
		INSERT INTO participants (id, session_id, user_id)
		VALUES ($1, $2, $3)
	`, uuid.NewString(), sessionID, userID)

	return err
}
