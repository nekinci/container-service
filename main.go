package main

import (
	"fmt"
	"github.com/nekinci/paas/api"
	"github.com/nekinci/paas/application"
	"github.com/nekinci/paas/proxy"
	"os"
)

func main() {

	fmt.Println("Application starting...")
	fmt.Printf("Process Id: %d, Parent Process Id: %d\n", os.Getpid(), os.Getppid())

	context := application.NewContext()

	s := proxy.NewServer(context)
	channel := make(chan bool)
	go func() {
		s.ListenAndServeL7("127.0.0.1:7888")
		channel <- true
	}()

	go api.ListenAndServe(context)
	<-channel
}
