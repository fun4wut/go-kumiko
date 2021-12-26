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
			"a": "2",
		})
	})
	server.Get("/rua/:id/eqqq", func(ctx *kumiko.Ctx) {
		v := func() string {
			if v, ok := ctx.GetPathParam("id"); ok {
				return v
			} else {
				return ""
			}
		}()
		ctx.WriteText(http.StatusOK, "221333"+v)
	})

	log.Fatal(server.Run("0.0.0.0:4001"))
}
