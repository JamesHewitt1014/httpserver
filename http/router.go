package http

type Router map[route]HttpFunction

type route struct {
	Path   string
	Method string
}

type HttpFunction func(request *Request) *Response

// Registers a new route, if the route already exists it will be overwritten
func (r Router) RegisterRoute(method string, path string, fn HttpFunction) {
	newRoute := route{
		Path:   path,
		Method: method,
	}
	r[newRoute] = fn
}

//NOTE: I could do a check here to see if httpMethod is one of the allowed methods?

// Calls matching Handler function for the request
func (r Router) dispatch(request *Request) *Response {
	route := route{
		Path:   request.Path,
		Method: request.HttpMethod,
	}

	fn, exists := r[route]
	if !exists {
		return r.responseFromError(ERROR_PATH_NOT_FOUND)
	}

	response := fn(request)

	return response
}

// Creates a response from a Go error
func (r Router) responseFromError(err error) *Response {
	httpError := errorAsHttpError(err)
	message := []byte(httpError.Error())
	return CreateResponse(httpError.responseStatus, message)
}

var ERROR_PATH_NOT_FOUND = HttpError{
	responseStatus: StatusNotFound,
	errorMessage:   "Resource not found",
}

// TODO:
// Add support for params (these need to not be part of the path - should they be handled in the parsing or the routing?)
