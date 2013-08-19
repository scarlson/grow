package main

import (
	"encoding/json"
	"fmt"
	"github.com/scarlson/grow"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Config struct {
    // config data gets loaded from local json file
	RedditSecret string
	RedditId     string
	UserAgent    string
}

func (self *Config) load(path string) error {
    // load process reads local json file to fill config struct
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &self)
	if err != nil {
        return err
	}
	return nil
}


// html templates for displaying data and required user actions
var notAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
<form action="/authorize" method="POST"><input type="submit" value="Authorize"/></form>
</body></html>
`))

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
{{.}}
</body></html>
`))


var conf = &Config{}

func Init() {
    // perform config load and initialize grow library with required oauth fields
	err := conf.load("./config.json")
    if err != nil {
        panic(err)
    }
	grow.Config(conf.UserAgent, "identity,read", conf.RedditId, conf.RedditSecret)
}


/* ===========================================================================
                             HTTP HANDLERS
=========================================================================== */

func handleRoot(w http.ResponseWriter, r *http.Request) {
    // display index template for user auth button
	notAuthenticatedTemplate.Execute(w, nil)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    // pass http handlers to grow library for oauth redirect
    grow.Authorize(w, r)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// callback from reddit, process oauth response
    user, _ := grow.Authorized(w, r)
	data := make(map[string]interface{})
    data["user"] = user.Name
    otro, _ := grow.GetUser("kn0thing")
    data["k"] = otro.Name
    sub, _ := grow.GetSubreddit("redditdev")
    data["sub"] = sub.Display_name
    data["comments"], _ = user.Comments()
    userInfoTemplate.Execute(w, data)
}

func main() {
    Init()
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/authorize", handleAuthorize)
	http.HandleFunc("/login", handleOAuth2Callback)
	http.ListenAndServe(":10000", nil)
}
