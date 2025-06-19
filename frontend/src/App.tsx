import { NavBar } from "@/components/NavBar"
import { HomeView } from "@/components/views/HomeView"
import { SettingsView } from "@/components/views/SettingsView"
import { SetupView } from "@/components/views/SetupView"
import { ThemeProvider } from "@/components/providers/ThemeProvider"
import React from "react"
import { GetConfig } from "wailsjs/go/main/App"
import { config } from "wailsjs/go/models"

type View = "home" | "settings"

function App() {
  const [currentView, setCurrentView] = React.useState<View>("settings")
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
      const hasWorkspace = cfg.WorkspaceRoot && cfg.WorkspaceRoot.trim() !== ""
      setRequiresSetup(!hasWorkspace)
      
      if (hasWorkspace) {
        setCurrentView("home")
      }
    } catch (error) {
      console.error("Failed to load configuration:", error)
      setRequiresSetup(true)
      setCurrentView("settings")
    }
  }

  const navigateToHome = () => {
    if (!requiresSetup) {
      setCurrentView("home")
    }
  }

  const navigateToSettings = () => {
    setCurrentView("settings")
  }

  const onConfigurationComplete = () => {
    setRequiresSetup(false)
    setCurrentView("home")
    loadConfiguration() // Reload config
  }

  const renderCurrentView = () => {
    if (requiresSetup) {
      return (
        <SetupView 
          config={appConfig}
          onConfigurationComplete={onConfigurationComplete}
        />
      )
    }

    switch (currentView) {
      case "home":
        return <HomeView />
      case "settings":
        return <SettingsView config={appConfig} />
      default:
        return <HomeView />
    }
  }

  return (
    <ThemeProvider defaultTheme="system">
      <div className="min-h-screen bg-background flex flex-col">
        <NavBar
          onNavigateToHome={navigateToHome}
          onNavigateToSettings={navigateToSettings}
          requiresSetup={requiresSetup}
        />
        {renderCurrentView()}
      </div>
    </ThemeProvider>
  )
}

export default App
