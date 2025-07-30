package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"httpserver/http"
)


func main() {
	router := http.Router{}
	router.RegisterRoute("GET", "/test", method1)
	router.RegisterRoute("PUT", "/example/2", method2)
	router.RegisterRoute("GET", "/example/2", method3)

	server := http.CreateServer(router)
	const port int = 8080
	err := server.Start(port)
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
	defer server.Close()
	log.Println("Server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Program shutdown")
}

func method1(request *http.Request) *http.Response {
	request.Print()
	content := "Hello, Handler!"
	response := http.CreateResponse(http.StatusOk, []byte(content))
	response.Print()
	return response
}

func method2(request *http.Request) *http.Response {
	request.Print()
	res := http.Ok()
	res.Print()
	return res
}

func method3(request *http.Request) *http.Response {
	request.Print()
	content := `<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
</head>
<body>
	<h1>Hello World!</h1> 
	<p>This is some content.</p>
</body>
</html>`
	response := http.CreateResponse(http.StatusOk, []byte(content))
	response.SetHeader("Content-Type", "text/html")
	response.Print()
	return response
}
