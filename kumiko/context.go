package kumiko

import (
	"encoding/json"
	"net/http"
)

// H HashMap
type H map[string]interface{}

type Ctx struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newCtx(w http.ResponseWriter, req *http.Request) *Ctx {
	return &Ctx{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
func (c Ctx) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c Ctx) WriteStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c Ctx) WriteHeader(key string, val string) {
	c.Writer.Header().Set(key, val)
}
func (c Ctx) WriteJson(code int, obj interface{}) {
	c.WriteHeader("Content-Type", "application/json")
	c.WriteStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c Ctx) WriteText(code int, txt string) {
	c.WriteHeader("content-type", "text/plain")
	c.WriteStatus(code)
	_, err := c.Writer.Write([]byte(txt))
	if err != nil {
		return
	}
}
