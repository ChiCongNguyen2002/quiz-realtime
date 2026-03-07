package quiz

type Repository interface {
	GetQuestionsByQuizID(quizID string) ([]Question, error)
}
