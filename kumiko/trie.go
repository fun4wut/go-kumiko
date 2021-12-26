package kumiko

import "strings"

type trieNode struct {
	pattern   string               // 叶子节点，待匹配的完整路由 如 /doc/p/rua，如果非叶子节点则为空字符串
	part      string               // 当前的路由部分, 如 p
	children  map[string]*trieNode // 子节点 part -> trieNode的映射
	wildNodes []*trieNode
	isWild    bool
}

func (root *trieNode) find(parts []string, height int) *trieNode {
	if len(parts) == height || strings.HasPrefix(root.part, "*") {
		// 如果pattern为空字符串，说明该节点非终结节点，查找失败
		if root.pattern == "" {
			return nil
		}
		return root
	}
	part := parts[height]
	// 因为有通配符匹配的存在，所以在搜索时，需要找到所有匹配的节点，并对所有匹配的节点都做一遍遍历
	nodes := append(root.wildNodes)
	if n, ok := root.children[part]; ok {
		nodes = append(nodes, n)
	}
	// 通配符节点+符合的普通节点就是要遍历的子节点了
	for _, n := range nodes {
		if res := n.find(parts, height+1); res != nil {
			return res
		}
	}
	return nil
}

func (root *trieNode) insert(parts []string, height int) {
	// 匹配到最后一段了，该root就是正确的位置
	if len(parts) == height {
		root.pattern = "/" + strings.Join(parts, "/")
		return
	}

	part := parts[height]
	// 插入时不需要考虑通配符，所以查到一个就行
	node, ok := root.children[part]
	if !ok {
		node = &trieNode{
			part:      part,
			children:  map[string]*trieNode{},
			wildNodes: []*trieNode{},
			isWild:    strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*"),
		}
		root.children[part] = node
		if node.isWild {
			root.wildNodes = append(root.wildNodes, node)
		}
	}
	node.insert(parts, height+1)
}
