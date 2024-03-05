package model

type EntrySys struct {
	EnvironmentSys
	ContentType *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"contentType,omitempty"`
	FirstPublishedAt string   `json:"firstPublishedAt,omitempty"`
	PublishedCounter int      `json:"publishedCounter,omitempty"`
	PublishedAt      string   `json:"publishedAt,omitempty"`
	PublishedBy      *BaseSys `json:"publishedBy,omitempty"`
	PublishedVersion int      `json:"publishedVersion,omitempty"`
	ArchivedAt       string   `json:"archivedAt,omitempty"`
	ArchivedBy       *BaseSys `json:"archivedBy,omitempty"`
	ArchivedVersion  int      `json:"archivedVersion,omitempty"`
}

type Entry struct {
	Locale string                 `json:"locale"`
	Sys    *EntrySys              `json:"sys"`
	Fields map[string]interface{} `json:"fields"`
}

func (entry *Entry) GetVersion() int {
	version := 1
	if entry.Sys != nil {
		version = entry.Sys.Version
	}

	return version
}

func (entry *Entry) IsNew() bool {
	return entry.Sys == nil || entry.Sys.ID == ""
}

func (entry *Entry) IsPublished() bool {
	return entry.Sys.PublishedVersion > 0 && entry.Sys.Version == entry.Sys.PublishedVersion+1
}
