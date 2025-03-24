package contentful

import (
	"bytes"
	"fmt"
	"strconv"
)

// AppUploadService service
type AppUploadService struct {
	service
	BaseURL string
}

type AppUpload struct {
	Sys *Sys `json:"sys"`
}

// Create creates a new access token
func (service *AppUploadService) Create(organizationID string, bundleData []byte) (*AppUpload, error) {
	orgBaseURL := service.c.BaseURL

	defer func() { service.c.BaseURL = orgBaseURL }()

	service.c.BaseURL = service.BaseURL

	path := fmt.Sprintf("/organizations/%s/app_uploads", organizationID)
	method := "POST"
	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bundleData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.Itoa(len(bundleData)))
	var asset AppUpload
	if err = service.c.do(req, &asset); err != nil {
		return nil, err
	}

	return &asset, nil
}
