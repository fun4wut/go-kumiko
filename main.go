package main

import (
	"awesomeProject/kotoha"
	"awesomeProject/kumiko"
	"flag"
	"fmt"
	"log"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *kotoha.Group {
	return kotoha.AddGroup("scores", 1<<11, kotoha.GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startApiServer(apiAddr string, cacheGroup *kotoha.Group) {
	server := kumiko.New()
	server.Get("/api", func(ctx *kumiko.Ctx) {
		key := ctx.GetQuery("key")
		log.Println("[API] handling api, key is", key)
		view, _ := cacheGroup.Get(key)
		ctx.WriteHeader("Content-Type", "application/octet-stream")
		ctx.Writer.Write(view.ByteSlice())
	})
	log.Println("api run at", apiAddr)
	log.Fatal(server.Run(apiAddr[7:]))
}

func startCacheServer(addr string, remoteAddrs []string, cacheGroup *kotoha.Group) {
	server := kumiko.New()
	peers := kotoha.NewHttpPool(addr)
	peers.AddPeer(remoteAddrs...)
	cacheGroup.RegisterPeerPicker(peers)
	g1 := server.Group(peers.BasePath)
	g1.Get("/:groupname/:key", kotoha.HandleGet())
	log.Println("cache server run at", addr)
	log.Fatal(server.Run(addr[7:]))
}

func main() {
	var port int
	var isApi bool
	flag.IntVar(&port, "p", 8001, "cache server port")
	flag.BoolVar(&isApi, "api", false, "is api server")
	flag.Parse()
	cacheGroup := createGroup()

	apiAddr := "http://0.0.0.0:9999"
	addrMap := map[int]string{
		8001: "http://0.0.0.0:8001",
		8002: "http://0.0.0.0:8002",
		8003: "http://0.0.0.0:8003",
	}
	var addrs []string
	for _, a := range addrMap {
		addrs = append(addrs, a)
	}

	if isApi {
		go startApiServer(apiAddr, cacheGroup)
	}
	startCacheServer(addrMap[port], addrs, cacheGroup)
}
