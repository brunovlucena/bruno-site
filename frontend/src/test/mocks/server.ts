import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

// For testing, we'll use a mock URL that doesn't require real connection
const API_BASE_URL = 'http://test-api.local'

export const handlers = [
  // Health check
  http.get(`${API_BASE_URL}/health`, () => {
    return HttpResponse.json({
      status: 'healthy',
      timestamp: Date.now()
    })
  }),

  // Projects endpoint
  http.get(`${API_BASE_URL}/api/projects`, () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Test Project',
        description: 'A test project description',
        short_description: 'Test project',
        type: 'web',
        icon: 'test-icon',
        github_url: 'https://github.com/test',
        live_url: 'https://test.com',
        technologies: ['React', 'TypeScript'],
        active: true
      }
    ])
  }),

  // About endpoint
  http.get(`${API_BASE_URL}/api/about`, () => {
    return HttpResponse.json({
      description: 'Test about description',
      highlights: [
        {
          icon: 'test-icon',
          text: 'Test highlight'
        }
      ]
    })
  }),

  // Contact endpoint
  http.get(`${API_BASE_URL}/api/contact`, () => {
    return HttpResponse.json({
      email: 'test@example.com',
      location: 'Test Location',
      linkedin: 'https://linkedin.com/test',
      github: 'https://github.com/test',
      availability: 'Available'
    })
  }),

  // Skills endpoint
  http.get(`${API_BASE_URL}/api/skills`, () => {
    return HttpResponse.json([
      {
        id: 1,
        name: 'React',
        category: 'Frontend',
        proficiency: 90,
        icon: 'react-icon',
        order: 1,
        active: true
      }
    ])
  }),

  // Experience endpoint
  http.get(`${API_BASE_URL}/api/experiences`, () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Software Engineer',
        company: 'Test Company',
        description: 'Test experience description',
        start_date: '2020-01',
        end_date: '2023-01',
        current: false,
        order: 1,
        technologies: ['React', 'TypeScript'],
        active: true
      }
    ])
  }),

  // Analytics tracking endpoint
  http.post(`${API_BASE_URL}/api/analytics/track`, () => {
    return HttpResponse.json({
      status: 'tracked',
      project_id: 1
    })
  }),

  // Metrics endpoint
  http.get(`${API_BASE_URL}/metrics`, () => {
    return new HttpResponse('# HELP test_metric Test metric\n# TYPE test_metric counter\ntest_metric 1\n', {
      headers: {
        'Content-Type': 'text/plain'
      }
    })
  })
]

export const server = setupServer(...handlers)
