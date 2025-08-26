import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import App from './App'

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />)
    // Basic test to ensure the app renders without errors
    expect(document.body).toBeInTheDocument()
  })

  it('renders the main app container', () => {
    render(<App />)
    // Check if the main app container exists
    const appElement = document.querySelector('#root')
    expect(appElement).toBeInTheDocument()
  })
})
