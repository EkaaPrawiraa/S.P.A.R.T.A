package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"S.P.A.R.T.A/backend/configs"
	"S.P.A.R.T.A/backend/internal/client"
	httpHandler "S.P.A.R.T.A/backend/internal/delivery/http/handler"
	"S.P.A.R.T.A/backend/internal/delivery/http/route"

	// repositories
	postgresRepo "S.P.A.R.T.A/backend/internal/repository/postgres"

	// usecases
	ucImpl "S.P.A.R.T.A/backend/internal/usecase"

	"S.P.A.R.T.A/backend/pkg/database"
)

func main() {

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
	openaiClient := client.NewOpenAIClient(cfg.OpenAIKey)

	// =========================
	// Repositories
	// =========================
	userRepo := postgresRepo.NewUserRepository(db)
	exerciseRepo := postgresRepo.NewExerciseRepository(db)
	workoutRepo := postgresRepo.NewWorkoutRepository(db)
	splitRepo := postgresRepo.NewSplitRepository(db)
	nutritionRepo := postgresRepo.NewNutritionRepository(db)
	plannerRepo := postgresRepo.NewPlannerRepository(db)

	// =========================
	// Usecases
	// =========================
	workoutUC := ucImpl.NewWorkoutUsecase(workoutRepo)
	splitUC := ucImpl.NewSplitUsecase(splitRepo)
	nutritionUC := ucImpl.NewNutritionUsecase(nutritionRepo)
	plannerUC := ucImpl.NewPlannerUsecase(plannerRepo, openaiClient)
	exerciseUC := ucImpl.NewExerciseUsecase(exerciseRepo)

	// =========================
	// Handlers
	// =========================
	workoutHandler := httpHandler.NewWorkoutHandler(workoutUC)
	splitHandler := httpHandler.NewSplitHandler(splitUC)
	nutritionHandler := httpHandler.NewNutritionHandler(nutritionUC)
	plannerHandler := httpHandler.NewPlannerHandler(plannerUC)
	exerciseHandler := httpHandler.NewExerciseHandler(exerciseUC)

	// =========================
	// Router
	// =========================
	var router *gin.Engine = route.SetupRouter(
		workoutHandler,
		splitHandler,
		nutritionHandler,
		plannerHandler,
		exerciseHandler,
	)

	// =========================
	// Start Server
	// =========================
	log.Println("server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
