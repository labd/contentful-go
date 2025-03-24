package client

import (
	"log/slog"

	"github.com/labd/contentful-go/service/common"
)

type ClientConfig struct {
	URL        *string
	HTTPClient common.HttpClient
	Debug      bool
	UserAgent  *string
	Token      string
	Logger     *slog.Logger
}
