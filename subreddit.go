package grow

// Type: t5_
type subredditThing struct {
	Id   string
	Name string
	Kind string
	Data Subreddit
}

type Subreddit struct {
	Accounts_active    float64
	Description        interface{}
	Description_html   interface{}
	Display_name       interface{}
	Header_img         interface{}
	Header_size        []interface{}
	Header_title       interface{}
	Over18             bool
	Public_description interface{}
	Subscribers        float64
	Title              interface{}
	Url                interface{}
}

func Subscribe() error {
    return nil
}

func About() error {

}
