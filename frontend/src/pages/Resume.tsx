import React, { useState, useEffect } from 'react'

interface Experience {
  id: number
  title: string
  company: string
  start_date: string
  end_date: string | null
  current: boolean
  description: string
  technologies: string[]
  order: number
}

const Resume: React.FC = () => {
  const [experience, setExperience] = useState<Experience[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchExperience = async () => {
      try {
        const response = await fetch('/api/v1/content/experience')
        if (!response.ok) {
          throw new Error('Failed to fetch experience data')
        }
        const data = await response.json()
        const sortedData = (data || []).sort((a, b) => {
          // Sort by start_date in descending order (most recent first)
          const dateA = new Date(a.start_date)
          const dateB = new Date(b.start_date)
          return dateB.getTime() - dateA.getTime()
        })
        setExperience(sortedData)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred')
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
                      <strong>Technologies:</strong> {exp.technologies.join(', ')}
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
