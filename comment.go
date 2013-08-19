package grow

// Type: t1_
type commentThing struct {
	Id   string
	Name string
	Kind string
	Data Comment
}

type Comment struct {
	Approved_by            interface{}
	Author                 interface{}
	Author_flair_css_class interface{}
	Author_flair_text      interface{}
	Banned_by              interface{}
	Body                   interface{}
	Body_html              interface{}
	Edited                 interface{}
	Gilded                 float64
	Likes                  bool
	Link_id                interface{}
	Link_title             interface{}
	Num_reports            float64
	Parent_id              interface{}
	Score_hidden           bool
	Subreddit              interface{}
	Subreddit_id           interface{}
	Distinguished          interface{}
	Ups                    float64
	Downs                  float64
}

func (self *Comment) Upvote() error {
	return nil
}

func (self *Comment) Downvote() error {
	return nil
}
