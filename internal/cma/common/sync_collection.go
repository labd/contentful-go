package common

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/labd/contentful-go/pkgs/common"
	"github.com/labd/contentful-go/service/cma"
	common2 "github.com/labd/contentful-go/service/common"
)

var _ cma.SyncCollection = &SyncCollection{}

type SyncCollection struct {
	common.InterfaceCollection
	path    string
	ctx     context.Context
	client  common2.RestClient
	page    uint16
	syncUrl *string
}

func (s SyncCollection) Next() (*common.InterfaceCollection, error) {

	queryValues := s.Query.Values()

	if s.syncUrl != nil {
		parsedUrl, err := url.Parse(*s.syncUrl)

		if err != nil {
			return nil, err
		}

		queryValues = parsedUrl.Query()
	}

	res, err := s.client.Get(s.ctx, s.path, queryValues, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&s.InterfaceCollection)
	if err != nil {
		return nil, err
	}
	if len(s.InterfaceCollection.NextPageUrl) > 0 {
		nextPageCollection, err := s.getNextPage(s.InterfaceCollection.NextPageUrl)
		if err != nil {
			return nil, err
		}

		s.InterfaceCollection.Items = append(s.InterfaceCollection.Items, nextPageCollection.Items...)
		s.InterfaceCollection.NextSyncUrl = nextPageCollection.NextSyncUrl
	}

	return &s.InterfaceCollection, nil
}

func (s SyncCollection) getNextPage(nextPageUrl string) (*common.InterfaceCollection, error) {

	parsedUrl, err := url.Parse(nextPageUrl)

	if err != nil {
		return nil, err
	}

	res, err := s.client.Get(s.ctx, s.path, parsedUrl.Query(), nil)
	if err != nil {
		return nil, err
	}

	collection := common.InterfaceCollection{}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&collection)
	if err != nil {
		return nil, err
	}

	if len(collection.NextPageUrl) > 0 {
		nextPageCollection, err := s.getNextPage(collection.NextPageUrl)
		if err != nil {
			return nil, err
		}

		collection.Items = append(collection.Items, nextPageCollection.Items...)
		collection.NextSyncUrl = nextPageCollection.NextSyncUrl
	}

	return &collection, nil
}

func (s SyncCollection) GetQuery() *common.Query {
	return s.Query
}

// NewSyncCollection initializes a new collection
func NewSyncCollection(options *CollectionOptions, syncUrl *string) *SyncCollection {

	return &SyncCollection{
		InterfaceCollection: *common.NewInterfaceCollection(&common.CollectionOptions{}),
		ctx:                 options.Ctx,
		path:                options.Path,
		client:              options.Client,
		page:                1,
		syncUrl:             syncUrl,
	}
}
