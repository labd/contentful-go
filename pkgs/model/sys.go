package model

type BaseSys struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	LinkType string `json:"linkType,omitempty"`
	Version  int    `json:"version,omitempty"`
}

type EnvironmentSys struct {
	SpaceSys
	Environment *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"environment,omitempty"`
}

type CreatedSys struct {
	BaseSys
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	UpdatedBy *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"updatedBy,omitempty"`
	CreatedBy *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"createdBy,omitempty"`
}

type SpaceSys struct {
	CreatedSys
	Space *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"space,omitempty"`
}

type PublishSys struct {
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

//// Sys model
//type Sys struct {
//	BaseSys
//	SpaceSys
//	//Revision    int `json:"revision,omitempty"`
//	//ContentType *struct {
//	//	Sys Sys `json:"sys,omitempty"`
//	//} `json:"contentType,omitempty"`
//	//FirstPublishedAt string `json:"firstPublishedAt,omitempty"`
//	//PublishedCounter int    `json:"publishedCounter,omitempty"`
//	//PublishedAt      string `json:"publishedAt,omitempty"`
//	//PublishedBy      *Sys   `json:"publishedBy,omitempty"`
//	//PublishedVersion int    `json:"publishedVersion,omitempty"`
//	//ArchivedAt       string `json:"archivedAt,omitempty"`
//	//ArchivedBy       *Sys   `json:"archivedBy,omitempty"`
//	//ArchivedVersion  int    `json:"archivedVersion,omitempty"`
//}
