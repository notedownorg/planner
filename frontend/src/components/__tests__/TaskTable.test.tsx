import { render, screen, fireEvent, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import '@testing-library/jest-dom'
import { jest } from '@jest/globals'
import { TaskTable } from '../TaskTable'
import { habits } from 'wailsjs/go/models'

// Mock the Wails API functions
jest.mock('wailsjs/go/main/App', () => ({
    GetCurrentWeekHabits: jest.fn(),
    ToggleHabit: jest.fn(),
    AddHabit: jest.fn(),
    RemoveHabit: jest.fn(),
    ReorderHabits: jest.fn(),
}))

import * as App from 'wailsjs/go/main/App'

describe('TaskTable', () => {
    const mockHabits = new habits.WeeklyHabits({
        year: 2024,
        week_number: 1,
        habits: {
            Exercise: { name: 'Exercise', completed: false, order: 0 },
            Read: { name: 'Read', completed: false, order: 1 },
            Meditate: { name: 'Meditate', completed: true, order: 2 },
        },
        day_status: {},
    })

    beforeEach(() => {
        jest.clearAllMocks()
        ;(App.GetCurrentWeekHabits as jest.MockedFunction<typeof App.GetCurrentWeekHabits>).mockResolvedValue(mockHabits)
        ;(App.ToggleHabit as jest.MockedFunction<typeof App.ToggleHabit>).mockResolvedValue(undefined)
        ;(App.AddHabit as jest.MockedFunction<typeof App.AddHabit>).mockResolvedValue(undefined)
        ;(App.RemoveHabit as jest.MockedFunction<typeof App.RemoveHabit>).mockResolvedValue(undefined)
        ;(App.ReorderHabits as jest.MockedFunction<typeof App.ReorderHabits>).mockResolvedValue(undefined)
    })

    describe('Rendering', () => {
        it('should render loading state initially', () => {
            render(<TaskTable />)
            expect(screen.getByText('Loading habits...')).toBeInTheDocument()
        })

        it('should render habits after loading', async () => {
            await act(async () => {
                render(<TaskTable />)
            })

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
                expect(screen.getByText('Read')).toBeInTheDocument()
                expect(screen.getByText('Meditate')).toBeInTheDocument()
            })
        })

        it('should render incomplete habits before completed ones', async () => {
            render(<TaskTable />)

            await waitFor(() => {
                const habits = screen.getAllByText(/Exercise|Read|Meditate/)
                expect(habits[0]).toHaveTextContent('Exercise')
                expect(habits[1]).toHaveTextContent('Read')
                expect(habits[2]).toHaveTextContent('Meditate')
            })
        })

        it('should show separator between incomplete and completed habits', async () => {
            render(<TaskTable />)

            await waitFor(() => {
                const separator = document.querySelector('.w-px.h-6.bg-gray-300')
                expect(separator).toBeInTheDocument()
            })
        })

        it('should show add button after incomplete habits', async () => {
            render(<TaskTable />)

            await waitFor(() => {
                const addButton = screen.getByTitle('Add new habit')
                expect(addButton).toBeInTheDocument()
            })
        })

        it('should render error state when loading fails', async () => {
            ;(App.GetCurrentWeekHabits as jest.MockedFunction<typeof App.GetCurrentWeekHabits>).mockRejectedValue(new Error('Network error'))

            await act(async () => {
                render(<TaskTable />)
            })

            await waitFor(() => {
                expect(screen.getByText('Failed to load habits')).toBeInTheDocument()
                expect(screen.getByText('Retry')).toBeInTheDocument()
            })
        })
    })

    describe('Habit Completion', () => {
        it('should toggle habit completion status', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const exerciseCheckbox = screen.getAllByTitle('Mark as complete')[0]
            await user.click(exerciseCheckbox)

            expect(App.ToggleHabit).toHaveBeenCalledWith('Exercise')
            expect(App.GetCurrentWeekHabits).toHaveBeenCalledTimes(2) // Initial load + reload after toggle
        })

        it('should show correct emoji for completion status', async () => {
            render(<TaskTable />)

            await waitFor(() => {
                const incompleteEmojis = screen.getAllByText('⭕')
                const completeEmojis = screen.getAllByText('✅')

                expect(incompleteEmojis).toHaveLength(2) // Exercise and Read
                expect(completeEmojis).toHaveLength(1) // Meditate
            })
        })
    })

    describe('Adding Habits', () => {
        it('should show input field when add button is clicked', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByTitle('Add new habit')).toBeInTheDocument()
            })

            await user.click(screen.getByTitle('Add new habit'))

            expect(screen.getByPlaceholderText('Add habit...')).toBeInTheDocument()
        })

        it('should add new habit when form is submitted', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByTitle('Add new habit')).toBeInTheDocument()
            })

            await user.click(screen.getByTitle('Add new habit'))
            const input = screen.getByPlaceholderText('Add habit...')

            await user.type(input, 'New Habit')
            await user.keyboard('{Enter}')

            expect(App.AddHabit).toHaveBeenCalledWith('New Habit')
        })

        it('should cancel adding habit on Escape', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByTitle('Add new habit')).toBeInTheDocument()
            })

            await user.click(screen.getByTitle('Add new habit'))
            await user.keyboard('{Escape}')

            expect(screen.getByTitle('Add new habit')).toBeInTheDocument()
            expect(screen.queryByPlaceholderText('Add habit...')).not.toBeInTheDocument()
        })

        it('should show checkmark button when text is entered', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByTitle('Add new habit')).toBeInTheDocument()
            })

            await user.click(screen.getByTitle('Add new habit'))
            const input = screen.getByPlaceholderText('Add habit...')

            await user.type(input, 'Test')

            expect(screen.getByTitle('Add habit')).toBeInTheDocument()
            expect(screen.getAllByText('✅')).toContain(screen.getByTitle('Add habit'))
        })
    })

    describe('Editing Habits', () => {
        it('should enter edit mode on double click', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!
            await user.dblClick(habitPill)

            const input = screen.getByDisplayValue('Exercise')
            expect(input).toBeInTheDocument()
            expect(input).toHaveFocus()
        })

        it('should save edited habit name', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!
            await user.dblClick(habitPill)

            const input = screen.getByDisplayValue('Exercise')
            await user.clear(input)
            await user.type(input, 'Workout')
            await user.keyboard('{Enter}')

            expect(App.AddHabit).toHaveBeenCalledWith('Workout')
            expect(App.RemoveHabit).toHaveBeenCalledWith('Exercise')
        })

        it('should cancel edit on Escape', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!
            await user.dblClick(habitPill)

            const input = screen.getByDisplayValue('Exercise')
            await user.clear(input)
            await user.type(input, 'Workout')
            await user.keyboard('{Escape}')

            expect(screen.getByText('Exercise')).toBeInTheDocument()
            expect(App.AddHabit).not.toHaveBeenCalled()
        })
    })

    describe('Deleting Habits', () => {
        it('should show delete button on hover', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!
            await user.hover(habitPill)

            const deleteButtons = screen.getAllByTitle('Remove habit')
            expect(deleteButtons.length).toBeGreaterThan(0)
        })

        it('should delete habit when delete button is clicked', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!
            await user.hover(habitPill)

            const deleteButtons = screen.getAllByTitle('Remove habit')
            const deleteButton = deleteButtons.find(btn => 
                btn.closest('[data-testid]')?.getAttribute('data-testid')?.includes('Exercise') ||
                btn.parentElement?.textContent?.includes('Exercise')
            ) || deleteButtons[0]
            await user.click(deleteButton)

            expect(App.RemoveHabit).toHaveBeenCalledWith('Exercise')
        })
    })

    describe('Drag and Drop', () => {
        it('should show drag handle on hover', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!
            await user.hover(habitPill)

            const dragHandles = screen.getAllByTitle('Drag to reorder')
            expect(dragHandles.length).toBeGreaterThan(0)
        })

        it('should handle drag and drop reordering', async () => {
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const exercisePill = screen.getByText('Exercise').closest('div')!
            const readPill = screen.getByText('Read').closest('div')!

            // Simulate drag start
            fireEvent.dragStart(exercisePill, {
                dataTransfer: {
                    setData: jest.fn(),
                },
            })

            // Simulate drag over
            fireEvent.dragOver(readPill, {
                preventDefault: jest.fn(),
            })

            // Simulate drop
            fireEvent.drop(readPill, {
                preventDefault: jest.fn(),
            })

            await waitFor(() => {
                expect(App.ReorderHabits).toHaveBeenCalledWith(['Read', 'Exercise', 'Meditate'])
            })
        })
    })

    describe('Responsive Layout', () => {
        it('should expand pill padding on hover', async () => {
            const user = userEvent.setup()
            render(<TaskTable />)

            await waitFor(() => {
                expect(screen.getByText('Exercise')).toBeInTheDocument()
            })

            const habitPill = screen.getByText('Exercise').closest('div')!

            expect(habitPill).toHaveClass('px-3')

            await user.hover(habitPill)

            expect(habitPill).toHaveClass('px-7')
        })

        it('should show placeholder text when no habits exist', async () => {
            ;(App.GetCurrentWeekHabits as jest.MockedFunction<typeof App.GetCurrentWeekHabits>).mockResolvedValue(new habits.WeeklyHabits({
                year: 2024,
                week_number: 1,
                habits: {},
                day_status: {},
            }))

            const user = userEvent.setup()
            await act(async () => {
                render(<TaskTable />)
            })

            await waitFor(() => {
                expect(screen.getByTitle('Add new habit')).toBeInTheDocument()
            })

            await user.click(screen.getByTitle('Add new habit'))

            expect(screen.getByPlaceholderText('Add your first habit...')).toBeInTheDocument()
        })
    })
})
