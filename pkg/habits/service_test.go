package habits

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notedownorg/planner/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestService(t *testing.T) (*Service, string) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "habits-test-*")
	require.NoError(t, err)

	cfg := &config.Config{
		WorkspaceRoot: tempDir,
		PeriodicNotes: config.PeriodicNotes{
			WeeklySubdir:     "weekly",
			WeeklyNameFormat: "YYYY-[W]WW",
		},
	}

	service := NewService(cfg)
	return service, tempDir
}

func TestGetWeeklyFilePath(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name       string
		year       int
		weekNumber int
		expected   string
	}{
		{
			name:       "Week 1",
			year:       2024,
			weekNumber: 1,
			expected:   filepath.Join(tempDir, "weekly", "2024-W01.md"),
		},
		{
			name:       "Week 10",
			year:       2024,
			weekNumber: 10,
			expected:   filepath.Join(tempDir, "weekly", "2024-W10.md"),
		},
		{
			name:       "Week 52",
			year:       2023,
			weekNumber: 52,
			expected:   filepath.Join(tempDir, "weekly", "2023-W52.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetWeeklyFilePath(tt.year, tt.weekNumber)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadWeeklyHabits_NewFile(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Load habits for a week that doesn't exist yet
	habits, err := service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)

	assert.Equal(t, 2024, habits.Year)
	assert.Equal(t, 1, habits.WeekNumber)
	assert.NotNil(t, habits.Habits)
	assert.NotNil(t, habits.DayStatus)
	assert.Len(t, habits.Habits, 0) // Should be empty with no default habits
}

func TestSaveAndLoadWeeklyHabits(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create test habits
	testHabits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Exercise": {Name: "Exercise", Completed: false, Order: 0},
			"Read":     {Name: "Read", Completed: true, Order: 1},
			"Meditate": {Name: "Meditate", Completed: false, Order: 2},
		},
		DayStatus: make(map[string]bool),
	}

	// Save habits
	err := service.SaveWeeklyHabits(testHabits)
	require.NoError(t, err)

	// Verify file was created
	filePath := service.GetWeeklyFilePath(2024, 1)
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// Load habits back
	loadedHabits, err := service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)

	// Verify loaded habits match saved habits
	assert.Equal(t, testHabits.Year, loadedHabits.Year)
	assert.Equal(t, testHabits.WeekNumber, loadedHabits.WeekNumber)
	assert.Len(t, loadedHabits.Habits, 3)

	// Check individual habits
	exercise := loadedHabits.Habits["Exercise"]
	assert.NotNil(t, exercise)
	assert.Equal(t, "Exercise", exercise.Name)
	assert.False(t, exercise.Completed)
	assert.Equal(t, 0, exercise.Order)

	read := loadedHabits.Habits["Read"]
	assert.NotNil(t, read)
	assert.Equal(t, "Read", read.Name)
	assert.True(t, read.Completed)
	// Read's order might be different due to sorting completed items at the end
}

func TestToggleHabit(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create initial habits
	testHabits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Exercise": {Name: "Exercise", Completed: false, Order: 0},
		},
		DayStatus: make(map[string]bool),
	}
	err := service.SaveWeeklyHabits(testHabits)
	require.NoError(t, err)

	// Toggle habit
	err = service.ToggleHabit(2024, 1, "Exercise")
	require.NoError(t, err)

	// Verify habit was toggled
	habits, err := service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.True(t, habits.Habits["Exercise"].Completed)

	// Toggle back
	err = service.ToggleHabit(2024, 1, "Exercise")
	require.NoError(t, err)

	habits, err = service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.False(t, habits.Habits["Exercise"].Completed)
}

func TestAddHabit(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Add first habit
	err := service.AddHabit(2024, 1, "Exercise")
	require.NoError(t, err)

	habits, err := service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.Len(t, habits.Habits, 1)
	assert.NotNil(t, habits.Habits["Exercise"])
	assert.Equal(t, 0, habits.Habits["Exercise"].Order)

	// Add second habit
	err = service.AddHabit(2024, 1, "Read")
	require.NoError(t, err)

	habits, err = service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.Len(t, habits.Habits, 2)
	assert.NotNil(t, habits.Habits["Read"])
	assert.Equal(t, 1, habits.Habits["Read"].Order)

	// Try to add duplicate habit (should be ignored)
	err = service.AddHabit(2024, 1, "Exercise")
	require.NoError(t, err)

	habits, err = service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.Len(t, habits.Habits, 2) // Still only 2 habits
}

func TestRemoveHabit(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create initial habits
	testHabits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Exercise": {Name: "Exercise", Completed: false, Order: 0},
			"Read":     {Name: "Read", Completed: false, Order: 1},
		},
		DayStatus: make(map[string]bool),
	}
	err := service.SaveWeeklyHabits(testHabits)
	require.NoError(t, err)

	// Remove habit
	err = service.RemoveHabit(2024, 1, "Exercise")
	require.NoError(t, err)

	// Verify habit was removed
	habits, err := service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.Len(t, habits.Habits, 1)
	assert.Nil(t, habits.Habits["Exercise"])
	assert.NotNil(t, habits.Habits["Read"])
}

func TestReorderHabits(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create initial habits
	testHabits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Exercise": {Name: "Exercise", Completed: false, Order: 0},
			"Read":     {Name: "Read", Completed: false, Order: 1},
			"Meditate": {Name: "Meditate", Completed: false, Order: 2},
		},
		DayStatus: make(map[string]bool),
	}
	err := service.SaveWeeklyHabits(testHabits)
	require.NoError(t, err)

	// Reorder habits
	newOrder := []string{"Read", "Meditate", "Exercise"}
	err = service.ReorderHabits(2024, 1, newOrder)
	require.NoError(t, err)

	// Verify new order
	habits, err := service.LoadWeeklyHabits(2024, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, habits.Habits["Read"].Order)
	assert.Equal(t, 1, habits.Habits["Meditate"].Order)
	assert.Equal(t, 2, habits.Habits["Exercise"].Order)
}

func TestMarkdownPersistence_OrderPreservation(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create habits with specific order
	testHabits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Third":  {Name: "Third", Completed: false, Order: 2},
			"First":  {Name: "First", Completed: false, Order: 0},
			"Second": {Name: "Second", Completed: false, Order: 1},
		},
		DayStatus: make(map[string]bool),
	}
	err := service.SaveWeeklyHabits(testHabits)
	require.NoError(t, err)

	// Read the file content
	filePath := service.GetWeeklyFilePath(2024, 1)
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	// Verify order in file
	fileContent := string(content)
	firstPos := indexOf(fileContent, "First")
	secondPos := indexOf(fileContent, "Second")
	thirdPos := indexOf(fileContent, "Third")

	assert.True(t, firstPos < secondPos, "First should come before Second in file")
	assert.True(t, secondPos < thirdPos, "Second should come before Third in file")
}

func TestMarkdownPersistence_CompletedAtEnd(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create habits with mixed completion status
	testHabits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Incomplete1": {Name: "Incomplete1", Completed: false, Order: 0},
			"Complete1":   {Name: "Complete1", Completed: true, Order: 1},
			"Incomplete2": {Name: "Incomplete2", Completed: false, Order: 2},
			"Complete2":   {Name: "Complete2", Completed: true, Order: 3},
		},
		DayStatus: make(map[string]bool),
	}
	err := service.SaveWeeklyHabits(testHabits)
	require.NoError(t, err)

	// Read the file content
	filePath := service.GetWeeklyFilePath(2024, 1)
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	// Verify completed tasks are at the end
	fileContent := string(content)
	incomplete1Pos := indexOf(fileContent, "Incomplete1")
	incomplete2Pos := indexOf(fileContent, "Incomplete2")
	complete1Pos := indexOf(fileContent, "Complete1")
	complete2Pos := indexOf(fileContent, "Complete2")

	assert.True(t, incomplete1Pos < complete1Pos, "Incomplete habits should come before completed ones")
	assert.True(t, incomplete2Pos < complete1Pos, "Incomplete habits should come before completed ones")
	assert.True(t, incomplete1Pos < complete2Pos, "Incomplete habits should come before completed ones")
	assert.True(t, incomplete2Pos < complete2Pos, "Incomplete habits should come before completed ones")
}

func TestGetDefaultHabitsFromPreviousWeek(t *testing.T) {
	service, tempDir := createTestService(t)
	defer os.RemoveAll(tempDir)

	// Create habits for week 1
	week1Habits := &WeeklyHabits{
		Year:       2024,
		WeekNumber: 1,
		Habits: map[string]*Habit{
			"Exercise": {Name: "Exercise", Completed: false, Order: 0},
			"Read":     {Name: "Read", Completed: true, Order: 1},
		},
		DayStatus: make(map[string]bool),
	}
	err := service.SaveWeeklyHabits(week1Habits)
	require.NoError(t, err)

	// Load habits for week 2 (should get defaults from week 1)
	week2Habits, err := service.LoadWeeklyHabits(2024, 2)
	require.NoError(t, err)

	// Verify habits were copied but not completed
	assert.Len(t, week2Habits.Habits, 2)
	assert.NotNil(t, week2Habits.Habits["Exercise"])
	assert.NotNil(t, week2Habits.Habits["Read"])
	assert.False(t, week2Habits.Habits["Exercise"].Completed)
	assert.False(t, week2Habits.Habits["Read"].Completed) // Should be reset to incomplete
}

// Helper function
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}