package repository

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
)

type SplitRepository interface {
	CreateTemplate(ctx context.Context, tpl *split.SplitTemplate) error
	UpdateTemplate(ctx context.Context, tpl *split.SplitTemplate) error
	ActivateTemplate(ctx context.Context, userID string, templateID string) error
	GetTemplateByID(ctx context.Context, id string) (*split.SplitTemplate, error)
	GetUserTemplates(ctx context.Context, userID string) ([]split.SplitTemplate, error)
	GetSplitDayByID(ctx context.Context, id string) (*split.SplitDay, error)
}
