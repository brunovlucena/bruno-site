import React, { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { projectsApi } from '../services/api'
import { Project } from '../types'
import { 
  SiReact, 
  SiTypescript, 
  SiVite, 
  SiTailwindcss, 
  SiReactrouter,
  SiGo,
  SiPostgresql,
  SiRedis,
  SiKubernetes,
  SiFlux,
  SiHelm,
  SiNginx,
  SiDocker,
  SiGithub,
  SiGithubactions,
  SiPulumi,
  SiPrometheus,
  SiGrafana,
  SiAmazon,
  SiGooglecloud
} from 'react-icons/si'
import { FaShieldAlt, FaLock, FaChartBar, FaCloud, FaRobot, FaRocket } from 'react-icons/fa'
import { BiCertification } from 'react-icons/bi'

const Home: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchParams] = useSearchParams()

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        setLoading(true)
        const fetchedProjects = await projectsApi.getAll()
        setProjects(fetchedProjects)
        setError(null)
      } catch (err) {
        console.error('Failed to fetch projects:', err)
        setError('Failed to load projects')
      } finally {
        setLoading(false)
      }
    }

    fetchProjects()
  }, [])

  // Handle scrolling to projects section when section parameter is present
  useEffect(() => {
    const section = searchParams.get('section')
    if (section === 'projects') {
      const projectsElement = document.getElementById('projects')
      if (projectsElement) {
        projectsElement.scrollIntoView({ behavior: 'smooth' })
      }
    }
  }, [searchParams])

  const getIconComponent = (iconName: string) => {
    const iconMap: { [key: string]: React.ComponentType<any> } = {
      'react': SiReact,
      'typescript': SiTypescript,
      'vite': SiVite,
      'tailwind': SiTailwindcss,
      'router': SiReactrouter,
      'go': SiGo,
      'postgresql': SiPostgresql,
      'redis': SiRedis,
      'kubernetes': SiKubernetes,
      'flux': SiFlux,
      'helm': SiHelm,
      'nginx': SiNginx,
      'docker': SiDocker,
      'github': SiGithub,
      'githubactions': SiGithubactions,
      'pulumi': SiPulumi,
      'prometheus': SiPrometheus,
      'grafana': SiGrafana,
      'shield': FaShieldAlt,
      'certification': BiCertification,
    }
    return iconMap[iconName.toLowerCase()] || SiGithub
  }

  return (
    <div className="home">
      <section className="hero">
        <div className="container">
          <h1>Senior Cloud Native Infrastructure Engineer</h1>
          <p>SRE • DevSecOps • AI/ML Infrastructure • Platform Engineering</p>
        </div>
      </section>
      
      <section id="about" className="section">
        <div className="container">
          <div className="about-header">
            <div className="about-intro">
              <div className="about-image-container">
                <img 
                  src="/assets/eu.png" 
                  alt="Bruno Lucena" 
                  className="about-image"
                />
              </div>
              <h2>About Me</h2>
              <p>Senior Cloud Native Infrastructure Engineer with 12+ years of experience designing, building, and scaling mission-critical cloud-native platforms. Expert in Kubernetes ecosystem, multi-cloud architectures (AWS/GCP), and modern observability stacks. Proven track record in Site Reliability Engineering (SRE), DevSecOps practices, and leading high-performing infrastructure teams. Specialized in AI/ML infrastructure, LLMOps, and building resilient systems that handle millions of requests. Passionate about automation, security-first approaches, and driving innovation in cloud-native technologies. Delivers enterprise-grade solutions with 99.9%+ uptime, optimized performance, and comprehensive security postures.</p>
            </div>
          </div>
          
          <div className="skills-grid">
            <div className="skill-tag">
              <SiGo className="skill-icon" style={{ color: '#00ADD8' }} />
              <span>Go</span>
            </div>
            <div className="skill-tag">
              <SiKubernetes className="skill-icon" style={{ color: '#326CE5' }} />
              <span>Kubernetes</span>
            </div>
            <div className="skill-tag">
              <SiDocker className="skill-icon" style={{ color: '#2496ED' }} />
              <span>Docker</span>
            </div>
            <div className="skill-tag">
              <SiPostgresql className="skill-icon" style={{ color: '#336791' }} />
              <span>PostgreSQL</span>
            </div>
            <div className="skill-tag">
              <SiRedis className="skill-icon" style={{ color: '#DC382D' }} />
              <span>Redis</span>
            </div>
            <div className="skill-tag">
              <SiPrometheus className="skill-icon" style={{ color: '#E6522C' }} />
              <span>Prometheus</span>
            </div>
            <div className="skill-tag">
              <SiGrafana className="skill-icon" style={{ color: '#F46800' }} />
              <span>Grafana</span>
            </div>
            <div className="skill-tag">
              <SiPulumi className="skill-icon" style={{ color: '#00B4D8' }} />
              <span>Pulumi</span>
            </div>
            <div className="skill-tag">
              <SiAmazon className="skill-icon" style={{ color: '#FF9900' }} />
              <span>AWS</span>
            </div>
            <div className="skill-tag">
              <SiGooglecloud className="skill-icon" style={{ color: '#4285F4' }} />
              <span>GCP</span>
            </div>
            <div className="skill-tag">
              <SiGithub className="skill-icon" style={{ color: '#181717' }} />
              <span>GitHub</span>
            </div>
            <div className="skill-tag">
              <SiGithubactions className="skill-icon" style={{ color: '#2088FF' }} />
              <span>GitHub Actions</span>
            </div>
          </div>
        </div>
      </section>

      <section id="projects" className="section">
        <div className="container">
          <h2>Homelab</h2>
          {loading && (
            <div className="loading">
              <p>Loading homelab projects...</p>
            </div>
          )}
          
          {error && (
            <div className="error">
              <p>Error: {error}</p>
            </div>
          )}
          
          {!loading && !error && projects.length === 0 && (
            <div className="no-projects">
              <p>No homelab projects available at the moment.</p>
            </div>
          )}
          
          {!loading && !error && projects.length > 0 && (
            <div className="projects-grid">
              {projects.map((project) => {
                const IconComponent = getIconComponent(project.technologies[0] || 'github')
                return (
                  <div key={project.id} className="project-card">
                    <div className="project-header">
                      <IconComponent className="project-icon" />
                      <h3>{project.title}</h3>
                    </div>
                    <p className="project-description">{project.description}</p>
                    
                    {/* YouTube Video Embed */}
                    {project.video_url && (
                      <div className="project-video">
                        <iframe
                          src={project.video_url}
                          title={`${project.title} Video`}
                          frameBorder="0"
                          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                          allowFullScreen
                        ></iframe>
                      </div>
                    )}
                    
                    <div className="project-meta">
                      <span className="project-type">{project.type}</span>
                      <span className="project-modules">{project.modules} modules</span>
                    </div>
                    {project.github_url && (
                      <a 
                        href={project.github_url} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        className="project-link"
                      >
                        View Homelab Project
                      </a>
                    )}
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </section>

      <footer className="footer">
        <div className="container">
          <div className="footer-content">
            <div className="footer-tech">
              <p>This site was built with:</p>
              <div className="footer-icons">
                <SiReact className="footer-icon" style={{ color: '#61DAFB' }} />
                <SiTypescript className="footer-icon" style={{ color: '#3178C6' }} />
                <SiVite className="footer-icon" style={{ color: '#646CFF' }} />
                <SiTailwindcss className="footer-icon" style={{ color: '#06B6D4' }} />
              </div>
            </div>
            <div className="footer-links">
              <a href="https://github.com/brunovlucena" target="_blank" rel="noopener noreferrer">
                <SiGithub className="footer-link-icon" />
                GitHub
              </a>
              <a href="https://www.linkedin.com/in/bvlucena" target="_blank" rel="noopener noreferrer">
                <SiGithub className="footer-link-icon" />
                LinkedIn
              </a>
            </div>
          </div>
        </div>
      </footer>

    </div>
  )
}

export default Home
