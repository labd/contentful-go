package common

import (
	"context"
	"encoding/json"
	"github.com/flaconi/contentful-go/pkgs/common"
	"github.com/flaconi/contentful-go/service/cma"
	common2 "github.com/flaconi/contentful-go/service/common"
)

var _ cma.NextableCollection[any, any] = &Collection[any, any]{}

// CollectionOptions holds init options
type CollectionOptions struct {
	Path   string
	Ctx    context.Context
	Client common2.RestClient
}

type Collection[Items any, Includes any] struct {
	common.Collection[Items, Includes]
	path   string
	ctx    context.Context
	client common2.RestClient
	page   uint16
}

// NewCollection initializes a new collection
func NewCollection[Items any, Includes any](options *CollectionOptions) *Collection[Items, Includes] {
	return &Collection[Items, Includes]{
		Collection: *common.NewCollection[Items, Includes](&common.CollectionOptions{}),
		ctx:        options.Ctx,
		path:       options.Path,
		client:     options.Client,
		page:       1,
	}
}

func (col Collection[Items, Includes]) Next() (*common.Collection[Items, Includes], error) {
	// setup query params
	skip := uint16(col.Limit) * (col.page - 1)
	col.Query.Skip(skip)

	res, err := col.client.Get(col.ctx, col.path, col.Query.Values(), nil)

	if err != nil {
		return nil, err
	}

	col.page++

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&col.Collection)
	if err != nil {
		return nil, err
	}
	return &col.Collection, nil
}
