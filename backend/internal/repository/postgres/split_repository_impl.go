package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type splitRepository struct {
	db *sql.DB
}

func NewSplitRepository(db *sql.DB) domainrepo.SplitRepository {
	return &splitRepository{db: db}
}

func (r *splitRepository) CreateTemplate(ctx context.Context, tpl *split.SplitTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO split_templates(id,user_id,name,description,created_by,focus_muscle,is_active,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		tpl.ID, tpl.UserID, tpl.Name, tpl.Description, tpl.CreatedBy, tpl.FocusMuscle, tpl.IsActive, tpl.CreatedAt)
	return err
}

func (r *splitRepository) UpdateTemplate(ctx context.Context, tpl *split.SplitTemplate) error {
	return nil
}

func (r *splitRepository) GetTemplateByID(ctx context.Context, id string) (*split.SplitTemplate, error) {
	return nil, nil
}

func (r *splitRepository) GetUserTemplates(ctx context.Context, userID string) ([]split.SplitTemplate, error) {
	return nil, nil
}
