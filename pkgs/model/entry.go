package model

type Entry struct {
	Locale string                 `json:"locale"`
	Sys    *PublishSys            `json:"sys"`
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
