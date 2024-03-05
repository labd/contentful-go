package model

type Asset struct {
	Locale string
	Sys    *EnvironmentSys `json:"sys,omitempty"`
	Fields *AssetFields    `json:"fields,omitempty"`
}

func (asset *Asset) GetVersion() int {
	version := 1
	if asset.Sys != nil {
		version = asset.Sys.Version
	}

	return version
}

func (asset *Asset) IsNew() bool {
	return asset.Sys == nil || asset.Sys.ID == ""
}

type AssetFields struct {
	Title       map[string]string `json:"title,omitempty"`
	Description map[string]string `json:"description,omitempty"`
	File        map[string]*File  `json:"file,omitempty"`
}

// File represents a Contentful File
type File struct {
	URL         string       `json:"url,omitempty"`
	UploadURL   string       `json:"upload,omitempty"`
	UploadFrom  *UploadFrom  `json:"uploadFrom,omitempty"`
	Details     *FileDetails `json:"details,omitempty"`
	FileName    string       `json:"fileName,omitempty"`
	ContentType string       `json:"contentType,omitempty"`
}

// UploadFrom model
type UploadFrom struct {
	Sys *BaseSys `json:"sys,omitempty"`
}

// FileDetails model
type FileDetails struct {
	Size  int          `json:"size,omitempty"`
	Image *ImageFields `json:"image,omitempty"`
}

// ImageFields model
type ImageFields struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}
