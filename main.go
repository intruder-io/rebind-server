package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	flag "github.com/spf13/pflag"
)

// https://stackoverflow.com/a/33881296
var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		// Also set "Connection: close"
		w.Header().Set("Connection", "close")

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {
	// Flags
	var port = flag.IntP("port", "p", 8080, "port to listen on")
	var assetsDir = flag.StringP("assets-dir", "a", "./assets", "assets directory")
	flag.Parse()

	var server http.Server

	// Create a static fileserver with 1 API endpopint
	mux := http.NewServeMux()
	mux.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Shutting down server")
		server.Shutdown(context.Background())
	})
	fs := http.FileServer(http.Dir(*assetsDir))
	mux.Handle("/", NoCache(fs))
	server = http.Server{
		Handler: mux,
	}
	listener, err := net.Listen("tcp6", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	server.Serve(listener)
}
