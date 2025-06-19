import React, { useState, useEffect } from 'react'
import { habits as habitTypes } from 'wailsjs/go/models'
import { GetCurrentWeekHabits, ToggleHabit, AddHabit, RemoveHabit, ReorderHabits } from 'wailsjs/go/main/App'
import { Button } from './ui/button'
import { GripVertical, X } from 'lucide-react'

interface HabitPillProps {
    habit: habitTypes.Habit
    index: number
    onToggle: (habitName: string) => void
    onDelete: (habitName: string) => void
    onEdit: (oldName: string, newName: string) => void
    onDragStart: (index: number) => void
    onDragOver: (index: number) => void
    onDrop: () => void
    isDragging: boolean
    dragOverIndex: number | null
    draggedIndex: number | null
}

const HabitPill = ({ 
    habit, 
    index, 
    onToggle, 
    onDelete, 
    onEdit, 
    onDragStart, 
    onDragOver, 
    onDrop,
    isDragging,
    dragOverIndex,
    draggedIndex
}: HabitPillProps) => {
    const [isHovered, setIsHovered] = useState(false)
    const [isEditing, setIsEditing] = useState(false)
    const [editValue, setEditValue] = useState(habit.name)

    const handleDoubleClick = () => {
        setIsEditing(true)
        setEditValue(habit.name)
    }

    const handleEditSave = () => {
        if (editValue.trim() && editValue.trim() !== habit.name) {
            onEdit(habit.name, editValue.trim())
        }
        setIsEditing(false)
    }

    const handleEditCancel = () => {
        setEditValue(habit.name)
        setIsEditing(false)
    }

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleEditSave()
        } else if (e.key === 'Escape') {
            handleEditCancel()
        }
    }

    // Check if this is a valid drop target
    const isValidDropTarget = draggedIndex !== null && draggedIndex !== index
    
    const pillClasses = `
        relative inline-flex items-center gap-2 py-2 rounded-full transition-all cursor-pointer
        ${isHovered && !isEditing ? 'px-7' : 'px-3'}
        ${habit.completed 
            ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-300' 
            : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'
        }
        ${isDragging ? 'opacity-50' : 'hover:scale-105'}
        ${dragOverIndex === index && isValidDropTarget ? 'ring-2 ring-blue-400' : ''}
    `

    return (
        <div
            className={pillClasses}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
            onDoubleClick={handleDoubleClick}
            draggable={isHovered && !isEditing}
            onDragStart={(e) => {
                e.dataTransfer.setData('text/plain', index.toString())
                onDragStart(index)
            }}
            onDragOver={(e) => {
                e.preventDefault()
                onDragOver(index)
            }}
            onDrop={(e) => {
                e.preventDefault()
                onDrop()
            }}
        >
            {/* Drag Handle - Absolute positioned */}
            {isHovered && !isEditing && (
                <div 
                    className="absolute left-1 top-1/2 transform -translate-y-1/2"
                    title="Drag to reorder"
                >
                    <GripVertical className="h-3 w-3 text-gray-400 cursor-grab" />
                </div>
            )}

            {/* Completion Status */}
            <button
                onClick={(e) => {
                    e.stopPropagation()
                    onToggle(habit.name)
                }}
                className="flex-shrink-0 text-lg leading-none transition-all hover:scale-110"
                title={habit.completed ? "Mark as incomplete" : "Mark as complete"}
            >
                {habit.completed ? '✅' : '⭕'}
            </button>

            {/* Habit Name */}
            {isEditing ? (
                <input
                    type="text"
                    value={editValue}
                    onChange={(e) => setEditValue(e.target.value)}
                    onBlur={handleEditSave}
                    onKeyDown={handleKeyDown}
                    className="bg-transparent border-none outline-none text-sm font-medium min-w-0 flex-1"
                    autoFocus
                />
            ) : (
                <span 
                    className="text-sm font-medium select-none"
                    title="Double-click to edit"
                >
                    {habit.name}
                </span>
            )}

            {/* Delete Button - Absolute positioned */}
            {isHovered && !isEditing && (
                <button
                    onClick={(e) => {
                        e.stopPropagation()
                        onDelete(habit.name)
                    }}
                    className="absolute right-1 top-1/2 transform -translate-y-1/2 text-red-500 hover:text-red-700 transition-colors"
                    title="Remove habit"
                >
                    <X className="h-3 w-3" />
                </button>
            )}
        </div>
    )
}

export const TaskTable = () => {
    const [weeklyHabits, setWeeklyHabits] = useState<habitTypes.WeeklyHabits | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [newHabitName, setNewHabitName] = useState('')
    const [draggedIndex, setDraggedIndex] = useState<number | null>(null)
    const [dragOverIndex, setDragOverIndex] = useState<number | null>(null)
    const [isAddingHabit, setIsAddingHabit] = useState(false)

    // Load habits on component mount
    useEffect(() => {
        loadHabits()
    }, [])

    const loadHabits = async () => {
        try {
            setLoading(true)
            setError(null)
            const habits = await GetCurrentWeekHabits()
            setWeeklyHabits(habits)
        } catch (err) {
            console.error('Failed to load habits:', err)
            setError('Failed to load habits')
        } finally {
            setLoading(false)
        }
    }

    const getSortedHabits = () => {
        if (!weeklyHabits?.habits) return []
        
        return Object.values(weeklyHabits.habits)
            .sort((a, b) => {
                // First sort by completion status (incomplete first)
                if (a.completed !== b.completed) {
                    return a.completed ? 1 : -1
                }
                // Then by order within each group
                return (a.order || 0) - (b.order || 0)
            })
    }

    const handleToggleHabit = async (habitName: string) => {
        try {
            await ToggleHabit(habitName)
            await loadHabits()
        } catch (err) {
            console.error('Failed to toggle habit:', err)
            setError('Failed to update habit')
        }
    }

    const handleAddHabit = async (habitName: string) => {
        if (!habitName.trim()) return
        
        try {
            await AddHabit(habitName.trim())
            await loadHabits()
            setNewHabitName('')
            setIsAddingHabit(false)
        } catch (err) {
            console.error('Failed to add habit:', err)
            setError('Failed to add habit')
        }
    }

    const handleRemoveHabit = async (habitName: string) => {
        try {
            await RemoveHabit(habitName)
            await loadHabits()
        } catch (err) {
            console.error('Failed to remove habit:', err)
            setError('Failed to remove habit')
        }
    }

    const handleEditHabit = async (oldName: string, newName: string) => {
        if (oldName === newName) return
        
        try {
            // Add new habit and remove old one
            await AddHabit(newName)
            await RemoveHabit(oldName)
            await loadHabits()
        } catch (err) {
            console.error('Failed to edit habit:', err)
            setError('Failed to edit habit')
        }
    }

    const handleDragStart = (index: number) => {
        setDraggedIndex(index)
    }

    const handleDragOver = (index: number) => {
        setDragOverIndex(index)
    }

    const handleDrop = async () => {
        if (draggedIndex === null || dragOverIndex === null || draggedIndex === dragOverIndex) {
            setDraggedIndex(null)
            setDragOverIndex(null)
            return
        }

        const sortedHabits = getSortedHabits()
        const reorderedHabits = [...sortedHabits]
        
        // Move the dragged item to the new position
        const [draggedItem] = reorderedHabits.splice(draggedIndex, 1)
        reorderedHabits.splice(dragOverIndex, 0, draggedItem)
        
        // Update order for all habits based on new positions
        const habitUpdates = reorderedHabits.map((habit, index) => ({
            ...habit,
            order: index
        }))
        
        // Extract names in new order
        const newOrder = habitUpdates.map(habit => habit.name)
        
        try {
            await ReorderHabits(newOrder)
            await loadHabits()
        } catch (err) {
            console.error('Failed to reorder habits:', err)
            setError('Failed to reorder habits')
        }

        setDraggedIndex(null)
        setDragOverIndex(null)
    }

    const handleSubmitNewHabit = (e: React.FormEvent) => {
        e.preventDefault()
        if (newHabitName.trim()) {
            handleAddHabit(newHabitName.trim())
        }
    }

    const handleCancelAdd = () => {
        setNewHabitName('')
        setIsAddingHabit(false)
    }

    if (loading) {
        return (
            <div className="border border-border rounded-lg p-3">
                <div className="text-sm text-muted-foreground">Loading habits...</div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="border border-border rounded-lg p-3">
                <div className="text-sm text-red-500 mb-2">{error}</div>
                <Button onClick={loadHabits} variant="outline" size="sm">
                    Retry
                </Button>
            </div>
        )
    }

    const sortedHabits = getSortedHabits()

    // Separate incomplete and completed habits
    const incompleteHabits = sortedHabits.filter(h => !h.completed)
    const completedHabits = sortedHabits.filter(h => h.completed)

    return (
        <div className="p-3">
            <div className="flex flex-wrap gap-2 items-center">
                {/* Incomplete Habit Pills */}
                {incompleteHabits.map((habit, index) => (
                    <HabitPill
                        key={habit.name}
                        habit={habit}
                        index={sortedHabits.indexOf(habit)}
                        onToggle={handleToggleHabit}
                        onDelete={handleRemoveHabit}
                        onEdit={handleEditHabit}
                        onDragStart={handleDragStart}
                        onDragOver={handleDragOver}
                        onDrop={handleDrop}
                        isDragging={draggedIndex === sortedHabits.indexOf(habit)}
                        dragOverIndex={dragOverIndex}
                        draggedIndex={draggedIndex}
                    />
                ))}

                {/* Add New Habit - positioned after incomplete habits */}
                {isAddingHabit ? (
                    <form onSubmit={handleSubmitNewHabit} className="inline-flex items-center gap-1">
                        <input
                            type="text"
                            value={newHabitName}
                            onChange={(e) => setNewHabitName(e.target.value)}
                            onKeyDown={(e) => {
                                if (e.key === 'Escape') {
                                    handleCancelAdd()
                                }
                            }}
                            onBlur={() => {
                                if (!newHabitName.trim()) {
                                    handleCancelAdd()
                                }
                            }}
                            placeholder={sortedHabits.length === 0 ? "Add your first habit..." : "Add habit..."}
                            className="px-3 py-2 text-sm rounded-full bg-gray-100 dark:bg-gray-800 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-400 w-32"
                            autoFocus
                        />
                        {newHabitName.trim() && (
                            <button
                                type="submit"
                                className="flex-shrink-0 text-lg leading-none text-green-600 hover:scale-110 transition-all"
                                title="Add habit"
                            >
                                ✅
                            </button>
                        )}
                    </form>
                ) : (
                    <button
                        onClick={() => setIsAddingHabit(true)}
                        className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:scale-105 transition-all"
                        title="Add new habit"
                    >
                        <span className="text-lg leading-none">+</span>
                    </button>
                )}

                {/* Separator between incomplete and completed */}
                {completedHabits.length > 0 && (
                    <div className="w-px h-6 bg-gray-300 dark:bg-gray-600 mx-1" />
                )}

                {/* Completed Habit Pills */}
                {completedHabits.map((habit) => (
                    <HabitPill
                        key={habit.name}
                        habit={habit}
                        index={sortedHabits.indexOf(habit)}
                        onToggle={handleToggleHabit}
                        onDelete={handleRemoveHabit}
                        onEdit={handleEditHabit}
                        onDragStart={handleDragStart}
                        onDragOver={handleDragOver}
                        onDrop={handleDrop}
                        isDragging={draggedIndex === sortedHabits.indexOf(habit)}
                        dragOverIndex={dragOverIndex}
                        draggedIndex={draggedIndex}
                    />
                ))}
            </div>
        </div>
    )
}