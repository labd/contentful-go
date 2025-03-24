package preview_api_keys

import (
	"context"
	"encoding/json"
	"fmt"

	cma2 "github.com/labd/contentful-go/internal/cma/common"
	"github.com/labd/contentful-go/pkgs/model"
	"github.com/labd/contentful-go/service/cma"
	"github.com/labd/contentful-go/service/common"
)

var _ cma.PreviewApiKeys = &previewApiKeys{}

type previewApiKeys struct {
	client   common.RestClient
	basePath string
}

func (a *previewApiKeys) Get(ctx context.Context, apiKeyID string) (result *model.PreviewAPIKey, err error) {

	res, err := a.client.Get(ctx, fmt.Sprintf("%s/%s", a.basePath, apiKeyID), nil, nil)

	if err != nil {
		return nil, err
	}
	var apiKey model.PreviewAPIKey

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&apiKey)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (a *previewApiKeys) List(ctx context.Context) cma.NextableCollection[*model.PreviewAPIKey, any] {

	return cma2.NewCollection[*model.PreviewAPIKey, any](&cma2.CollectionOptions{
		Path:   a.basePath,
		Client: a.client,
		Ctx:    ctx,
	})
}

func NewPreviewApiKeysService(client common.RestClient) cma.PreviewApiKeys {
	return &previewApiKeys{
		client:   client,
		basePath: "/preview_api_keys",
	}
}
