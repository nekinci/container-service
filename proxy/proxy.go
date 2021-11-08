package proxy

import (
	"errors"
	"fmt"
	"github.com/nekinci/paas/application"
	"github.com/nekinci/paas/garbagecollector"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	InvalidHostHeaderError     error = errors.New("Invalid host header!\n")
	HostHeaderReadTimeoutError error = errors.New("While reading host header from connection occurred timeout\n")
)

type Proxy struct {
	ctx *application.Context
}

func NewServer(ctx *application.Context) Proxy {
	proxy := Proxy{ctx: ctx}
	return proxy
}

func (p Proxy) ListenAndServeL7(addr string) error {

	go garbagecollector.ScheduleCollect(p.ctx)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		hostName := request.Host
		println(hostName + request.URL.Path)
		app := p.ctx.Get(hostName)

		if app == nil {
			file, err := os.ReadFile("./resources/no-available.html")
			if err != nil {
				writer.Write([]byte("No available!"))
			} else {
				writer.Write(file)
			}
			return
		}
		serveProxy(fmt.Sprintf("http://0.0.0.0:%s", app.GetPort()), writer, request)
	})

	fmt.Printf("%v", http.ListenAndServe(addr, nil))
	return nil
}

func serveProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}

//func (proxy Proxy) ListenAndServeL4(addr string) error {
//
//	listen, err := net.Listen("tcp", addr)
//
//	if err != nil {
//		panic(err)
//		return err
//	}
//
//	for {
//		conn, err := listen.Accept()
//		if err != nil {
//			log.Printf("An error occurred while accepting connection: %v", err)
//		}
//		_ = conn
//
//		go func() {
//			defer conn.Close()
//
//			r1, w1 := net.Pipe()
//			r2, w2 := net.Pipe()
//			defer w1.Close()
//			defer w2.Close()
//			defer conn.Close()
//
//			_, _, _, _ = r1, w1, r2, w2
//
//			go func() {
//				written, err2 := io.Copy(io.MultiWriter(w2, os.Stdout), conn)
//				if err2 != nil {
//					log.Printf("err2:: %v", err)
//					return
//				}
//				log.Printf("MW::Written: %d", written)
//
//			}()
//
//
//			//hostHeader, err := proxy.GetHostWithTimeout(r1, 10)
//			if err != nil {
//				panic(err)
//			}
//			//hostHeader = getPrefix(hostHeader)
//
//			c := proxy.ctx.Get("frontend")
//			if c == nil {
//				return
//			}
//			port := c.GetPort()
//			dial, _ := net.Dial(c.GetProtocol(), "0.0.0.0"+":"+port)
//			defer dial.Close()
//
//			go func() {
//				io.Copy(dial, r2)
//			}()
//
//			io.Copy(conn, dial)
//			println("CIKIYOR:..")
//
//		}()
//	}
//
//	return nil
//}
//
//
//func (proxy *Proxy) GetHostWithTimeout(reader io.Reader, timeout int) (string, error) {
//	var error error = nil
//	var hostHeader string = ""
//	containerChan := make(chan bool)
//	go func() {
//		timer := time.NewTimer(time.Second * time.Duration(timeout))
//		<- timer.C
//		error = HostHeaderReadTimeoutError
//		containerChan <- true
//	}()
//
//	go func() {
//		h, err := proxy.GetHost(reader)
//		hostHeader = h
//		error = err
//		containerChan <- true
//	}()
//
//	<- containerChan
//	return hostHeader, error
//}
//
//func (proxy *Proxy) GetHost(reader io.Reader) (string, error) {
//	hostHeader := ""
//	scanner := bufio.NewScanner(reader)
//	for scanner.Scan() {
//		readiedString := scanner.Text()
//		if strings.Contains(readiedString, "Host:") || strings.Contains(readiedString, "Host :") {
//			hostHeader = strings.Replace(readiedString, "Host:", "", 1)
//			hostHeader = strings.Replace(hostHeader, "Host :", "", 1)
//			hostHeader = strings.TrimSpace(hostHeader)
//			break
//		}
//	}
//	if hostHeader == "" {
//		return "", InvalidHostHeaderError
//	}
//
//	return hostHeader, nil
//}
