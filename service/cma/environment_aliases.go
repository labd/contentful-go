package cma

import (
	"context"
	"github.com/flaconi/contentful-go/pkgs/model"
)

type EnvironmentAliases interface {
	Get(ctx context.Context, environmentId string) (*model.EnvironmentAlias, error)

	List(ctx context.Context) NextableCollection[*model.EnvironmentAlias, any]

	Upsert(ctx context.Context, env *model.EnvironmentAlias) error

	Delete(ctx context.Context, env *model.EnvironmentAlias) error
}
