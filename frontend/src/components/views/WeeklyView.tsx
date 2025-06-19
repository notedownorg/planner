import { config as configTypes } from 'wailsjs/go/models'
import { TaskTable } from '../TaskTable'

interface WeeklyViewProps {
    config: configTypes.Config | null
}

export function WeeklyView({ config }: WeeklyViewProps) {
    const enabledComponents = config?.WeeklyView?.EnabledComponents || {
        HabitTracker: true,
    }

    // Get current week dates
    const getCurrentWeekDates = () => {
        const today = new Date()
        const dayOfWeek = today.getDay()
        const mondayOffset = dayOfWeek === 0 ? -6 : 1 - dayOfWeek // Convert Sunday (0) to -6, Monday to 0
        
        const monday = new Date(today)
        monday.setDate(today.getDate() + mondayOffset)
        
        const sunday = new Date(monday)
        sunday.setDate(monday.getDate() + 6)
        
        return { monday, sunday }
    }

    const getWeekInfo = () => {
        const { monday } = getCurrentWeekDates()
        
        // Calculate ISO week number
        const getISOWeekNumber = (date: Date) => {
            const target = new Date(date.valueOf())
            const dayNumber = (date.getDay() + 6) % 7
            target.setDate(target.getDate() - dayNumber + 3)
            const firstThursday = target.valueOf()
            target.setMonth(0, 1)
            if (target.getDay() !== 4) {
                target.setMonth(0, 1 + ((4 - target.getDay()) + 7) % 7)
            }
            return 1 + Math.ceil((firstThursday - target.valueOf()) / 604800000)
        }
        
        const weekNumber = getISOWeekNumber(monday)
        const year = monday.getFullYear()
        
        return { weekNumber, year }
    }

    return (
        <div className="flex-1 max-w-7xl mx-auto py-8 px-8">
            <div className="space-y-4">
                <div>
                    <h1 className="text-3xl font-bold text-foreground">
                        Week {getWeekInfo().weekNumber.toString().padStart(2, '0')} 
                        <span className="text-muted-foreground ml-2 font-normal">{getWeekInfo().year}</span>
                    </h1>
                </div>

                {enabledComponents.HabitTracker && <TaskTable />}
            </div>
        </div>
    )
}