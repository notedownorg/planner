package markdown

import "strings"

// NodeType represents the type of markdown node
type NodeType string

const (
	NodeDocument  NodeType = "document"
	NodeHeading   NodeType = "heading"
	NodeParagraph NodeType = "paragraph"
	NodeTask      NodeType = "task"
	NodeList      NodeType = "list"
	NodeListItem  NodeType = "list_item"
	NodeText      NodeType = "text"
)

// Node is the base interface for all markdown nodes
type Node interface {
	Type() NodeType
	Children() []Node
	AddChild(Node)
	ClearChildren()
	Parent() Node
	SetParent(Node)
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	NodeType NodeType
	children []Node
	parent   Node
}

func (n *BaseNode) Type() NodeType { return n.NodeType }
func (n *BaseNode) Children() []Node { return n.children }
func (n *BaseNode) AddChild(child Node) {
	n.children = append(n.children, child)
	child.SetParent(n)
}
func (n *BaseNode) ClearChildren() {
	// Clear parent references from existing children
	for _, child := range n.children {
		child.SetParent(nil)
	}
	// Clear the children slice
	n.children = []Node{}
}

// AddChildNode is a helper that adds a child and sets parent correctly
func AddChildNode(parent Node, child Node) {
	parent.AddChild(child)
}
func (n *BaseNode) Parent() Node { return n.parent }
func (n *BaseNode) SetParent(parent Node) { n.parent = parent }

// Document represents the root of a markdown document tree
type Document struct {
	BaseNode
}

func NewDocument() *Document {
	return &Document{
		BaseNode: BaseNode{
			NodeType: NodeDocument,
			children: []Node{},
		},
	}
}

func (d *Document) AddChild(child Node) {
	d.children = append(d.children, child)
	child.SetParent(d)
}

// Heading represents an ATX heading node
type Heading struct {
	BaseNode
	Level int    // 1-6
	Title string
}

func NewHeading(level int, title string) *Heading {
	return &Heading{
		BaseNode: BaseNode{
			NodeType: NodeHeading,
			children: []Node{},
		},
		Level: level,
		Title: title,
	}
}

func (h *Heading) AddChild(child Node) {
	h.children = append(h.children, child)
	child.SetParent(h)
}

// Task represents a task/checkbox item
type Task struct {
	BaseNode
	Checked bool
	Content string
}

func NewTask(checked bool, content string) *Task {
	return &Task{
		BaseNode: BaseNode{
			NodeType: NodeTask,
			children: []Node{},
		},
		Checked: checked,
		Content: content,
	}
}

// Paragraph represents a paragraph of text
type Paragraph struct {
	BaseNode
	Content string
}

func NewParagraph(content string) *Paragraph {
	return &Paragraph{
		BaseNode: BaseNode{
			NodeType: NodeParagraph,
			children: []Node{},
		},
		Content: content,
	}
}

// List represents an unordered or ordered list
type List struct {
	BaseNode
	Ordered bool
}

func NewList(ordered bool) *List {
	return &List{
		BaseNode: BaseNode{
			NodeType: NodeList,
			children: []Node{},
		},
		Ordered: ordered,
	}
}

// ListItem represents an item in a list
type ListItem struct {
	BaseNode
	Content string
}

func NewListItem(content string) *ListItem {
	return &ListItem{
		BaseNode: BaseNode{
			NodeType: NodeListItem,
			children: []Node{},
		},
		Content: content,
	}
}

// Text represents raw text content
type Text struct {
	BaseNode
	Content string
}

func NewText(content string) *Text {
	return &Text{
		BaseNode: BaseNode{
			NodeType: NodeText,
			children: []Node{},
		},
		Content: content,
	}
}

// Utility functions

// FindHeadings recursively finds all heading nodes in the tree
func FindHeadings(node Node) []*Heading {
	var headings []*Heading
	
	if h, ok := node.(*Heading); ok {
		headings = append(headings, h)
	}
	
	for _, child := range node.Children() {
		headings = append(headings, FindHeadings(child)...)
	}
	
	return headings
}

// FindHeadingByTitle finds a heading by title (case-insensitive)
func FindHeadingByTitle(node Node, title string) *Heading {
	lowerTitle := strings.ToLower(title)
	headings := FindHeadings(node)
	
	for _, h := range headings {
		if strings.ToLower(h.Title) == lowerTitle {
			return h
		}
	}
	
	return nil
}

// FindTasks recursively finds all task nodes in the tree
func FindTasks(node Node) []*Task {
	var tasks []*Task
	
	if t, ok := node.(*Task); ok {
		tasks = append(tasks, t)
	}
	
	for _, child := range node.Children() {
		tasks = append(tasks, FindTasks(child)...)
	}
	
	return tasks
}

// GetNodeDepth returns the depth of a node in the tree (0 for root)
func GetNodeDepth(node Node) int {
	depth := 0
	current := node
	
	for current.Parent() != nil {
		depth++
		current = current.Parent()
	}
	
	return depth
}