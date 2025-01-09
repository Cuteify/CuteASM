package parser

type Block interface {
	//Parse(p *Parser)
}

type Node struct {
	Value    Block
	Father   *Node
	Children []*Node
}

func (n *Node) AddChild(node *Node) {
	n.Children = append(n.Children, node)
	node.Father = n
}
