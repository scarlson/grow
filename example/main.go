package main

import (
    "encoding/json"
    //"fmt"
    //"github.com/scarlson/grow"
    "grow"
    "html/template"
    "io/ioutil"
    "net/http"
)

// config data gets loaded from local json file
type Config struct {
    RedditSecret string
    RedditId     string
    UserAgent    string
}

// load process reads local json file to fill config struct
func (self *Config) load(path string) error {
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

// perform config load and initialize grow library with required oauth fields
func Init() {
    err := conf.load("./config.json")
    if err != nil {
        panic(err)
    }
    grow.Config(conf.UserAgent, "identity,submit,read,privatemessages", conf.RedditId, conf.RedditSecret)
}

/* ===========================================================================
                             HTTP HANDLERS
=========================================================================== */

// display index template for user auth button
func handleRoot(w http.ResponseWriter, r *http.Request) {
    notAuthenticatedTemplate.Execute(w, nil)
}

// pass http handlers to grow library for oauth redirect
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    grow.Authorize(w, r)
}

// callback from reddit, process oauth response
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    user, _ := grow.Authorized(w, r)
    data := make(map[string]interface{})
    data["user"] = user.Name
    //a, b := grow.SubmitComment(nil, "ohshit")
    //data["noauth"], data["noautherr"] = string(a), b
    otro, otroerr := grow.GetUser("kn0thing")
    data["k"], data["kerr"] = otro.Name, otroerr
    sub, suberr := grow.GetSubreddit("redditdev")
    data["sub"], data["suberr"] = sub.Name, suberr
    data["comments"], data["commerr"] = user.Comments()
    //d, comperr := grow.Compose("kamoi", "test", "test")
    //data["compose"], data["comperr"] = string(d), comperr
    link, linkerr := grow.LinkInfo("t5_2qizd", 100, "")
    data["link"], data["linkerr"] = string(link), linkerr
    userInfoTemplate.Execute(w, data)
}

func main() {
    Init()
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/authorize", handleAuthorize)
    http.HandleFunc("/login", handleOAuth2Callback)
    http.ListenAndServe(":10000", nil)
}
