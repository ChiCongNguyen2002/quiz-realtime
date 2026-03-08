package main

import (
	"log"
	"quiz-realtime/api/handler/http"

	appConfig "quiz-realtime/configs"
	notification "quiz-realtime/internal/infrastructure/notification"
	pgRepo "quiz-realtime/internal/infrastructure/repository/postgres"
	redisRepo "quiz-realtime/internal/infrastructure/repository/redis"
	ws "quiz-realtime/internal/infrastructure/websocket"
	appQuiz "quiz-realtime/internal/service/quiz"
	"quiz-realtime/pkg/database"
	redisClient "quiz-realtime/pkg/redis"
)

func main() {
	cfg, err := appConfig.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	masterCfg := database.Config{
		Host:            cfg.Database.Master.Host,
		Port:            cfg.Database.Master.Port,
		User:            cfg.Database.Master.User,
		Password:        cfg.Database.Master.Password,
		Name:            cfg.Database.Master.Name,
		SSLMode:         cfg.Database.Master.SSLMode,
		MaxOpenConns:    cfg.Database.Master.MaxOpenConns,
		MaxIdleConns:    cfg.Database.Master.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.Master.ConnMaxLifetime,
	}

	replicaCfgs := make([]database.Config, 0, len(cfg.Database.Replicas))
	for _, r := range cfg.Database.Replicas {
		replicaCfgs = append(replicaCfgs, database.Config{
			Host:            r.Host,
			Port:            r.Port,
			User:            r.User,
			Password:        r.Password,
			Name:            r.Name,
			SSLMode:         r.SSLMode,
			MaxOpenConns:    r.MaxOpenConns,
			MaxIdleConns:    r.MaxIdleConns,
			ConnMaxLifetime: r.ConnMaxLifetime,
		})
	}

	dbGroup, err := database.NewDBGroup(masterCfg, replicaCfgs)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer dbGroup.Close()

	// Redis with connection pool
	redis := redisClient.NewRedis(redisClient.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Repositories
	quizRepository := pgRepo.NewQuizRepository(dbGroup.MasterDB())
	scoreRepository := pgRepo.NewScoreRepository(dbGroup.MasterDB())
	leaderboardRepository := redisRepo.NewLeaderboardRepository(redis)
	sessionRepository := pgRepo.NewSessionRepository(dbGroup.MasterDB())
	participantRepository := pgRepo.NewParticipantRepository(dbGroup.MasterDB())

	// Services
	broadcaster := notification.NewWebsocketBroadcaster(hub)

	quizService := appQuiz.NewService(
		quizRepository,
		sessionRepository,
		participantRepository,
		scoreRepository,
		leaderboardRepository,
		broadcaster,
	)

	// Handlers
	handler := &http.QuizHandler{
		Service: quizService,
	}

	router := http.SetupRouter(handler, hub)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
