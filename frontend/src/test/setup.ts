import '@testing-library/jest-dom'
import { vi } from 'vitest'

// Disable MSW for now to avoid localhost issues
// import { server } from './mocks/server'
// beforeAll(() => server.listen({ onUnhandledRequest: 'bypass' }))
// afterEach(() => server.resetHandlers())
// afterAll(() => server.close())

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  unobserve() {}
  root: Element | null = null
  rootMargin: string = ''
  thresholds: ReadonlyArray<number> = []
  takeRecords() { return [] }
}

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  unobserve() {}
}

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})
