import React from 'react'

const Home: React.FC = () => {
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
              <p>Senior Cloud Native Infrastructure Engineer with 12+ years of experience architecting and maintaining robust cloud-native solutions on AWS, GCP, and bare-metal environments. Expert in Kubernetes orchestration, observability systems, and automation using Terraform, Pulumi, and CI/CD pipelines.</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default Home
