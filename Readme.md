
# Http-Server Library
This project is a simple and very limited Http Server library for Go. I built this to further my understanding of the HTTP protocol and of how HTTP server libraries work.

Features:
* Provides a basic interface for the http server & routing
* Accepts HTTP/1.1 requests via a TCP connection
* Parses the HTTP/1.1 request
* Responds with valid HTTP/1.1 responses

# Usage
Create server using:
`server := http.CreateServer()`
Register as many routes you want - using the HttpMethod, request path, and the function you wish to apply to that route
Note that this function must take a `http.Request` as its parameter and return a `http.Response`.
`server.RegisterRoute("GET", "/example/of/path", functionName)`
Start the server by using the start function and providing a port (i.e. 80)
`err := server.Start(80)`

# Example Server
An example server using the library can be found under `example/main.go`. 

To run this example server run `go run ./example` for the project root directory.

The server will attempt to run on port `:8080` - if successful you will be able to interact with the server using `cUrl` or another HttpClient (i.e. Postman). Check the code in `example/main.go` to see which paths are setup.

# See for future extensions
* Caching - RFC 9111
* TCP implementation - [Beej's Guide to Network Programming](https://beej.us/guide/bgnet/)

# Resources
* Relevant RFC's are 9112 (HTTP/1.1 protocol) and 9110 (HTTP Semantics)
* Go standard library - `net/http`
* [Creating a mini HTTP server from "Scratch" (Python)](https://medium.com/@sakhawy/creating-an-http-server-from-scratch-ed41ef83314b)

