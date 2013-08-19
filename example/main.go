package main

import (
	"encoding/json"
	"fmt"
	"grow"
	"html/template"
	"io/ioutil"
	"net/http"
)

/* ===========================================================================
                        CONFIGS AND TEMPLATES
=========================================================================== */

type Config struct {
	RedditSecret string
	RedditId     string
	UserAgent    string
}

func (self *Config) load(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &self)
	if err != nil {
		fmt.Print("bad json ", err)
	}
	return nil
}

var conf = &Config{}

func Init() {
	err := conf.load("./config.json")
    if err != nil {
        panic(err)
    }
	grow.Config(conf.UserAgent, "identity,read", conf.RedditId, conf.RedditSecret)
}

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

/* ===========================================================================
                             HTTP HANDLERS
=========================================================================== */

func handleRoot(w http.ResponseWriter, r *http.Request) {
	notAuthenticatedTemplate.Execute(w, nil)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    grow.Authorize(w, r)
}

// callback from Reddit
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
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
