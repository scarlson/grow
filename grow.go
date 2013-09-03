package grow

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/* ===========================================================================
                         PACKAGE VARS AND INTERFACES
=========================================================================== */

var config = &oauth.Config{}
var requestURL = "https://oauth.reddit.com"
var state = ""
var UserAgent = ""
var AuthedUser = Account{}
var transport = &oauth.Transport{
	Config:    config,
	Transport: &http.Transport{},
}

/* ===========================================================================
                         PRIVATE HELPER FUNCS
=========================================================================== */

// random string generator
func randomString(l int) string {
	bites := make([]byte, l)
	for i := 0; i < l; i++ {
		bites[i] = byte(randInt(65, 90))
	}
	return string(bites)
}

// these funcs are one of those things that should be built into the core library, imo
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

/* ===========================================================================
                         IDENTITY SCOPE
=========================================================================== */

// api/me, hopefully returns a user with a modhash because oauth me doesn't
func NoauthMe() (*Account, error) {
	url := fmt.Sprintf("%s%s", requestURL, "/api/me.json")
	req, err := noauthRequest("GET", url)
	thing := &accountThing{}
	if err != nil {
		fmt.Printf("\nErr1: %v\n%v\n", err, string(req))
		return &Account{}, err
	}
	err = json.Unmarshal(req, &thing)
	if err != nil {
		fmt.Printf("\nErr2: %v\n%v\n", err, string(req))
		return &Account{}, err
	}
	acc := thing.Data
	fmt.Printf("\n\nNoAuthed User: %+v\nModhash: %+v\n", acc, acc.Modhash)
	return &acc, nil
}

// api/v1/me, returns the authed user as an Account object,
// should eventually refactor to return AccountThing instead
func Me() (*Account, error) {
	if !strings.Contains(config.Scope, "identity") {
		return &Account{}, nil // TODO: out of scope error
	}
	client := transport.Client()
	p := fmt.Sprintf("%s%s", requestURL, "/api/v1/me")
	req, err := http.NewRequest("GET", p, nil)

	// build required headers
	req.Header.Add("User-Agent", UserAgent)
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
	//fmt.Println(string(contents))
	account := &Account{}
	// cast the contents into an account object -- why is this not an account thing?
	err = json.Unmarshal(contents, account)
	AuthedUser = *account
	//fmt.Printf("\n\nAuthed User: %+v\nModhash: %+v\n", AuthedUser, AuthedUser.Modhash)
	return account, err
}

/* ===========================================================================
                         EDIT SCOPE
=========================================================================== */

// api/del, delete a comment
func DelComment(le *Comment) error {
	if !strings.Contains(config.Scope, "edit") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/editusertext, edit a comment or self post
func EditUserText() error {
	if !strings.Contains(config.Scope, "edit") {
		return nil // TODO: out of scope error
	}
	return nil
}

/* ===========================================================================
                         SUBMIT SCOPE
=========================================================================== */

// api/comment, reply to a commentable thing
func SubmitComment(le commentable, text string) ([]byte, error) {
	if !strings.Contains(config.Scope, "submit") {
		return nil, nil // TODO: out of scope error
	}
	v := url.Values{}
	v.Set("api_type", "json")
	v.Set("text", text)
	v.Set("thing_id", "t3_1lm2hf")
	v.Set("uh", AuthedUser.Modhash)
	return oauthPostRequest("/api/compose", v)
}

// api/submit, submit a link to a subreddit
func SubmitLink() error {
	if !strings.Contains(config.Scope, "submit") {
		return nil // TODO: out of scope error
	}
	return nil
}

/* ===========================================================================
                         READ SCOPE
=========================================================================== */

// api/info, fetch a link or list of links by url
func LinkInfo(id string, url string, limit int) error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/multi/mine
func MultiMine() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/multi/multipath
func MultiMultipath() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/multi/multipath/description
func GetMultipathDescription() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/multi/multipath/description
func PutMultipathDescription() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/multi/multipath/r/srname
func MultipathSubreddit() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// comments/article
func ArticleComments(article string, comment string, context string, depth string, limit int, sort string) error {
	url := fmt.Sprintf("%s%s", requestURL, "/comments/article")
	contents, err := oauthGetRequest(url)
	_, _ = contents, err
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// controversial
func Controversial() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// hot
func Hot() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// listing
func GetListing() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// new
func New() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// top
func Top() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// sort
func Sort() error {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	return nil
}

// fetch a subreddit's about.json, return a subreddit object, uses OAuth
// should refactor to return SubredditThing instead
func GetSubreddit(name string) (Subreddit, error) {
	if !strings.Contains(config.Scope, "read") {
		return Subreddit{}, nil // TODO: out of scope error
	}
	suburl := fmt.Sprintf("/r/%s/about.json", name)
	contents, err := oauthGetRequest(suburl)
	subt := &subredditThing{}
	err = json.Unmarshal(contents, subt)
	sub := subt.Data
	return sub, err
}

/* ===========================================================================
                         VOTE SCOPE
=========================================================================== */

// api/vote, upvote/downvote a link/comment
// http://www.reddit.com/dev/api#POST_api_vote
func Vote(le votable, dir int) error {
	if !strings.Contains(config.Scope, "vote") {
		return nil // TODO: out of scope error
	}
	return nil
}

/* ===========================================================================
                         MODPOSTS SCOPE
=========================================================================== */

// api/marknsfw, report a comment/link as NSFW
func MarkNSFW() error {
	if !strings.Contains(config.Scope, "modposts") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/set_contest_mode, set an object's "contest" mode
func SetContest(val bool) error {
	if !strings.Contains(config.Scope, "modposts") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/set_subreddit_sticky, set a link as subreddit's sticky
func SetSubredditSticky() error {
	if !strings.Contains(config.Scope, "modposts") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/unmarknsfw, remove a NSFW mark -- should probably just make this a toggle?
func UnmarkNSFW() error {
	if !strings.Contains(config.Scope, "modposts") {
		return nil // TODO: out of scope error
	}
	return nil
}

/* ===========================================================================
                          MOD CONFIG SCOPE
=========================================================================== */

/* ===========================================================================
                          MOD FLAIR SCOPE
=========================================================================== */

/* ===========================================================================
                          MOD LOG SCOPE
=========================================================================== */

/* ===========================================================================
                          MY SUBREDDITS SCOPE
=========================================================================== */

/* ===========================================================================
                          PRIVATE MESSAGES SCOPE
=========================================================================== */

// api/block
func Block() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/compose
/*
api_type - the string json
captcha - the user's response to the CAPTCHA challenge
iden - the identifier of the CAPTCHA challenge
subject
text - raw markdown text
to - the name of an existing user
uh - a valid modhash
*/
func Compose(to string, subject string, body string) ([]byte, error) {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil, nil // TODO: out of scope error
	}
	v := url.Values{}
	v.Set("api_type", "json")
	v.Set("subject", subject)
	v.Set("text", body)
	v.Set("to", to)
	v.Set("uh", AuthedUser.Modhash)
	return oauthPostRequest("/api/compose", v)
}

// api/read_message
func ReadMessage() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

// api/unread_message
func UnreadMessage() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

// message/inbox
func Inbox() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

// message/sent
func Sent() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

// message/unread
func Unread() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

// message/where
func Where() error {
	if !strings.Contains(config.Scope, "privatemessages") {
		return nil // TODO: out of scope error
	}
	return nil
}

/* ===========================================================================
                          UNAUTHED HANDLERS
=========================================================================== */

// fetch a user's about.json and return its account object, doesn't use OAuth
// refactor to accountthing?
func GetUser(name string) (Account, error) {
	url := fmt.Sprintf("http://reddit.com/user/%s/about.json", name)
	req, err := noauthRequest("GET", url)
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

/* ===========================================================================
                          REQUEST HANDLERS
=========================================================================== */

// send a non tokenized request for non-API restricted data, usually an about.json or some such
func noauthRequest(method string, url string) ([]byte, error) {
	client := transport.Client()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// send an oauthed request using a tokenized transport, data returned will depend on authed user
func oauthPostRequest(path string, data url.Values) ([]byte, error) {
	// is there a better way to handle post requests?
	client := transport.Client()
	p := fmt.Sprintf("%s%s", requestURL, path)
	req, err := http.NewRequest("POST", p, strings.NewReader(data.Encode()))

	// build required headers
	req.Header.Add("User-Agent", UserAgent)
	access_token := fmt.Sprintf("bearer %s", transport.Token.AccessToken)
	req.Header.Add("Authorization", access_token)
	fmt.Printf("OAuth Post: %v, %v", req, data)
	// send the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// send an oauthed request using a tokenized transport, data returned will depend on authed user
func oauthGetRequest(path string) ([]byte, error) {
	client := transport.Client()
	p := fmt.Sprintf("%s%s", requestURL, path)
	req, err := http.NewRequest("GET", p, nil)

	// build required headers
	req.Header.Add("User-Agent", UserAgent)
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
	UserAgent = useragent
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
