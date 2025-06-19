import { Button } from '@/components/ui/button'
import React from 'react'
import { SaveConfig, SelectWorkspaceDirectory, ValidateWorkspacePath } from 'wailsjs/go/main/App'
import { config as configTypes } from 'wailsjs/go/models'

interface SettingsViewProps {
    config: configTypes.Config | null
}

export function SettingsView({ config }: SettingsViewProps) {
    const [workspacePath, setWorkspacePath] = React.useState('')
    const [weeklySubdir, setWeeklySubdir] = React.useState('_periodic/weekly')
    const [weeklyNameFormat, setWeeklyNameFormat] = React.useState('YYYY-[W]WW')
    const [weeklyComponents, setWeeklyComponents] = React.useState({
        HabitTracker: true,
    })
    const [isValidating, setIsValidating] = React.useState(false)
    const [validationError, setValidationError] = React.useState('')
    const [isSaving, setIsSaving] = React.useState(false)

    React.useEffect(() => {
        if (config) {
            setWorkspacePath(config.WorkspaceRoot || '')
            setWeeklySubdir(config.PeriodicNotes?.WeeklySubdir || '_periodic/weekly')
            setWeeklyNameFormat(config.PeriodicNotes?.WeeklyNameFormat || 'YYYY-[W]WW')
            setWeeklyComponents(
                config.WeeklyView?.EnabledComponents || {
                    HabitTracker: true,
                }
            )
        }
    }, [config])

    const handleBrowseWorkspace = async () => {
        try {
            const selectedPath = await SelectWorkspaceDirectory()
            if (selectedPath) {
                setWorkspacePath(selectedPath)
                setValidationError('')
            }
        } catch (error) {
            console.error('Failed to select workspace:', error)
        }
    }

    const validateWorkspace = async (path: string) => {
        if (!path.trim()) {
            setValidationError('Please select a workspace directory')
            return false
        }

        setIsValidating(true)
        try {
            await ValidateWorkspacePath(path)
            setValidationError('')
            return true
        } catch {
            setValidationError(
                'Invalid workspace directory. Please ensure the path exists and is writable.'
            )
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
                WorkspaceRoot: workspacePath,
                PeriodicNotes: {
                    WeeklySubdir: weeklySubdir,
                    WeeklyNameFormat: weeklyNameFormat,
                },
                WeeklyView: {
                    EnabledComponents: weeklyComponents,
                },
            })

            await SaveConfig(updatedConfig)
        } catch (error) {
            console.error('Failed to save configuration:', error)
            setValidationError('Failed to save configuration')
        } finally {
            setIsSaving(false)
        }
    }

    const canSave = workspacePath.trim() !== '' && !isValidating && !isSaving

    return (
        <div className="flex-1 max-w-2xl mx-auto py-8 px-8">
            <div className="space-y-8">
                <div>
                    <h1 className="text-3xl font-bold text-foreground">Settings</h1>
                    <p className="text-muted-foreground mt-2">
                        Configure your application preferences.
                    </p>
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

                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold text-foreground">Periodic Notes</h2>
                        <div>
                            <div className="flex items-center gap-2 mb-2">
                                <label className="block text-sm font-medium text-foreground">
                                    Weekly Notes Subdirectory
                                </label>
                                <button
                                    type="button"
                                    className="text-muted-foreground hover:text-foreground relative group"
                                    aria-label="Path relative to workspace root where weekly notes will be stored"
                                >
                                    <svg
                                        className="w-4 h-4"
                                        fill="currentColor"
                                        viewBox="0 0 20 20"
                                    >
                                        <path
                                            fillRule="evenodd"
                                            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                                            clipRule="evenodd"
                                        />
                                    </svg>
                                    <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 text-xs text-background bg-foreground rounded opacity-0 group-hover:opacity-100 transition-opacity duration-150 delay-75 pointer-events-none whitespace-nowrap z-10">
                                        Path relative to workspace root where weekly notes will be
                                        stored
                                    </div>
                                </button>
                            </div>
                            <input
                                type="text"
                                value={weeklySubdir}
                                onChange={(e) => setWeeklySubdir(e.target.value)}
                                className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
                                placeholder="_periodic/weekly"
                            />
                        </div>

                        <div>
                            <div className="flex items-center gap-2 mb-2">
                                <label className="block text-sm font-medium text-foreground">
                                    Weekly Note Name Format
                                </label>
                                <button
                                    type="button"
                                    className="text-muted-foreground hover:text-foreground relative group"
                                    aria-label="Format for weekly note filenames"
                                >
                                    <svg
                                        className="w-4 h-4"
                                        fill="currentColor"
                                        viewBox="0 0 20 20"
                                    >
                                        <path
                                            fillRule="evenodd"
                                            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                                            clipRule="evenodd"
                                        />
                                    </svg>
                                    <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 text-xs text-background bg-foreground rounded opacity-0 group-hover:opacity-100 transition-opacity duration-150 delay-75 pointer-events-none z-10 w-80 text-center">
                                        Format for weekly note filenames. Use YYYY for year, WW for
                                        week number, [W] for literal 'W'. Example: YYYY-[W]WW
                                        becomes 2024-W01
                                    </div>
                                </button>
                            </div>
                            <input
                                type="text"
                                value={weeklyNameFormat}
                                onChange={(e) => setWeeklyNameFormat(e.target.value)}
                                className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
                                placeholder="YYYY-[W]WW"
                            />
                        </div>
                    </div>

                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold text-foreground">Weekly View</h2>
                        <div>
                            <div className="flex items-center gap-2 mb-3">
                                <label className="block text-sm font-medium text-foreground">
                                    Enabled Components
                                </label>
                                <button
                                    type="button"
                                    className="text-muted-foreground hover:text-foreground relative group"
                                    aria-label="Choose which components appear in your weekly view"
                                >
                                    <svg
                                        className="w-4 h-4"
                                        fill="currentColor"
                                        viewBox="0 0 20 20"
                                    >
                                        <path
                                            fillRule="evenodd"
                                            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                                            clipRule="evenodd"
                                        />
                                    </svg>
                                    <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 text-xs text-background bg-foreground rounded opacity-0 group-hover:opacity-100 transition-opacity duration-150 delay-75 pointer-events-none whitespace-nowrap z-10">
                                        Choose which components appear in your weekly view
                                    </div>
                                </button>
                            </div>
                            <div className="space-y-3">
                                <label className="flex items-center space-x-2 cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={weeklyComponents.HabitTracker}
                                        onChange={(e) =>
                                            setWeeklyComponents((prev) => ({
                                                ...prev,
                                                HabitTracker: e.target.checked,
                                            }))
                                        }
                                        className="rounded border-border"
                                    />
                                    <span className="text-sm text-foreground">Habit Tracker</span>
                                </label>
                            </div>
                        </div>
                    </div>

                    <Button onClick={handleSaveConfig} disabled={!canSave}>
                        {isSaving ? 'Saving...' : 'Save Settings'}
                    </Button>
                </div>
            </div>
        </div>
    )
}
