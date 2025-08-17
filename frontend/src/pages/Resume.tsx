import React from 'react'

const Resume: React.FC = () => {
  return (
    <div className="resume">
      <div className="container">
        <h1>Bruno Lucena</h1>
        <h2>Senior SRE/DevSecOps Engineer</h2>
        
        <section className="resume-section">
          <h3>Experience</h3>
          <div className="experience-item">
            <h4>Senior Infrastructure Engineer</h4>
            <p>12+ years of experience in cloud-native infrastructure</p>
          </div>
        </section>
        
        <section className="resume-section">
          <h3>Skills</h3>
          <ul>
            <li>Kubernetes & Docker</li>
            <li>AWS, GCP, Azure</li>
            <li>Terraform, Pulumi</li>
            <li>CI/CD Pipelines</li>
            <li>Observability & Monitoring</li>
          </ul>
        </section>
      </div>
    </div>
  )
}

export default Resume
