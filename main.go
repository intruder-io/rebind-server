package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
)

// BlockListListener is a net.Listener that blocks connections from a list of IPs. New IPs to block can be added with the Block method.
type BlockListListener struct {
	listener       net.Listener
	blockTime      time.Duration
	blockList      map[string]time.Time
	blockListMutex sync.RWMutex
}

func NewBlockListListener(network, addr string, blockTime time.Duration) (*BlockListListener, error) {
	l, err := net.Listen(network, addr)
	ret := &BlockListListener{
		listener:       l,
		blockTime:      blockTime,
		blockList:      make(map[string]time.Time),
		blockListMutex: sync.RWMutex{},
	}
	return ret, err
}

func (l *BlockListListener) Accept() (net.Conn, error) {
	for {
		conn, err := l.listener.Accept()
		addr := conn.RemoteAddr().(*net.TCPAddr).IP.String()
		if t, blocked := l.blockList[addr]; blocked && t.After(time.Now()) {
			conn.Close()
			continue
		}
		return conn, err
	}
}

func (l *BlockListListener) Close() error {
	return l.listener.Close()
}

func (l *BlockListListener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *BlockListListener) Block(addr string) {
	l.blockListMutex.Lock()
	defer l.blockListMutex.Unlock()
	l.blockList[addr] = time.Now().Add(l.blockTime)
}

func main() {
	// Flags
	var port = flag.IntP("port", "p", 8080, "port to listen on")
	var blockTime = flag.DurationP("block-time", "b", 5*time.Second, "time to block")
	var assetsDir = flag.StringP("assets-dir", "a", "./assets", "assets directory")
	flag.Parse()

	// Listener
	listener, err := NewBlockListListener("tcp6", fmt.Sprintf(":%d", *port), *blockTime)
	if err != nil {
		log.Fatal(err)
	}

	// Create a static fileserver with 1 API endpopint
	mux := http.NewServeMux()
	mux.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		addr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Blocking %s\n", addr)
		listener.Block(addr)
	})
	fs := http.FileServer(http.Dir(*assetsDir))
	mux.Handle("/", fs)
	server := http.Server{
		Handler: mux,
	}

	server.Serve(listener)
}
