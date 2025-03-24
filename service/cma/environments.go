package cma

import (
	"context"

	"github.com/labd/contentful-go/pkgs/model"
)

type Environments interface {
	Get(ctx context.Context, environmentId string) (*model.Environment, error)

	List(ctx context.Context) NextableCollection[*model.Environment, any]

	Upsert(ctx context.Context, env *model.Environment, sourceEnv *string) error

	Delete(ctx context.Context, env *model.Environment) error
}
