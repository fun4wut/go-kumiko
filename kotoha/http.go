package kotoha

import (
	"awesomeProject/kumiko"
	"log"
	"net/http"
)

const defaultBasePath = "/__kotoha"

type HttpPool struct {
	Self     string
	BasePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		Self:     self,
		BasePath: defaultBasePath,
	}
}

func HandleGet() kumiko.HandlerFn {
	return func(ctx *kumiko.Ctx) {
		log.Println("Hit path")
		groupName, _ := ctx.GetPathParam("groupname")
		key, _ := ctx.GetPathParam("key")
		group := GetGroup(groupName)
		v, err := group.Get(key)
		if err != nil {
			ctx.WriteText(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
		ctx.Writer.Write(v.ByteSlice())
	}
}
