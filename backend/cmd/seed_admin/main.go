package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"S.P.A.R.T.A/backend/configs"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/user"
	postgresRepo "S.P.A.R.T.A/backend/internal/repository/postgres"
	"S.P.A.R.T.A/backend/pkg/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	cfg := configs.LoadConfig()
	if strings.ToLower(strings.TrimSpace(cfg.AppEnv)) != "local" {
		log.Fatal("seed_admin is only allowed when APP_ENV=local")
	}

	name := strings.TrimSpace(os.Getenv("SEED_ADMIN_NAME"))
	email := strings.TrimSpace(strings.ToLower(os.Getenv("SEED_ADMIN_EMAIL")))
	pass := strings.TrimSpace(os.Getenv("SEED_ADMIN_PASSWORD"))
	if name == "" {
		name = "Local Admin"
	}
	if email == "" || pass == "" {
		log.Fatal("missing SEED_ADMIN_EMAIL or SEED_ADMIN_PASSWORD")
	}

	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgresRepo.NewUserRepository(db)
	ctx := context.Background()
	if existing, err := userRepo.GetByEmail(ctx, email); err == nil && existing != nil {
		fmt.Println("admin user already exists:", existing.ID)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().UTC()
	u := &user.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "admin",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := userRepo.Create(ctx, u); err != nil {
		log.Fatal(err)
	}

	fmt.Println("created admin user:", u.ID)
}
