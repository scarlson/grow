package grow

// generic interface for any object that's votable
type votable interface {
	Upvote() error
	Downvote() error
}

// generic interface for any object that's commentable
type commentable interface {
}
