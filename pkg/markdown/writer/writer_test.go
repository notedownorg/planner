package writer

import (
	"strings"
	"testing"

	"github.com/notedownorg/planner/pkg/markdown"
	"github.com/notedownorg/planner/pkg/markdown/reader"
	"github.com/stretchr/testify/assert"
)

func TestWriteDocument(t *testing.T) {
	tests := []struct {
		name     string
		build    func() *markdown.Document
		expected string
	}{
		{
			name: "Empty document",
			build: func() *markdown.Document {
				return markdown.NewDocument()
			},
			expected: "",
		},
		{
			name: "Single heading",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()
				doc.AddChild(markdown.NewHeading(1, "Title"))
				return doc
			},
			expected: "# Title",
		},
		{
			name: "Heading with paragraph",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()
				h1 := markdown.NewHeading(1, "Title")
				h1.AddChild(markdown.NewParagraph("This is content."))
				doc.AddChild(h1)
				return doc
			},
			expected: `# Title

This is content.`,
		},
		{
			name: "Multiple headings",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()
				doc.AddChild(markdown.NewHeading(1, "First"))
				doc.AddChild(markdown.NewHeading(2, "Second"))
				doc.AddChild(markdown.NewHeading(3, "Third"))
				return doc
			},
			expected: `# First

## Second

### Third`,
		},
		{
			name: "Nested headings",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()
				h1 := markdown.NewHeading(1, "Level 1")
				h2 := markdown.NewHeading(2, "Level 2")
				h3 := markdown.NewHeading(3, "Level 3")

				h1.AddChild(h2)
				h2.AddChild(h3)
				doc.AddChild(h1)
				return doc
			},
			expected: `# Level 1

## Level 2

### Level 3`,
		},
		{
			name: "Task list",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()
				h1 := markdown.NewHeading(1, "Tasks")
				h1.AddChild(markdown.NewTask(true, "Completed"))
				h1.AddChild(markdown.NewTask(false, "Incomplete"))
				doc.AddChild(h1)
				return doc
			},
			expected: `# Tasks

- [x] Completed
- [ ] Incomplete`,
		},
		{
			name: "List with items",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()
				list := markdown.NewList(false)
				list.AddChild(markdown.NewListItem("First"))
				list.AddChild(markdown.NewListItem("Second"))
				doc.AddChild(list)
				return doc
			},
			expected: `- First
- Second`,
		},
		{
			name: "Complex document",
			build: func() *markdown.Document {
				doc := markdown.NewDocument()

				// Week heading
				week := markdown.NewHeading(1, "Week 51")
				doc.AddChild(week)

				// Habits section
				habits := markdown.NewHeading(2, "Habits")
				habits.AddChild(markdown.NewTask(true, "Exercise"))
				habits.AddChild(markdown.NewTask(true, "Meditation"))
				habits.AddChild(markdown.NewTask(false, "Reading"))
				week.AddChild(habits)

				// Notes section
				notes := markdown.NewHeading(2, "Notes")
				notes.AddChild(markdown.NewParagraph("Productive week overall."))
				week.AddChild(notes)

				return doc
			},
			expected: `# Week 51

## Habits

- [x] Exercise
- [x] Meditation
- [ ] Reading

## Notes

Productive week overall.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := tt.build()
			result := WriteDocument(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteHeading(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		title    string
		expected string
	}{
		{
			name:     "Level 1 heading",
			level:    1,
			title:    "Main Title",
			expected: "# Main Title\n",
		},
		{
			name:     "Level 2 heading",
			level:    2,
			title:    "Subtitle",
			expected: "## Subtitle\n",
		},
		{
			name:     "Level 6 heading",
			level:    6,
			title:    "Deep Heading",
			expected: "###### Deep Heading\n",
		},
		{
			name:     "Invalid level defaults to 1",
			level:    7,
			title:    "Invalid",
			expected: "# Invalid\n",
		},
		{
			name:     "Zero level defaults to 1",
			level:    0,
			title:    "Zero",
			expected: "# Zero\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			WriteHeading(&builder, tt.level, tt.title)
			assert.Equal(t, tt.expected, builder.String())
		})
	}
}

func TestWriteNode(t *testing.T) {
	tests := []struct {
		name     string
		node     markdown.Node
		expected string
	}{
		{
			name:     "Single heading",
			node:     markdown.NewHeading(2, "Test"),
			expected: "## Test",
		},
		{
			name:     "Paragraph",
			node:     markdown.NewParagraph("Test content."),
			expected: "Test content.",
		},
		{
			name:     "Task",
			node:     markdown.NewTask(true, "Completed"),
			expected: "- [x] Completed",
		},
		{
			name: "Heading with content",
			node: func() markdown.Node {
				h := markdown.NewHeading(1, "Title")
				h.AddChild(markdown.NewParagraph("Content"))
				return h
			}(),
			expected: `# Title

Content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WriteNode(tt.node)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateTaskList(t *testing.T) {
	tasks := map[string]bool{
		"Task 1": true,
		"Task 2": false,
		"Task 3": true,
	}

	list := CreateTaskList(tasks)
	assert.Equal(t, markdown.NodeList, list.Type())
	assert.False(t, list.Ordered)
	assert.Len(t, list.Children(), 3)

	// Check that all tasks are created properly
	taskCount := 0
	checkedCount := 0
	for _, child := range list.Children() {
		if task, ok := child.(*markdown.Task); ok {
			taskCount++
			if task.Checked {
				checkedCount++
			}
		}
	}
	assert.Equal(t, 3, taskCount)
	assert.Equal(t, 2, checkedCount)
}

func TestUpdateOrCreateHeading(t *testing.T) {
	tests := []struct {
		name          string
		existingTitle string
		newTitle      string
		newLevel      int
		shouldUpdate  bool
	}{
		{
			name:          "Create new heading",
			existingTitle: "",
			newTitle:      "New Heading",
			newLevel:      2,
			shouldUpdate:  false,
		},
		{
			name:          "Update existing heading",
			existingTitle: "Existing",
			newTitle:      "Existing",
			newLevel:      3,
			shouldUpdate:  true,
		},
		{
			name:          "Case insensitive update",
			existingTitle: "Test Heading",
			newTitle:      "test heading",
			newLevel:      2,
			shouldUpdate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := markdown.NewDocument()

			// Add existing heading if needed
			if tt.existingTitle != "" {
				doc.AddChild(markdown.NewHeading(1, tt.existingTitle))
			}

			initialCount := len(doc.Children())
			heading := UpdateOrCreateHeading(doc, tt.newLevel, tt.newTitle)

			assert.NotNil(t, heading)
			assert.Equal(t, tt.newLevel, heading.Level)

			if tt.shouldUpdate {
				assert.Len(t, doc.Children(), initialCount)
			} else {
				assert.Len(t, doc.Children(), initialCount+1)
			}
		})
	}
}

func TestWriteList(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		ordered  bool
		expected string
	}{
		{
			name:     "Empty list",
			items:    []string{},
			ordered:  false,
			expected: "",
		},
		{
			name:     "Unordered list",
			items:    []string{"First", "Second", "Third"},
			ordered:  false,
			expected: "- First\n- Second\n- Third\n",
		},
		{
			name:     "Ordered list",
			items:    []string{"Step 1", "Step 2", "Step 3"},
			ordered:  true,
			expected: "1. Step 1\n2. Step 2\n3. Step 3\n",
		},
		{
			name:     "Single item unordered",
			items:    []string{"Only item"},
			ordered:  false,
			expected: "- Only item\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WriteList(tt.items, tt.ordered)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteCheckboxList(t *testing.T) {
	tests := []struct {
		name  string
		items map[string]bool
		check func(t *testing.T, result string)
	}{
		{
			name:  "Empty list",
			items: map[string]bool{},
			check: func(t *testing.T, result string) {
				assert.Empty(t, result)
			},
		},
		{
			name: "All checked",
			items: map[string]bool{
				"Task 1": true,
				"Task 2": true,
			},
			check: func(t *testing.T, result string) {
				assert.Contains(t, result, "- [x] Task 1")
				assert.Contains(t, result, "- [x] Task 2")
				assert.NotContains(t, result, "- [ ]")
			},
		},
		{
			name: "Mixed checked/unchecked",
			items: map[string]bool{
				"Done": true,
				"Todo": false,
			},
			check: func(t *testing.T, result string) {
				assert.Contains(t, result, "- [x] Done")
				assert.Contains(t, result, "- [ ] Todo")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WriteCheckboxList(tt.items)
			tt.check(t, result)
		})
	}
}

func TestWriteTable(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		rows     [][]string
		validate func(t *testing.T, result string)
	}{
		{
			name:    "Basic table",
			headers: []string{"Name", "Age"},
			rows: [][]string{
				{"Alice", "30"},
				{"Bob", "25"},
			},
			validate: func(t *testing.T, result string) {
				assert.Contains(t, result, "| Name | Age |")
				assert.Contains(t, result, "| --- | --- |")
				assert.Contains(t, result, "| Alice | 30 |")
				assert.Contains(t, result, "| Bob | 25 |")
			},
		},
		{
			name:    "Single column",
			headers: []string{"Item"},
			rows: [][]string{
				{"Apple"},
				{"Banana"},
			},
			validate: func(t *testing.T, result string) {
				assert.Contains(t, result, "| Item |")
				assert.Contains(t, result, "| --- |")
				assert.Contains(t, result, "| Apple |")
				assert.Contains(t, result, "| Banana |")
			},
		},
		{
			name:    "Empty rows",
			headers: []string{"Col1", "Col2"},
			rows:    [][]string{},
			validate: func(t *testing.T, result string) {
				assert.Contains(t, result, "| Col1 | Col2 |")
				assert.Contains(t, result, "| --- | --- |")
				// No data rows
				lines := strings.Split(strings.TrimSpace(result), "\n")
				assert.Len(t, lines, 2) // Header + separator only
			},
		},
		{
			name:    "Row with fewer columns",
			headers: []string{"A", "B", "C"},
			rows: [][]string{
				{"1", "2", "3"},
				{"4", "5"}, // Missing third column
			},
			validate: func(t *testing.T, result string) {
				assert.Contains(t, result, "| 1 | 2 | 3 |")
				assert.Contains(t, result, "| 4 | 5 |")
				assert.NotContains(t, result, "| 4 | 5 | |") // Should not add empty column
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WriteTable(tt.headers, tt.rows)
			tt.validate(t, result)
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that we can parse and write back markdown without losing structure
	tests := []struct {
		name     string
		markdown string
	}{
		{
			name: "Basic document",
			markdown: `# Title

## Section 1

Content here.

## Section 2

More content.`,
		},
		{
			name: "Document with tasks",
			markdown: `# Tasks

## Today

- [x] Morning routine
- [ ] Code review
- [x] Team meeting

## Tomorrow

- [ ] Deploy update
- [ ] Write docs`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the markdown
			doc, err := reader.ParseMarkdown(tt.markdown)
			assert.NoError(t, err)

			// Write it back
			output := WriteDocument(doc)

			// The output should be equivalent (though whitespace might differ)
			assert.NotEmpty(t, output)

			// Re-parse to verify structure is maintained
			doc2, err := reader.ParseMarkdown(output)
			assert.NoError(t, err)

			// Compare heading counts
			headings1 := markdown.FindHeadings(doc)
			headings2 := markdown.FindHeadings(doc2)
			assert.Len(t, headings2, len(headings1))

			// Compare task counts
			tasks1 := markdown.FindTasks(doc)
			tasks2 := markdown.FindTasks(doc2)
			assert.Len(t, tasks2, len(tasks1))
		})
	}
}
