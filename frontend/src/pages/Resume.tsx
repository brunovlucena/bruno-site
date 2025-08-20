import React, { useState, useEffect } from 'react'
import { apiClient } from '../services/api'
import { Experience } from '../services/api'

// Technology icons mapping
const technologyIcons: { [key: string]: string } = {
  // Cloud Platforms
  'AWS': '☁️',
  'GCP': '☁️',
  'Google Cloud Platform': '☁️',
  'Azure': '☁️',
  'OpenStack': '☁️',
  
  // Kubernetes & Containerization
  'Kubernetes': '☸️',
  'Docker': '🐳',
  'EKS': '☸️',
  'Kops': '☸️',
  'Bare-metal': '🖥️',
  
  // Infrastructure as Code
  'Terraform': '🏗️',
  'Pulumi': '☁️',
  'Ansible': '🤖',
  'Chef': '👨‍🍳',
  'Saltstack': '🧂',
  
  // CI/CD & DevOps
  'CI/CD': '🔄',
  'GitHub Actions': '🐙',
  'GitLab CI/CD': '🦊',
  'Jenkins': '🤖',
  'ArgoCD': '🚀',
  'Flux': '⚡',
  'GitOps': '📦',
  
  // Monitoring & Observability
  'Prometheus': '📊',
  'Grafana': '📈',
  'Loki': '📝',
  'Tempo': '⏱️',
  'Thanos': '⚡',
  'ELK': '🦷',
  'EFK': '🦷',
  'OpenTelemetry': '👁️',
  'Monitoring': '📊',
  'Logging': '📝',
  'Tracing': '🔍',
  'Alerting': '🚨',
  'Metrics': '📈',
  
  // Programming Languages
  'Go': '🐹',
  'Golang': '🐹',
  'Python': '🐍',
  'Bash': '💻',
  'JavaScript': '📗',
  'TypeScript': '📘',
  
  // Databases & Messaging
  'PostgreSQL': '🐘',
  'Redis': '🔴',
  'RabbitMQ': '🐰',
  'MongoDB': '🍃',
  'Kafka': '📨',
  
  // Distributed Systems
  'Mesos': '🕷️',
  'Consul': '🏛️',
  'Linkerd': '🔗',
  'Distributed Systems': '🌐',
  
  // Serverless & Platforms
  'Serverless': '⚡',
  'AWS Lambda': '⚡',
  'Knative': '☸️',
  'CloudEvents': '☁️',
  
  // Security
  'Security': '🔒',
  'Compliance': '📋',
  'Network Security': '🛡️',
  'VPN': '🔐',
  
  // AI/ML
  'RAG': '🤖',
  'Vertex AI': '🧠',
  'Machine Learning': '🤖',
  
  // Networking
  'Load Balancing': '⚖️',
  'API Gateway': '🚪',
  'Service Mesh': '🕸️',
  
  // Management
  'Team Leadership': '👥',
  'People Management': '👥',
  'Project Management': '📊',
  'Agile/Scrum': '🔄',
  
  // General
  'Infrastructure': '🏗️',
  'Automation': '⚙️',
  'Operations': '🔧',
  'Cloud Operations': '☁️',
  'Infrastructure as Code': '🏗️',
  'Cloud Migration': '🔄',
  'VMware ESXi': '💻',
  'Collaboration': '🤝',
  'Troubleshooting': '🔧',
  'Observability': '👁️',
  'Cloud Native Infrastructure': '☸️',
  'Cloud-Native Infrastructure': '☸️',
  'Problem-Solving': '🧩'
}

const Resume: React.FC = () => {
  const [experience, setExperience] = useState<Experience[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchExperience = async () => {
      try {
        const data = await apiClient.getExperiences()
        console.log('Raw experience data:', data)
        
        // Remove duplicates based on unique combination of title, company, and start_date
        const uniqueData = data.reduce((acc, current) => {
          const key = `${current.title}-${current.company}-${current.start_date}`
          const exists = acc.find(item => 
            `${item.title}-${item.company}-${item.start_date}` === key
          )
          if (!exists) {
            acc.push(current)
          }
          return acc
        }, [] as Experience[])
        
        console.log('Unique experience data:', uniqueData)
        
        const sortedData = uniqueData.sort((a, b) => {
          // Sort by order first (highest first), then by start_date (most recent first)
          if (a.order !== b.order) {
            return b.order - a.order
          }
          const dateA = new Date(a.start_date)
          const dateB = new Date(b.start_date)
          return dateB.getTime() - dateA.getTime()
        })
        
        console.log('Sorted experience data:', sortedData)
        console.log('Technologies check:', sortedData.map(exp => ({ title: exp.title, technologies: exp.technologies, length: exp.technologies?.length })))
        setExperience(sortedData)
        setError(null)
      } catch (err) {
        console.error('Failed to fetch experiences:', err)
        setError('Failed to fetch experience data')
      } finally {
        setLoading(false)
      }
    }

    fetchExperience()
  }, [])

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'long' 
    })
  }

  const formatPeriod = (startDate: string, endDate: string | null, current: boolean) => {
    const start = formatDate(startDate)
    if (current) {
      return `${start} - Present`
    }
    if (endDate) {
      const end = formatDate(endDate)
      return `${start} - ${end}`
    }
    return start
  }

  const getTechnologyIcon = (technology: string): string => {
    return technologyIcons[technology] || '🔧'
  }

  if (loading) {
    return (
      <div className="resume">
        <div className="container">
          <h1>Bruno Lucena</h1>
          <h2>Senior Cloud Native Infrastructure Engineer</h2>
          <div className="loading">
            <p>Loading professional experience from database...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="resume">
        <div className="container">
          <h1>Bruno Lucena</h1>
          <h2>Senior Cloud Native Infrastructure Engineer</h2>
          <div className="error">
            <p>Error loading experience: {error}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="resume">
      <div className="container">
        <h1>Bruno Lucena</h1>
        <h2>Senior Cloud Native Infrastructure Engineer</h2>
        
        <section className="resume-section">
          <h3>Professional Experience</h3>
          {experience.length === 0 ? (
            <div className="no-experience">
              <p>No experience data available.</p>
            </div>
          ) : (
            <div className="experience-items">
              {experience.map((exp) => (
                <div key={exp.id} className="experience-item">
                  <div className="experience-header">
                    <h4 className="experience-title">{exp.title}</h4>
                    <span className="experience-company">
                      {exp.company === 'Crealytics' && (
                        <a href="https://www.crealytics.com/" target="_blank" rel="noopener noreferrer" className="company-link">
                          {exp.company}
                        </a>
                      )}
                      {exp.company === 'Tempest Security Intelligence' && (
                        <a href="https://www.tempest.com.br/" target="_blank" rel="noopener noreferrer" className="company-link">
                          {exp.company}
                        </a>
                      )}
                      {exp.company === 'Mobimeo' && (
                        <a href="https://mobimeo.com/en/home-page/" target="_blank" rel="noopener noreferrer" className="company-link">
                          {exp.company}
                        </a>
                      )}
                      {exp.company === 'Notifi' && (
                        <a href="http://notifi.network/" target="_blank" rel="noopener noreferrer" className="company-link">
                          {exp.company}
                        </a>
                      )}
                      {exp.company === 'Namecheap, Inc' && (
                        <a href="https://www.namecheap.com/" target="_blank" rel="noopener noreferrer" className="company-link">
                          {exp.company}
                        </a>
                      )}
                      {exp.company === 'Lesara' && (
                        <a href="https://www.linkedin.com/company/lesara/" target="_blank" rel="noopener noreferrer" className="company-link">
                          {exp.company}
                        </a>
                      )}
                      {!['Crealytics', 'Tempest Security Intelligence', 'Mobimeo', 'Notifi', 'Namecheap, Inc', 'Lesara'].includes(exp.company) && exp.company}
                    </span>
                    <span className="experience-period">
                      {formatPeriod(exp.start_date, exp.end_date, exp.current)}
                    </span>
                  </div>
                  <div className="experience-description">
                    <p>{exp.description}</p>
                  </div>
                  {exp.technologies && exp.technologies.length > 0 && (
                    <div className="experience-technologies">
                      <strong>Technologies used at {exp.company}:</strong>
                      <div className="technology-icons">
                        {exp.technologies.map((tech, index) => (
                          <span key={index} className="technology-icon" title={tech}>
                            {getTechnologyIcon(tech)}
                          </span>
                        ))}
                      </div>
                      <div className="technology-names">
                        {exp.technologies.join(', ')}
                      </div>
                    </div>
                                    )}
                  {(!exp.technologies || exp.technologies.length === 0) && (
                    <div className="experience-technologies">
                      <strong>DEBUG: No technologies for {exp.company}</strong>
                      <div>Technologies: {JSON.stringify(exp.technologies)}</div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  )
}

export default Resume
