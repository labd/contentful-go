package app_installations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cma2 "github.com/flaconi/contentful-go/internal/cma/common"
	"github.com/flaconi/contentful-go/pkgs/model"
	"github.com/flaconi/contentful-go/service/cma"
	"github.com/flaconi/contentful-go/service/common"
	"net/http"
	"strings"
)

var _ cma.AppInstallations = &appInstallationsService{}

type appInstallationsService struct {
	client   common.RestClient
	basePath string
}

func (a appInstallationsService) Get(ctx context.Context, appDefinitionId string) (*model.AppInstallation, error) {
	res, err := a.client.Get(ctx, fmt.Sprintf("%s/%s", a.basePath, appDefinitionId), nil, nil)

	if err != nil {
		return nil, err
	}
	var appInstallation model.AppInstallation

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&appInstallation)
	if err != nil {
		return nil, err
	}
	return &appInstallation, nil
}

func (a appInstallationsService) List(ctx context.Context) cma.NextableCollection[*model.AppInstallation, any] {
	return cma2.NewCollection[*model.AppInstallation, any](&cma2.CollectionOptions{
		Path:   a.basePath,
		Client: a.client,
		Ctx:    ctx,
	})
}

func (a appInstallationsService) Upsert(ctx context.Context, appInstallation *model.AppInstallation) error {
	bytesArray, err := json.Marshal(appInstallation)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	if len(appInstallation.Terms) > 0 {
		headers.Set("X-Contentful-Marketplace", strings.Join(appInstallation.Terms, ","))
	}

	res, err := a.client.Put(ctx, fmt.Sprintf("%s/%s", a.basePath, appInstallation.Sys.AppDefinition.Sys.ID), nil, headers, bytes.NewReader(bytesArray))

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&appInstallation)
}

func (a appInstallationsService) Delete(ctx context.Context, appInstallation *model.AppInstallation) error {
	headers := make(http.Header)

	_, err := a.client.Delete(ctx, fmt.Sprintf("%s/%s", a.basePath, appInstallation.Sys.AppDefinition.Sys.ID), nil, headers)

	return err
}

func NewAppInstallationsService(client common.RestClient) cma.AppInstallations {
	return &appInstallationsService{
		client:   client,
		basePath: "/app_installations",
	}
}
