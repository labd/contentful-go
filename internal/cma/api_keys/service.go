package api_keys

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

var _ cma.ApiKeys = &apiKeysService{}

type apiKeysService struct {
	client   common.RestClient
	basePath string
}

func (a *apiKeysService) Get(ctx context.Context, apiKeyID string) (result *model.APIKey, err error) {

	res, err := a.client.Get(ctx, fmt.Sprintf("%s/%s", a.basePath, apiKeyID), nil, nil)

	if err != nil {
		return nil, err
	}
	var apiKey model.APIKey

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&apiKey)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (a *apiKeysService) List(ctx context.Context) cma.NextableCollection[*model.APIKey, any] {
	return cma2.NewCollection[*model.APIKey, any](&cma2.CollectionOptions{
		Path:   a.basePath,
		Client: a.client,
		Ctx:    ctx,
	})
}

// Upsert updates or creates a new api key entity
func (a *apiKeysService) Upsert(ctx context.Context, apiKey *model.APIKey) error {
	bytesArray, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(apiKey.GetVersion()))

	var res *http.Response

	if apiKey.IsNew() {
		res, err = a.client.Post(ctx, a.basePath, nil, headers, bytes.NewReader(bytesArray))
	} else {
		res, err = a.client.Put(ctx, fmt.Sprintf("%s/%s", a.basePath, apiKey.Sys.ID), nil, headers, bytes.NewReader(bytesArray))
	}

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&apiKey)
}

func (a *apiKeysService) Delete(ctx context.Context, apiKey *model.APIKey) error {

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(apiKey.GetVersion()))

	_, err := a.client.Delete(ctx, fmt.Sprintf("%s/%s", a.basePath, apiKey.Sys.ID), nil, headers)

	return err
}

func NewApiKeysService(client common.RestClient) cma.ApiKeys {
	return &apiKeysService{
		client:   client,
		basePath: "/api_keys",
	}
}
