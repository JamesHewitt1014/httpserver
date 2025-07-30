package http

import (
	"fmt"
	"net"
)

type HttpServer struct {
	handler 	 Handler
	listener net.Listener
	open     bool
}
// NOTE: open represents an actively listening server - could use an atomic bool or maybe even a channel instead if worried about race conditions

type Handler interface{
	dispatch(*Request) *Response
	responseFromError(error) *Response
}
// Note: Instead of returning a response, taking in a writer like the one in the Go std library would allow for text streaming responses. - Could add a write func?

func CreateServer(handler Handler) *HttpServer {
	return &HttpServer{
		handler:	  handler,
		listener: nil,
		open:     false,
	}	
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
		response = s.handler.responseFromError(err)
	} else {
		response = s.handler.dispatch(request)
	}

	response.Write(connection)
}

func (s *HttpServer) Close() error {
	err := s.listener.Close()
	s.open = false
	return err
}
