
# Http-Server Library
This project is a simple and very limited HTTP Server library for Go. I built this to further my understanding of the HTTP protocol and the inner workings of HTTP server libraries.

Features:
* Provides a basic interface for the http server & routing
* Accepts HTTP/1.1 requests via a TCP connection
* Parses the HTTP/1.1 request
* Responds with valid HTTP/1.1 responses

# Usage
The http server requires a handler. The user can make their own handler as long as it conforms to the handler interface. The default handler is `http.Router`, this struct routes http requests to associated functions based on the requests path. Note that these functions must take a `http.Request` as the parameter and return a `http.Response`.

```
router := http.Router{}
server.RegisterRoute("GET", "/example/of/path", functionName1)
server.RegisterRoute("GET", "/different/path", functionName2)
```

Create the server using:

`server := http.CreateServer(router)`

Start the server by using the start function and providing a port (i.e. 80)

`err := server.Start(80)`

# Example Server
An example server using the library can be found under `example/main.go`. 

To run this example server use the command `go run ./example` from the project root directory.

The server will attempt to run on port `:8080` - if successful you will be able to interact with the server using `cUrl` or another HttpClient (i.e. Postman). Check the code in `example/main.go` to see which paths are setup.

# See for future extensions
* Caching - RFC 9111
* TCP implementation - [Beej's Guide to Network Programming](https://beej.us/guide/bgnet/)

# Resources
* Relevant RFC's are 9112 (HTTP/1.1 protocol) and 9110 (HTTP Semantics)
* Go standard library - `net/http`
* [Creating a mini HTTP server from "Scratch" (Python)](https://medium.com/@sakhawy/creating-an-http-server-from-scratch-ed41ef83314b)

