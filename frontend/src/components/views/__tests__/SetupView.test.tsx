import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { SetupView } from '../SetupView'
import { jest } from '@jest/globals'
import * as AppModule from 'wailsjs/go/main/App'

// Mock the Wails modules
jest.mock('wailsjs/go/main/App')
jest.mock('wailsjs/go/models')

describe('SetupView', () => {
  const mockOnConfigurationComplete = jest.fn()
  
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders setup form', () => {
    render(
      <SetupView 
        config={null} 
        onConfigurationComplete={mockOnConfigurationComplete} 
      />
    )
    
    expect(screen.getByText('Welcome to Notedown Planner')).toBeInTheDocument()
    expect(screen.getByText('Please configure your workspace to get started.')).toBeInTheDocument()
    expect(screen.getByText('Workspace Directory *')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Browse' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Complete Setup' })).toBeInTheDocument()
  })

  it('handles workspace directory selection', async () => {
    const user = userEvent.setup()
    const mockPath = '/Users/test/workspace'
    ;(AppModule.SelectWorkspaceDirectory as jest.MockedFunction<typeof AppModule.SelectWorkspaceDirectory>).mockResolvedValue(mockPath)
    
    render(
      <SetupView 
        config={null} 
        onConfigurationComplete={mockOnConfigurationComplete} 
      />
    )
    
    await user.click(screen.getByRole('button', { name: 'Browse' }))
    
    await waitFor(() => {
      expect(screen.getByDisplayValue(mockPath)).toBeInTheDocument()
    })
  })

  it('validates workspace path before saving', async () => {
    const user = userEvent.setup()
    const invalidPath = '/invalid/path'
    ;(AppModule.SelectWorkspaceDirectory as jest.MockedFunction<typeof AppModule.SelectWorkspaceDirectory>).mockResolvedValue(invalidPath)
    ;(AppModule.ValidateWorkspacePath as jest.MockedFunction<typeof AppModule.ValidateWorkspacePath>).mockRejectedValue(new Error('Invalid path'))
    
    render(
      <SetupView 
        config={null} 
        onConfigurationComplete={mockOnConfigurationComplete} 
      />
    )
    
    // Select an invalid path
    await user.click(screen.getByRole('button', { name: 'Browse' }))
    await waitFor(() => {
      expect(screen.getByDisplayValue(invalidPath)).toBeInTheDocument()
    })
    
    // Try to save with invalid path
    await user.click(screen.getByRole('button', { name: 'Complete Setup' }))
    
    await waitFor(() => {
      expect(screen.getByText(/Invalid workspace directory/)).toBeInTheDocument()
    })
  })

  it('saves configuration successfully', async () => {
    const user = userEvent.setup()
    const mockPath = '/Users/test/workspace'
    ;(AppModule.SelectWorkspaceDirectory as jest.MockedFunction<typeof AppModule.SelectWorkspaceDirectory>).mockResolvedValue(mockPath)
    ;(AppModule.ValidateWorkspacePath as jest.MockedFunction<typeof AppModule.ValidateWorkspacePath>).mockResolvedValue(undefined)
    ;(AppModule.SaveConfig as jest.MockedFunction<typeof AppModule.SaveConfig>).mockResolvedValue(undefined)
    
    render(
      <SetupView 
        config={null} 
        onConfigurationComplete={mockOnConfigurationComplete} 
      />
    )
    
    await user.click(screen.getByRole('button', { name: 'Browse' }))
    await waitFor(() => {
      expect(screen.getByDisplayValue(mockPath)).toBeInTheDocument()
    })
    
    await user.click(screen.getByRole('button', { name: 'Complete Setup' }))
    
    await waitFor(() => {
      expect(AppModule.SaveConfig).toHaveBeenCalledWith({
        WorkspaceRoot: mockPath,
      })
      expect(mockOnConfigurationComplete).toHaveBeenCalled()
    })
  })
})