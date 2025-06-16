package gfey

import "strings"

type node struct {
	pattern  string  // 待匹配的路由
	part     string  // 路由部分
	children []*node // node*型的切片
	isWild   bool    // 判断是否为精准匹配，part含有':'或'*'时为true
}

// matchChild: 匹配子节点的part, 返回第一个符合的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChild: 匹配子节点的part，返回所有符合的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert: 递归查找匹配的叶子节点，若不存在则新建路径节点，并在叶子节点保存完整的pattern
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归到叶子节点，保存完整pattern；只有叶子节点保存完整的pattern，中间节点的pattern为""
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]       // 当前层级的路径片段
	child := n.matchChild(part) // 查找是否已有对应子节点
	if child == nil {           // 没有则新建
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1) // 递归插入下一级
}

// search: 返回满足pattern的叶子节点，若不存在返回nil
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件：到达路径末尾或遇到通配符节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
