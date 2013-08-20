package grow

// Type: t3_, reddit link wrapper
type linkThing struct {
	Id   string
	Name string
	Kind string
	Data Link
}


// data payload from a link wrapper
type Link struct {
	Author                 interface{}
	Author_flair_css_class interface{}
	Autho_flair_text       interface{}
	Clicked                bool
	Domain                 interface{}
	Hidden                 bool
	Is_self                bool
	Likes                  bool
	Link_flair_css_class   interface{}
	Link_flair_text        interface{}
	Media                  struct{}
	Media_embed            struct{}
	Num_comments           float64
	Over_18                bool
	Permalink              interface{}
	Saved                  bool
	Score                  float64
	Selftext               interface{}
	Selftext_html          interface{}
	Subreddit              interface{}
	Subreddit_id           interface{}
	Thumbnail              interface{}
	Title                  interface{}
	Url                    interface{}
	Edited                 float64
	Distinguished          interface{}
	Ups                    float64
	Downs                  float64
}

// upvote the link for the authed user
func (self *Link) Upvote() error {
	return nil
}

// downvote the link for the authed user
func (self *Link) Downvote() error {
	return nil
}
