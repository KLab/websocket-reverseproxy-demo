//package gorouter
package main

import (
    "flag"
    "fmt"
    "github.com/methane/rproxy"
    "github.com/garyburd/redigo/redis"
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

    if !strings.HasPrefix(path, "/chatroom/") {
        return
    }

    con := redisPool.Get()
    defer con.Close()
    /*
    if err != nil {
        log.Println("Can't connect to redis: ", err.String())
        req.URL.Host = backend
        return
    }
    */
    host, err := redis.String(con.Do("GET", path))
    if err == redis.ErrNil || host == "" {
        log.Println("Can't resolve for path: ", path)
        return
    }
    if err != nil {
        log.Println("Can't resolve for path: ", path)
        log.Println(" error=", err)
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
        MaxIdle: 3,
        IdleTimeout: 10 * time.Second,
        Dial: func() (conn redis.Conn, err error) {
            conn, err = redis.Dial("tcp", "127.0.0.1:6379")
            return
        },
    };

    proxy := &rproxy.ReverseProxy{Director: director}
    http.ListenAndServe(listenAddr, proxy)
}