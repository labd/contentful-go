package locales

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

var _ cma.Locales = &localeService{}

type localeService struct {
	client   common.RestClient
	basePath string
}

func (l localeService) Get(ctx context.Context, localeId string) (*model.Locale, error) {
	res, err := l.client.Get(ctx, fmt.Sprintf("%s/%s", l.basePath, localeId), nil, nil)

	if err != nil {
		return nil, err
	}
	var locale model.Locale

	err = locale.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return &locale, nil
}

func (l localeService) List(ctx context.Context) cma.NextableCollection[*model.Locale, any] {
	return cma2.NewCollection[*model.Locale, any](&cma2.CollectionOptions{
		Path:   l.basePath,
		Client: l.client,
		Ctx:    ctx,
	})
}

func (l localeService) Upsert(ctx context.Context, locale *model.Locale) error {
	bytesArray, err := json.Marshal(locale)
	if err != nil {
		return err
	}

	headers := make(http.Header)

	headers.Set("X-Contentful-Version", strconv.Itoa(locale.GetVersion()))

	var res *http.Response

	if locale.IsNew() {
		res, err = l.client.Post(ctx, l.basePath, nil, headers, bytes.NewReader(bytesArray))
	} else {
		res, err = l.client.Put(ctx, fmt.Sprintf("%s/%s", l.basePath, locale.Sys.ID), nil, headers, bytes.NewReader(bytesArray))
	}

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&locale)
}

func (l localeService) Delete(ctx context.Context, locale *model.Locale) error {
	headers := make(http.Header)

	_, err := l.client.Delete(ctx, fmt.Sprintf("%s/%s", l.basePath, locale.Sys.ID), nil, headers)

	return err
}

func NewLocaleService(client common.RestClient) cma.Locales {
	return &localeService{
		client:   client,
		basePath: "/locales",
	}
}
