import '@testing-library/jest-dom'
import { jest } from '@jest/globals'

// Mock Wails runtime
global.runtime = {
  EventsOn: jest.fn(),
  EventsOff: jest.fn(),
  EventsEmit: jest.fn(),
  Environment: jest.fn(),
  WindowShow: jest.fn(),
  WindowHide: jest.fn(),
}

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // deprecated
    removeListener: jest.fn(), // deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
})

// Clean up after each test
afterEach(() => {
  jest.clearAllMocks()
})