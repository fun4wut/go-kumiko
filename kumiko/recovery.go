package kumiko

import (
	"log"
)

func Recovery() HandlerFn {
	return func(ctx *Ctx) {
		// 错误恢复，捕获错误。
		defer func() {
			var err interface{}
			if err = recover(); err != nil {
				log.Println("Found Error", err)
			}
		}()
		ctx.Next()
	}
}
