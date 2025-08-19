import { test, expect } from '@playwright/test'

test.describe('Home Page', () => {
  test('should load the home page', async ({ page }) => {
    await page.goto('/')
    
    // Check if the page loads
    await expect(page).toHaveTitle(/Bruno/)
    
    // Check for main content
    await expect(page.locator('main')).toBeVisible()
  })

  test('should display projects section', async ({ page }) => {
    await page.goto('/')
    
    // Wait for projects to load
    await page.waitForSelector('[data-testid="projects-section"]', { timeout: 10000 })
    
    // Check if projects are displayed
    const projects = page.locator('[data-testid="project-card"]')
    await expect(projects.first()).toBeVisible()
  })

  test('should display about section', async ({ page }) => {
    await page.goto('/')
    
    // Scroll to about section
    await page.locator('[data-testid="about-section"]').scrollIntoViewIfNeeded()
    
    // Check if about content is displayed
    await expect(page.locator('[data-testid="about-description"]')).toBeVisible()
  })

  test('should display contact information', async ({ page }) => {
    await page.goto('/')
    
    // Scroll to contact section
    await page.locator('[data-testid="contact-section"]').scrollIntoViewIfNeeded()
    
    // Check if contact info is displayed
    await expect(page.locator('[data-testid="contact-email"]')).toBeVisible()
  })

  test('should track project views when clicking on projects', async ({ page }) => {
    await page.goto('/')
    
    // Wait for projects to load
    await page.waitForSelector('[data-testid="project-card"]', { timeout: 10000 })
    
    // Click on first project
    await page.locator('[data-testid="project-card"]').first().click()
    
    // Check if tracking request was made
    await page.waitForRequest(request => 
      request.url().includes('/api/v1/analytics/track') && 
      request.method() === 'POST'
    )
  })

  test('should be responsive on mobile', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 })
    await page.goto('/')
    
    // Check if mobile navigation works
    await expect(page.locator('[data-testid="mobile-menu"]')).toBeVisible()
  })

  test('should handle API errors gracefully', async ({ page }) => {
    // Mock API failure
    await page.route('**/api/v1/projects', route => 
      route.fulfill({ status: 500, body: 'Internal Server Error' })
    )
    
    await page.goto('/')
    
    // Check if error state is handled
    await expect(page.locator('[data-testid="error-message"]')).toBeVisible()
  })

  test('should have proper accessibility', async ({ page }) => {
    await page.goto('/')
    
    // Check for proper heading structure
    const headings = page.locator('h1, h2, h3, h4, h5, h6')
    await expect(headings.first()).toBeVisible()
    
    // Check for alt text on images
    const images = page.locator('img')
    for (let i = 0; i < await images.count(); i++) {
      const alt = await images.nth(i).getAttribute('alt')
      expect(alt).toBeTruthy()
    }
  })
})
