package habits

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/notedownorg/planner/pkg/config"
	"github.com/notedownorg/planner/pkg/markdown"
	"github.com/notedownorg/planner/pkg/markdown/reader"
	"github.com/notedownorg/planner/pkg/markdown/writer"
)

// Service manages habit tracking with markdown persistence
type Service struct {
	config *config.Config
}

// NewService creates a new habit service
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// GetWeeklyFilePath generates the file path for a specific week's notes
func (s *Service) GetWeeklyFilePath(year int, weekNumber int) string {
	// Format: YYYY-[W]WW -> 2024-W01
	weekStr := fmt.Sprintf("%d-W%02d", year, weekNumber)
	filename := fmt.Sprintf("%s.md", weekStr)

	return filepath.Join(s.config.WorkspaceRoot, s.config.PeriodicNotes.WeeklySubdir, filename)
}

// getCurrentWeekInfo gets the current ISO week information
func getCurrentWeekInfo() (int, int) {
	now := time.Now()
	year, week := now.ISOWeek()
	return year, week
}

// LoadWeeklyHabits loads habits for a specific week
func (s *Service) LoadWeeklyHabits(year int, weekNumber int) (*WeeklyHabits, error) {
	filePath := s.GetWeeklyFilePath(year, weekNumber)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, create new with default habits
		return s.createNewWeeklyHabits(year, weekNumber)
	}

	// Read existing file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read weekly file: %w", err)
	}

	// Parse markdown
	doc, err := reader.ParseMarkdown(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Extract habits from markdown
	return s.extractHabitsFromDocument(doc, year, weekNumber)
}

// SaveWeeklyHabits saves habits to markdown file
func (s *Service) SaveWeeklyHabits(habits *WeeklyHabits) error {
	filePath := s.GetWeeklyFilePath(habits.Year, habits.WeekNumber)

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read existing content if file exists
	var doc *markdown.Document
	if _, err := os.Stat(filePath); err == nil {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read existing file: %w", err)
		}

		doc, err = reader.ParseMarkdown(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse existing markdown: %w", err)
		}
	} else {
		// Create new document
		doc = s.createNewWeeklyDocument(habits.Year, habits.WeekNumber)
	}

	// Update or create habits section
	if err := s.updateHabitsSection(doc, habits); err != nil {
		return fmt.Errorf("failed to update habits section: %w", err)
	}

	// Write back to file
	content := writer.WriteDocument(doc)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// createNewWeeklyHabits creates a new weekly habits structure with defaults
func (s *Service) createNewWeeklyHabits(year int, weekNumber int) (*WeeklyHabits, error) {
	// Get default habits from previous week or config
	defaultHabits := s.getDefaultHabitsWithoutRecursion(year, weekNumber)

	habits := &WeeklyHabits{
		Year:       year,
		WeekNumber: weekNumber,
		Habits:     make(map[string]*Habit),
		DayStatus:  make(map[string]bool),
	}

	// Initialize with default habits (all uncompleted)
	for i, habitName := range defaultHabits {
		habits.Habits[habitName] = &Habit{
			Name:      habitName,
			Completed: false,
			Order:     i,
		}
	}

	return habits, nil
}

// getDefaultHabits gets default habits from previous week or returns empty
func (s *Service) getDefaultHabits(year int, weekNumber int) []string {
	// Try to get habits from previous week
	prevYear, prevWeek := year, weekNumber-1
	if prevWeek < 1 {
		prevYear--
		prevWeek = 52 // Approximate - could be 53 in some years
	}

	prevHabits, err := s.LoadWeeklyHabits(prevYear, prevWeek)
	if err == nil && len(prevHabits.Habits) > 0 {
		var habits []string
		for habitName := range prevHabits.Habits {
			habits = append(habits, habitName)
		}
		return habits
	}

	// Return empty list - no default habits
	return []string{}
}

// getDefaultHabitsWithoutRecursion gets default habits from previous week file directly without recursion
func (s *Service) getDefaultHabitsWithoutRecursion(year int, weekNumber int) []string {
	// Try to get habits from previous week
	prevYear, prevWeek := year, weekNumber-1
	if prevWeek < 1 {
		prevYear--
		prevWeek = 52 // Approximate - could be 53 in some years
	}

	// Check if previous week file exists and read it directly
	filePath := s.GetWeeklyFilePath(prevYear, prevWeek)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// No previous week file, return empty
		return []string{}
	}

	// Read existing file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}
	}

	// Parse markdown
	doc, err := reader.ParseMarkdown(string(content))
	if err != nil {
		return []string{}
	}

	// Extract habits from markdown
	prevHabits, err := s.extractHabitsFromDocument(doc, prevYear, prevWeek)
	if err != nil || len(prevHabits.Habits) == 0 {
		return []string{}
	}

	var habits []string
	for habitName := range prevHabits.Habits {
		habits = append(habits, habitName)
	}
	return habits
}

// createNewWeeklyDocument creates a new markdown document for a week
func (s *Service) createNewWeeklyDocument(year int, weekNumber int) *markdown.Document {
	doc := markdown.NewDocument()

	// Add main heading
	weekTitle := fmt.Sprintf("Week %02d", weekNumber)
	weekHeading := markdown.NewHeading(1, weekTitle)
	doc.AddChild(weekHeading)

	return doc
}

// extractHabitsFromDocument extracts habits from a markdown document
func (s *Service) extractHabitsFromDocument(doc *markdown.Document, year int, weekNumber int) (*WeeklyHabits, error) {
	habits := &WeeklyHabits{
		Year:       year,
		WeekNumber: weekNumber,
		Habits:     make(map[string]*Habit),
		DayStatus:  make(map[string]bool),
	}

	// Find the Habits heading
	habitsHeading := markdown.FindHeadingByTitle(doc, "Habits")
	if habitsHeading == nil {
		// No habits section found, return empty but valid structure
		return habits, nil
	}

	// Extract tasks from the habits section
	tasks := markdown.FindTasks(habitsHeading)
	for i, task := range tasks {
		// Use the entire task content as the habit name
		if len(task.Content) > 0 {
			habitName := task.Content
			habits.Habits[habitName] = &Habit{
				Name:      habitName,
				Completed: task.Checked,
				Order:     i, // Use index to preserve order from file
			}
		}
	}

	return habits, nil
}

// updateHabitsSection updates the habits section in a markdown document
func (s *Service) updateHabitsSection(doc *markdown.Document, habits *WeeklyHabits) error {
	// Find or create the main week heading
	weekTitle := fmt.Sprintf("Week %02d", habits.WeekNumber)
	weekHeading := markdown.FindHeadingByTitle(doc, weekTitle)
	if weekHeading == nil {
		weekHeading = markdown.NewHeading(1, weekTitle)
		doc.AddChild(weekHeading)
	}

	// Find or create the Habits heading
	habitsHeading := markdown.FindHeadingByTitle(weekHeading, "Habits")
	if habitsHeading == nil {
		habitsHeading = markdown.NewHeading(2, "Habits")
		weekHeading.AddChild(habitsHeading)
	} else {
		// Clear existing children to avoid duplicates when updating
		habitsHeading.ClearChildren()
	}

	// Sort habits by order before adding them
	var sortedHabits []*Habit
	for _, habit := range habits.Habits {
		sortedHabits = append(sortedHabits, habit)
	}
	// Sort by completion status first, then by Order field
	sort.Slice(sortedHabits, func(i, j int) bool {
		// Incomplete habits come first
		if sortedHabits[i].Completed != sortedHabits[j].Completed {
			return !sortedHabits[i].Completed
		}
		// Within same completion status, sort by order
		return sortedHabits[i].Order < sortedHabits[j].Order
	})

	// Add habit tasks in sorted order
	for _, habit := range sortedHabits {
		task := markdown.NewTask(habit.Completed, habit.Name)
		habitsHeading.AddChild(task)
	}

	return nil
}

// GetCurrentWeekHabits gets habits for the current week
func (s *Service) GetCurrentWeekHabits() (*WeeklyHabits, error) {
	year, week := getCurrentWeekInfo()
	return s.LoadWeeklyHabits(year, week)
}

// ToggleHabit toggles the completion status of a habit
func (s *Service) ToggleHabit(year int, weekNumber int, habitName string) error {
	habits, err := s.LoadWeeklyHabits(year, weekNumber)
	if err != nil {
		return err
	}

	if habit, exists := habits.Habits[habitName]; exists {
		habit.Completed = !habit.Completed
	}

	return s.SaveWeeklyHabits(habits)
}

// AddHabit adds a new habit to the current week
func (s *Service) AddHabit(year int, weekNumber int, habitName string) error {
	habits, err := s.LoadWeeklyHabits(year, weekNumber)
	if err != nil {
		return err
	}

	// Only add if not already exists
	if _, exists := habits.Habits[habitName]; !exists {
		// Set order as the next highest order
		maxOrder := -1
		for _, habit := range habits.Habits {
			if habit.Order > maxOrder {
				maxOrder = habit.Order
			}
		}
		habits.Habits[habitName] = &Habit{
			Name:      habitName,
			Completed: false,
			Order:     maxOrder + 1,
		}
	}

	return s.SaveWeeklyHabits(habits)
}

// RemoveHabit removes a habit from the current week
func (s *Service) RemoveHabit(year int, weekNumber int, habitName string) error {
	habits, err := s.LoadWeeklyHabits(year, weekNumber)
	if err != nil {
		return err
	}

	delete(habits.Habits, habitName)

	return s.SaveWeeklyHabits(habits)
}

// ReorderHabits reorders habits according to the provided habit names array
func (s *Service) ReorderHabits(year int, weekNumber int, habitNames []string) error {
	habits, err := s.LoadWeeklyHabits(year, weekNumber)
	if err != nil {
		return err
	}

	// Update order for each habit based on position in array
	for i, habitName := range habitNames {
		if habit, exists := habits.Habits[habitName]; exists {
			habit.Order = i
		}
	}

	return s.SaveWeeklyHabits(habits)
}
