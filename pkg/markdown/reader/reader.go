package reader

import (
	"regexp"
	"strings"

	"github.com/notedownorg/planner/pkg/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var taskRegex = regexp.MustCompile(`^\s*-\s*\[([ xX])\]\s*(.*)$`)

// ParseMarkdown parses markdown content and builds a tree structure
func ParseMarkdown(content string) (*markdown.Document, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.TaskList),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	source := []byte(content)
	
	astDoc := md.Parser().Parse(text.NewReader(source))
	
	// Create the document root
	doc := markdown.NewDocument()
	
	// Build the tree by walking the AST
	err := buildTree(doc, astDoc, source)
	
	return doc, err
}

// buildTree recursively builds our tree from the goldmark AST
func buildTree(parent markdown.Node, astNode ast.Node, source []byte) error {
	var headingStack []*markdown.Heading
	
	// Process each child of the current AST node
	for child := astNode.FirstChild(); child != nil; child = child.NextSibling() {
		node, err := convertASTNode(child, source)
		if err != nil {
			return err
		}
		
		if node != nil {
			//fmt.Printf("Created node: %T\n", node)
			switch n := node.(type) {
			case *markdown.Heading:
				// Pop headings from stack that are at same or higher level
				for len(headingStack) > 0 && headingStack[len(headingStack)-1].Level >= n.Level {
					headingStack = headingStack[:len(headingStack)-1]
				}
				
				// Add to appropriate parent
				if len(headingStack) > 0 {
					headingStack[len(headingStack)-1].AddChild(n)
				} else {
					parent.AddChild(n)
				}
				
				// Push to stack
				headingStack = append(headingStack, n)
				
			default:
				// Add to current heading if exists, otherwise to parent
				if len(headingStack) > 0 {
					headingStack[len(headingStack)-1].AddChild(node)
				} else {
					parent.AddChild(node)
				}
			}
			
			// Recursively process children only for non-text nodes
			if _, isText := node.(*markdown.Text); !isText {
				if err := buildTree(node, child, source); err != nil {
					return err
				}
			}
		} else {
			// If we didn't create a node for this AST node, process its children
			// directly under the current parent
			if err := buildTree(parent, child, source); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// convertASTNode converts a goldmark AST node to our markdown node
func convertASTNode(astNode ast.Node, source []byte) (markdown.Node, error) {
	// Debug: print the AST node type
	//fmt.Printf("Converting AST node: %T\n", astNode)
	
	switch node := astNode.(type) {
	case *ast.Heading:
		title := extractNodeText(node, source)
		return markdown.NewHeading(node.Level, title), nil
		
	case *ast.Paragraph:
		// Check if this paragraph contains a task
		text := extractNodeText(node, source)
		if task := parseTask(text); task != nil {
			return task, nil
		}
		return markdown.NewParagraph(text), nil
		
	case *ast.List:
		ordered := node.IsOrdered()
		return markdown.NewList(ordered), nil
		
	case *ast.ListItem:
		// Check if this list item contains a task checkbox
		hasTaskCheckBox := false
		var isChecked bool
		
		// Walk through all descendants to find TaskCheckBox
		ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if entering {
				if taskBox, ok := n.(*extast.TaskCheckBox); ok {
					hasTaskCheckBox = true
					isChecked = taskBox.IsChecked
					return ast.WalkStop, nil
				}
			}
			return ast.WalkContinue, nil
		})
		
		// Extract the content (excluding the checkbox)
		content := extractListItemText(node, source)
		
		if hasTaskCheckBox {
			return markdown.NewTask(isChecked, content), nil
		}
		
		return markdown.NewListItem(content), nil
		
	case *extast.TaskCheckBox:
		// These are handled in the ListItem case
		return nil, nil
		
	case *ast.TextBlock:
		// Handle text blocks
		content := extractNodeText(node, source)
		return markdown.NewText(content), nil
		
	case *ast.Text:
		// Skip text nodes that are direct children of headings
		// as the heading text is already extracted in the heading creation
		if parent := astNode.Parent(); parent != nil {
			if _, ok := parent.(*ast.Heading); ok {
				return nil, nil
			}
		}
		content := string(node.Segment.Value(source))
		return markdown.NewText(content), nil
		
	default:
		// For other node types, we'll process their children but not create a node
		return nil, nil
	}
}

// extractNodeText extracts all text content from a node
func extractNodeText(node ast.Node, source []byte) string {
	var text strings.Builder
	
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if textNode, ok := n.(*ast.Text); ok {
				text.Write(textNode.Segment.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})
	
	return strings.TrimSpace(text.String())
}

// extractListItemText extracts text content from a list item, excluding task checkbox
func extractListItemText(node ast.Node, source []byte) string {
	var text strings.Builder
	
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			// Skip TaskCheckBox nodes
			if _, ok := n.(*extast.TaskCheckBox); ok {
				return ast.WalkSkipChildren, nil
			}
			if textNode, ok := n.(*ast.Text); ok {
				text.Write(textNode.Segment.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})
	
	return strings.TrimSpace(text.String())
}

// parseTask checks if text is a task and returns a Task node if it is
func parseTask(text string) *markdown.Task {
	matches := taskRegex.FindStringSubmatch(text)
	if len(matches) == 3 {
		checked := matches[1] != " "
		content := matches[2]
		return markdown.NewTask(checked, content)
	}
	return nil
}

