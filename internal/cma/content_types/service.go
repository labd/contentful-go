package content_types

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

var _ cma.ContentTypes = &contentTypeService{}

type contentTypeService struct {
	client   common.RestClient
	basePath string
}

func (c contentTypeService) ListActivated(ctx context.Context) cma.NextableCollection[*model.ContentType, any] {
	return cma2.NewCollection[*model.ContentType, any](&cma2.CollectionOptions{
		Path:   "/public/content_types",
		Client: c.client,
		Ctx:    ctx,
	})
}

func (c contentTypeService) Get(ctx context.Context, contentTypeId string) (*model.ContentType, error) {
	res, err := c.client.Get(ctx, fmt.Sprintf("%s/%s", c.basePath, contentTypeId), nil, nil)

	if err != nil {
		return nil, err
	}
	var contentType model.ContentType

	err = contentType.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return &contentType, nil
}

func (c contentTypeService) List(ctx context.Context) cma.NextableCollection[*model.ContentType, any] {
	return cma2.NewCollection[*model.ContentType, any](&cma2.CollectionOptions{
		Path:   c.basePath,
		Client: c.client,
		Ctx:    ctx,
	})
}

func (c contentTypeService) Upsert(ctx context.Context, contentType *model.ContentType) error {
	bytesArray, err := json.Marshal(contentType)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(contentType.GetVersion()))

	var res *http.Response
	var path string

	if contentType.IsNew() {
		path = fmt.Sprintf("%s/%s", c.basePath, contentType.Name)
	} else {
		path = fmt.Sprintf("%s/%s", c.basePath, contentType.Sys.ID)
	}

	res, err = c.client.Put(ctx, path, nil, headers, bytes.NewReader(bytesArray))

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&contentType)
}

func (c contentTypeService) Delete(ctx context.Context, contentType *model.ContentType) error {
	headers := make(http.Header)

	_, err := c.client.Delete(ctx, fmt.Sprintf("%s/%s", c.basePath, contentType.Sys.ID), nil, headers)

	return err
}

func (c contentTypeService) Activate(ctx context.Context, contentType *model.ContentType) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(contentType.GetVersion()))

	res, err := c.client.Put(ctx, fmt.Sprintf("%s/%s/published", c.basePath, contentType.Sys.ID), nil, headers, nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&contentType)
}

func (c contentTypeService) Deactivate(ctx context.Context, contentType *model.ContentType) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(contentType.GetVersion()))

	res, err := c.client.Delete(ctx, fmt.Sprintf("%s/%s/published", c.basePath, contentType.Sys.ID), nil, headers)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&contentType)
}

func NewContentTypeService(client common.RestClient) cma.ContentTypes {
	return &contentTypeService{
		client:   client,
		basePath: "/content_types",
	}
}
