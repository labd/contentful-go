package cma

import (
	"context"

	"github.com/labd/contentful-go/pkgs/model"
)

type Locales interface {
	Get(ctx context.Context, localeId string) (*model.Locale, error)

	List(ctx context.Context) NextableCollection[*model.Locale, any]

	Upsert(ctx context.Context, locale *model.Locale) error

	Delete(ctx context.Context, locale *model.Locale) error
}
