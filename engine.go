package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	port       string
	socketPath string
	ip         string
)

func init() {
	flag.StringVar(&port, "p", "", "port to serve on")
	flag.StringVar(&socketPath, "s", "", "unix socket to proxy requests")
	flag.StringVar(&ip, "i", "", "ip to listen on")
	flag.Parse()
}

func writeError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// Proxy request from the http server to docker's
// unix socket
func proxy(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c := httputil.NewClientConn(conn, nil)
	defer c.Close()

	res, err := c.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer res.Body.Close()

	copyHeader(w.Header(), res.Header)
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Println(err)
	}
}

func main() {
	var err error
	if port == "" {
		err = fmt.Errorf("Please provide a port to serve on")
	} else if socketPath == "" {
		err = fmt.Errorf("Please provide a socket to proxy to")
	}
	if err != nil {
		writeError(err)
	}

	http.HandleFunc("/", proxy)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil); err != nil {
		writeError(err)
	}
}
