package postgres

import (
	"quiz-realtime/internal/domain/quiz"

	"gorm.io/gorm"
)

type QuizRepository struct {
	DB *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *QuizRepository {
	return &QuizRepository{DB: db}
}

func (r *QuizRepository) GetQuestionsByQuizID(quizID string) ([]quiz.Question, error) {
	var questions []quiz.Question

	err := r.DB.Where("quiz_id = ?", quizID).Find(&questions).Error

	return questions, err
}
