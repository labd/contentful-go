package cma

import (
	"context"
	"github.com/flaconi/contentful-go/pkgs/model"
)

type Entries interface {
	Get(ctx context.Context, entryId string) (*model.Entry, error)

	List(ctx context.Context) NextableCollection[*model.Entry, any]

	Upsert(ctx context.Context, contentTypeID string, entry *model.Entry) error

	Delete(ctx context.Context, entry *model.Entry) error

	Publish(ctx context.Context, entry *model.Entry) error

	Unpublish(ctx context.Context, entry *model.Entry) error

	Archive(ctx context.Context, entry *model.Entry) error

	Unarchive(ctx context.Context, entry *model.Entry) error

	ListPublished(ctx context.Context) NextableCollection[*model.Entry, any]
}
