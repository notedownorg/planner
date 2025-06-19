package habits

import "time"

// Habit represents a single habit with a name, completion status, and order
type Habit struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
}

// WeeklyHabits represents habits for a specific week
type WeeklyHabits struct {
	Year        int                `json:"year"`
	WeekNumber  int                `json:"week_number"`
	Habits      map[string]*Habit  `json:"habits"` // key is habit name
	DayStatus   map[string]bool    `json:"day_status"` // tracks which days have been marked
}

// HabitDay represents habits for a specific day
type HabitDay struct {
	Date   time.Time          `json:"date"`
	Habits map[string]*Habit  `json:"habits"` // key is habit name
}

// HabitConfig represents the configuration for habits
type HabitConfig struct {
	DefaultHabits []string `json:"default_habits"` // list of habit names
}