package environment_aliases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	cma2 "github.com/labd/contentful-go/internal/cma/common"
	"github.com/labd/contentful-go/pkgs/model"
	"github.com/labd/contentful-go/service/cma"
	"github.com/labd/contentful-go/service/common"
)

var _ cma.EnvironmentAliases = &environmentAliasService{}

type environmentAliasService struct {
	client   common.RestClient
	basePath string
}

func (e environmentAliasService) Get(ctx context.Context, environmentId string) (*model.EnvironmentAlias, error) {
	res, err := e.client.Get(ctx, fmt.Sprintf("%s/%s", e.basePath, environmentId), nil, nil)

	if err != nil {
		return nil, err
	}
	var entry model.EnvironmentAlias

	if decodeError := entry.Decode(res.Body); decodeError != nil {
		return nil, decodeError
	}
	return &entry, nil
}

func (e environmentAliasService) List(ctx context.Context) cma.NextableCollection[*model.EnvironmentAlias, any] {
	return cma2.NewCollection[*model.EnvironmentAlias, any](&cma2.CollectionOptions{
		Path:   e.basePath,
		Client: e.client,
		Ctx:    ctx,
	})
}

func (e environmentAliasService) Upsert(ctx context.Context, env *model.EnvironmentAlias) error {
	bytesArray, err := json.Marshal(env)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(env.GetVersion()))

	var res *http.Response

	path := fmt.Sprintf("%s/%s", e.basePath, env.Sys.ID)

	res, err = e.client.Put(ctx, path, nil, headers, bytes.NewReader(bytesArray))

	if err != nil {
		return err
	}

	return env.Decode(res.Body)
}

func (e environmentAliasService) Delete(ctx context.Context, env *model.EnvironmentAlias) error {
	headers := make(http.Header)

	_, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s", e.basePath, env.Sys.ID), nil, headers)

	return err
}

func NewEnvironmentAliasService(client common.RestClient) cma.EnvironmentAliases {
	return &environmentAliasService{
		client:   client,
		basePath: "/environment_aliases",
	}
}
