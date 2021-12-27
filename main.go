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
	{
		g1 := server.Group("/rua")
		g1.Get("/:id/faq", func(ctx *kumiko.Ctx) {
			v := func() string {
				if v, ok := ctx.GetPathParam("id"); ok {
					return v
				} else {
					return ""
				}
			}()
			ctx.WriteText(http.StatusOK, "221333"+v)
		})
	}

	log.Fatal(server.Run("0.0.0.0:4001"))
}
