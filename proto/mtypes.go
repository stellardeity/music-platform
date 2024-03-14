package proto

type WsCmd struct {
	Cmd string `json:"cmd"`
}

type WsMessage struct {
	WsCmd
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
