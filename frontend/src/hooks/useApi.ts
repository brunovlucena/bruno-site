import { useQuery, useMutation, useQueryClient } from 'react-query'
import { apiClient, Project, Skill, Experience, AboutData, ContactData, Content } from '../services/api'

// =============================================================================
// 🎯 PROJECTS HOOKS
// =============================================================================

export const useProjects = () => {
  return useQuery('projects', apiClient.getProjects, {
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 10 * 60 * 1000, // 10 minutes
  })
}

export const useProject = (id: number) => {
  return useQuery(['project', id], () => apiClient.getProject(id), {
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}

export const useCreateProject = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.createProject, {
    onSuccess: () => {
      queryClient.invalidateQueries('projects')
    },
  })
}

export const useUpdateProject = () => {
  const queryClient = useQueryClient()
  
  return useMutation(
    ({ id, project }: { id: number; project: Partial<Project> }) =>
      apiClient.updateProject(id, project),
    {
      onSuccess: (_, { id }) => {
        queryClient.invalidateQueries('projects')
        queryClient.invalidateQueries(['project', id])
      },
    }
  )
}

export const useDeleteProject = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.deleteProject, {
    onSuccess: () => {
      queryClient.invalidateQueries('projects')
    },
  })
}

// =============================================================================
// 🛠️ SKILLS HOOKS
// =============================================================================

export const useSkills = () => {
  return useQuery('skills', apiClient.getSkills, {
    staleTime: 10 * 60 * 1000, // 10 minutes
    cacheTime: 20 * 60 * 1000, // 20 minutes
  })
}

export const useSkill = (id: number) => {
  return useQuery(['skill', id], () => apiClient.getSkill(id), {
    enabled: !!id,
    staleTime: 10 * 60 * 1000,
  })
}

export const useCreateSkill = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.createSkill, {
    onSuccess: () => {
      queryClient.invalidateQueries('skills')
    },
  })
}

export const useUpdateSkill = () => {
  const queryClient = useQueryClient()
  
  return useMutation(
    ({ id, skill }: { id: number; skill: Partial<Skill> }) =>
      apiClient.updateSkill(id, skill),
    {
      onSuccess: (_, { id }) => {
        queryClient.invalidateQueries('skills')
        queryClient.invalidateQueries(['skill', id])
      },
    }
  )
}

export const useDeleteSkill = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.deleteSkill, {
    onSuccess: () => {
      queryClient.invalidateQueries('skills')
    },
  })
}

// =============================================================================
// 💼 EXPERIENCES HOOKS
// =============================================================================

export const useExperiences = () => {
  return useQuery('experiences', apiClient.getExperiences, {
    staleTime: 10 * 60 * 1000,
    cacheTime: 20 * 60 * 1000,
  })
}

export const useExperience = (id: number) => {
  return useQuery(['experience', id], () => apiClient.getExperience(id), {
    enabled: !!id,
    staleTime: 10 * 60 * 1000,
  })
}

export const useCreateExperience = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.createExperience, {
    onSuccess: () => {
      queryClient.invalidateQueries('experiences')
    },
  })
}

export const useUpdateExperience = () => {
  const queryClient = useQueryClient()
  
  return useMutation(
    ({ id, experience }: { id: number; experience: Partial<Experience> }) =>
      apiClient.updateExperience(id, experience),
    {
      onSuccess: (_, { id }) => {
        queryClient.invalidateQueries('experiences')
        queryClient.invalidateQueries(['experience', id])
      },
    }
  )
}

export const useDeleteExperience = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.deleteExperience, {
    onSuccess: () => {
      queryClient.invalidateQueries('experiences')
    },
  })
}

// =============================================================================
// 📄 CONTENT HOOKS
// =============================================================================

export const useContent = () => {
  return useQuery('content', apiClient.getContent, {
    staleTime: 5 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
  })
}

export const useContentByType = (type: string) => {
  return useQuery(['content', type], () => apiClient.getContentByType(type), {
    enabled: !!type,
    staleTime: 5 * 60 * 1000,
  })
}

export const useCreateContent = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.createContent, {
    onSuccess: () => {
      queryClient.invalidateQueries('content')
    },
  })
}

export const useUpdateContent = () => {
  const queryClient = useQueryClient()
  
  return useMutation(
    ({ id, content }: { id: number; content: Partial<Content> }) =>
      apiClient.updateContent(id, content),
    {
      onSuccess: (_, { id }) => {
        queryClient.invalidateQueries('content')
        queryClient.invalidateQueries(['content', id])
      },
    }
  )
}

export const useDeleteContent = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.deleteContent, {
    onSuccess: () => {
      queryClient.invalidateQueries('content')
    },
  })
}

// =============================================================================
// 👤 ABOUT HOOKS
// =============================================================================

export const useAbout = () => {
  return useQuery('about', apiClient.getAbout, {
    staleTime: 10 * 60 * 1000,
    cacheTime: 20 * 60 * 1000,
  })
}

export const useUpdateAbout = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.updateAbout, {
    onSuccess: () => {
      queryClient.invalidateQueries('about')
    },
  })
}

// =============================================================================
// 📞 CONTACT HOOKS
// =============================================================================

export const useContact = () => {
  return useQuery('contact', apiClient.getContact, {
    staleTime: 10 * 60 * 1000,
    cacheTime: 20 * 60 * 1000,
  })
}

export const useUpdateContact = () => {
  const queryClient = useQueryClient()
  
  return useMutation(apiClient.updateContact, {
    onSuccess: () => {
      queryClient.invalidateQueries('contact')
    },
  })
}

// =============================================================================
// 🏥 HEALTH HOOKS
// =============================================================================

export const useHealthCheck = () => {
  return useQuery('health', apiClient.healthCheck, {
    refetchInterval: 30 * 1000, // Check every 30 seconds
    refetchIntervalInBackground: true,
  })
}

