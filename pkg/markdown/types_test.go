package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		create   func() Node
		expected NodeType
	}{
		{
			name:     "Document node",
			create:   func() Node { return NewDocument() },
			expected: NodeDocument,
		},
		{
			name:     "Heading node",
			create:   func() Node { return NewHeading(1, "Test") },
			expected: NodeHeading,
		},
		{
			name:     "Paragraph node",
			create:   func() Node { return NewParagraph("Test content") },
			expected: NodeParagraph,
		},
		{
			name:     "Task node",
			create:   func() Node { return NewTask(true, "Test task") },
			expected: NodeTask,
		},
		{
			name:     "List node",
			create:   func() Node { return NewList(false) },
			expected: NodeList,
		},
		{
			name:     "ListItem node",
			create:   func() Node { return NewListItem("Test item") },
			expected: NodeListItem,
		},
		{
			name:     "Text node",
			create:   func() Node { return NewText("Test text") },
			expected: NodeText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.create()
			assert.Equal(t, tt.expected, node.Type())
			assert.Empty(t, node.Children())
			assert.Nil(t, node.Parent())
		})
	}
}

func TestNodeRelationships(t *testing.T) {
	doc := NewDocument()
	heading := NewHeading(1, "Test Heading")
	para := NewParagraph("Test paragraph")
	
	// Test adding children
	doc.AddChild(heading)
	heading.AddChild(para)
	
	// Test parent/child relationships
	assert.Equal(t, doc, heading.Parent())
	assert.Equal(t, heading, para.Parent())
	assert.Len(t, doc.Children(), 1)
	assert.Len(t, heading.Children(), 1)
	assert.Equal(t, heading, doc.Children()[0])
	assert.Equal(t, para, heading.Children()[0])
}

func TestFindHeadings(t *testing.T) {
	tests := []struct {
		name     string
		build    func() Node
		expected []string // Expected heading titles
	}{
		{
			name: "No headings",
			build: func() Node {
				doc := NewDocument()
				doc.AddChild(NewParagraph("Some text"))
				return doc
			},
			expected: []string{},
		},
		{
			name: "Single heading",
			build: func() Node {
				doc := NewDocument()
				doc.AddChild(NewHeading(1, "Title"))
				return doc
			},
			expected: []string{"Title"},
		},
		{
			name: "Multiple headings at same level",
			build: func() Node {
				doc := NewDocument()
				doc.AddChild(NewHeading(1, "First"))
				doc.AddChild(NewHeading(1, "Second"))
				doc.AddChild(NewHeading(1, "Third"))
				return doc
			},
			expected: []string{"First", "Second", "Third"},
		},
		{
			name: "Nested headings",
			build: func() Node {
				doc := NewDocument()
				h1 := NewHeading(1, "Level 1")
				h2 := NewHeading(2, "Level 2")
				h3 := NewHeading(3, "Level 3")
				
				doc.AddChild(h1)
				h1.AddChild(h2)
				h2.AddChild(h3)
				
				return doc
			},
			expected: []string{"Level 1", "Level 2", "Level 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.build()
			headings := FindHeadings(node)
			
			assert.Len(t, headings, len(tt.expected))
			for i, expected := range tt.expected {
				assert.Equal(t, expected, headings[i].Title)
			}
		})
	}
}

func TestFindHeadingByTitle(t *testing.T) {
	tests := []struct {
		name        string
		searchTitle string
		found       bool
		foundTitle  string
	}{
		{
			name:        "Find existing heading",
			searchTitle: "Test Heading",
			found:       true,
			foundTitle:  "Test Heading",
		},
		{
			name:        "Case insensitive search",
			searchTitle: "test heading",
			found:       true,
			foundTitle:  "Test Heading",
		},
		{
			name:        "Not found",
			searchTitle: "Non-existent",
			found:       false,
		},
	}

	// Build test document
	doc := NewDocument()
	doc.AddChild(NewHeading(1, "Test Heading"))
	doc.AddChild(NewHeading(2, "Another Heading"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heading := FindHeadingByTitle(doc, tt.searchTitle)
			
			if tt.found {
				assert.NotNil(t, heading)
				assert.Equal(t, tt.foundTitle, heading.Title)
			} else {
				assert.Nil(t, heading)
			}
		})
	}
}

func TestFindTasks(t *testing.T) {
	tests := []struct {
		name     string
		build    func() Node
		expected int // Expected number of tasks
		checked  []bool
	}{
		{
			name: "No tasks",
			build: func() Node {
				doc := NewDocument()
				doc.AddChild(NewParagraph("Some text"))
				return doc
			},
			expected: 0,
			checked:  []bool{},
		},
		{
			name: "Tasks in list",
			build: func() Node {
				doc := NewDocument()
				list := NewList(false)
				list.AddChild(NewTask(true, "Completed task"))
				list.AddChild(NewTask(false, "Incomplete task"))
				doc.AddChild(list)
				return doc
			},
			expected: 2,
			checked:  []bool{true, false},
		},
		{
			name: "Tasks under heading",
			build: func() Node {
				doc := NewDocument()
				heading := NewHeading(1, "Tasks")
				heading.AddChild(NewTask(true, "Task 1"))
				heading.AddChild(NewTask(true, "Task 2"))
				heading.AddChild(NewTask(false, "Task 3"))
				doc.AddChild(heading)
				return doc
			},
			expected: 3,
			checked:  []bool{true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.build()
			tasks := FindTasks(node)
			
			assert.Len(t, tasks, tt.expected)
			for i, checked := range tt.checked {
				assert.Equal(t, checked, tasks[i].Checked)
			}
		})
	}
}

func TestGetNodeDepth(t *testing.T) {
	tests := []struct {
		name     string
		build    func() Node
		expected int
	}{
		{
			name: "Root node depth",
			build: func() Node {
				return NewDocument()
			},
			expected: 0,
		},
		{
			name: "Direct child depth",
			build: func() Node {
				doc := NewDocument()
				heading := NewHeading(1, "Test")
				doc.AddChild(heading)
				return heading
			},
			expected: 1,
		},
		{
			name: "Nested node depth",
			build: func() Node {
				doc := NewDocument()
				h1 := NewHeading(1, "Level 1")
				h2 := NewHeading(2, "Level 2")
				para := NewParagraph("Deep content")
				
				doc.AddChild(h1)
				h1.AddChild(h2)
				h2.AddChild(para)
				
				return para
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.build()
			depth := GetNodeDepth(node)
			assert.Equal(t, tt.expected, depth)
		})
	}
}

func TestHeadingNode(t *testing.T) {
	tests := []struct {
		name  string
		level int
		title string
	}{
		{name: "H1", level: 1, title: "Main Title"},
		{name: "H2", level: 2, title: "Subtitle"},
		{name: "H6", level: 6, title: "Deepest Level"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heading := NewHeading(tt.level, tt.title)
			assert.Equal(t, tt.level, heading.Level)
			assert.Equal(t, tt.title, heading.Title)
			assert.Equal(t, NodeHeading, heading.Type())
		})
	}
}

func TestTaskNode(t *testing.T) {
	tests := []struct {
		name    string
		checked bool
		content string
	}{
		{name: "Checked task", checked: true, content: "Completed item"},
		{name: "Unchecked task", checked: false, content: "TODO item"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask(tt.checked, tt.content)
			assert.Equal(t, tt.checked, task.Checked)
			assert.Equal(t, tt.content, task.Content)
			assert.Equal(t, NodeTask, task.Type())
		})
	}
}

func TestBaseNode_ClearChildren(t *testing.T) {
	// Create a parent node
	parent := &BaseNode{
		NodeType: NodeHeading,
		children: []Node{},
	}

	// Create child nodes
	child1 := &BaseNode{NodeType: NodeText}
	child2 := &BaseNode{NodeType: NodeTask}
	child3 := &BaseNode{NodeType: NodeParagraph}

	// Add children
	parent.AddChild(child1)
	parent.AddChild(child2)
	parent.AddChild(child3)

	// Verify children were added
	assert.Len(t, parent.Children(), 3)
	assert.Equal(t, parent, child1.Parent())
	assert.Equal(t, parent, child2.Parent())
	assert.Equal(t, parent, child3.Parent())

	// Clear children
	parent.ClearChildren()

	// Verify children were cleared
	assert.Len(t, parent.Children(), 0)
	assert.Nil(t, child1.Parent())
	assert.Nil(t, child2.Parent())
	assert.Nil(t, child3.Parent())
}

func TestComplexDocumentStructure(t *testing.T) {
	// Create a complex document structure
	doc := NewDocument()
	
	// Main heading
	mainHeading := NewHeading(1, "Week 01")
	doc.AddChild(mainHeading)
	
	// Habits section
	habitsHeading := NewHeading(2, "Habits")
	mainHeading.AddChild(habitsHeading)
	
	// Add habits as tasks
	habit1 := NewTask(false, "Exercise")
	habit2 := NewTask(true, "Read 30 minutes")
	habit3 := NewTask(false, "Meditate")
	
	habitsHeading.AddChild(habit1)
	habitsHeading.AddChild(habit2)
	habitsHeading.AddChild(habit3)
	
	// Notes section
	notesHeading := NewHeading(2, "Notes")
	mainHeading.AddChild(notesHeading)
	
	notesPara := NewParagraph("Some notes for the week")
	notesHeading.AddChild(notesPara)
	
	// Verify structure
	assert.Len(t, doc.Children(), 1)
	assert.Len(t, mainHeading.Children(), 2)
	assert.Len(t, habitsHeading.Children(), 3)
	assert.Len(t, notesHeading.Children(), 1)
	
	// Test finding specific nodes
	foundHabits := FindHeadingByTitle(doc, "Habits")
	assert.NotNil(t, foundHabits)
	
	tasks := FindTasks(foundHabits)
	assert.Len(t, tasks, 3)
	
	// Test clearing habits
	foundHabits.ClearChildren()
	assert.Len(t, foundHabits.Children(), 0)
	
	// Verify parent references were cleared
	assert.Nil(t, habit1.Parent())
	assert.Nil(t, habit2.Parent())
	assert.Nil(t, habit3.Parent())
}