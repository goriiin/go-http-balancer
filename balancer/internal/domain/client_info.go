package domain

type ClientInfo struct {
	ClientID   string `json:"client_id"`
	Capacity   int64  `json:"capacity"`
	Rate       int64  `json:"rate"`
	CustomData string `json:"custom_data,omitempty"`
}
