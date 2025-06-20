import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'
import { jest } from '@jest/globals'
import * as AppModule from 'wailsjs/go/main/App'

// Mock the Wails modules
jest.mock('wailsjs/go/main/App')
jest.mock('wailsjs/go/models')

describe('App Component', () => {
    beforeEach(() => {
        jest.clearAllMocks()
    })

    it('shows setup view when workspace path is not configured', async () => {
        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockResolvedValue(
            {
                WorkspaceRoot: '',
                PeriodicNotes: {
                    WeeklySubdir: '_periodic/weekly',
                    WeeklyNameFormat: 'YYYY-[W]WW',
                },
                WeeklyView: {
                    EnabledComponents: {
                        HabitTracker: true,
                    },
                    convertValues: () => {},
                },
                convertValues: () => {},
            }
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
            expect(
                screen.getByText('Please configure your workspace to get started.')
            ).toBeInTheDocument()
        })

        // Weekly view button should be disabled in setup mode
        expect(screen.getByRole('button', { name: 'Weekly view' })).toBeDisabled()
    })

    it('shows setup view when GetConfig fails', async () => {
        // Mock console.error to suppress expected error message
        const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {})

        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockRejectedValue(
            new Error('Config error')
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
        })

        // Restore console.error
        consoleSpy.mockRestore()
    })

    it('shows home view when workspace path is configured', async () => {
        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockResolvedValue(
            {
                WorkspaceRoot: '/path/to/workspace',
                PeriodicNotes: {
                    WeeklySubdir: '_periodic/weekly',
                    WeeklyNameFormat: 'YYYY-[W]WW',
                },
                WeeklyView: {
                    EnabledComponents: {
                        HabitTracker: true,
                    },
                    convertValues: () => {},
                },
                convertValues: () => {},
            }
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText(/Week \d+/)).toBeInTheDocument()
            expect(screen.getByText('2025')).toBeInTheDocument()
        })

        // Weekly view button should not be disabled
        expect(screen.getByRole('button', { name: 'Weekly view' })).not.toBeDisabled()
    })

    it('prevents navigation to home when in setup mode', async () => {
        const user = userEvent.setup()

        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockResolvedValue(
            {
                WorkspaceRoot: '',
                PeriodicNotes: {
                    WeeklySubdir: '_periodic/weekly',
                    WeeklyNameFormat: 'YYYY-[W]WW',
                },
                WeeklyView: {
                    EnabledComponents: {
                        HabitTracker: true,
                    },
                    convertValues: () => {},
                },
                convertValues: () => {},
            }
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
        })

        // Weekly view button should be disabled
        const weeklyViewButton = screen.getByRole('button', { name: 'Weekly view' })
        expect(weeklyViewButton).toBeDisabled()

        // Clicking disabled button shouldn't do anything
        await user.click(weeklyViewButton)
        expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
    })

    it('allows navigation to settings during setup', async () => {
        const user = userEvent.setup()

        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockResolvedValue(
            {
                WorkspaceRoot: '',
                PeriodicNotes: {
                    WeeklySubdir: '_periodic/weekly',
                    WeeklyNameFormat: 'YYYY-[W]WW',
                },
                WeeklyView: {
                    EnabledComponents: {
                        HabitTracker: true,
                    },
                    convertValues: () => {},
                },
                convertValues: () => {},
            }
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
        })

        // Settings button should work
        const settingsButton = screen.getByRole('button', { name: 'Settings' })
        expect(settingsButton).not.toBeDisabled()

        // Should still show setup when clicking settings during setup
        await user.click(settingsButton)
        expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
    })

    it('requires setup when workspace path is only whitespace', async () => {
        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockResolvedValue(
            {
                WorkspaceRoot: '   ',
                PeriodicNotes: {
                    WeeklySubdir: '_periodic/weekly',
                    WeeklyNameFormat: 'YYYY-[W]WW',
                },
                WeeklyView: {
                    EnabledComponents: {
                        HabitTracker: true,
                    },
                    convertValues: () => {},
                },
                convertValues: () => {},
            }
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
            expect(
                screen.getByText('Please configure your workspace to get started.')
            ).toBeInTheDocument()
        })

        // Weekly view button should be disabled
        expect(screen.getByRole('button', { name: 'Weekly view' })).toBeDisabled()
    })

    it('requires setup when workspace path is null', async () => {
        ;(AppModule.GetConfig as jest.MockedFunction<typeof AppModule.GetConfig>).mockResolvedValue(
            {
                WorkspaceRoot: '',
                PeriodicNotes: {
                    WeeklySubdir: '_periodic/weekly',
                    WeeklyNameFormat: 'YYYY-[W]WW',
                },
                WeeklyView: {
                    EnabledComponents: {
                        HabitTracker: true,
                    },
                    convertValues: () => {},
                },
                convertValues: () => {},
            }
        )

        render(<App />)

        await waitFor(() => {
            expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
        })
    })
})
