import { Button } from "@/components/ui/button"
import React from "react"
import { SaveConfig, SelectWorkspaceDirectory, ValidateWorkspacePath } from "wailsjs/go/main/App"
import { config as configTypes } from "wailsjs/go/models"

interface SettingsViewProps {
  config: configTypes.Config | null
}

export function SettingsView({ config }: SettingsViewProps) {
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
    } catch (error) {
      console.error("Failed to save configuration:", error)
      setValidationError("Failed to save configuration")
    } finally {
      setIsSaving(false)
    }
  }

  const canSave = workspacePath.trim() !== "" && !isValidating && !isSaving

  return (
    <div className="flex-1 max-w-2xl mx-auto py-8 px-8">
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Settings</h1>
          <p className="text-muted-foreground mt-2">Configure your application preferences.</p>
        </div>

        <div className="space-y-6">
          <div className="space-y-4">
            <h2 className="text-xl font-semibold text-foreground">Workspace</h2>
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">
                Workspace Directory
              </label>
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={workspacePath}
                  onChange={(e) => setWorkspacePath(e.target.value)}
                  className="flex-1 px-3 py-2 border border-border rounded-md bg-background text-foreground"
                  readOnly
                />
                <Button onClick={handleBrowseWorkspace} variant="outline">
                  Change
                </Button>
              </div>
              {validationError && (
                <p className="text-red-500 text-sm mt-1">{validationError}</p>
              )}
            </div>
          </div>

          <Button 
            onClick={handleSaveConfig} 
            disabled={!canSave}
          >
            {isSaving ? "Saving..." : "Save Settings"}
          </Button>
        </div>
      </div>
    </div>
  )
}