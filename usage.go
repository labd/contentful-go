package contentful

// Usage model
type Usage struct {
	Sys           *Sys              `json:"sys"`
	UnitOfMeasure string            `json:"unitOfMeasure"`
	Metric        string            `json:"metric"`
	DateRange     DateRange         `json:"dateRange"`
	TotalUsage    int               `json:"usage"`
	UsagePerDay   map[string]string `json:"usagePerDay"`
}

// DateRange model
type DateRange struct {
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}
