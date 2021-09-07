package main

import (
	"bufio"
	"fmt"
	"github.com/nekinci/paas/container"
	"github.com/nekinci/paas/garbagecollector"
	"github.com/nekinci/paas/specification"
	"github.com/nekinci/paas/util"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main(){

	fmt.Printf("My process Id: %d, My Parent Process Id: %d\n", os.Getpid(), os.Getppid())
	stream, err := os.ReadFile("example.yml")
	if err != nil {
		panic(err)
	}

	app, err := specification.NewApplication(stream)

	if err != nil{
		panic(err)
	}

	ctx := container.NewContext()
	go garbagecollector.ScheduleCollect(ctx)
	newContainer := container.NewContainer(app, ctx)
	err = newContainer.Run()

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:7888")

	if err != nil {
		log.Fatalf("Application not started %v", err)
	}


	for  {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Connection not accepted : %v", err)
		}

		go func() {
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

			path := request.URL.Path
			path = path[1:]

			c := ctx.Get(path)

			if c == nil {
				noAvailable := util.NewNoAvailable()
				noAvailable.Write(conn)
				return
			}

			if c.Status != container.RUNNING {
				noLongerAvailable := util.NewNoLongerAvailable()
				noLongerAvailable.Write(conn)
				return
			}


			addr := c.IPV4 + ":" + strconv.Itoa(*c.BindingPort)
			dial, err := net.Dial("tcp", addr)

			go func() {
				io.Copy(conn, dial)
			}()

			io.Copy(dial, r2)

		}()
	}

}

//byteList, _ := ioutil.ReadAll(conn)
//
//go func() {
//	n, err:=conn.Write(byteList)
//	conn.Close()
//	//	n, err := io.Copy(conn, dial)
//	println(n)
//	if err != nil {
//		log.Printf("An error occurred: %v", err)
//	}
//}()
//
//_, err := io.Copy(dial, bytes.NewReader(byteList))
//if err != nil {
//log.Printf("An error occurred : %v", err)
//}