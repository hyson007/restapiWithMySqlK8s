package message

type Message struct {
	Content string
}

func New(content string) Message {
	return Message{content}
}
