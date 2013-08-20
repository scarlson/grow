package grow

// Type: t1_, reddit comment wrapper
type commentThing struct {
	Id   string
	Name string
	Kind string
	Data Comment
}

// data payload from a comment wrapper
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

// Upvote this comment for the authed user
func (self *Comment) Upvote() error {
	return nil
}

// Downvote this comment for the authed user
func (self *Comment) Downvote() error {
	return nil
}
