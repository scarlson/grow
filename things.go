package grow

//  Base thing used to determine what we're working with
type Thing struct {
	Id   string
	Name string
	Kind string
	Data interface{}
}
