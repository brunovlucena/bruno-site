import axios from 'axios';
import { Project, Skill, Experience, Content, AnalyticsData } from '@/types';

// Use relative API URL - nginx will proxy /api/* requests to the backend
const API_BASE_URL = '/api';

const baseURL = API_BASE_URL;

const api = axios.create({
  baseURL: baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 5000, // 5 second timeout
});

// Request interceptor for tracking
api.interceptors.request.use((config) => {
  // Add user agent and other tracking info
  config.headers['User-Agent'] = navigator.userAgent;
  config.headers['Referer'] = window.location.href;
  return config;
});

// Response interceptor for better error handling
api.interceptors.response.use(
  (response) => {
    // Check if response is empty or malformed
    if (!response.data) {
      console.warn('Empty response from API:', response.config?.url);
      return Promise.reject({
        message: 'Empty response from server',
        status: response.status,
        url: response.config?.url
      });
    }
    return response;
  },
  (error) => {
    console.warn('API Error:', {
      url: error.config?.url,
      status: error.response?.status,
      message: error.message
    });
    
    // Return a rejected promise with a more descriptive error
    return Promise.reject({
      message: error.response?.data?.message || error.message || 'Network error',
      status: error.response?.status,
      url: error.config?.url
    });
  }
);

// Projects API
export const projectsApi = {
  getAll: async (): Promise<Project[]> => {
    try {
      console.log('📡 API: Calling /v1/projects...');
      const response = await api.get('/v1/projects');
      console.log('✅ API: Projects response received:', {
        status: response.status,
        dataType: typeof response.data,
        isArray: Array.isArray(response.data),
        dataLength: Array.isArray(response.data) ? response.data.length : 0,
        sample: Array.isArray(response.data) && response.data.length > 0 ? { id: response.data[0].id, title: response.data[0].title } : null
      });
      return response.data || [];
    } catch (error) {
      console.error('❌ API: Failed to load projects:', error);
      return [];
    }
  },

  getById: async (id: number): Promise<Project> => {
    const response = await api.get(`/v1/projects/${id}`);
    return response.data;
  },

  create: async (project: Omit<Project, 'id' | 'created_at' | 'updated_at'>): Promise<Project> => {
    const response = await api.post('/v1/projects', project);
    return response.data;
  },

  update: async (id: number, updates: Partial<Project>): Promise<void> => {
    await api.put(`/v1/projects/${id}`, updates);
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/v1/projects/${id}`);
  },
};

// Skills API
export const skillsApi = {
  getAll: async (): Promise<Skill[]> => {
    try {
      console.log('📡 API: Calling /v1/content/skills...');
      const response = await api.get('/v1/content/skills');
      console.log('✅ API: Skills response received:', {
        status: response.status,
        dataLength: response.data?.length || 0
      });
      return response.data || [];
    } catch (error) {
      console.error('❌ API: Failed to load skills:', error);
      return [];
    }
  },

  update: async (skills: Skill[]): Promise<void> => {
    await api.put('/v1/content/skills', { skills });
  },
};

// Experience API
export const experienceApi = {
  getAll: async (): Promise<Experience[]> => {
    try {
      console.log('📡 API: Calling /v1/content/experience...');
      const response = await api.get('/v1/content/experience');
      console.log('✅ API: Experience response received:', {
        status: response.status,
        dataLength: response.data?.length || 0
      });
      return response.data || [];
    } catch (error) {
      console.error('❌ API: Failed to load experience:', error);
      return [];
    }
  },

  update: async (experience: Experience[]): Promise<void> => {
    await api.put('/v1/content/experience', { experience });
  },
};

// Content API
export const contentApi = {
  getAbout: async (): Promise<Content> => {
    try {
      const response = await api.get('/v1/about');
      return response.data || { key: 'about', value: { description: '' } };
    } catch (error) {
      console.warn('Failed to load about content:', error);
      return { key: 'about', value: { description: '' } };
    }
  },

  updateAbout: async (content: Content['value']): Promise<void> => {
    await api.put('/v1/about', content);
  },

  getContact: async (): Promise<Content> => {
    const response = await api.get('/v1/contact');
    return response.data;
  },

  updateContact: async (content: Content['value']): Promise<void> => {
    await api.put('/v1/contact', content);
  },
};

// Analytics API
export const analyticsApi = {
  trackVisit: async (data: { project_id?: number; ip: string; user_agent: string; referrer?: string }): Promise<void> => {
    await api.post('/v1/analytics/visit', data);
  },
};

// Health check
export const healthApi = {
  check: async (): Promise<{ status: string; time: string; version: string }> => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api; 