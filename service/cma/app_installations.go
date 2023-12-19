package cma

import (
	"context"
	"github.com/flaconi/contentful-go/pkgs/model"
)

type AppInstallations interface {
	Get(ctx context.Context, appDefinitionId string) (*model.AppInstallation, error)

	List(ctx context.Context) NextableCollection[*model.AppInstallation, any]

	Upsert(ctx context.Context, appInstallation *model.AppInstallation) error

	Delete(ctx context.Context, appInstallation *model.AppInstallation) error
}
