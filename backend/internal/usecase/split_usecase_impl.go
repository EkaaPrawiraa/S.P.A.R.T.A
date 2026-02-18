package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
)

type splitUsecase struct {
	repo domainrepo.SplitRepository
}

func NewSplitUsecase(repo domainrepo.SplitRepository) domainuc.SplitUsecase {
	return &splitUsecase{repo: repo}
}

func (u *splitUsecase) CreateTemplate(ctx context.Context, tpl *split.SplitTemplate) error {
	return u.repo.CreateTemplate(ctx, tpl)
}

func (u *splitUsecase) GetTemplate(ctx context.Context, id string) (*split.SplitTemplate, error) {
	return u.repo.GetTemplateByID(ctx, id)
}

func (u *splitUsecase) GetUserTemplates(ctx context.Context, userID string) ([]split.SplitTemplate, error) {
	return u.repo.GetUserTemplates(ctx, userID)
}
