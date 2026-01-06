package socket

type ReadMessage struct {
	UserID  string `json:"-"`
	Channel string `json:"channel"`
	Content string `json:"content"`
}

func (m *ReadMessage) isForAuthentication() bool {
	return m.Channel == "watergun:auth"
}
