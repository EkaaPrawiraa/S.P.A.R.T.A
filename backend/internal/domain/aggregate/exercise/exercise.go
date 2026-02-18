package exercise

import "time"

type Exercise struct {
    ID               string
    Name             string
    PrimaryMuscle    string
    SecondaryMuscles []string
    Equipment        string
    CreatedAt        time.Time
    Media            []ExerciseMedia
}

type ExerciseMedia struct {
    ID           string
    MediaType    string
    MediaURL     string
    ThumbnailURL *string
    CreatedAt    time.Time
}
