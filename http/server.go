package http

import (
	"fmt"
	"net"
)

type HttpServer struct {
	router 	 *router
	listener net.Listener
	open     bool
}
// NOTE: 
// 1. open represents an actively listening server - could use an atomic bool or maybe even a channel instead if worried about race conditions
// 2. Could add an insertable error handler to HttpServer: for setting how default errors are handled (i.e. change content-type to html, and display "500 server error" page)
// 3. Routing is currently limited, an alternative might be to make Router an interface and to allow end-users to define routing behaviour. This would probably require it to be decoupled more from the HttpServer struct.

func CreateServer() *HttpServer {
	router := newRouter()
	return &HttpServer{
		router:	  router,
		listener: nil,
		open:     false,
	}	
}

// Registers a new route, if the route already exists it will be overwritten
func (s *HttpServer) RegisterRoute(method string, path string, fn Handler){
	s.router.RegisterRoute(method, path, fn)	
}

// Starts the HttpServer. It will begin listening for Http Requests
func (s *HttpServer) Start(port int) (error) {
	address := fmt.Sprintf(":%d", port)
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listener = tcpListener
	s.open = true

	go s.listen()

	return nil
}

func (s *HttpServer) listen() {
	defer s.Close()

	for s.open {
		connection, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(connection)
	}
}

func (s *HttpServer) handleConnection(connection net.Conn) {
	defer connection.Close()

	var response *Response

	request, err := requestFromStream(connection)
	if err != nil {
		response = ResponseFromError(err)
		response.Write(connection)
		return
	}

	response, err = s.router.dispatch(request)
	if err != nil {
		response = ResponseFromError(err)
	}

	response.Write(connection)
}

func (s *HttpServer) Close() error {
	err := s.listener.Close()
	s.open = false
	return err
}
