package writer

import (
	"fmt"
	"strings"

	"github.com/notedownorg/planner/pkg/markdown"
)

// WriteDocument converts a Document tree back to markdown format
func WriteDocument(doc *markdown.Document) string {
	var builder strings.Builder
	writeNode(&builder, doc, 0)
	return strings.TrimSpace(builder.String())
}

// writeNode recursively writes a node and its children
func writeNode(builder *strings.Builder, node markdown.Node, depth int) {
	switch n := node.(type) {
	case *markdown.Document:
		// Write all children
		for i, child := range n.Children() {
			if i > 0 {
				builder.WriteString("\n")
			}
			writeNode(builder, child, depth)
		}
		
	case *markdown.Heading:
		// Write the heading
		WriteHeading(builder, n.Level, n.Title)
		
		// Write children with proper spacing
		for i, child := range n.Children() {
			// Add spacing logic based on previous and current child types
			if i > 0 {
				prev := n.Children()[i-1]
				
				// Check if previous sibling ends with a task
				prevEndsWithTask := false
				if prevHeading, ok := prev.(*markdown.Heading); ok {
					if len(prevHeading.Children()) > 0 {
						lastChild := prevHeading.Children()[len(prevHeading.Children())-1]
						if _, isTask := lastChild.(*markdown.Task); isTask {
							prevEndsWithTask = true
						}
					}
				} else if _, isTask := prev.(*markdown.Task); isTask {
					prevEndsWithTask = true
				}
				
				// Determine spacing based on node types and previous ending
				if _, isPrevTask := prev.(*markdown.Task); isPrevTask {
					if _, isCurrentTask := child.(*markdown.Task); isCurrentTask {
						// Consecutive tasks, don't add extra spacing (task already has \n)
					} else if _, isCurrentHeading := child.(*markdown.Heading); isCurrentHeading {
						// Task to heading, add one extra spacing
						builder.WriteString("\n")
					} else {
						// Task to other node, no extra spacing (task already has \n)
					}
				} else {
					// Previous was not a task, but check if it ended with a task
					if _, isCurrentHeading := child.(*markdown.Heading); isCurrentHeading {
						if prevEndsWithTask {
							// Previous section ended with task, only add one newline
							builder.WriteString("\n")
						} else {
							// Normal spacing between headings
							builder.WriteString("\n\n")
						}
					} else {
						builder.WriteString("\n")
					}
				}
			} else {
				builder.WriteString("\n")
			}
			writeNode(builder, child, depth+1)
		}
		
	case *markdown.Paragraph:
		builder.WriteString(n.Content)
		
		// Write children
		for _, child := range n.Children() {
			writeNode(builder, child, depth+1)
		}
		
	case *markdown.Task:
		// Write task without extra indentation
		if n.Checked {
			builder.WriteString("- [x] ")
		} else {
			builder.WriteString("- [ ] ")
		}
		builder.WriteString(n.Content)
		builder.WriteString("\n")
		
		// Don't write children for tasks as they contain duplicate content
		
	case *markdown.List:
		// Write children (list items)
		for i, child := range n.Children() {
			if i > 0 {
				// No extra spacing between list items
			}
			writeNode(builder, child, depth+1)
		}
		
	case *markdown.ListItem:
		// Write list item with proper indentation
		if depth > 0 {
			builder.WriteString(strings.Repeat("  ", depth-1))
		}
		
		// Check if parent is ordered
		if parent, ok := n.Parent().(*markdown.List); ok && parent.Ordered {
			// For ordered lists, we'd need to track the index
			builder.WriteString("1. ")
		} else {
			builder.WriteString("- ")
		}
		
		builder.WriteString(n.Content)
		builder.WriteString("\n")
		
		// Write children
		for _, child := range n.Children() {
			writeNode(builder, child, depth+1)
		}
		
	case *markdown.Text:
		builder.WriteString(n.Content)
		
		// Write children (though text nodes typically don't have children)
		for _, child := range n.Children() {
			writeNode(builder, child, depth+1)
		}
	}
}

// needsSpacing determines if a node needs spacing before it
func needsSpacing(node markdown.Node) bool {
	switch node.(type) {
	case *markdown.Heading, *markdown.Paragraph, *markdown.List:
		return true
	default:
		return false
	}
}

// needsExtraSpacing determines if extra spacing is needed between two consecutive nodes
func needsExtraSpacing(prev, current markdown.Node) bool {
	// Add extra spacing between tasks and headings
	if _, isPrevTask := prev.(*markdown.Task); isPrevTask {
		if _, isCurrentHeading := current.(*markdown.Heading); isCurrentHeading {
			return true
		}
	}
	return false
}

// WriteHeading writes an ATX heading to the builder
func WriteHeading(builder *strings.Builder, level int, title string) {
	if level < 1 || level > 6 {
		level = 1 // Default to H1 if invalid level
	}
	
	// Write the # characters
	builder.WriteString(strings.Repeat("#", level))
	builder.WriteString(" ")
	builder.WriteString(title)
	builder.WriteString("\n")
}

// WriteNode writes a single node and its children to markdown format
func WriteNode(node markdown.Node) string {
	var builder strings.Builder
	writeNode(&builder, node, 0)
	return strings.TrimSpace(builder.String())
}

// CreateHeadingWithContent creates a heading with content children
func CreateHeadingWithContent(level int, title string, children ...markdown.Node) *markdown.Heading {
	heading := markdown.NewHeading(level, title)
	for _, child := range children {
		heading.AddChild(child)
	}
	return heading
}

// AddNodeToHeading adds a node as a child of a heading
func AddNodeToHeading(heading *markdown.Heading, node markdown.Node) {
	heading.AddChild(node)
}

// CreateTaskList creates a list of tasks
func CreateTaskList(tasks map[string]bool) *markdown.List {
	list := markdown.NewList(false)
	
	for content, checked := range tasks {
		task := markdown.NewTask(checked, content)
		list.AddChild(task)
	}
	
	return list
}

// UpdateOrCreateHeading finds a heading by title and updates it, or creates a new one
func UpdateOrCreateHeading(doc *markdown.Document, level int, title string) *markdown.Heading {
	heading := markdown.FindHeadingByTitle(doc, title)
	
	if heading != nil {
		heading.Level = level
		return heading
	}
	
	// Create new heading
	newHeading := markdown.NewHeading(level, title)
	doc.AddChild(newHeading)
	return newHeading
}

// WriteList writes a markdown list from a slice of strings
func WriteList(items []string, ordered bool) string {
	var builder strings.Builder
	
	for i, item := range items {
		if ordered {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
		} else {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
	}
	
	return builder.String()
}

// WriteCheckboxList writes a markdown checkbox list
func WriteCheckboxList(items map[string]bool) string {
	var builder strings.Builder
	
	for item, checked := range items {
		if checked {
			builder.WriteString(fmt.Sprintf("- [x] %s\n", item))
		} else {
			builder.WriteString(fmt.Sprintf("- [ ] %s\n", item))
		}
	}
	
	return builder.String()
}

// WriteTable writes a simple markdown table
func WriteTable(headers []string, rows [][]string) string {
	var builder strings.Builder
	
	// Write headers
	builder.WriteString("| ")
	for _, header := range headers {
		builder.WriteString(header)
		builder.WriteString(" | ")
	}
	builder.WriteString("\n")
	
	// Write separator
	builder.WriteString("|")
	for range headers {
		builder.WriteString(" --- |")
	}
	builder.WriteString("\n")
	
	// Write rows
	for _, row := range rows {
		builder.WriteString("| ")
		for i, cell := range row {
			if i < len(headers) { // Ensure we don't exceed header count
				builder.WriteString(cell)
				builder.WriteString(" | ")
			}
		}
		builder.WriteString("\n")
	}
	
	return builder.String()
}