package grow

// generic interface for any object that's votable
type votable interface {
	Upvote() error
	Downvote() error
}
