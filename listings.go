package grow

// Kind: Listing, generic listing wrapper - houses abritrary data types
type listingThing struct {
	Kind string
	Data Listing
}

// generic listing object, arbitrary data types
type Listing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []interface{}
}

// Array of reddit account wrapper objects
type AccountListing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []accountThing
}

// Array of reddit subreddit wrapper objects
type SubredditListing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []subredditThing
}

// Array of reddit comment wrapper objects
type commentListingThing struct {
	Kind string
	Data CommentListing
}

// Array of reddit comment payloads
type CommentListing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []commentThing
}
