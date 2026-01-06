package socket

type ReadMessage struct {
	UserID  string `json:"-"`
	Room    string `json:"room"`
	Content string `json:"content"`
}

func (m *ReadMessage) isForAuthentication() bool {
	return m.Room == "watergun:auth"
}
