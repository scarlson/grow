package grow

// Kind: Listing
type listingThing struct {
	Kind string
	Data Listing
}

type Listing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []interface{}
}

type AccountListing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []accountThing
}

type SubredditListing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []subredditThing
}

type commentListingThing struct {
	Kind string
	Data CommentListing
}

type CommentListing struct {
	Before   interface{}
	After    interface{}
	Modhash  interface{}
	Children []commentThing
}
