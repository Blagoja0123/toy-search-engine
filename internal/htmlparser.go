package internal

import (
	"fmt"
	"strings"
)

type htmlTreeNode struct {
	tag      string
	content  []string
	children []*htmlTreeNode
	parent   *htmlTreeNode
}

func newHTMLNode(tag string, content []string, parent *htmlTreeNode) *htmlTreeNode {
	return &htmlTreeNode{
		tag:      tag,
		content:  content,
		children: []*htmlTreeNode{},
		parent:   parent,
	}
}

func (n *htmlTreeNode) addChild(child *htmlTreeNode) {
	n.children = append(n.children, child)
}

type tree struct {
	root *htmlTreeNode
}

func newHTMLTree() *tree {
	return &tree{nil}
}

func (t *tree) setRoot(node *htmlTreeNode) {
	t.root = node
}

func tokenize(n *htmlTreeNode) []string {
	var tokens []string
	//if len(n.children) == 0 {
	//	return []string{""}
	//}
	if n != nil {

		var fmtContent []string
		for _, val := range n.content {
			tmpVal := strings.Split(val, " ")
			fmtContent = append(fmtContent, tmpVal...)
		}
		tokens = append(tokens, fmtContent...)
		for _, child := range n.children {
			tokens = append(tokens, tokenize(child)...)
		}
	}
	return tokens
}

func printHTMLTree(node *htmlTreeNode, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "\t"
	}

	fmt.Printf("%s<%s>%s\n", indent, node.tag, node.content)

	for _, child := range node.children {
		printHTMLTree(child, depth+1)
	}

	fmt.Printf("%s</%s>\n", indent, node.tag)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (n *htmlTreeNode) addContent(content string) {
	n.content = append(n.content, content)
}

func stripFile(current *htmlTreeNode, htmlFile []string, index int) {
	line := htmlFile[index]
	if line == "</body>" {
		return
	}
	if line[1] == '/' {
		stripFile(current.parent, htmlFile, index+1)
	}
	if line[0] != '<' {
		current.addContent(line)
		stripFile(current, htmlFile, index+1)

	}
	if line[0] == '<' && line[1] != '/' && line[1] != '!' {
		tag := strings.ReplaceAll(line, "<", "")
		tag = strings.ReplaceAll(tag, ">", "")
		tempNode := newHTMLNode(tag, []string{""}, current)
		current.addChild(tempNode)
		stripFile(tempNode, htmlFile, index+1)
	}
}

func formatFile(htmlFile []string) []string {
	var res []string
	for i := range htmlFile {
		tempStr := htmlFile[i]
		tempStr = strings.Replace(tempStr, "\t", "", -1)
		tempStr = strings.TrimSpace(tempStr)
		res = append(res, tempStr)
	}
	return res
}

func mapTokens(t *tree) map[string]int {
	mp := make(map[string]int)

	tokens := tokenize(t.root)
	for _, line := range tokens {
		token := strings.Fields(line)
		for _, indToken := range token {
			if _, ok := mp[indToken]; ok {
				mp[indToken]++
			} else {
				mp[indToken] = 0
			}
		}
	}
	return mp
}
