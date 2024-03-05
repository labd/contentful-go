package cma

import (
	"context"
	"github.com/flaconi/contentful-go/pkgs/model"
)

type ContentTypes interface {
	Get(ctx context.Context, typeId string) (*model.ContentType, error)

	List(ctx context.Context) NextableCollection[*model.ContentType, any]

	Upsert(ctx context.Context, contentType *model.ContentType) error

	Delete(ctx context.Context, contentType *model.ContentType) error

	Activate(ctx context.Context, contentType *model.ContentType) error

	Deactivate(ctx context.Context, contentType *model.ContentType) error

	ListActivated(ctx context.Context) NextableCollection[*model.ContentType, any]
}
