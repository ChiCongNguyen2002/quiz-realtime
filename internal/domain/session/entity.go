package session

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	ID     string `gorm:"primaryKey;column:id"`
	QuizID string `gorm:"column:quiz_id"`
}

func (Session) TableName() string {
	return "sessions"
}

type Participant struct {
	gorm.Model
	SessionID string `gorm:"column:session_id"`
	UserID    string `gorm:"column:user_id"`
}

func (Participant) TableName() string {
	return "participants"
}
