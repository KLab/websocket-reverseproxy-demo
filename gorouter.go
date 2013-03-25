//package gorouter
package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/methane/rproxy"
	"log"
	"net/http"
	"strings"
	"time"
)

var redisPool *redis.Pool
var backend string

// リクエストを振り分ける関数.
func director(req *http.Request) {
	path := req.URL.Path
	req.URL.Scheme = "http"
	req.URL.Host = backend

	if !strings.HasPrefix(path, "/chatsocket/") {
		return
	}

	con := redisPool.Get()
	defer con.Close()

	host, err := redis.String(con.Do("GET", path))
	if err != nil || host == "" {
		log.Println("Can't resolve for path: ", path)
		if err != redis.ErrNil {
			log.Println(" error=", err)
		}
		return
	}
	log.Println("Resolved backend: ", host, path)
	req.URL.Host = host
}

func main() {
	port := flag.Int("port", 8090, "port number")
	flag.StringVar(&backend, "backend", "127.0.0.1:8888", "Default backend")
	flag.Parse()
	listenAddr := fmt.Sprintf(":%v", *port)
	log.Println("Listening on: ", listenAddr)

	redisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 10 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", "127.0.0.1:6379")
			return
		},
	}

	proxy := &rproxy.ReverseProxy{Director: director}
	http.ListenAndServe(listenAddr, proxy)
}
