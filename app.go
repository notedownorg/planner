package main

import (
	"context"
	"fmt"
	"time"

	"github.com/notedownorg/planner/pkg/config"
	"github.com/notedownorg/planner/pkg/habits"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	habitService *habits.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	
	// Initialize habit service with config
	cfg, err := config.Load()
	if err != nil {
		// Use default config if loading fails
		cfg = config.NewConfigWithDefaults()
	}
	a.habitService = habits.NewService(cfg)
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetConfig returns the current configuration
func (a *App) GetConfig() (*config.Config, error) {
	return config.Load()
}

// SaveConfig saves the configuration to disk
func (a *App) SaveConfig(cfg *config.Config) error {
	return config.Save(cfg)
}

// SelectWorkspaceDirectory opens a native directory picker
func (a *App) SelectWorkspaceDirectory() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Select Workspace Directory",
	}

	result, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return "", err
	}

	return result, nil
}

// ValidateWorkspacePath validates if the given path is suitable for a workspace
func (a *App) ValidateWorkspacePath(path string) error {
	return config.ValidateWorkspacePath(path)
}

// GetCurrentWeekHabits returns habits for the current week
func (a *App) GetCurrentWeekHabits() (*habits.WeeklyHabits, error) {
	if a.habitService == nil {
		return nil, fmt.Errorf("habit service not initialized")
	}
	return a.habitService.GetCurrentWeekHabits()
}

// ToggleHabit toggles the completion status of a habit for the current week
func (a *App) ToggleHabit(habitName string) error {
	if a.habitService == nil {
		return fmt.Errorf("habit service not initialized")
	}
	
	// Get current week info
	year, week := getCurrentWeekInfo()
	return a.habitService.ToggleHabit(year, week, habitName)
}

// AddHabit adds a new habit to the current week
func (a *App) AddHabit(habitName string) error {
	if a.habitService == nil {
		return fmt.Errorf("habit service not initialized")
	}
	
	// Get current week info
	year, week := getCurrentWeekInfo()
	return a.habitService.AddHabit(year, week, habitName)
}

// RemoveHabit removes a habit from the current week
func (a *App) RemoveHabit(habitName string) error {
	if a.habitService == nil {
		return fmt.Errorf("habit service not initialized")
	}
	
	// Get current week info
	year, week := getCurrentWeekInfo()
	return a.habitService.RemoveHabit(year, week, habitName)
}

// ReorderHabits reorders habits for the current week
func (a *App) ReorderHabits(habitNames []string) error {
	if a.habitService == nil {
		return fmt.Errorf("habit service not initialized")
	}
	
	// Get current week info
	year, week := getCurrentWeekInfo()
	return a.habitService.ReorderHabits(year, week, habitNames)
}

// getCurrentWeekInfo gets the current ISO week information
func getCurrentWeekInfo() (int, int) {
	// This should match the frontend implementation
	// Using time.Now().ISOWeek() for consistency
	now := time.Now()
	year, week := now.ISOWeek()
	return year, week
}
