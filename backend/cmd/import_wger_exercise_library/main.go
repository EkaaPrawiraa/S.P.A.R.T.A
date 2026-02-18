package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"S.P.A.R.T.A/backend/configs"
	"S.P.A.R.T.A/backend/internal/client"
	"S.P.A.R.T.A/backend/pkg/database"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	_ = godotenv.Load()

	cfg := configs.LoadConfig()

	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed connect database:", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	wger := client.NewWgerClient()

	pageSize := 100
	offset := 0
	languageCode := "en"

	wgerUUIDToExerciseID := map[string]string{}
	insertedExercises := 0
	insertedImages := 0
	insertedVideos := 0

	for {
		page, err := wger.ListExerciseInfo(ctx, pageSize, offset, languageCode)
		if err != nil {
			log.Fatal("failed fetch wger exerciseinfo:", err)
		}

		for _, ex := range page.Results {
			name := pickWgerName(ex)
			if strings.TrimSpace(name) == "" {
				continue
			}

			existingID, err := findExerciseIDByName(ctx, db, name)
			if err != nil {
				log.Fatal("failed lookup exercise by name:", err)
			}

			exerciseID := existingID
			if exerciseID == "" {
				exerciseID = deterministicUUID("wger:exercise:" + ex.UUID)
				createdAt := time.Now().UTC()
				primaryMuscle := strings.ToLower(strings.TrimSpace(ex.Category.Name))
				secondary := wgerSecondaryMuscles(ex)
				equipment := wgerEquipment(ex)

				ok, err := insertExerciseIfNotExists(ctx, db, exerciseID, name, primaryMuscle, secondary, equipment, createdAt)
				if err != nil {
					log.Fatal("failed insert exercise:", err)
				}
				if ok {
					insertedExercises++
				}
			}

			wgerUUIDToExerciseID[ex.UUID] = exerciseID

			img, ok := pickMainImage(ex.Images)
			if ok && strings.TrimSpace(img.Image) != "" {
				createdAt := time.Now().UTC()
				mediaID := deterministicUUID("wger:exercise_media:image:" + img.UUID)
				inserted, err := insertExerciseMediaIfNotExists(ctx, db, mediaID, exerciseID, "image", img.Image, nil, createdAt)
				if err != nil {
					log.Fatal("failed insert exercise image media:", err)
				}
				if inserted {
					insertedImages++
				}
			}
		}

		if page.Next == nil || strings.TrimSpace(*page.Next) == "" {
			break
		}
		offset += pageSize
	}

	videoOffset := 0
	videoPageSize := 50
	for {
		page, err := wger.ListVideos(ctx, videoPageSize, videoOffset)
		if err != nil {
			log.Fatal("failed fetch wger videos:", err)
		}

		for _, v := range page.Results {
			if strings.TrimSpace(v.Video) == "" || strings.TrimSpace(v.ExerciseUUID) == "" {
				continue
			}
			exerciseID := wgerUUIDToExerciseID[v.ExerciseUUID]
			if exerciseID == "" {
				continue
			}

			createdAt := time.Now().UTC()
			mediaID := deterministicUUID("wger:exercise_media:video:" + v.UUID)
			inserted, err := insertExerciseMediaIfNotExists(ctx, db, mediaID, exerciseID, "video", v.Video, nil, createdAt)
			if err != nil {
				log.Fatal("failed insert exercise video media:", err)
			}
			if inserted {
				insertedVideos++
			}
		}

		if page.Next == nil || strings.TrimSpace(*page.Next) == "" {
			break
		}
		videoOffset += videoPageSize
	}

	log.Printf(
		"import_wger_exercise_library done: exercises_inserted=%d images_inserted=%d videos_inserted=%d",
		insertedExercises,
		insertedImages,
		insertedVideos,
	)
}

func deterministicUUID(input string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(input)).String()
}

func pickWgerName(ex client.WgerExerciseInfo) string {
	for _, tr := range ex.Translations {
		if strings.TrimSpace(tr.Name) != "" {
			return strings.TrimSpace(tr.Name)
		}
	}
	return ""
}

func wgerEquipment(ex client.WgerExerciseInfo) string {
	if len(ex.Equipment) == 0 {
		return ""
	}
	names := make([]string, 0, len(ex.Equipment))
	for _, e := range ex.Equipment {
		if strings.TrimSpace(e.Name) != "" {
			names = append(names, strings.ToLower(strings.TrimSpace(e.Name)))
		}
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

func wgerSecondaryMuscles(ex client.WgerExerciseInfo) []string {
	set := map[string]struct{}{}
	for _, m := range ex.Muscles {
		name := strings.ToLower(strings.TrimSpace(m.Name))
		if name != "" {
			set[name] = struct{}{}
		}
	}
	for _, m := range ex.MusclesSecondary {
		name := strings.ToLower(strings.TrimSpace(m.Name))
		if name != "" {
			set[name] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func pickMainImage(images []client.WgerExerciseImage) (client.WgerExerciseImage, bool) {
	if len(images) == 0 {
		return client.WgerExerciseImage{}, false
	}
	for _, img := range images {
		if img.IsMain {
			return img, true
		}
	}
	return images[0], true
}

func findExerciseIDByName(ctx context.Context, db *sql.DB, name string) (string, error) {
	row := db.QueryRowContext(ctx, `SELECT id FROM exercises WHERE lower(name) = lower($1) LIMIT 1`, name)
	var id string
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func insertExerciseIfNotExists(ctx context.Context, db *sql.DB, id, name, primaryMuscle string, secondaryMuscles []string, equipment string, createdAt time.Time) (bool, error) {
	res, err := db.ExecContext(ctx,
		`INSERT INTO exercises(id,name,primary_muscle,secondary_muscles,equipment,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (id) DO NOTHING`,
		id,
		name,
		primaryMuscle,
		secondaryMuscles,
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
