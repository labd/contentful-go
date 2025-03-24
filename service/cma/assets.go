package cma

import (
	"context"

	"github.com/labd/contentful-go/pkgs/model"
)

type Assets interface {
	Get(ctx context.Context, assetId string) (*model.Asset, error)

	List(ctx context.Context) NextableCollection[*model.Asset, any]

	Upsert(ctx context.Context, asset *model.Asset) error

	Process(ctx context.Context, asset *model.Asset) error

	Delete(ctx context.Context, asset *model.Asset) error

	Publish(ctx context.Context, asset *model.Asset) error

	Unpublish(ctx context.Context, asset *model.Asset) error

	Archive(ctx context.Context, asset *model.Asset) error

	Unarchive(ctx context.Context, asset *model.Asset) error

	ListPublished(ctx context.Context) NextableCollection[*model.Asset, any]
}
