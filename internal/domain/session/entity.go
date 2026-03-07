package session

type Session struct {
	ID     string
	QuizID string
}

type Participant struct {
	ID        string
	SessionID string
	UserID    string
}
