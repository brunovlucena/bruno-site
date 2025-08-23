import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'

// Simple test component for now
const TestComponent = () => {
  return (
    <div>
      <h1>Test Component</h1>
      <p>This is a test</p>
    </div>
  )
}

describe('Basic Component Test', () => {
  it('renders without crashing', () => {
    render(<TestComponent />)
    expect(screen.getByText('Test Component')).toBeInTheDocument()
  })

  it('displays test content', () => {
    render(<TestComponent />)
    expect(screen.getByText('This is a test')).toBeInTheDocument()
  })
})
