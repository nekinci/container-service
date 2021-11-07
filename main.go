package main

import (
	"fmt"
	"github.com/nekinci/paas/api"
	"github.com/nekinci/paas/application"
	"github.com/nekinci/paas/proxy"
	"github.com/nekinci/paas/specification"
	"io/ioutil"
	"os"
)

func main() {

	fmt.Println("Application starting...")
	fmt.Printf("Process Id: %d, Parent Process Id: %d\n", os.Getpid(), os.Getppid())

	context := application.NewContext()

	s := proxy.NewServer(context)
	channel := make(chan bool)
	go func() {
		s.ListenAndServeL7("0.0.0.0:443")
		channel <- true
	}()

	stream, err := ioutil.ReadFile("example.yml")
	if err != nil {
		panic(err)
	}

	tryItApp, specErr := specification.NewApplication(stream)
	tryItApp.Email = "superuser@containerdemo.live"
	if specErr != nil {
		panic(err)
	}

	context.Handle(tryItApp, false)

	go api.ListenAndServe(context)
	<-channel
}
