package quiz

import (
	"log"

	domainLeaderboard "quiz-realtime/internal/domain/leaderboard"
	domainQuiz "quiz-realtime/internal/domain/quiz"
	domainSession "quiz-realtime/internal/domain/session"
	dto "quiz-realtime/internal/dto/quiz"
)

type LeaderboardBroadcaster interface {
	BroadcastLeaderboardUpdated(resp dto.SubmitAnswerResponse) error
}

type Service struct {
	QuizRepo        domainQuiz.Repository
	SessionRepo     domainSession.SessionRepository
	ParticipantRepo domainSession.ParticipantRepository
	ScoreRepo       domainLeaderboard.ScoreRepository
	LeaderboardRepo domainLeaderboard.Repository
	Broadcaster     LeaderboardBroadcaster
}

func NewService(
	quizRepo domainQuiz.Repository,
	sessionRepo domainSession.SessionRepository,
	participantRepo domainSession.ParticipantRepository,
	scoreRepo domainLeaderboard.ScoreRepository,
	leaderboardRepo domainLeaderboard.Repository,
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

func (s *Service) SubmitAnswer(sessionID string, req dto.SubmitAnswerRequest) (dto.SubmitAnswerResponse, error) {
	sess, err := s.SessionRepo.GetByID(sessionID)
	if err != nil {
		return dto.SubmitAnswerResponse{}, err
	}

	questions, err := s.QuizRepo.GetQuestionsByQuizID(sess.QuizID)
	if err != nil {
		return dto.SubmitAnswerResponse{}, err
	}

	correctMap := map[string]string{}
	for _, q := range questions {
		correctMap[q.ID] = q.CorrectAnswer
	}

	score := 0
	for _, ans := range req.Answers {
		if correctMap[ans.QuestionID] == ans.Answer {
			score++
		}
	}

	if err := s.ScoreRepo.SaveScore(sessionID, req.UserID, score); err != nil {
		return dto.SubmitAnswerResponse{}, err
	}

	if err := s.LeaderboardRepo.UpdateScore(sessionID, req.UserID, score); err != nil {
		log.Println("leaderboard redis update failed:", err)
	}

	leaderboard, err := s.getLeaderboardInternal(sessionID)
	if err != nil {
		return dto.SubmitAnswerResponse{}, err
	}

	resp := dto.SubmitAnswerResponse{
		SessionID:   sessionID,
		UserID:      req.UserID,
		Score:       score,
		Leaderboard: leaderboard,
	}

	if s.Broadcaster != nil {
		if err := s.Broadcaster.BroadcastLeaderboardUpdated(resp); err != nil {
			log.Println("broadcast leaderboard failed:", err)
		}
	}

	return resp, nil
}

func (s *Service) getLeaderboardInternal(sessionID string) ([]domainLeaderboard.Entry, error) {
	if s.LeaderboardRepo != nil {
		entries, err := s.LeaderboardRepo.GetLeaderboard(sessionID)
		if err == nil {
			return entries, nil
		}
	}

	if s.ScoreRepo != nil {
		return s.ScoreRepo.GetTopScores(sessionID, 10)
	}

	return []domainLeaderboard.Entry{}, nil
}

func (s *Service) GetLeaderboard(sessionID string) (dto.GetLeaderboardResponse, error) {
	entries, err := s.getLeaderboardInternal(sessionID)
	if err != nil {
		return dto.GetLeaderboardResponse{}, err
	}

	return dto.GetLeaderboardResponse{
		SessionID:   sessionID,
		Leaderboard: entries,
	}, nil
}

func (s *Service) CreateSession(req dto.CreateSessionRequest) (dto.CreateSessionResponse, error) {
	sess, err := s.SessionRepo.Create(req.QuizID)
	if err != nil {
		return dto.CreateSessionResponse{}, err
	}

	return dto.CreateSessionResponse{
		SessionID: sess.ID,
		QuizID:    sess.QuizID,
	}, nil
}

func (s *Service) JoinSession(sessionID string, req dto.JoinSessionRequest) (dto.JoinSessionResponse, error) {
	if err := s.ParticipantRepo.AddParticipant(sessionID, req.UserID); err != nil {
		return dto.JoinSessionResponse{}, err
	}

	return dto.JoinSessionResponse{
		SessionID: sessionID,
		UserID:    req.UserID,
	}, nil
}
