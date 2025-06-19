import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { SettingsView } from '../SettingsView'
import { jest } from '@jest/globals'
import * as AppModule from 'wailsjs/go/main/App'

// Mock the Wails modules
jest.mock('wailsjs/go/main/App')
jest.mock('wailsjs/go/models')

describe('SettingsView', () => {
    const mockConfig = {
        WorkspaceRoot: '/Users/test/workspace',
    }

    beforeEach(() => {
        jest.clearAllMocks()
    })

    it('renders settings form with config data', () => {
        render(<SettingsView config={mockConfig} />)

        expect(screen.getByText('Settings')).toBeInTheDocument()
        expect(screen.getByText('Configure your application preferences.')).toBeInTheDocument()
        expect(screen.getByText('Workspace')).toBeInTheDocument()
        expect(screen.getByDisplayValue('/Users/test/workspace')).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'Change' })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'Save Settings' })).toBeInTheDocument()
    })

    it('renders with empty config', () => {
        render(<SettingsView config={null} />)

        expect(screen.getByText('Settings')).toBeInTheDocument()
        expect(screen.getByDisplayValue('')).toBeInTheDocument()
    })

    it('handles workspace directory change', async () => {
        const user = userEvent.setup()
        const newPath = '/Users/test/new-workspace'
        ;(
            AppModule.SelectWorkspaceDirectory as jest.MockedFunction<
                typeof AppModule.SelectWorkspaceDirectory
            >
        ).mockResolvedValue(newPath)

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Change' }))

        await waitFor(() => {
            expect(screen.getByDisplayValue(newPath)).toBeInTheDocument()
        })
    })

    it('does not change workspace when SelectWorkspaceDirectory returns empty', async () => {
        const user = userEvent.setup()
        ;(
            AppModule.SelectWorkspaceDirectory as jest.MockedFunction<
                typeof AppModule.SelectWorkspaceDirectory
            >
        ).mockResolvedValue('')

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Change' }))

        // Should keep original value
        expect(screen.getByDisplayValue('/Users/test/workspace')).toBeInTheDocument()
    })

    it('handles SelectWorkspaceDirectory error gracefully', async () => {
        const user = userEvent.setup()
        const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {})
        ;(
            AppModule.SelectWorkspaceDirectory as jest.MockedFunction<
                typeof AppModule.SelectWorkspaceDirectory
            >
        ).mockRejectedValue(new Error('Directory selection failed'))

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Change' }))

        // Should keep original value and not crash
        expect(screen.getByDisplayValue('/Users/test/workspace')).toBeInTheDocument()

        consoleSpy.mockRestore()
    })

    it('validates workspace before saving', async () => {
        const user = userEvent.setup()
        ;(
            AppModule.ValidateWorkspacePath as jest.MockedFunction<
                typeof AppModule.ValidateWorkspacePath
            >
        ).mockRejectedValue(new Error('Invalid path'))

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Save Settings' }))

        await waitFor(() => {
            expect(
                screen.getByText(
                    'Invalid workspace directory. Please ensure the path exists and is writable.'
                )
            ).toBeInTheDocument()
        })
    })

    it('validates that save button is disabled for empty workspace path', () => {
        render(<SettingsView config={{ WorkspaceRoot: '' }} />)

        // Save button should be disabled when path is empty
        const saveButton = screen.getByRole('button', { name: 'Save Settings' })
        expect(saveButton).toBeDisabled()
    })

    it('saves settings successfully', async () => {
        const user = userEvent.setup()
        ;(
            AppModule.ValidateWorkspacePath as jest.MockedFunction<
                typeof AppModule.ValidateWorkspacePath
            >
        ).mockResolvedValue(undefined)
        ;(
            AppModule.SaveConfig as jest.MockedFunction<typeof AppModule.SaveConfig>
        ).mockResolvedValue(undefined)

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Save Settings' }))

        await waitFor(() => {
            expect(AppModule.SaveConfig).toHaveBeenCalledWith({
                WorkspaceRoot: '/Users/test/workspace',
            })
        })
    })

    it('handles save configuration error', async () => {
        const user = userEvent.setup()
        const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {})
        ;(
            AppModule.ValidateWorkspacePath as jest.MockedFunction<
                typeof AppModule.ValidateWorkspacePath
            >
        ).mockResolvedValue(undefined)
        ;(
            AppModule.SaveConfig as jest.MockedFunction<typeof AppModule.SaveConfig>
        ).mockRejectedValue(new Error('Save failed'))

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Save Settings' }))

        await waitFor(() => {
            expect(screen.getByText('Failed to save configuration')).toBeInTheDocument()
        })

        consoleSpy.mockRestore()
    })

    it('disables save button when path is empty', () => {
        render(<SettingsView config={{ WorkspaceRoot: '' }} />)

        expect(screen.getByRole('button', { name: 'Save Settings' })).toBeDisabled()
    })

    it('disables save button when validating', async () => {
        const user = userEvent.setup()
        // Mock ValidateWorkspacePath to never resolve (simulating ongoing validation)
        ;(
            AppModule.ValidateWorkspacePath as jest.MockedFunction<
                typeof AppModule.ValidateWorkspacePath
            >
        ).mockImplementation(() => new Promise(() => {}))

        render(<SettingsView config={mockConfig} />)

        // Start validation
        await user.click(screen.getByRole('button', { name: 'Save Settings' }))

        // Button should be disabled during validation
        expect(screen.getByRole('button', { name: 'Save Settings' })).toBeDisabled()
    })

    it('shows validation loading state', async () => {
        const user = userEvent.setup()
        // Mock ValidateWorkspacePath to never resolve (simulating ongoing validation)
        ;(
            AppModule.ValidateWorkspacePath as jest.MockedFunction<
                typeof AppModule.ValidateWorkspacePath
            >
        ).mockImplementation(() => new Promise(() => {}))

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Save Settings' }))

        // Should show validation state
        expect(screen.getByRole('button', { name: 'Save Settings' })).toBeDisabled()
    })

    it('clears validation error when workspace path changes', async () => {
        const user = userEvent.setup()
        const newPath = '/Users/test/new-workspace'

        // First, trigger a validation error
        ;(
            AppModule.ValidateWorkspacePath as jest.MockedFunction<
                typeof AppModule.ValidateWorkspacePath
            >
        ).mockRejectedValue(new Error('Invalid path'))

        render(<SettingsView config={mockConfig} />)

        await user.click(screen.getByRole('button', { name: 'Save Settings' }))

        await waitFor(() => {
            expect(
                screen.getByText(
                    'Invalid workspace directory. Please ensure the path exists and is writable.'
                )
            ).toBeInTheDocument()
        })

        // Now change the workspace path
        ;(
            AppModule.SelectWorkspaceDirectory as jest.MockedFunction<
                typeof AppModule.SelectWorkspaceDirectory
            >
        ).mockResolvedValue(newPath)

        await user.click(screen.getByRole('button', { name: 'Change' }))

        await waitFor(() => {
            expect(
                screen.queryByText(
                    'Invalid workspace directory. Please ensure the path exists and is writable.'
                )
            ).not.toBeInTheDocument()
        })
    })
})
