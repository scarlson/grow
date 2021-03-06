package grow

// Type: t5_, reddit subreddit wrapper
type subredditThing struct {
	Id   string
	Name string
	Kind string
	Data Subreddit
}

// data payload for subreddit wrapper
type Subreddit struct {
	Accounts_active    float64
	Description        string
	Description_html   interface{}
	Display_name       interface{}
	Header_img         interface{}
	Header_size        []interface{}
	Header_title       interface{}
	Id                 string
	Name               string
	Over18             bool
	Public_description string
	Subscribers        float64
	Title              string
	Url                interface{}
}

// subscribe to this subreddit for the authed user
//
// is there an Unsubscribe?  Or does Subscribe toggle?
func Subscribe() error {
	return nil
}

// retrieve subreddit about.json, needs to return something not an error
func About() error {
	return nil
}
