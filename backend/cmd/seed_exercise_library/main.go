package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"S.P.A.R.T.A/backend/configs"
	"S.P.A.R.T.A/backend/internal/client"
	"S.P.A.R.T.A/backend/pkg/database"
)

type seedExercise struct {
	Name          string
	PrimaryMuscle string
	Equipment     string
	WikiTitle     string
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	_ = godotenv.Load()

	cfg := configs.LoadConfig()

	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed connect database:", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	wiki := client.NewWikipediaClient()

	seeds := []seedExercise{
		{Name: "Bench Press", PrimaryMuscle: "chest", Equipment: "barbell", WikiTitle: "Bench press"},
		{Name: "Squat", PrimaryMuscle: "legs", Equipment: "barbell", WikiTitle: "Squat (exercise)"},
		{Name: "Deadlift", PrimaryMuscle: "back", Equipment: "barbell", WikiTitle: "Deadlift"},
		{Name: "Overhead Press", PrimaryMuscle: "shoulders", Equipment: "barbell", WikiTitle: "Overhead press"},
		{Name: "Pull-up", PrimaryMuscle: "back", Equipment: "bodyweight", WikiTitle: "Pull-up (exercise)"},
		{Name: "Bent-Over Row", PrimaryMuscle: "back", Equipment: "barbell", WikiTitle: "Bent-over row"},
		{Name: "Biceps Curl", PrimaryMuscle: "arms", Equipment: "dumbbell", WikiTitle: "Biceps curl"},
		{Name: "Triceps Pushdown", PrimaryMuscle: "arms", Equipment: "cable", WikiTitle: "Triceps pushdown"},
		{Name: "Lat Pulldown", PrimaryMuscle: "back", Equipment: "machine", WikiTitle: "Lat pulldown"},
		{Name: "Leg Press", PrimaryMuscle: "legs", Equipment: "machine", WikiTitle: "Leg press"},
		{Name: "Leg Extension", PrimaryMuscle: "legs", Equipment: "machine", WikiTitle: "Leg extension"},
		{Name: "Leg Curl", PrimaryMuscle: "legs", Equipment: "machine", WikiTitle: "Leg curl"},
		{Name: "Calf Raise", PrimaryMuscle: "legs", Equipment: "machine", WikiTitle: "Calf raise"},
		{Name: "Plank", PrimaryMuscle: "core", Equipment: "bodyweight", WikiTitle: "Plank (exercise)"},
		{Name: "Crunch", PrimaryMuscle: "core", Equipment: "bodyweight", WikiTitle: "Crunch (exercise)"},
		{Name: "Romanian Deadlift", PrimaryMuscle: "legs", Equipment: "barbell", WikiTitle: "Romanian deadlift"},
		{Name: "Hip Thrust", PrimaryMuscle: "glutes", Equipment: "barbell", WikiTitle: "Hip thrust"},
		{Name: "Lunge", PrimaryMuscle: "legs", Equipment: "bodyweight", WikiTitle: "Lunge (exercise)"},
		{Name: "Push-up", PrimaryMuscle: "chest", Equipment: "bodyweight", WikiTitle: "Push-up"},
		{Name: "Dumbbell Fly", PrimaryMuscle: "chest", Equipment: "dumbbell", WikiTitle: "Fly (exercise)"},
	}

	insertedExercises := 0
	insertedMedia := 0

	for _, seed := range seeds {
		exID := uuid.NewSHA1(uuid.NameSpaceOID, []byte("exercise:"+seed.Name)).String()
		createdAt := time.Now().UTC()

		ok, err := insertExerciseIfNotExists(ctx, db, exID, seed.Name, seed.PrimaryMuscle, seed.Equipment, createdAt)
		if err != nil {
			log.Fatal("failed insert exercise:", err)
		}
		if ok {
			insertedExercises++
		}

		thumbURL := ""
		if seed.WikiTitle != "" {
			url, err := wiki.GetPageThumbnailURL(ctx, seed.WikiTitle)
			if err == nil {
				thumbURL = url
			}
		}
		if thumbURL == "" {
			continue
		}

		mediaID := uuid.NewSHA1(uuid.NameSpaceOID, []byte("exercise_media:image:"+seed.Name)).String()
		ok, err = insertExerciseMediaIfNotExists(ctx, db, mediaID, exID, "image", thumbURL, nil, createdAt)
		if err != nil {
			log.Fatal("failed insert exercise media:", err)
		}
		if ok {
			insertedMedia++
		}
	}

	log.Printf("seed_exercise_library done: exercises_inserted=%d media_inserted=%d", insertedExercises, insertedMedia)
}

func insertExerciseIfNotExists(ctx context.Context, db *sql.DB, id, name, primaryMuscle, equipment string, createdAt time.Time) (bool, error) {
	res, err := db.ExecContext(ctx,
		`INSERT INTO exercises(id,name,primary_muscle,secondary_muscles,equipment,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (id) DO NOTHING`,
		id,
		name,
		primaryMuscle,
		[]string{},
		equipment,
		createdAt,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func insertExerciseMediaIfNotExists(ctx context.Context, db *sql.DB, id, exerciseID, mediaType, mediaURL string, thumbnailURL *string, createdAt time.Time) (bool, error) {
	res, err := db.ExecContext(ctx,
		`INSERT INTO exercise_media(id,exercise_id,media_type,media_url,thumbnail_url,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (id) DO NOTHING`,
		id,
		exerciseID,
		mediaType,
		mediaURL,
		thumbnailURL,
		createdAt,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
