package quiz

import (
	"quiz-realtime/internal/domain/leaderboard"
	"quiz-realtime/internal/domain/quiz"
	"quiz-realtime/internal/domain/session"
	dto "quiz-realtime/internal/dto/quiz"
)

type LeaderboardBroadcaster interface {
	BroadcastLeaderboardUpdated(resp dto.SubmitAnswerResponse) error
}

type Service struct {
	QuizRepo        quiz.Repository
	SessionRepo     session.SessionRepository
	ParticipantRepo session.ParticipantRepository
	ScoreRepo       leaderboard.ScoreRepository
	LeaderboardRepo leaderboard.Repository
	Broadcaster     LeaderboardBroadcaster
}

func NewService(
	quizRepo quiz.Repository,
	sessionRepo session.SessionRepository,
	participantRepo session.ParticipantRepository,
	scoreRepo leaderboard.ScoreRepository,
	leaderboardRepo leaderboard.Repository,
	broadcaster LeaderboardBroadcaster,
) *Service {
	return &Service{
		QuizRepo:        quizRepo,
		SessionRepo:     sessionRepo,
		ParticipantRepo: participantRepo,
		ScoreRepo:       scoreRepo,
		LeaderboardRepo: leaderboardRepo,
		Broadcaster:     broadcaster,
	}
}

func (s *Service) SubmitAnswer(sessionID string, req dto.SubmitAnswerRequest) (*dto.SubmitAnswerResponse, error) {
	sess, err := s.SessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, err
	}

	questions, err := s.QuizRepo.GetQuestionsByQuizID(sess.QuizID)
	if err != nil {
		return nil, err
	}

	score := s.calculateScore(questions, req.Answers)

	if err := s.ScoreRepo.SaveScore(sessionID, req.UserID, score); err != nil {
		return nil, err
	}

	if err := s.LeaderboardRepo.UpdateScore(sessionID, req.UserID, score); err != nil {
		return nil, err
	}

	entries, err := s.getLeaderboardEntries(sessionID)
	if err != nil {
		return nil, err
	}

	resp := &dto.SubmitAnswerResponse{
		SessionID:   sessionID,
		UserID:      req.UserID,
		Score:       score,
		Leaderboard: entries,
	}

	if s.Broadcaster != nil {
		_ = s.Broadcaster.BroadcastLeaderboardUpdated(*resp)
	}

	return resp, nil
}

func (s *Service) calculateScore(questions []quiz.Question, answers []dto.Answer) int {
	correctMap := make(map[string]string, len(questions))
	for _, q := range questions {
		correctMap[q.ID] = q.CorrectAnswer
	}

	score := 0
	for _, ans := range answers {
		if correctMap[ans.QuestionID] == ans.Answer {
			score++
		}
	}
	return score
}

func (s *Service) getLeaderboardEntries(sessionID string) ([]leaderboard.Entry, error) {
	if s.LeaderboardRepo != nil {
		return s.LeaderboardRepo.GetLeaderboard(sessionID)
	}

	if s.ScoreRepo != nil {
		return s.ScoreRepo.GetTopScores(sessionID, 10)
	}

	return nil, nil
}

func (s *Service) GetLeaderboard(sessionID string) (*dto.GetLeaderboardResponse, error) {
	entries, err := s.getLeaderboardEntries(sessionID)
	if err != nil {
		return nil, err
	}

	return &dto.GetLeaderboardResponse{
		SessionID:   sessionID,
		Leaderboard: entries,
	}, nil
}

func (s *Service) CreateSession(req dto.CreateSessionRequest) (*dto.CreateSessionResponse, error) {
	sess, err := s.SessionRepo.Create(req.QuizID)
	if err != nil {
		return nil, err
	}

	return &dto.CreateSessionResponse{
		SessionID: sess.ID,
		QuizID:    sess.QuizID,
	}, nil
}

func (s *Service) JoinSession(sessionID string, req dto.JoinSessionRequest) (*dto.JoinSessionResponse, error) {
	if err := s.ParticipantRepo.AddParticipant(sessionID, req.UserID); err != nil {
		return nil, err
	}

	return &dto.JoinSessionResponse{
		SessionID: sessionID,
		UserID:    req.UserID,
	}, nil
}
