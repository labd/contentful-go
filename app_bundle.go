package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// AppBundleService service
type AppBundleService service

type appBundleCreation struct {
	Upload  upload `json:"upload"`
	Comment string `json:"comment,omitempty"`
}

type upload struct {
	Sys Sys `json:"sys"`
}

type AppBundle struct {
	Sys     *Sys         `json:"sys"`
	Comment string       `json:"comment"`
	Files   []BundleFile `json:"files"`
}

type BundleFile struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	MD5  string `json:"md5"`
}

// Create creates a new app bundle
func (service *AppBundleService) Create(organizationID string, definitionId string, comment string, uploadId string) (*AppBundle, error) {
	var uploadData = appBundleCreation{
		Upload:  upload{Sys: Sys{ID: uploadId, Type: "Link", LinkType: "AppUpload"}},
		Comment: comment,
	}

	bytesArray, err := json.Marshal(uploadData)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/organizations/%s/app_definitions/%s/app_bundles", organizationID, definitionId)
	method := "POST"
	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, err
	}

	var asset AppBundle
	if err = service.c.do(req, &asset); err != nil {
		return nil, err
	}

	return &asset, nil
}
