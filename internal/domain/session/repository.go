package session

type SessionRepository interface {
	Create(quizID string) (*Session, error)
	GetByID(id string) (*Session, error)
}

type ParticipantRepository interface {
	AddParticipant(sessionID string, userID string) error
}
