package messages

// easyjson:json
type UrlMessage struct {
	Domain string `json:"domain"`
	Path string `json:"path"`
}

// easyjson:json
type LocationMessage struct {
	Location string `json:"location"`
}
