import { Button } from '@/components/ui/button'
import { useTheme } from '@/components/providers/ThemeProvider'
import { Moon, Sun, Settings, CalendarRange } from 'lucide-react'

interface NavBarProps {
    onNavigateToSettings: () => void
    onNavigateToWeeklyView: () => void
    requiresSetup: boolean
}

export function NavBar({ onNavigateToSettings, onNavigateToWeeklyView, requiresSetup }: NavBarProps) {
    const { theme, setTheme } = useTheme()

    return (
        <nav className="w-full border-b bg-background">
            <div className="flex h-16 items-center px-4 justify-between">
                <div className="flex items-center">
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={onNavigateToWeeklyView}
                        disabled={requiresSetup}
                    >
                        <CalendarRange className="h-[1.2rem] w-[1.2rem]" />
                        <span className="sr-only">Weekly view</span>
                    </Button>
                </div>
                <div className="flex items-center space-x-2">
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
                    >
                        <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                        <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                        <span className="sr-only">Toggle theme</span>
                    </Button>
                    <Button variant="ghost" size="icon" onClick={onNavigateToSettings}>
                        <Settings className="h-[1.2rem] w-[1.2rem]" />
                        <span className="sr-only">Settings</span>
                    </Button>
                </div>
            </div>
        </nav>
    )
}
