package grow

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/* ===========================================================================
                         PACKAGE VARS AND INTERFACES
=========================================================================== */

var config = &oauth.Config{}
var requestURL = "https://oauth.reddit.com"
var global_delay = 2 * time.Second // reddit throttles requests to 30/min, we pop a work queue every 2s
var lastRequest time.Time
var state string
var UserAgent string
var AuthedUser = Account{}
var transport = &oauth.Transport{
	Config:    config,
	Transport: &http.Transport{},
}

//requestChan := make(chan []byte)

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

func redditQueue(r *http.Request) ([]byte, error) {
	// process the queue requests every global_delay seconds
	// there should be a better way to do this without blocking
	client := transport.Client()
	for time.Now().Before(lastRequest.Add(global_delay)) {
		time.Sleep(time.Millisecond * 100)
	}
	resp, err := client.Do(r)
	fmt.Printf("\n\nQueue req: %v\n\n", r)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	lastRequest = time.Now()
	return ioutil.ReadAll(resp.Body)
}

/* ===========================================================================
                         IDENTITY SCOPE
=========================================================================== */

// api/me, hopefully returns a user with a modhash because oauth me doesn't
func NoauthMe() (*Account, error) {
	urls := fmt.Sprintf("%s%s", requestURL, "/api/me.json")
	req, err := noauthRequest("GET", urls)
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
func SubmitComment(le interface{}, text string) ([]byte, error) {
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
func LinkInfo(id string, limit int, turl string) ([]byte, error) {
	if !strings.Contains(config.Scope, "read") {
		return nil, nil // TODO: out of scope error
	}
	v := url.Values{}
	if id != "" {
		v.Set("id", id)
	}
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 1
	}
	if turl != "" {
		v.Set("url", turl)
	}
	v.Set("limit", strconv.Itoa(limit))
	durl := fmt.Sprintf("/api/info?%s", v.Encode())
	return oauthGetRequest(durl)
}

// api/multi/mine
func MultiMine() ([]byte, error) {
	if !strings.Contains(config.Scope, "read") {
		return nil, nil // TODO: out of scope error
	}
	return oauthGetRequest("/api/multi/mine")
}

// api/multi/multipath
func MultiMultipath(multipath string) ([]byte, error) {
	if !strings.Contains(config.Scope, "read") {
		return nil, nil // TODO: out of scope error
	}
	urls := fmt.Sprintf("/api/multi/%s", multipath)
	return oauthGetRequest(urls)
}

// api/multi/multipath/description
func GetMultipathDescription() ([]byte, error) {
	if !strings.Contains(config.Scope, "read") {
		return nil // TODO: out of scope error
	}
	urls := fmt.Sprintf("/api/multi/%s/description", multipath)
	return oauthGetRequest(urls)
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
	urls := fmt.Sprintf("%s%s", requestURL, "/comments/article")
	contents, err := oauthGetRequest(urls)
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
	urls := fmt.Sprintf("http://reddit.com/user/%s/about.json", name)
	req, err := noauthRequest("GET", urls)
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
func noauthRequest(method string, urls string) ([]byte, error) {
	fmt.Printf("\n\n%vting: %v\n\n", method, urls)
	req, err := http.NewRequest(method, urls, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	return redditQueue(req)
}

// send an oauthed request using a tokenized transport, data returned will depend on authed user
func oauthPostRequest(path string, data url.Values) ([]byte, error) {
	// is there a better way to handle post requests?
	p := fmt.Sprintf("%s%s", requestURL, path)
	fmt.Printf("\n\nOAuth POSTting: %v\n\n", p)
	req, err := http.NewRequest("POST", p, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// build required headers
	req.Header.Add("User-Agent", UserAgent)
	access_token := fmt.Sprintf("bearer %s", transport.Token.AccessToken)
	req.Header.Add("Authorization", access_token)

	return redditQueue(req)
}

// send an oauthed request using a tokenized transport, data returned will depend on authed user
func oauthGetRequest(path string) ([]byte, error) {
	p := fmt.Sprintf("%s%s", requestURL, path)
	fmt.Printf("\n\nOAuth GETting: %v\n\n", p)
	req, err := http.NewRequest("GET", p, nil)
	if err != nil {
		return nil, err
	}

	// build required headers
	req.Header.Add("User-Agent", UserAgent)
	access_token := fmt.Sprintf("bearer %s", transport.Token.AccessToken)
	req.Header.Add("Authorization", access_token)

	return redditQueue(req)
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
	urls := config.AuthCodeURL(state)
	http.Redirect(w, r, urls, http.StatusFound)
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
