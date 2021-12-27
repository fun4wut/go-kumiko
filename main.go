package main

import (
	"awesomeProject/kumiko"
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := kumiko.New()
	server.Use(func(ctx *kumiko.Ctx) {
		fmt.Println("part1")
		ctx.Next()
		fmt.Println("part4")
	})
	server.Use(func(ctx *kumiko.Ctx) {
		fmt.Println("part2")
		ctx.Next()
		fmt.Println("part3")
	})
	server.Get("/", func(ctx *kumiko.Ctx) {
		fmt.Println("handling main")
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
