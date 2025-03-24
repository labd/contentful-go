package cma

import (
	"context"

	"github.com/labd/contentful-go/pkgs/model"
)

type PreviewApiKeys interface {
	Get(ctx context.Context, apiKeyID string) (*model.PreviewAPIKey, error)

	List(ctx context.Context) NextableCollection[*model.PreviewAPIKey, any]
}
