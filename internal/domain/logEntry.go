package domain

type LogEntry struct {
	Timestamp       int64               `json:"timestamp"`
	Message         string              `json:"message"`
	Service         string              `json:"service"`
	Method          string              `json:"method"`
	Path            string              `json:"path"`
	Query           string              `json:"query"`
	Request         string              `json:"request"`
	Response        string              `json:"response"`
	Status          int                 `json:"status"`
	RequestHeaders  map[string][]string `json:"requestHeaders"`
	ResponseHeaders map[string][]string `json:"responseHeaders"`
}
