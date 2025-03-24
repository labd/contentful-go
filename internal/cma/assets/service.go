package assets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-multierror"
	cma2 "github.com/labd/contentful-go/internal/cma/common"
	"github.com/labd/contentful-go/pkgs/model"
	"github.com/labd/contentful-go/service/cma"
	"github.com/labd/contentful-go/service/common"
)

var _ cma.Assets = &assetService{}

type assetService struct {
	client   common.RestClient
	basePath string
}

func (e assetService) Get(ctx context.Context, assetId string) (*model.Asset, error) {
	res, err := e.client.Get(ctx, fmt.Sprintf("%s/%s", e.basePath, assetId), nil, nil)

	if err != nil {
		return nil, err
	}
	var entry model.Asset

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (e assetService) List(ctx context.Context) cma.NextableCollection[*model.Asset, any] {
	return cma2.NewCollection[*model.Asset, any](&cma2.CollectionOptions{
		Path:   e.basePath,
		Client: e.client,
		Ctx:    ctx,
	})
}

func (e assetService) Upsert(ctx context.Context, asset *model.Asset) error {
	bytesArray, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	var res *http.Response

	if asset.IsNew() {
		res, err = e.client.Post(ctx, e.basePath, nil, headers, bytes.NewReader(bytesArray))
	} else {
		res, err = e.client.Put(ctx, fmt.Sprintf("%s/%s", e.basePath, asset.Sys.ID), nil, headers, bytes.NewReader(bytesArray))
	}

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&asset)
}

func (e assetService) Process(ctx context.Context, asset *model.Asset) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	var err error

	for locale := range asset.Fields.File {
		_, resultErr := e.client.Put(ctx, fmt.Sprintf("%s/%s/files/%s/process", e.basePath, asset.Sys.ID, locale), nil, headers, nil)
		if resultErr != nil {
			err = multierror.Append(err, resultErr)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (e assetService) Delete(ctx context.Context, asset *model.Asset) error {
	headers := make(http.Header)

	_, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s", e.basePath, asset.Sys.ID), nil, headers)

	return err
}

func (e assetService) Publish(ctx context.Context, asset *model.Asset) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	res, err := e.client.Put(ctx, fmt.Sprintf("%s/%s/published", e.basePath, asset.Sys.ID), nil, headers, nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&asset)
}

func (e assetService) Unpublish(ctx context.Context, asset *model.Asset) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	res, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s/published", e.basePath, asset.Sys.ID), nil, headers)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&asset)
}

func (e assetService) ListPublished(ctx context.Context) cma.NextableCollection[*model.Asset, any] {
	return cma2.NewCollection[*model.Asset, any](&cma2.CollectionOptions{
		Path:   "/public/assets",
		Client: e.client,
		Ctx:    ctx,
	})
}

func (e assetService) Archive(ctx context.Context, asset *model.Asset) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	res, err := e.client.Put(ctx, fmt.Sprintf("%s/%s/archived", e.basePath, asset.Sys.ID), nil, headers, nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&asset)
}

func (e assetService) Unarchive(ctx context.Context, asset *model.Asset) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	res, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s/archived", e.basePath, asset.Sys.ID), nil, headers)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&asset)
}

func NewAssetService(client common.RestClient) cma.Assets {
	return &assetService{
		client:   client,
		basePath: "/assets",
	}
}
