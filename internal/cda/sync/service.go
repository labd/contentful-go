package sync

import (
	"context"

	cma2 "github.com/labd/contentful-go/internal/cma/common"
	"github.com/labd/contentful-go/service/cda"
	"github.com/labd/contentful-go/service/cma"
	"github.com/labd/contentful-go/service/common"
)

var _ cda.Sync = &syncService{}

type syncService struct {
	client   common.RestClient
	basePath string
}

func (s syncService) GetFromSyncUrl(ctx context.Context, syncUrl string) cma.SyncCollection {
	return cma2.NewSyncCollection(&cma2.CollectionOptions{
		Path:   s.basePath,
		Client: s.client,
		Ctx:    ctx,
	}, &syncUrl)
}

func (s syncService) Init(ctx context.Context, syncType cda.SyncType, contentType *string) cma.SyncCollection {
	collection := cma2.NewSyncCollection(&cma2.CollectionOptions{
		Path:   s.basePath,
		Client: s.client,
		Ctx:    ctx,
	}, nil)

	collection.GetQuery().Limit(1000).Equal("type", syncType.String()).Equal("initial", "true")
	if syncType == cda.Entry {
		collection.GetQuery().Limit(100).ContentType(*contentType).Equal("type", cda.Entry.String())
	}

	return collection
}

func NewSyncService(client common.RestClient) cda.Sync {
	return syncService{
		client:   client,
		basePath: "/sync",
	}
}
