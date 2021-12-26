package main

import (
	"awesomeProject/kumiko"
	"log"
	"net/http"
)

func main() {
	server := kumiko.New()
	server.Get("/", func(ctx *kumiko.Ctx) {
		ctx.WriteJson(http.StatusOK, kumiko.H{
			"a": "1",
		})
	})
	log.Fatal(server.Run("0.0.0.0:4001"))
}
