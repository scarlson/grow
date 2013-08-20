package grow

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

/* ===========================================================================
                         PACKAGE VARS AND INTERFACES
=========================================================================== */

var config = &oauth.Config{}
var requestURL = "https://oauth.reddit.com"
var state = ""
var transport = &oauth.Transport{
	Config:    config,
	Transport: &http.Transport{},
}

// generic interface for any object that's votable
type votable interface {
	Upvote() error
	Downvote() error
}

// generic interface for any object that's commentable
type commentable interface {
}

/* ===========================================================================
                         PRIVATE HELPER FUNCS
=========================================================================== */

// random string generator
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

// these funcs are one of those things that should be built into the core library, imo
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

/* ===========================================================================
                         PUBLIC FUNCS
=========================================================================== */

// this is redundant to the comment method that should do the same.  Then again,
// it should be a commentable interface so you can comment on links as well
func SubmitComment() (*Comment, error) {
	return nil, nil
}

// reddit's api/v1/me handler, returns the authed user as an Account object,
// should eventually refactor to return AccountThing instead
func Me() (*Account, error) {
	client := transport.Client()
	p := fmt.Sprintf("%s%s", requestURL, "/api/v1/me")
	req, err := http.NewRequest("GET", p, nil)

	// build required headers
	req.Header.Add("User-Agent", "useragent")
	access_token := fmt.Sprintf("bearer %s", transport.Token.AccessToken)
	req.Header.Add("Authorization", access_token)

	// send the request
	res, err := client.Do(req)

	if err != nil {
		return &Account{}, err
	}
	defer res.Body.Close()

	// read the contents from the http response
	contents, err := ioutil.ReadAll(res.Body)
	account := &Account{}
	// cast the contents into an account objects -- why is this not an account thing?
	err = json.Unmarshal(contents, account)
	return account, err
}

// fetch a user's about.json and return its account object, doesn't use OAuth
func GetUser(name string) (Account, error) {
	url := fmt.Sprintf("http://reddit.com/user/%s/about.json", name)
	req, err := noauthRequest("GET", url, "")
	thing := &accountThing{}
	if err != nil {
		return Account{}, err
	}
	err = json.Unmarshal(req, &thing)
	if err != nil {
		return Account{}, err
	}
	acc := thing.Data
	return acc, nil
}

// fetch a subreddit's about.json, return a subreddit object, uses OAuth -- should it not?
func GetSubreddit(name string) (Subreddit, error) {
	suburl := fmt.Sprintf("/r/%s/about.json", name)
	contents, err := oauthRequest("GET", suburl, "")
	subt := &subredditThing{}
	err = json.Unmarshal(contents, subt)
	sub := subt.Data
	return sub, err
}

/* ===========================================================================
                          REQUEST HANDLERS
=========================================================================== */

// send a non tokenized request for non-API restricted data, usually an about.json or some such
func noauthRequest(method string, url string, user string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", user)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// send an oauthed request using a tokenized transport, data returned will depend on authed user
func oauthRequest(method string, path string, user string) ([]byte, error) {
	client := transport.Client()
	p := fmt.Sprintf("%s%s", requestURL, path)
	req, err := http.NewRequest(method, p, nil)

	// build required headers
	req.Header.Add("User-Agent", user)
	access_token := fmt.Sprintf("bearer %s", transport.Token.AccessToken)
	req.Header.Add("Authorization", access_token)

	// send the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// setup the oauth client according to the user's app data from reddit
func Config(useragent string, scope string, redditid string, redditsecret string) bool {
	config.ClientId = redditid
	config.ClientSecret = redditsecret
	config.Scope = scope
	config.AuthURL = "https://ssl.reddit.com/api/v1/authorize"
	config.TokenURL = "https://ssl.reddit.com/api/v1/access_token"
	config.RedirectURL = "http://redditbank.com/login"
	return true
}

// send the authorize redirect to reddit for a user to authorize our app's access to their account
func Authorize(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())
	state = randomString(10)
	url := config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// oauth callback from reddit
func Authorized(w http.ResponseWriter, r *http.Request) (*Account, error) {
	//Get the code from the response
	code := r.FormValue("code")
	check_state := r.FormValue("state")

	if check_state != state {
		panic("states don't match!")
	}

	// Set up a Transport using the config.
	trans := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			req.SetBasicAuth(config.ClientId, config.ClientSecret)
			return nil, nil
		},
	}
	token_trans := &oauth.Transport{
		Config:    config,
		Transport: trans,
	}

	token, err := token_trans.Exchange(code)
	if err != nil {
		panic(err)
	}

	// add token to the global transport so every request can use it
	transport.Token = token

	// finally return the authed user, should probably just return some
	// error or nil, allow the user to decide if they need Me()
	return Me()
}
