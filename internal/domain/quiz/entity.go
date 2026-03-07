package quiz

type Quiz struct {
	ID    string
	Title string
}

type Question struct {
	ID            string
	QuizID        string
	Content       string
	CorrectAnswer string
}

type UserAnswer struct {
	UserID     string
	QuestionID string
	Answer     string
}
