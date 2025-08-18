import React from 'react'
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
import { FaShieldAlt } from 'react-icons/fa'
import { BiCertification } from 'react-icons/bi'

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
