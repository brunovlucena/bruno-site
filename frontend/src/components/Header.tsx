import React from 'react'
import { Link } from 'react-router-dom'

const Header: React.FC = () => {
  return (
    <header className="header">
      <div className="header-container">
        <div className="header-left">
                    <Link to="/" className="logo">
            <span className="logo-text">Bruno Lucena</span>
          </Link>
          
          {/* Homelab Link */}
          <Link to="/?section=projects" className="homelab-link">
            Homelab
          </Link>
        </div>
        
        <div className="header-actions">
          <a href="#" className="header-link">
            <div className="header-link-icon">🔍</div>
            <span>Search</span>
          </a>
          <Link to="/resume" className="header-link">
            <div className="header-link-icon">📖</div>
            <span>Resume</span>
          </Link>
          <a href="https://www.linkedin.com/in/bvlucena" className="header-link" target="_blank" rel="noopener noreferrer">
            <div className="header-link-icon">🔗</div>
            <span>LinkedIn</span>
          </a>
          <a href="https://github.com/brunovlucena" className="header-link" target="_blank" rel="noopener noreferrer">
            <div className="header-link-icon">🐙</div>
            <span>GitHub</span>
          </a>
          <a href="/contact" className="header-link">
            <span>Contact</span>
          </a>
        </div>
      </div>
    </header>
  )
}

export default Header
