import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import App from '../App'

// Mock the API service
vi.mock('../services/api', () => ({
  getProjects: vi.fn(() => Promise.resolve([
    {
      id: 1,
      title: 'Test Project',
      description: 'A test project',
      short_description: 'Test',
      type: 'web',
      modules: 3,
      icon: 'test-icon',
      url: 'https://test.com',
      technologies: ['React', 'TypeScript'],
      active: true
    }
  ])),
  getAbout: vi.fn(() => Promise.resolve({
    description: 'Test about description',
    highlights: [
      { icon: 'test-icon', text: 'Test highlight' }
    ]
  })),
  getContact: vi.fn(() => Promise.resolve({
    email: 'test@example.com',
    location: 'Test Location',
    linkedin: 'https://linkedin.com/test',
    github: 'https://github.com/test',
    availability: 'Available'
  }))
}))

describe('App Component', () => {
  it('renders without crashing', () => {
    render(<App />)
    expect(screen.getByRole('main')).toBeInTheDocument()
  })

  it('displays loading state initially', () => {
    render(<App />)
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })

  it('loads and displays projects', async () => {
    render(<App />)
    
    await waitFor(() => {
      expect(screen.getByText('Test Project')).toBeInTheDocument()
    })
  })

  it('displays about information', async () => {
    render(<App />)
    
    await waitFor(() => {
      expect(screen.getByText('Test about description')).toBeInTheDocument()
    })
  })

  it('displays contact information', async () => {
    render(<App />)
    
    await waitFor(() => {
      expect(screen.getByText('test@example.com')).toBeInTheDocument()
    })
  })

  it('handles API errors gracefully', async () => {
    // Mock API error
    const { getProjects } = await import('../services/api')
    vi.mocked(getProjects).mockRejectedValueOnce(new Error('API Error'))
    
    render(<App />)
    
    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument()
    })
  })
})
