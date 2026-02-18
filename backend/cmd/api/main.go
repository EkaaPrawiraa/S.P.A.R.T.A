package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"S.P.A.R.T.A/backend/configs"
	"S.P.A.R.T.A/backend/internal/ai/orchestrator"
	"S.P.A.R.T.A/backend/internal/client"
	httpHandler "S.P.A.R.T.A/backend/internal/delivery/http/handler"
	"S.P.A.R.T.A/backend/internal/delivery/http/route"

	// repositories
	postgresRepo "S.P.A.R.T.A/backend/internal/repository/postgres"
	redisRepo "S.P.A.R.T.A/backend/internal/repository/redis"

	// usecases
	ucImpl "S.P.A.R.T.A/backend/internal/usecase"

	"S.P.A.R.T.A/backend/pkg/cache"
	"S.P.A.R.T.A/backend/pkg/database"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// Best-effort load .env for local development.
	_ = godotenv.Load()

	// =========================
	// Load Config
	// =========================
	cfg := configs.LoadConfig()

	// =========================
	// Database Connection
	// =========================
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed connect database:", err)
	}

	// =========================
	// External Clients
	// =========================
	openaiClient := client.NewOpenAIClient(cfg.OpenAIKey, cfg.OpenAIModel, cfg.OpenAIBase)

	redisClient, err := cache.NewRedisClient(cache.RedisConfig{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		log.Fatal("failed connect redis:", err)
	}

	// =========================
	// Repositories
	// =========================
	exerciseRepo := postgresRepo.NewExerciseRepository(db)
	workoutRepo := postgresRepo.NewWorkoutRepository(db)
	splitRepo := postgresRepo.NewSplitRepository(db)
	nutritionRepo := postgresRepo.NewNutritionRepository(db)
	plannerRepo := postgresRepo.NewPlannerRepository(db)
	userRepo := postgresRepo.NewUserRepository(db)
	adminInviteRepo := postgresRepo.NewAdminInviteRepository(db)
	motivationRepo := redisRepo.NewMotivationRepository(redisClient)
	exerciseCacheRepo := redisRepo.NewExerciseCacheRepository(redisClient)

	// =========================
	// Usecases
	// =========================
	workoutUC := ucImpl.NewWorkoutUsecase(workoutRepo)
	splitUC := ucImpl.NewSplitUsecase(splitRepo)
	nutritionUC := ucImpl.NewNutritionUsecase(nutritionRepo)
	exerciseUC := ucImpl.NewExerciseUsecase(exerciseRepo, exerciseCacheRepo)

	aiOrchestrator := orchestrator.NewOrchestrator(openaiClient)
	aiCoachUC := ucImpl.NewAICoachUsecase(aiOrchestrator, splitRepo, exerciseRepo, plannerRepo, workoutRepo, nutritionRepo, motivationRepo)
	plannerUC := ucImpl.NewPlannerUsecase(plannerRepo, aiCoachUC)
	authUC := ucImpl.NewAuthUsecase(userRepo, adminInviteRepo, cfg.JWTSecret)
	adminUC := ucImpl.NewAdminUsecase(adminInviteRepo)

	// =========================
	// Handlers
	// =========================
	workoutHandler := httpHandler.NewWorkoutHandler(workoutUC)
	splitHandler := httpHandler.NewSplitHandler(splitUC)
	nutritionHandler := httpHandler.NewNutritionHandler(nutritionUC)
	plannerHandler := httpHandler.NewPlannerHandler(plannerUC)
	exerciseHandler := httpHandler.NewExerciseHandler(exerciseUC)
	aiCoachHandler := httpHandler.NewAICoachHandler(aiCoachUC)
	authHandler := httpHandler.NewAuthHandler(authUC)
	adminHandler := httpHandler.NewAdminHandler(adminUC)

	// =========================
	// Router
	// =========================
	var router *gin.Engine = route.SetupRouter(
		workoutHandler,
		splitHandler,
		nutritionHandler,
		plannerHandler,
		exerciseHandler,
		aiCoachHandler,
		authHandler,
		adminHandler,
		cfg.JWTSecret,
	)

	// =========================
	// Start Server
	// =========================
	addr := ":" + cfg.Port
	log.Println("server running on", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
