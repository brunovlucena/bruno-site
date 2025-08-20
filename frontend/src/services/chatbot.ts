import { apiClient } from './api';

export interface ChatbotResponse {
  text: string;
  suggestions?: string[];
  action?: 'navigate' | 'show_project' | 'show_contact';
  data?: any;
}

export class ChatbotService {
  private static instance: ChatbotService;
  private projects: any[] = [];
  private skills: any[] = [];
  private experience: any[] = [];
  private about: any = null;

  private constructor() {}

  static getInstance(): ChatbotService {
    if (!ChatbotService.instance) {
      ChatbotService.instance = new ChatbotService();
    }
    return ChatbotService.instance;
  }

  async initialize(): Promise<void> {
    try {
      console.log('🔄 ChatbotService: Starting initialization...');
      
      // Load data in parallel with better error handling
      console.log('📡 ChatbotService: Loading projects...');
      const [projectsData, skillsData, experienceData, aboutData] = await Promise.allSettled([
        apiClient.getProjects(),
        apiClient.getSkills(),
        apiClient.getAbout(),
        apiClient.getExperiences()
      ]);

      console.log('📊 ChatbotService: Processing results...');
      console.log('  - Projects status:', projectsData.status);
      console.log('  - Skills status:', skillsData.status);
      console.log('  - Experience status:', experienceData.status);
      console.log('  - About status:', aboutData.status);

      this.projects = projectsData.status === 'fulfilled' ? (projectsData.value as unknown as any[]) : [];
      this.skills = skillsData.status === 'fulfilled' ? (skillsData.value as unknown as any[]) : [];
      this.experience = experienceData.status === 'fulfilled' ? (experienceData.value as unknown as any[]) : [];
      this.about = aboutData.status === 'fulfilled' ? (aboutData.value as unknown as any) : { key: 'about', value: { description: '' } };

      console.log('✅ ChatbotService: Data loaded successfully:', {
        projects: this.projects.length,
        skills: this.skills.length,
        experience: this.experience.length,
        hasAbout: !!this.about?.value
      });

      if (this.projects.length > 0) {
        console.log('📋 ChatbotService: Sample project:', this.projects[0]);
      }
    } catch (error) {
      console.error('❌ ChatbotService: Failed to load data:', error);
      // Initialize with empty data to prevent errors
      this.projects = [];
      this.skills = [];
      this.experience = [];
      this.about = { value: '' };
    }
  }

  async processMessage(userInput: string): Promise<ChatbotResponse> {
    const input = userInput.toLowerCase().trim();
    
    // Experience and work history
    if (this.matchesKeywords(input, ['experience', 'work', 'job', 'career', 'background'])) {
      return this.handleExperienceQuery(input);
    }
    
    // Projects
    if (this.matchesKeywords(input, ['project', 'work', 'site', 'github'])) {
      return this.handleProjectsQuery(input);
    }
    
    // Skills and technologies
    if (this.matchesKeywords(input, ['skill', 'technology', 'tech', 'stack', 'tools', 'languages'])) {
      return this.handleSkillsQuery(input);
    }
    
    // Contact information
    if (this.matchesKeywords(input, ['contact', 'email', 'reach', 'get in touch', 'hire', 'available'])) {
      return this.handleContactQuery(input);
    }
    
    // Resume
    if (this.matchesKeywords(input, ['resume', 'cv', 'education', 'certification'])) {
      return this.handleResumeQuery(input);
    }
    
    // About
    if (this.matchesKeywords(input, ['about', 'who', 'introduce', 'tell me about'])) {
      return this.handleAboutQuery(input);
    }
    
    // Greetings
    if (this.matchesKeywords(input, ['hello', 'hi', 'hey', 'good morning', 'good afternoon'])) {
      return {
        text: "Hello! I'm Bruno's AI assistant. I can help you learn more about his experience, projects, skills, and how to get in touch. What would you like to know?",
        suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his skills?', 'How can I contact him?']
      };
    }
    
    // Help
    if (this.matchesKeywords(input, ['help', 'what can you do', 'commands', 'options'])) {
      return {
        text: "I can help you with information about Bruno's professional background. Here are some things you can ask me about:",
        suggestions: ['Experience & Work History', 'Projects & Site', 'Skills & Technologies', 'Contact Information', 'Resume & Education']
      };
    }
    
    // Default response
    return {
      text: "That's an interesting question! Bruno has a diverse background in cloud infrastructure, AI/ML, and DevOps. Could you be more specific about what you'd like to know? I can help with his experience, projects, skills, or contact information.",
      suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his skills?', 'How can I contact him?']
    };
  }

  private matchesKeywords(input: string, keywords: string[]): boolean {
    return keywords.some(keyword => input.includes(keyword));
  }

  private handleExperienceQuery(input: string): ChatbotResponse {
    if (this.experience.length === 0) {
      return {
        text: "Bruno has 12+ years of experience in SRE, DevSecOps, and AI Engineering. He's worked with major cloud providers (AWS, GCP, Azure), Kubernetes, and has extensive experience in infrastructure automation and AI/ML technologies.",
        suggestions: ['Tell me about specific roles', 'What companies has he worked for?', 'Show me his skills']
      };
    }

    const recentExperience = (this.experience as any[]).slice(0, 3);
    const experienceText = recentExperience.map(exp => 
      `${exp.title} at ${exp.company} (${exp.period})`
    ).join(', ');

    return {
      text: `Bruno has extensive experience including: ${experienceText}. He specializes in cloud-native infrastructure, AI/ML, and DevOps automation. Would you like to know about specific roles or technologies?`,
      suggestions: ['Tell me about specific roles', 'What are his key skills?', 'Show me his projects']
    };
  }

  private handleProjectsQuery(input: string): ChatbotResponse {
    if (this.projects.length === 0) {
      return {
        text: "Bruno has worked on several interesting projects including cloud-native infrastructure, AI/ML implementations, and DevOps automation. Some highlights include Kubernetes cluster management, CI/CD pipelines, and AI model deployment.",
        suggestions: ['Tell me about his experience', 'What are his skills?', 'How can I contact him?']
      };
    }

    const activeProjects = (this.projects as any[]).filter(p => p.active).slice(0, 3);
    const projectNames = activeProjects.map(p => p.title).join(', ');

    return {
      text: `Bruno has worked on various projects including: ${projectNames}. These cover areas like cloud infrastructure, AI/ML, and automation. Which area interests you most?`,
      suggestions: ['Tell me about specific projects', 'What technologies does he use?', 'Show me his experience']
    };
  }

  private handleSkillsQuery(input: string): ChatbotResponse {
    if (this.skills.length === 0) {
      return {
        text: "Bruno's key skills include Kubernetes, Docker, AWS/GCP/Azure, Terraform, Python, Go, AI/ML, CI/CD, monitoring, and security. He's also experienced with various AI frameworks and cloud-native technologies.",
        suggestions: ['Tell me about his experience', 'Show me his projects', 'How can I contact him?']
      };
    }

    const skillCategories = (this.skills as any[]).reduce((acc, skill) => {
      if (!acc[skill.category]) acc[skill.category] = [];
      acc[skill.category].push(skill.name);
      return acc;
    }, {} as Record<string, string[]>);

    const skillText = Object.entries(skillCategories)
      .map(([category, skills]) => `${category}: ${(skills as string[]).slice(0, 3).join(', ')}`)
      .join('; ');

    return {
      text: `Bruno's skills include: ${skillText}. He has expertise across cloud infrastructure, AI/ML, and DevOps. What specific technology would you like to know more about?`,
      suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his certifications?']
    };
  }

  private handleContactQuery(input: string): ChatbotResponse {
    return {
      text: "You can reach Bruno through LinkedIn, GitHub, or email. He's currently available for new opportunities and consulting work. Would you like me to provide specific contact information or discuss his availability?",
      action: 'show_contact',
      suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his skills?']
    };
  }

  private handleResumeQuery(input: string): ChatbotResponse {
    return {
      text: "Bruno's resume includes his extensive experience in cloud infrastructure, his work with major tech companies, and his expertise in AI/ML. You can view his detailed resume on the resume page, or I can highlight specific aspects of his background.",
      action: 'navigate',
      data: { path: '/resume' },
      suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his skills?']
    };
  }

  private handleAboutQuery(input: string): ChatbotResponse {
    if (this.about?.value?.description) {
      const aboutText = this.about.value.description.length > 200 
        ? this.about.value.description.substring(0, 200) + '...'
        : this.about.value.description;
      
      return {
        text: aboutText,
        suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his skills?']
      };
    }

    return {
      text: "Bruno is a Senior SRE/DevSecOps/AI Engineer with 12+ years of experience in cloud-native infrastructure, Kubernetes, and AI/ML technologies. He's passionate about building scalable, secure, and efficient systems.",
      suggestions: ['Tell me about his experience', 'Show me his projects', 'What are his skills?']
    };
  }

  getQuickSuggestions(): string[] {
    return [
      'Tell me about his experience',
      'Show me his projects', 
      'What are his skills?',
      'How can I contact him?',
      'Tell me about his background'
    ];
  }
}

export default ChatbotService.getInstance(); 