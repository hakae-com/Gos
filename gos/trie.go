package gos

import "strings"

type node struct {
	pattern  string  //待匹配路由
	part     string  //路由中的一部分
	children []*node //子节点
	isWild   bool    //是否精准匹配
}

func (m *node) matchChild(part string) *node {
	for _, child := range m.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (m *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range m.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (m *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		m.pattern = pattern
		return
	}
	part := parts[height]
	child := m.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: (part[0] == ':' || part[0] == '*')}
		m.children = append(m.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (m *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(m.part, "*") {
		if m.pattern == "" {
			return nil
		}
		return m
	}
	part := parts[height]
	children := m.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
