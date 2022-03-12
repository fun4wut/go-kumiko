package kotoha

// PeerGetter 向远程机器通过http获取数据的抽象
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

// PeerPicker 根据key值，找到对应的机器
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}
