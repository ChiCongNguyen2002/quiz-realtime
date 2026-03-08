package quiz

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	ID    string `gorm:"primaryKey;column:id"`
	Title string `gorm:"column:title"`
}

func (Quiz) TableName() string {
	return "quizzes"
}

type Question struct {
	gorm.Model
	ID            string `gorm:"primaryKey;column:id"`
	QuizID        string `gorm:"column:quiz_id"`
	Content       string `gorm:"column:content"`
	CorrectAnswer string `gorm:"column:correct_answer"`
}

func (Question) TableName() string {
	return "questions"
}

type UserAnswer struct {
	gorm.Model
	UserID     string `gorm:"column:user_id"`
	QuestionID string `gorm:"column:question_id"`
	Answer     string `gorm:"column:answer"`
}

func (UserAnswer) TableName() string {
	return "user_answers"
}
