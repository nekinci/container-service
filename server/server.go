package server

import (
	"bufio"
	"github.com/nekinci/paas/application"
	"github.com/nekinci/paas/specification"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	ctx application.Context
}

func NewServer() Server {
	server := Server{
		ctx: *application.NewContext(),
	}
	return server
}

func (server Server) ListenAndServe(addr string) error {

	listen, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("An error occurred while accepting connection: %v", err)
		}
		_ = conn

		go func() {
			defer conn.Close()
			r1, w1 := io.Pipe()
			r2, w2 := io.Pipe()

			go func() {
				io.Copy(io.MultiWriter(w1, w2), conn)
			}()

			request, err := http.ReadRequest(bufio.NewReader(r1))
			if err != nil {
				// No Http
				log.Fatalf("%v", err)
			}

			path := request.Host + request.URL.Path

			c := server.ctx.Get(path)
			if c == nil {
				return
			}

			port := c.GetPort()
			dial, err := net.Dial(c.GetProtocol(), "0.0.0.0"+":"+port)

			go func() {
				io.Copy(conn, dial)
			}()

			io.Copy(dial, r2)

		}()
	}

	return nil
}

func (server Server) Handle(app *specification.Specification) {
	server.ctx.RunApplication(*app)
}
