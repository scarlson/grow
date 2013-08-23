package grow

// Type: t4_, reddit message wrapper
type messageThing struct {
	Id   string
	Name string
	Kind string
	Data Message
}

// data payload for message wrapper
type Message struct {
	Author        interface{}
	Body          interface{}
	Body_html     interface{}
	Context       interface{}
	First_message struct{}
	Id            interface{}
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

// send constructed message on behalf of authed user
func (self *Message) Send() error {
	return nil
}
