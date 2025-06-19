# Markdown Package

This package provides a tree-based markdown parser and writer specifically designed for handling structured markdown documents with ATX headings, tasks, and hierarchical content.

## Tree Structure

The markdown package represents documents as a tree structure where each node can have children. This design allows for proper representation of markdown's hierarchical nature.

### Node Types

- **Document**: The root node of the tree
- **Heading**: ATX headings (# through ######) that can contain other content
- **Paragraph**: Text paragraphs
- **Task**: Checkbox items (- [ ] or - [x])
- **List**: Ordered or unordered lists
- **ListItem**: Individual items within a list
- **Text**: Raw text content

### Tree Example

```
Document
├── Heading (Level 1: "Week 51")
│   ├── Heading (Level 2: "Habits")
│   │   ├── Task (Checked: "Morning meditation")
│   │   ├── Task (Unchecked: "Exercise 30 min")
│   │   └── Task (Checked: "Read 20 pages")
│   └── Heading (Level 2: "Daily Notes")
│       └── Paragraph ("Productive day, completed project X")
└── Heading (Level 1: "Goals")
    └── List (Unordered)
        ├── ListItem ("Complete weekly review")
        └── ListItem ("Plan next sprint")
```

## Usage

### Reading Markdown

```go
import (
    "github.com/notedownorg/planner/pkg/markdown/reader"
)

// Parse markdown content into a tree
doc, err := reader.ParseMarkdown(markdownContent)
if err != nil {
    return err
}

// Find specific headings
habitSection := markdown.FindHeadingByTitle(doc, "Habits")

// Get all tasks in the document
tasks := markdown.FindTasks(doc)

// Get all level 2 headings
headings := markdown.FindHeadings(doc)
for _, h := range headings {
    if h.Level == 2 {
        fmt.Println(h.Title)
    }
}
```

### Writing Markdown

```go
import (
    "github.com/notedownorg/planner/pkg/markdown"
    "github.com/notedownorg/planner/pkg/markdown/writer"
)

// Create a document structure
doc := markdown.NewDocument()

// Add a main heading
weekHeading := markdown.NewHeading(1, "Week 51")
doc.AddChild(weekHeading)

// Add a habits section under the week
habitsHeading := markdown.NewHeading(2, "Habits")
weekHeading.AddChild(habitsHeading)

// Add tasks under habits
habitsHeading.AddChild(markdown.NewTask(true, "Morning meditation"))
habitsHeading.AddChild(markdown.NewTask(false, "Exercise 30 min"))

// Convert to markdown string
output := writer.WriteDocument(doc)
```

### Modifying Existing Documents

```go
// Find and update a heading
heading := markdown.FindHeadingByTitle(doc, "Habits")
if heading != nil {
    // Add a new task
    heading.AddChild(markdown.NewTask(false, "Drink 8 glasses of water"))
}

// Update or create a heading
habitHeading := writer.UpdateOrCreateHeading(doc, 2, "Daily Habits")

// Create a task list
tasks := map[string]bool{
    "Morning routine": true,
    "Evening review": false,
    "Meal prep": true,
}
taskList := writer.CreateTaskList(tasks)
habitHeading.AddChild(taskList)
```

## Design Principles

### 1. Tree-Based Structure
Unlike flat section-based parsers, this package maintains the full hierarchical structure of markdown documents. This allows for:
- Proper nesting of headings and content
- Maintaining parent-child relationships
- Easy manipulation of document structure

### 2. Heading-Centric Organization
ATX headings serve as the primary organizational structure. Content is associated with its parent heading, making it ideal for:
- Weekly planners with day/category sections
- Hierarchical note-taking systems
- Structured documentation

### 3. Task Support
First-class support for task lists (checkboxes) makes this package particularly suitable for:
- Habit tracking
- Todo lists
- Progress tracking

### 4. Goldmark Integration
The parser is built on top of the robust goldmark parser, providing:
- Standards-compliant markdown parsing
- Extension support (task lists, tables, etc.)
- Reliable AST traversal

## Use Cases

### Weekly Planner
```markdown
# Week 51

## Monday
- [x] Team standup
- [ ] Code review
- [ ] Update documentation

## Habits
- [x] Morning meditation
- [x] Exercise
- [ ] Reading

## Notes
Productive start to the week. Completed major refactoring.
```

### Project Documentation
```markdown
# Project Overview

## Architecture
System uses microservices pattern...

### Services
- [x] Auth service implemented
- [x] User service implemented
- [ ] Notification service pending

## API Documentation
### Endpoints
...
```

## Performance Considerations

1. **Tree Traversal**: Most operations (finding headings, tasks) require tree traversal. For large documents, consider caching frequently accessed nodes.

2. **Memory Usage**: The entire document is kept in memory as a tree structure. For very large documents, consider streaming approaches.

3. **Modification**: The tree structure makes modifications efficient as you can directly access and modify specific nodes without reparsing.

## Future Enhancements

Potential areas for expansion:
- Table support as tree nodes
- Code block nodes with language metadata
- Link and image nodes
- Frontmatter support
- Streaming parser for large documents
- Tree diffing for change tracking