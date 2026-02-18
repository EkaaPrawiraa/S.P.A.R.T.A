package workout

type WorkoutExplanation struct {
	Summary       string
	ExerciseNotes []WorkoutExerciseNote
}

type WorkoutExerciseNote struct {
	Name string
	Note string
}
