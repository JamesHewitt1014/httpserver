package http

type router struct {
	routes map[route]Handler
}
// Note: Router is setup like this so that it can be easily extracted from HttpServer if it needs to be decoupled later on

// type Router interface{
// 	dispatch() *Response
// }

type route struct {
	Path   string
	Method string
}

// This is a simple implementation that returns a Response struct
type Handler func(request *Request) *Response
// Note:
// Instead of returning a response, taking in a writer like the one in the Go std library would allow for text streaming responses.
// Example: type Handler func(writer io.Writer, request *Request)

func newRouter() *router{
	return &router{
		routes: map[route]Handler{},
	}
}

// Registers a new route, if the route already exists it will be overwritten
func (r *router) RegisterRoute(method string, path string, fn Handler) {
	newRoute := route{
		Path: path,
		Method: method,
	}
	r.routes[newRoute] = fn
}
//NOTE: I could do a check here to see if httpMethod is one of the allowed methods?

// Calls matching Handler function for the request
func (r *router) dispatch(request *Request) (*Response, error) {
	route := route{
		Path: request.Path,
		Method: request.HttpMethod,
	}

	fn, exists := r.routes[route]
	if !exists {
		return nil, ERROR_PATH_NOT_FOUND
	}

	response := fn(request)

	return response, nil
}

var ERROR_PATH_NOT_FOUND = HttpError{
		responseStatus: StatusNotFound,
		errorMsg: "Resource not found",
	}

// TODO:
// Add support for params (these need to not be part of the path - should they be handled in the parsing or the routing?)
