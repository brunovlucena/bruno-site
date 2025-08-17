import React from 'react'
import { Link } from 'react-router-dom'

const Header: React.FC = () => {
  return (
    <header className="header">
      <div className="header-container">
        <div className="header-brand">
          <Link to="/" className="logo">
            <span className="logo-text">Bruno Lucena</span>
          </Link>
        </div>
        
        <nav className="header-nav">
          <ul className="nav-menu">
            <li className="nav-item">
              <Link to="/" className="nav-link">Home</Link>
            </li>
            <li className="nav-item">
              <Link to="/resume" className="nav-link">Resume</Link>
            </li>
          </ul>
        </nav>
        
        <div className="header-actions">
          <a href="https://www.linkedin.com/in/bvlucena" className="header-link" target="_blank" rel="noopener noreferrer">
            <div className="header-link-icon">🔗</div>
            <span>LinkedIn</span>
          </a>
          <a href="https://github.com/brunovlucena" className="header-link" target="_blank" rel="noopener noreferrer">
            <div className="header-link-icon">🐙</div>
            <span>GitHub</span>
          </a>
        </div>
      </div>
    </header>
  )
}

export default Header
