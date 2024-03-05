package environments

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
	"strconv"
)

var _ cma.Environments = &environmentService{}

type environmentService struct {
	client   common.RestClient
	basePath string
}

func (e environmentService) Get(ctx context.Context, environmentId string) (*model.Environment, error) {
	res, err := e.client.Get(ctx, fmt.Sprintf("%s/%s", e.basePath, environmentId), nil, nil)

	if err != nil {
		return nil, err
	}
	var entry model.Environment

	err = entry.Decode(res.Body)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (e environmentService) List(ctx context.Context) cma.NextableCollection[*model.Environment, any] {
	return cma2.NewCollection[*model.Environment, any](&cma2.CollectionOptions{
		Path:   e.basePath,
		Client: e.client,
		Ctx:    ctx,
	})
}

func (e environmentService) Upsert(ctx context.Context, env *model.Environment, sourceEnv *string) error {
	bytesArray, err := json.Marshal(env)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(env.GetVersion()))

	if sourceEnv != nil {
		headers.Set("X-Contentful-Source-Environment", *sourceEnv)
	}

	var res *http.Response
	var path string

	if env.IsNew() {
		path = fmt.Sprintf("%s/%s", e.basePath, env.Name)
	} else {
		path = fmt.Sprintf("%s/%s", e.basePath, env.Sys.ID)
	}

	res, err = e.client.Put(ctx, path, nil, headers, bytes.NewReader(bytesArray))

	if err != nil {
		return err
	}

	return env.Decode(res.Body)
}

func (e environmentService) Delete(ctx context.Context, env *model.Environment) error {
	headers := make(http.Header)

	_, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s", e.basePath, env.Sys.ID), nil, headers)

	return err
}

func NewEnvironmentService(client common.RestClient) cma.Environments {
	return &environmentService{
		client:   client,
		basePath: "/environments",
	}
}
