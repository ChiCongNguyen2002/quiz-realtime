package main

import (
	"log"
	"quiz-realtime/api/handler/http"

	notification "quiz-realtime/internal/infrastructure/notification"
	pgRepo "quiz-realtime/internal/infrastructure/persistence/postgres"
	redisRepo "quiz-realtime/internal/infrastructure/persistence/redis"
	ws "quiz-realtime/internal/infrastructure/websocket"
	appQuiz "quiz-realtime/internal/service/quiz"
	appConfig "quiz-realtime/pkg/config"
	"quiz-realtime/pkg/database"
	redisClient "quiz-realtime/pkg/redis"
)

func main() {
	cfg, err := appConfig.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgres(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	redis := redisClient.NewRedis(cfg.Redis.Addr)

	hub := ws.NewHub()
	go hub.Run()

	quizRepository := pgRepo.NewQuizRepository(db)
	scoreRepository := pgRepo.NewScoreRepository(db)
	leaderboardRepository := redisRepo.NewLeaderboardRepository(redis)
	sessionRepository := pgRepo.NewSessionRepository(db)
	participantRepository := pgRepo.NewParticipantRepository(db)

	broadcaster := notification.NewWebsocketBroadcaster(hub)

	quizService := appQuiz.NewService(
		quizRepository,
		sessionRepository,
		participantRepository,
		scoreRepository,
		leaderboardRepository,
		broadcaster,
	)

	handler := &http.QuizHandler{
		Service:           quizService,
		WebsocketHub:      hub,
		BroadcastOnSubmit: true,
	}

	router := http.SetupRouter(handler, hub)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
