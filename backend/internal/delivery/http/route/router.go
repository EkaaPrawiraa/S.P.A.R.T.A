package route

import (
	"S.P.A.R.T.A/backend/internal/delivery/http/handler"
	"S.P.A.R.T.A/backend/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	workoutHandler *handler.WorkoutHandler,
	splitHandler *handler.SplitHandler,
	nutritionHandler *handler.NutritionHandler,
	plannerHandler *handler.PlannerHandler,
	exerciseHandler *handler.ExerciseHandler,
	AICoachHandler  *handler.AICoachHandler,
) *gin.Engine {

	r := gin.Default()
	
	authMW := middleware.NewAuthMiddleware("SUPER_SECRET")

	api := r.Group("/api/v1")

	// public routes
    api.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

	secured := api.Group("/")
    secured.Use(authMW.RequireAuth())

	// workouts
	workouts := api.Group("/workouts")
	{
		workouts.POST("", workoutHandler.CreateWorkoutSession)
		workouts.GET("/:id", workoutHandler.GetWorkoutSession)
		workouts.GET("/user/:user_id", workoutHandler.GetUserWorkoutSessions)
	}

	// splits
	splits := api.Group("/splits")
	{
		splits.POST("", splitHandler.CreateTemplate)
		splits.GET("/user/:user_id", splitHandler.GetUserTemplates)
	}

	// nutrition
	nutrition := api.Group("/nutrition")
	{
		nutrition.POST("", nutritionHandler.UpsertDailyNutrition)
		nutrition.GET("/user/:user_id", nutritionHandler.GetDailyNutrition)
	}

	// planner
	planner := api.Group("/planner")
	{
		planner.POST("/generate/:user_id", plannerHandler.GenerateRecommendation)
		planner.GET("/user/:user_id", plannerHandler.GetUserRecommendations)
	}

	// exercises
	exercises := api.Group("/exercises")
	{
		exercises.GET("", exerciseHandler.ListExercises)
		exercises.GET("/:id", exerciseHandler.GetExercise)
	}

	// ai
	ai := api.Group("/ai")
	{
		ai.POST("/ai/generate-split", AICoachHandler.GenerateSplit)
		ai.POST("/ai/overload", AICoachHandler.SuggestOverload)
	}

	return r
}
