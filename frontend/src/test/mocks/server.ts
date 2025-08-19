import { setupServer } from 'msw/node'
import { rest } from 'msw'

const API_BASE_URL = 'http://localhost:8080'

export const handlers = [
  // Health check
  rest.get(`${API_BASE_URL}/health`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        status: 'healthy',
        timestamp: Date.now()
      })
    )
  }),

  // Projects endpoint
  rest.get(`${API_BASE_URL}/api/v1/projects`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json([
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
    )
  }),

  // About endpoint
  rest.get(`${API_BASE_URL}/api/v1/about`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        description: 'Test about description',
        highlights: [
          {
            icon: 'test-icon',
            text: 'Test highlight'
          }
        ]
      })
    )
  }),

  // Contact endpoint
  rest.get(`${API_BASE_URL}/api/v1/contact`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        email: 'test@example.com',
        location: 'Test Location',
        linkedin: 'https://linkedin.com/test',
        github: 'https://github.com/test',
        availability: 'Available'
      })
    )
  }),

  // Skills endpoint
  rest.get(`${API_BASE_URL}/api/v1/skills`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json([
        {
          id: 1,
          name: 'React',
          category: 'Frontend',
          proficiency: 90,
          icon: 'react-icon',
          order: 1
        }
      ])
    )
  }),

  // Experience endpoint
  rest.get(`${API_BASE_URL}/api/v1/experience`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json([
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
    )
  }),

  // Analytics tracking endpoint
  rest.post(`${API_BASE_URL}/api/v1/analytics/track`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        status: 'tracked',
        project_id: 1
      })
    )
  }),

  // Metrics endpoint
  rest.get(`${API_BASE_URL}/metrics`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.set('Content-Type', 'text/plain'),
      ctx.body('# HELP test_metric Test metric\n# TYPE test_metric counter\ntest_metric 1\n')
    )
  })
]

export const server = setupServer(...handlers)
