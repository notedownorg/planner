package reader

import (
	"testing"

	"github.com/notedownorg/planner/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		validate func(t *testing.T, doc *markdown.Document)
	}{
		{
			name: "Empty document",
			content: "",
			validate: func(t *testing.T, doc *markdown.Document) {
				assert.NotNil(t, doc)
				assert.Empty(t, doc.Children())
			},
		},
		{
			name: "Single heading",
			content: "# Title",
			validate: func(t *testing.T, doc *markdown.Document) {
				assert.Len(t, doc.Children(), 1)
				heading, ok := doc.Children()[0].(*markdown.Heading)
				assert.True(t, ok)
				assert.Equal(t, 1, heading.Level)
				assert.Equal(t, "Title", heading.Title)
			},
		},
		{
			name: "Heading with paragraph",
			content: `# Title

This is a paragraph.`,
			validate: func(t *testing.T, doc *markdown.Document) {
				assert.Len(t, doc.Children(), 1)
				heading := doc.Children()[0].(*markdown.Heading)
				assert.Len(t, heading.Children(), 1)
				para, ok := heading.Children()[0].(*markdown.Paragraph)
				assert.True(t, ok)
				assert.Equal(t, "This is a paragraph.", para.Content)
			},
		},
		{
			name: "Nested headings",
			content: `# Level 1
## Level 2
### Level 3`,
			validate: func(t *testing.T, doc *markdown.Document) {
				headings := markdown.FindHeadings(doc)
				assert.Len(t, headings, 3)
				assert.Equal(t, 1, headings[0].Level)
				assert.Equal(t, 2, headings[1].Level)
				assert.Equal(t, 3, headings[2].Level)
			},
		},
		{
			name: "Task list",
			content: `# Tasks
- [x] Completed task
- [ ] Incomplete task
- [X] Another completed task`,
			validate: func(t *testing.T, doc *markdown.Document) {
				tasks := markdown.FindTasks(doc)
				assert.Len(t, tasks, 3)
				assert.True(t, tasks[0].Checked)
				assert.Equal(t, "Completed task", tasks[0].Content)
				assert.False(t, tasks[1].Checked)
				assert.Equal(t, "Incomplete task", tasks[1].Content)
				assert.True(t, tasks[2].Checked)
			},
		},
		{
			name: "Mixed content",
			content: `# Week 51

## Monday

Morning tasks:
- [x] Wake up early
- [ ] Exercise

### Notes
Had a productive day.

## Tuesday

Todo list pending.`,
			validate: func(t *testing.T, doc *markdown.Document) {
				// Check structure
				assert.Len(t, doc.Children(), 1) // Week 51
				week := doc.Children()[0].(*markdown.Heading)
				assert.Equal(t, "Week 51", week.Title)
				
				// Find Monday section
				monday := markdown.FindHeadingByTitle(doc, "Monday")
				assert.NotNil(t, monday)
				assert.Equal(t, 2, monday.Level)
				
				// Check tasks
				tasks := markdown.FindTasks(doc)
				assert.Len(t, tasks, 2)
				
				// Check Notes subsection
				notes := markdown.FindHeadingByTitle(doc, "Notes")
				assert.NotNil(t, notes)
				assert.Equal(t, 3, notes.Level)
			},
		},
		{
			name: "Lists",
			content: `# Lists

## Unordered
- Item 1
- Item 2
- Item 3

## Ordered
1. First
2. Second
3. Third`,
			validate: func(t *testing.T, doc *markdown.Document) {
				// Find all lists
				var lists []*markdown.List
				var findLists func(n markdown.Node)
				findLists = func(n markdown.Node) {
					if list, ok := n.(*markdown.List); ok {
						lists = append(lists, list)
					}
					for _, child := range n.Children() {
						findLists(child)
					}
				}
				findLists(doc)
				
				assert.Len(t, lists, 2)
				assert.False(t, lists[0].Ordered)
				assert.True(t, lists[1].Ordered)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := ParseMarkdown(tt.content)
			assert.NoError(t, err)
			tt.validate(t, doc)
		})
	}
}

func TestParseTask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *markdown.Task
	}{
		{
			name:  "Checked task with x",
			input: "- [x] Completed task",
			expected: &markdown.Task{
				Checked: true,
				Content: "Completed task",
			},
		},
		{
			name:  "Checked task with X",
			input: "- [X] Completed task",
			expected: &markdown.Task{
				Checked: true,
				Content: "Completed task",
			},
		},
		{
			name:  "Unchecked task",
			input: "- [ ] Incomplete task",
			expected: &markdown.Task{
				Checked: false,
				Content: "Incomplete task",
			},
		},
		{
			name:  "Task with leading spaces",
			input: "  - [x] Indented task",
			expected: &markdown.Task{
				Checked: true,
				Content: "Indented task",
			},
		},
		{
			name:     "Not a task - regular list item",
			input:    "- Regular item",
			expected: nil,
		},
		{
			name:     "Not a task - paragraph",
			input:    "This is just text",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := parseTask(tt.input)
			
			if tt.expected == nil {
				assert.Nil(t, task)
			} else {
				assert.NotNil(t, task)
				assert.Equal(t, tt.expected.Checked, task.Checked)
				assert.Equal(t, tt.expected.Content, task.Content)
			}
		})
	}
}

func TestComplexDocumentStructure(t *testing.T) {
	content := `# Project Plan

## Overview
This is the project overview.

## Tasks

### High Priority
- [x] Design architecture
- [x] Set up repository
- [ ] Implement core features

### Low Priority
- [ ] Add animations
- [ ] Optimize performance

## Notes

### Meeting Notes
Discussed timeline and deliverables.

### Ideas
- Consider using caching
- Add metrics collection`

	doc, err := ParseMarkdown(content)
	assert.NoError(t, err)

	// Verify document structure
	projectPlan := markdown.FindHeadingByTitle(doc, "Project Plan")
	assert.NotNil(t, projectPlan)
	assert.Equal(t, 1, projectPlan.Level)

	// Check all level 2 headings
	headings := markdown.FindHeadings(doc)
	level2Count := 0
	for _, h := range headings {
		if h.Level == 2 {
			level2Count++
		}
	}
	assert.Equal(t, 3, level2Count) // Overview, Tasks, Notes

	// Verify tasks are properly parsed
	tasks := markdown.FindTasks(doc)
	assert.Len(t, tasks, 5)
	
	// Count completed tasks
	completed := 0
	for _, task := range tasks {
		if task.Checked {
			completed++
		}
	}
	assert.Equal(t, 2, completed)
}