package main

import (
	"fmt"
	"github.com/nekinci/paas/server"
	"github.com/nekinci/paas/specification"
	"net"
	"net/http"
	"os"
)

// The container engine will serve own. :)

func main() {

	fmt.Printf("My process Id: %d, My Parent Process Id: %d\n", os.Getpid(), os.Getppid())

	listener, err := net.Listen("tcp", "172.20.10.9:8090")
	go func() {
		conn, _ := listener.Accept()
		conn.Write([]byte("Hey I am alive"))
	}()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(("Hello world!")))
	})

	go http.ListenAndServe("localhost:4040", nil)

	stream, err := os.ReadFile("example.yml")
	if err != nil {
		panic(err)
	}

	app, err := specification.NewApplication(stream)

	if err != nil {
		panic(err)
	}

	_ = app
	s := server.NewServer()
	go s.ListenAndServe("127.0.0.1:7888")
	s.Handle(app)

	for {
	}
}
