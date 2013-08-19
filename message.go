package grow

// Type: t4_
type messageThing struct {
	Id   string
	Name string
	Kind string
	Data Message
}

type Message struct {
	Author        interface{}
	Body          interface{}
	Body_html     interface{}
	Context       interface{}
	First_message struct{}
	Likes         bool
	Link_title    interface{}
	Name          interface{}
	New           bool
	Parent_id     interface{}
	Replies       interface{}
	Subject       interface{}
	Subreddit     interface{}
	Was_comment   bool
}

func (self *Message) Send() error {
	return nil
}
