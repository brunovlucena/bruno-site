import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

const API_BASE_URL = 'http://localhost:8080'

// For testing, we'll use a mock URL that doesn't require real connection
const TEST_API_BASE_URL = 'http://test-api.local'

export const handlers = [
  // Health check
  http.get(`${API_BASE_URL}/health`, () => {
    return HttpResponse.json({
      status: 'healthy',
      timestamp: Date.now()
    })
  }),

  // Projects endpoint
  http.get(`${API_BASE_URL}/api/v1/projects`, () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Test Project',
        description: 'A test project description',
        short_description: 'Test project',
        type: 'web',
        modules: 3,
        icon: 'test-icon',
        url: 'https://test.com',
        video_url: 'https://test.com/video',
        technologies: ['React', 'TypeScript'],
        active: true
      }
    ])
  }),

  // About endpoint
  http.get(`${API_BASE_URL}/api/v1/about`, () => {
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
  http.get(`${API_BASE_URL}/api/v1/contact`, () => {
    return HttpResponse.json({
      email: 'test@example.com',
      location: 'Test Location',
      linkedin: 'https://linkedin.com/test',
      github: 'https://github.com/test',
      availability: 'Available'
    })
  }),

  // Skills endpoint
  http.get(`${API_BASE_URL}/api/v1/skills`, () => {
    return HttpResponse.json([
      {
        id: 1,
        name: 'React',
        category: 'Frontend',
        proficiency: 90,
        icon: 'react-icon',
        order: 1
      }
    ])
  }),

  // Experience endpoint
  http.get(`${API_BASE_URL}/api/v1/experience`, () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'Software Engineer',
        company: 'Test Company',
        description: 'Test experience description',
        start_date: '2020-01',
        end_date: '2023-01',
        technologies: ['React', 'TypeScript'],
        active: true
      }
    ])
  }),

  // Analytics tracking endpoint
  http.post(`${API_BASE_URL}/api/v1/analytics/track`, () => {
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
