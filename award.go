package grow

// Type: t6_, reddit award wrapper
type awardThing struct {
	Id   string
	Name string
	Kind string
	Data Award
}

// data payload from award wrapper
type Award struct {
	Id   interface{}
	Name interface{}
}
