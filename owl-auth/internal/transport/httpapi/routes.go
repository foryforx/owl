package httpapi

type Router interface {
	Routes() []Route
}

type Route struct {
	Method  string
	Path    string
	Handler Handler
	IsPubic bool
}
