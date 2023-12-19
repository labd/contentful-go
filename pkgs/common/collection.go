package common

import (
	"github.com/flaconi/contentful-go/pkgs/model"
)

// CollectionOptions holds init options
type CollectionOptions struct {
	Limit uint16
}

type Collection[Items any, Includes any] struct {
	Query
	Sys      *model.Sys `json:"sys"`
	Total    int        `json:"total"`
	Skip     int        `json:"skip"`
	Limit    int        `json:"limit"`
	Items    []Items    `json:"items"`
	Includes Includes   `json:"includes"`
}

// NewCollection initializes a new collection
func NewCollection[Items any, Includes any](options *CollectionOptions) *Collection[Items, Includes] {
	query := NewQuery()
	query.Order("sys.createdAt", true)

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	return &Collection[Items, Includes]{
		Query: *query,
	}
}
