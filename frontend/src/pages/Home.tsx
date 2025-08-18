import React, { useState, useEffect } from 'react'
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
  SiGrafana
} from 'react-icons/si'
import { FaShieldAlt, FaLock, FaChartBar, FaCloud, FaRobot, FaRocket } from 'react-icons/fa'
import { BiCertification } from 'react-icons/bi'

const Home: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

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
          <h1>SRE/DevSecOps/AI Engineer</h1>
          <p>Senior Cloud Native Infrastructure Engineer with 12+ years of experience</p>
        </div>
      </section>
      
      <section id="about" className="section">
        <div className="container">
          <div className="about-header">
            <div className="about-intro">
              <h2>About Me</h2>
              <p>Senior Cloud Native Infrastructure Engineer with 12+ years of experience architecting and maintaining robust cloud-native solutions on AWS, GCP, and bare-metal environments. Expert in Kubernetes orchestration, observability systems, and automation using Terraform, Pulumi, and CI/CD pipelines. Proven track record in SRE practices, security analysis, and leading infrastructure teams. Passionate about solving complex infrastructure challenges and driving innovation in Agentic DevOps and LLMOps. Delivers high-availability, scalable solutions that minimize downtime and optimize system performance.</p>
            </div>
          </div>
          
          <div className="skills-grid">
            <div className="skill-tag">
              <FaLock className="skill-icon" />
              <span>IT Security</span>
            </div>
            <div className="skill-tag">
              <FaChartBar className="skill-icon" />
              <span>Project Management</span>
            </div>
            <div className="skill-tag">
              <SiKubernetes className="skill-icon" />
              <span>Kubernetes</span>
            </div>
            <div className="skill-tag">
              <FaCloud className="skill-icon" />
              <span>AWS & GCP</span>
            </div>
            <div className="skill-tag">
              <SiPrometheus className="skill-icon" />
              <span>Observability</span>
            </div>
            <div className="skill-tag">
              <FaRobot className="skill-icon" />
              <span>AI/LLMOps</span>
            </div>
            <div className="skill-tag">
              <FaLock className="skill-icon" />
              <span>Security</span>
            </div>
            <div className="skill-tag">
              <FaRocket className="skill-icon" />
              <span>Automation</span>
            </div>
          </div>
        </div>
      </section>

      <section id="projects" className="section">
        <div className="container">
          <h2>Projects</h2>
          {loading && (
            <div className="loading">
              <p>Loading projects...</p>
            </div>
          )}
          
          {error && (
            <div className="error">
              <p>Error: {error}</p>
            </div>
          )}
          
          {!loading && !error && projects.length === 0 && (
            <div className="no-projects">
              <p>No projects available at the moment.</p>
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
                        View Project
                      </a>
                    )}
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </section>

      <section id="technologies" className="section">
        <div className="container">
          <h2>Technologies Used</h2>
          <p>This portfolio site was built using modern cloud-native technologies:</p>
          
          <div className="tech-grid">
            <div className="tech-category">
              <h3>Frontend</h3>
              <ul>
                <li><SiReact className="tech-icon" /> React 19 with TypeScript 5.9</li>
                <li><SiVite className="tech-icon" /> Vite for fast builds</li>
                <li><SiTailwindcss className="tech-icon" /> Tailwind CSS for styling</li>
                <li><SiReactrouter className="tech-icon" /> React Router for navigation</li>
              </ul>
            </div>
            
            <div className="tech-category">
              <h3>Backend</h3>
              <ul>
                <li><SiGo className="tech-icon" /> Go (Golang) API</li>
                <li><SiPostgresql className="tech-icon" /> PostgreSQL database</li>
                <li><SiRedis className="tech-icon" /> Redis for caching</li>
                <li><FaShieldAlt className="tech-icon" /> JWT authentication</li>
              </ul>
            </div>
            
            <div className="tech-category">
              <h3>Infrastructure</h3>
              <ul>
                <li><SiKubernetes className="tech-icon" /> Kubernetes orchestration</li>
                <li><SiFlux className="tech-icon" /> Flux GitOps for deployment</li>
                <li><SiHelm className="tech-icon" /> Helm charts for packaging</li>
                <li><SiNginx className="tech-icon" /> Nginx ingress controller</li>
                <li><BiCertification className="tech-icon" /> Cert-manager for SSL</li>
              </ul>
            </div>
            
            <div className="tech-category">
              <h3>DevOps</h3>
              <ul>
                <li><SiDocker className="tech-icon" /> Docker containerization</li>
                <li><SiGithub className="tech-icon" /> GitHub Container Registry</li>
                <li><SiGithubactions className="tech-icon" /> GitHub Actions CI/CD</li>
                <li><SiPulumi className="tech-icon" /> Pulumi for IaC</li>
                <li><SiPrometheus className="tech-icon" /> Prometheus + Grafana monitoring</li>
              </ul>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default Home
