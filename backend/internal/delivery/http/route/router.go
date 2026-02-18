package route

import (
	"S.P.A.R.T.A/backend/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	workoutHandler *handler.WorkoutHandler,
	splitHandler *handler.SplitHandler,
	nutritionHandler *handler.NutritionHandler,
	plannerHandler *handler.PlannerHandler,
	exerciseHandler *handler.ExerciseHandler,
) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api")

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

	return r
}
