package model

type WebSocketMessage struct {
	Type    string            `json:"type"`
	Content map[string]string `json:"content"`
}
