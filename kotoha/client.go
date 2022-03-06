package kotoha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// http 客户端，负责往别的peer发送请求，并获取结果
type httpGetter struct {
	baseUrl string
}

func (h httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v/%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)
