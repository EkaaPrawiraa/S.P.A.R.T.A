package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
)

type SplitUsecase interface {
	CreateTemplate(ctx context.Context, tpl *split.SplitTemplate) error
	UpdateTemplate(ctx context.Context, tpl *split.SplitTemplate) error
	ActivateTemplate(ctx context.Context, userID string, templateID string) error
	GetTemplate(ctx context.Context, id string) (*split.SplitTemplate, error)
	GetUserTemplates(ctx context.Context, userID string) ([]split.SplitTemplate, error)
}
