package client

import (
	"github.com/flaconi/contentful-go/service/common"
	"log/slog"
)

type ClientConfig struct {
	URL        *string
	HTTPClient common.HttpClient
	Debug      bool
	UserAgent  *string
	Token      string
	Logger     *slog.Logger
}
