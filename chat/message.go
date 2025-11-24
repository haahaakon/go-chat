package chat

type Message struct {
	Userid  int
	Content string
}

func NewMessage(userid int, content string) *Message {
	message := Message{
		Userid:  userid,
		Content: content,
	}
	return &message
}
