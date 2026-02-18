package dto

type SplitTemplate struct {
	Name string
	Days []SplitDay
}

type SplitDay struct {
	DayName   string
	Focus     []string
	Exercises []SplitExercise
}

type SplitExercise struct {
	Name     string
	Sets     int
	RepRange string
	Priority string
}
