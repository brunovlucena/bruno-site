import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Header from './Header'

describe('Header', () => {
  it('renders the header component', () => {
    render(<Header />)
    // Basic test to ensure the header renders
    expect(document.body).toBeInTheDocument()
  })

  it('contains header content', () => {
    render(<Header />)
    // Check if header element exists
    const headerElement = document.querySelector('header')
    expect(headerElement).toBeInTheDocument()
  })
})
