package repository

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
)

type SplitRepository interface {
    CreateTemplate(ctx context.Context, tpl *split.SplitTemplate) error
    UpdateTemplate(ctx context.Context, tpl *split.SplitTemplate) error
    GetTemplateByID(ctx context.Context, id string) (*split.SplitTemplate, error)
    GetUserTemplates(ctx context.Context, userID string) ([]split.SplitTemplate, error)
}
