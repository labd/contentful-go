package app_definitions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cma2 "github.com/labd/contentful-go/internal/cma/common"
	"github.com/labd/contentful-go/pkgs/model"
	"github.com/labd/contentful-go/service/cma"
	"github.com/labd/contentful-go/service/common"
)

var _ cma.AppDefinitions = &appDefinitionService{}

type appDefinitionService struct {
	client   common.RestClient
	basePath string
}

func (a appDefinitionService) Get(ctx context.Context, appDefinitionId string) (*model.AppDefinition, error) {
	res, err := a.client.Get(ctx, fmt.Sprintf("%s/%s", a.basePath, appDefinitionId), nil, nil)

	if err != nil {
		return nil, err
	}
	var appInstallation model.AppDefinition

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&appInstallation)
	if err != nil {
		return nil, err
	}
	return &appInstallation, nil
}

func (a appDefinitionService) List(ctx context.Context) cma.NextableCollection[*model.AppDefinition, any] {
	return cma2.NewCollection[*model.AppDefinition, any](&cma2.CollectionOptions{
		Path:   a.basePath,
		Client: a.client,
		Ctx:    ctx,
	})
}

func (a appDefinitionService) Upsert(ctx context.Context, appDefinition *model.AppDefinition) error {
	bytesArray, err := json.Marshal(appDefinition)
	if err != nil {
		return err
	}

	var res *http.Response

	if appDefinition.IsNew() {
		res, err = a.client.Post(ctx, a.basePath, nil, nil, bytes.NewReader(bytesArray))
	} else {
		res, err = a.client.Put(ctx, fmt.Sprintf("%s/%s", a.basePath, appDefinition.Sys.ID), nil, nil, bytes.NewReader(bytesArray))
	}
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&appDefinition)
}

func (a appDefinitionService) Delete(ctx context.Context, appDefinition *model.AppDefinition) error {
	headers := make(http.Header)

	_, err := a.client.Delete(ctx, fmt.Sprintf("%s/%s", a.basePath, appDefinition.Sys.ID), nil, headers)

	return err
}

func NewAppDefinitionService(client common.RestClient) cma.AppDefinitions {
	return &appDefinitionService{
		client:   client,
		basePath: "/app_definitions",
	}
}
