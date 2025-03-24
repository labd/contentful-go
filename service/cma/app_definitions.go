package cma

import (
	"context"

	"github.com/labd/contentful-go/pkgs/model"
)

type AppDefinitions interface {
	Get(ctx context.Context, appDefinitionId string) (*model.AppDefinition, error)

	List(ctx context.Context) NextableCollection[*model.AppDefinition, any]

	Upsert(ctx context.Context, appDefinition *model.AppDefinition) error

	Delete(ctx context.Context, appDefinition *model.AppDefinition) error
}
