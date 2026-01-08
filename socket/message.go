package socket

type ReadMessage struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

type WriteMessage struct {
	SenderID string `json:"sender_id"`
	RoomID   string `json:"room_id"`
	Content  string `json:"content"`
}
