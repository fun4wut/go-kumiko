package kumiko

import "fmt"

func Recovery() HandlerFn {
	return func(ctx *Ctx) {
		// 错误恢复，捕获错误。
		defer func() {
			if p := recover(); &p != nil {
				fmt.Println("FoundErr", p)
			}
		}()
		ctx.Next()
	}
}
