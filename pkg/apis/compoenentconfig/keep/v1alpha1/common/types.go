package common

type Header struct {
	HeaderMap map[string]string `json:"header_map"`
}

type Router struct {
	RouterPath string `json:"routerPath,omitempty"`
}

type Message struct {
	Content string `json:"content"`
}

type WebSocketMessage struct {
	HeaderMap      *Header   `json:"headers"`
	Routers        *[]Router `json:"routers"`
	MessageContent *Message  `json:"messageContent"`
}
