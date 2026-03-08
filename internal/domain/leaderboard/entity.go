package leaderboard

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	SessionID string `gorm:"column:session_id;uniqueIndex:idx_session_user"`
	UserID    string `gorm:"column:user_id;uniqueIndex:idx_session_user"`
	Score     int    `gorm:"column:score"`
}

func (Score) TableName() string {
	return "scores"
}

type Entry struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
}
