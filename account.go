package grow

import (
	"encoding/json"
	"fmt"
)

// Type: t2_
type accountThing struct {
	Id   string
	Name string
	Kind string
	Data Account
}

type Account struct {
	Comment_karma      float64
	Created            interface{}
	Created_utc        interface{}
	Has_mail           interface{}
	Has_mod_mail       interface{}
	Has_verified_email bool
	Id                 interface{}
	Is_friend          bool
	Is_gold            bool
	Is_mod             bool
	Link_karma         float64
	Modhash            interface{}
	Name               interface{}
	Over_18            bool
}

func (user Account) Comments() ([]Comment, error) {
	url := fmt.Sprintf("http://reddit.com/user/%s/comments.json", user.Name)
	commentThing := &commentListingThing{}
	req, err := noauthRequest("GET", url, "asdf")
	if err != nil {
		return []Comment{}, err
	}
	err = json.Unmarshal(req, commentThing)
	comments := commentThing.Data.Children
	c := make([]Comment, len(comments))
	for i := range comments {
		c[i] = comments[i].Data
	}
	return c, err
}
