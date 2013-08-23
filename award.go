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

// build Award's fullname
func (le Award) Fullname() string {
	return le.Name.(string)
}
