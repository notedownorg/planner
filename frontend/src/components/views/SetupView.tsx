import { Button } from "@/components/ui/button"
import React from "react"
import { SaveConfig, SelectWorkspaceDirectory, ValidateWorkspacePath } from "wailsjs/go/main/App"
import { config as configTypes } from "wailsjs/go/models"

interface SetupViewProps {
  config: configTypes.Config | null
  onConfigurationComplete: () => void
}

export function SetupView({ config, onConfigurationComplete }: SetupViewProps) {
  const [workspacePath, setWorkspacePath] = React.useState("")
  const [isValidating, setIsValidating] = React.useState(false)
  const [validationError, setValidationError] = React.useState("")
  const [isSaving, setIsSaving] = React.useState(false)

  React.useEffect(() => {
    if (config) {
      setWorkspacePath(config.WorkspaceRoot || "")
    }
  }, [config])

  const handleBrowseWorkspace = async () => {
    try {
      const selectedPath = await SelectWorkspaceDirectory()
      if (selectedPath) {
        setWorkspacePath(selectedPath)
        setValidationError("")
      }
    } catch (error) {
      console.error("Failed to select workspace:", error)
    }
  }

  const validateWorkspace = async (path: string) => {
    if (!path.trim()) {
      setValidationError("Please select a workspace directory")
      return false
    }

    setIsValidating(true)
    try {
      await ValidateWorkspacePath(path)
      setValidationError("")
      return true
    } catch (error) {
      setValidationError("Invalid workspace directory. Please ensure the path exists and is writable.")
      return false
    } finally {
      setIsValidating(false)
    }
  }

  const handleSaveConfig = async () => {
    const isValid = await validateWorkspace(workspacePath)
    if (!isValid) return

    setIsSaving(true)
    try {
      const updatedConfig = configTypes.Config.createFrom({
        WorkspaceRoot: workspacePath
      })

      await SaveConfig(updatedConfig)
      onConfigurationComplete()
    } catch (error) {
      console.error("Failed to save configuration:", error)
      setValidationError("Failed to save configuration")
    } finally {
      setIsSaving(false)
    }
  }

  const canSave = workspacePath.trim() !== "" && !isValidating && !isSaving

  return (
    <div className="flex-1 flex items-center justify-center px-8">
      <div className="max-w-md w-full space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground">Welcome to Notedown Planner</h1>
          <p className="text-muted-foreground mt-2">Please configure your workspace to get started.</p>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-2">
              Workspace Directory *
            </label>
            <div className="flex space-x-2">
              <input
                type="text"
                value={workspacePath}
                onChange={(e) => setWorkspacePath(e.target.value)}
                placeholder="Select your workspace folder..."
                className="flex-1 px-3 py-2 border border-border rounded-md bg-background text-foreground"
                readOnly
              />
              <Button onClick={handleBrowseWorkspace} variant="outline">
                Browse
              </Button>
            </div>
            {validationError && (
              <p className="text-red-500 text-sm mt-1">❌ {validationError}</p>
            )}
            {isValidating && (
              <p className="text-muted-foreground text-sm mt-1">⏳ Validating workspace...</p>
            )}
          </div>

          <Button 
            onClick={handleSaveConfig} 
            disabled={!canSave}
            className="w-full"
          >
            {isSaving ? "Setting up..." : "Complete Setup"}
          </Button>
        </div>
      </div>
    </div>
  )
}