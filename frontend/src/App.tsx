import { NavBar } from '@/components/NavBar'
import { SettingsView } from '@/components/views/SettingsView'
import { SetupView } from '@/components/views/SetupView'
import { WeeklyView } from '@/components/views/WeeklyView'
import { ThemeProvider } from '@/components/providers/ThemeProvider'
import React from 'react'
import { GetConfig } from 'wailsjs/go/main/App'
import { config } from 'wailsjs/go/models'

type View = 'weekly' | 'settings'

function App() {
    const [currentView, setCurrentView] = React.useState<View>('settings')
    const [appConfig, setAppConfig] = React.useState<config.Config | null>(null)
    const [requiresSetup, setRequiresSetup] = React.useState(true)

    React.useEffect(() => {
        loadConfiguration()
    }, [])

    const loadConfiguration = async () => {
        try {
            const cfg = await GetConfig()
            setAppConfig(cfg)

            // Check if workspace is configured
            const hasWorkspace = cfg.WorkspaceRoot && cfg.WorkspaceRoot.trim() !== ''
            setRequiresSetup(!hasWorkspace)

            if (hasWorkspace) {
                setCurrentView('weekly')
            }
        } catch (error) {
            console.error('Failed to load configuration:', error)
            setRequiresSetup(true)
            setCurrentView('settings')
        }
    }

    const navigateToWeeklyView = () => {
        if (!requiresSetup) {
            setCurrentView('weekly')
        }
    }

    const navigateToSettings = () => {
        setCurrentView('settings')
    }

    const onConfigurationComplete = () => {
        setRequiresSetup(false)
        setCurrentView('weekly')
        loadConfiguration() // Reload config
    }

    const renderCurrentView = () => {
        if (requiresSetup) {
            return (
                <SetupView config={appConfig} onConfigurationComplete={onConfigurationComplete} />
            )
        }

        switch (currentView) {
            case 'weekly':
                return <WeeklyView config={appConfig} />
            case 'settings':
                return <SettingsView config={appConfig} />
            default:
                return <WeeklyView config={appConfig} />
        }
    }

    return (
        <ThemeProvider defaultTheme="system">
            <div className="min-h-screen bg-background flex flex-col">
                <NavBar
                    onNavigateToWeeklyView={navigateToWeeklyView}
                    onNavigateToSettings={navigateToSettings}
                    requiresSetup={requiresSetup}
                />
                {renderCurrentView()}
            </div>
        </ThemeProvider>
    )
}

export default App
