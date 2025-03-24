package entries

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

var _ cma.Entries = &entryService{}

type entryService struct {
	client   common.RestClient
	basePath string
}

func (e entryService) ListPublished(ctx context.Context) cma.NextableCollection[*model.Entry, any] {
	return cma2.NewCollection[*model.Entry, any](&cma2.CollectionOptions{
		Path:   "/public/entries",
		Client: e.client,
		Ctx:    ctx,
	})
}

func (e entryService) Get(ctx context.Context, entryId string) (*model.Entry, error) {
	res, err := e.client.Get(ctx, fmt.Sprintf("%s/%s", e.basePath, entryId), nil, nil)

	if err != nil {
		return nil, err
	}
	var entry model.Entry

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (e entryService) List(ctx context.Context) cma.NextableCollection[*model.Entry, any] {
	return cma2.NewCollection[*model.Entry, any](&cma2.CollectionOptions{
		Path:   e.basePath,
		Client: e.client,
		Ctx:    ctx,
	})
}

func (e entryService) Upsert(ctx context.Context, contentTypeID string, entry *model.Entry) error {
	bytesArray, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(entry.GetVersion()))
	headers.Set("X-Contentful-Content-Type", contentTypeID)

	var res *http.Response

	if entry.IsNew() {
		res, err = e.client.Post(ctx, e.basePath, nil, headers, bytes.NewReader(bytesArray))
	} else {
		res, err = e.client.Put(ctx, fmt.Sprintf("%s/%s", e.basePath, entry.Sys.ID), nil, headers, bytes.NewReader(bytesArray))
	}

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&entry)
}

func (e entryService) Delete(ctx context.Context, entry *model.Entry) error {
	headers := make(http.Header)

	_, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s", e.basePath, entry.Sys.ID), nil, headers)

	return err
}

func (e entryService) Publish(ctx context.Context, entry *model.Entry) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(entry.GetVersion()))

	res, err := e.client.Put(ctx, fmt.Sprintf("%s/%s/published", e.basePath, entry.Sys.ID), nil, headers, nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&entry)
}

func (e entryService) Unpublish(ctx context.Context, entry *model.Entry) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(entry.GetVersion()))

	res, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s/published", e.basePath, entry.Sys.ID), nil, headers)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&entry)
}

func (e entryService) Archive(ctx context.Context, entry *model.Entry) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(entry.GetVersion()))

	res, err := e.client.Put(ctx, fmt.Sprintf("%s/%s/archived", e.basePath, entry.Sys.ID), nil, headers, nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&entry)
}

func (e entryService) Unarchive(ctx context.Context, entry *model.Entry) error {
	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(entry.GetVersion()))

	res, err := e.client.Delete(ctx, fmt.Sprintf("%s/%s/archived", e.basePath, entry.Sys.ID), nil, headers)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&entry)
}

func NewEntriesService(client common.RestClient) cma.Entries {
	return &entryService{
		client:   client,
		basePath: "/entries",
	}
}
