package route

import (
	"time"

	"S.P.A.R.T.A/backend/internal/delivery/http/handler"
	"S.P.A.R.T.A/backend/internal/delivery/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	workoutHandler *handler.WorkoutHandler,
	splitHandler *handler.SplitHandler,
	nutritionHandler *handler.NutritionHandler,
	plannerHandler *handler.PlannerHandler,
	exerciseHandler *handler.ExerciseHandler,
	AICoachHandler *handler.AICoachHandler,
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	jwtSecret string,
) *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())

	// Allow browser-based local dev (frontend on :3000 calling backend on :8080).
	corsCfg := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type", "X-Request-Id"},
		ExposeHeaders: []string{
			"X-Request-Id",
		},
		MaxAge: 12 * time.Hour,
	}
	r.Use(cors.New(corsCfg))

	// Structured access logs + correlation.
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.MetricsMiddleware())

	authMW := middleware.NewAuthMiddleware(jwtSecret)
	adminMW := middleware.NewAdminMiddleware()

	api := r.Group("/api/v1")

	// public routes
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api.GET("/metrics", middleware.MetricsHandler())

	// auth (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	secured := api.Group("/")
	secured.Use(authMW.RequireAuth())

	// admin (secured + admin-only)
	admin := secured.Group("/admin")
	admin.Use(adminMW.RequireAdmin())
	{
		admin.POST("/invites", adminHandler.CreateInvite)
	}

	// workouts
	workouts := secured.Group("/workouts")
	{
		workouts.POST("", workoutHandler.CreateWorkoutSession)
		workouts.GET("/:id", workoutHandler.GetWorkoutSession)
		workouts.GET("/user/:user_id", workoutHandler.GetUserWorkoutSessions)
	}

	// splits
	splits := secured.Group("/splits")
	{
		splits.POST("", splitHandler.CreateTemplate)
		splits.GET("/:id", splitHandler.GetTemplate)
		splits.PUT("/:id", splitHandler.UpdateTemplate)
		splits.POST("/:id/activate", splitHandler.ActivateTemplate)
		splits.GET("/user/:user_id", splitHandler.GetUserTemplates)
	}

	// nutrition
	nutrition := secured.Group("/nutrition")
	{
		nutrition.POST("", nutritionHandler.UpsertDailyNutrition)
		nutrition.GET("/user/:user_id", nutritionHandler.GetDailyNutrition)
	}

	// planner
	planner := secured.Group("/planner")
	{
		planner.POST("/generate/:user_id", plannerHandler.GenerateRecommendation)
		planner.GET("/user/:user_id", plannerHandler.GetUserRecommendations)
	}

	// exercises
	exercises := secured.Group("/exercises")
	{
		exercises.POST("", exerciseHandler.CreateExercise)
		exercises.GET("", exerciseHandler.ListExercises)
		exercises.GET("/:id", exerciseHandler.GetExercise)
		exercises.POST("/:id/media", exerciseHandler.AddExerciseMedia)
	}

	// ai
	ai := secured.Group("/ai")
	{
		ai.POST("/generate-split", AICoachHandler.GenerateSplit)
		ai.POST("/overload", AICoachHandler.SuggestOverload)
		ai.POST("/workout", AICoachHandler.GenerateWorkoutPlan)
		ai.GET("/motivation", AICoachHandler.GetDailyMotivation)
		ai.GET("/coaching", AICoachHandler.GetCoachingSuggestions)
		ai.POST("/explain-workout", AICoachHandler.ExplainWorkoutPlan)
	}

	return r
}
