package grow

//  Base wrapper object used to determine what reddit object working with
type Thing struct {
	Id   string
	Name string
	Kind string
	Data interface{}
}
