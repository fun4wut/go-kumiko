package main

import (
	"awesomeProject/kotoha"
	"awesomeProject/kumiko"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	server := kumiko.New()
	server.Get("/", func(ctx *kumiko.Ctx) {
		log.Println("handling main")
		ctx.WriteJson(http.StatusOK, kumiko.H{
			"a": "2",
		})
	})
	kotoha.AddGroup("scores", 1<<10, kotoha.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	addr := "0.0.0.0:4001"
	peers := kotoha.NewHttpPool(addr)
	g1 := server.Group(peers.BasePath)
	g1.Get("/:groupname/:key", kotoha.HandleGet())
	log.Fatal(server.Run(peers.Self))
}
