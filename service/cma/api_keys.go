package cma

import (
	"context"

	"github.com/labd/contentful-go/pkgs/model"
)

type ApiKeys interface {
	Get(ctx context.Context, apiKeyID string) (*model.APIKey, error)

	List(ctx context.Context) NextableCollection[*model.APIKey, any]

	Upsert(ctx context.Context, apiKey *model.APIKey) error

	Delete(ctx context.Context, apiKey *model.APIKey) error
}
