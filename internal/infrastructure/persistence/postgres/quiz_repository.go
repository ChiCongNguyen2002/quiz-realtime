package postgres

import (
	"database/sql"

	"quiz-realtime/internal/domain/quiz"
)

type QuizRepository struct {
	DB *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{DB: db}
}

func (r *QuizRepository) GetQuestionsByQuizID(
	quizID string,
) ([]quiz.Question, error) {
	rows, err := r.DB.Query(`
		SELECT id, quiz_id, content, correct_answer
		FROM questions
		WHERE quiz_id = $1
	`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []quiz.Question

	for rows.Next() {
		var q quiz.Question

		if err := rows.Scan(
			&q.ID,
			&q.QuizID,
			&q.Content,
			&q.CorrectAnswer,
		); err != nil {
			return nil, err
		}

		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
