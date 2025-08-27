# 🔒 Security Scanning Guide

This document describes the comprehensive security scanning setup for the Bruno Site project, including automated scanning, manual scanning, and best practices.

## 📋 Overview

The Bruno Site project implements a multi-layered security scanning approach:

- **Automated Scanning**: GitHub Actions workflows run regularly
- **Manual Scanning**: Local tools for development and testing
- **Dependency Monitoring**: Dependabot for automatic updates
- **Secrets Detection**: GitLeaks and TruffleHog for sensitive data
- **Container Security**: Trivy for Docker image vulnerabilities
- **Code Security**: GoSec and ESLint for code analysis

## 🚀 Quick Start

### 1. Install Security Tools

```bash
# Install all required security scanning tools
make security-install-tools
```

### 2. Run Comprehensive Security Scan

```bash
# Run all security scans
make security-scan

# Or use the script directly
./scripts/security-scan.sh all
```

### 3. Run Specific Scans

```bash
# Scan only containers
make security-scan-containers

# Scan only dependencies
make security-scan-deps

# Scan only code
make security-scan-code

# Scan only for secrets
make security-scan-secrets
```

## 🔧 Available Commands

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make security-install-tools` | Install all security scanning tools |
| `make security-scan` | Run comprehensive security scan |
| `make security-scan-local` | Run basic local security scan |
| `make security-scan-containers` | Scan Docker containers |
| `make security-scan-deps` | Scan dependencies |
| `make security-scan-code` | Scan code for security issues |
| `make security-scan-secrets` | Scan for secrets |
| `make security-report` | Generate security report |

### Script Commands

```bash
# Basic usage
./scripts/security-scan.sh [scan_type] [output_dir] [verbose] [fail_on_critical]

# Examples
./scripts/security-scan.sh all ./security-reports false true
./scripts/security-scan.sh containers ./reports true false
./scripts/security-scan.sh dependencies
./scripts/security-scan.sh code
./scripts/security-scan.sh secrets
```

## 🛠️ Security Tools

### 1. Trivy (Container & Dependency Scanning)
- **Purpose**: Scan Docker images and dependencies for vulnerabilities
- **Installation**: `curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.48.0`
- **Usage**: `trivy image [image_name]` or `trivy fs [directory]`

### 2. GoSec (Go Code Security)
- **Purpose**: Static analysis of Go code for security issues
- **Installation**: `go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest`
- **Usage**: `gosec ./...`

### 3. GitLeaks (Secrets Detection)
- **Purpose**: Detect secrets and sensitive information in code
- **Installation**: `curl -sSfL https://raw.githubusercontent.com/gitleaks/gitleaks/master/install.sh | sh -s -- -b /usr/local/bin v8.18.0`
- **Usage**: `gitleaks detect --source .`

### 4. TruffleHog (Additional Secrets Detection)
- **Purpose**: Additional secrets detection with verification
- **Installation**: `pip install trufflehog`
- **Usage**: `trufflehog --only-verified .`

## 🔄 Automated Scanning

### GitHub Actions Workflow

The project includes a comprehensive GitHub Actions workflow (`.github/workflows/security-scan.yml`) that runs:

- **Scheduled**: Every Monday at 2 AM UTC
- **On Push**: To main and develop branches
- **On PR**: For all pull requests
- **Manual**: Via workflow dispatch

### Dependabot Configuration

Dependabot is configured to automatically:

- **Weekly Updates**: Regular dependency updates on Mondays
- **Daily Security**: Critical security updates daily
- **Multiple Ecosystems**: Go, Node.js, Docker, GitHub Actions

## 📊 Scan Results

### Output Files

Security scans generate the following files:

| File | Description |
|------|-------------|
| `trivy-*-results.json` | Container vulnerability results |
| `gosec-results.json` | Go code security issues |
| `gitleaks-results.json` | Detected secrets |
| `trufflehog-results.json` | Additional secrets findings |
| `npm-audit-results.json` | Node.js dependency vulnerabilities |
| `eslint-security-results.json` | JavaScript security issues |
| `security-report.md` | Comprehensive security report |

### Understanding Results

#### Critical Vulnerabilities
- **Action Required**: Immediate attention needed
- **Impact**: High risk of exploitation
- **Response**: Fix within 24 hours

#### High Vulnerabilities
- **Action Required**: Priority attention
- **Impact**: Significant risk
- **Response**: Fix within 72 hours

#### Medium Vulnerabilities
- **Action Required**: Plan for remediation
- **Impact**: Moderate risk
- **Response**: Fix within 1 week

#### Low Vulnerabilities
- **Action Required**: Monitor and plan
- **Impact**: Low risk
- **Response**: Fix within 1 month

## 🔧 Configuration

### GitLeaks Configuration

The `.gitleaks.toml` file configures:

- **Allowlist**: Files and patterns to ignore
- **Rules**: Custom detection patterns
- **Severity Levels**: Risk assessment for findings

### Security Headers

Recommended security headers for the application:

```nginx
# Nginx configuration
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

## 🚨 Security Alerts

### Critical Issues

When critical vulnerabilities are detected:

1. **Immediate Notification**: GitHub Actions will fail and alert
2. **Security Report**: Detailed report generated
3. **Remediation Plan**: Immediate action required
4. **Follow-up**: Verification of fixes

### Alert Channels

- **GitHub Issues**: Automatic issue creation for critical findings
- **Pull Request Comments**: Security summary on PRs
- **Security Tab**: Results uploaded to GitHub Security tab
- **Email Notifications**: Configured via GitHub settings

## 📈 Security Metrics

### Key Performance Indicators

| Metric | Target | Measurement |
|--------|--------|-------------|
| Vulnerability Density | < 1 per 1000 lines | Automated scanning |
| Time to Fix (Critical) | < 24 hours | Issue tracking |
| Time to Fix (High) | < 72 hours | Issue tracking |
| Security Test Coverage | > 90% | Coverage reports |
| Dependency Update Frequency | Weekly | Dependabot metrics |

### Reporting

Regular security reports include:

- **Vulnerability Summary**: Count by severity
- **Trend Analysis**: Changes over time
- **Remediation Status**: Progress on fixes
- **Risk Assessment**: Overall security posture

## 🔍 Manual Security Testing

### Local Development

```bash
# Quick security check during development
make security-scan-local

# Check specific component
make security-scan-containers
make security-scan-deps
```

### Pre-commit Checks

Add to your development workflow:

```bash
# Before committing code
make security-scan-secrets
make security-scan-code

# Before pushing
make security-scan
```

### Integration Testing

```bash
# Test security in running application
curl -I http://localhost:8080/health
curl -I http://localhost:3000

# Check security headers
curl -I -s http://localhost:8080 | grep -i security
```

## 🛡️ Best Practices

### Development

1. **Regular Scans**: Run security scans before commits
2. **Dependency Updates**: Keep dependencies updated
3. **Code Review**: Review security findings
4. **Documentation**: Document security decisions

### Deployment

1. **Pre-deployment Scan**: Scan before production deployment
2. **Security Gates**: Block deployment on critical issues
3. **Monitoring**: Monitor for new vulnerabilities
4. **Incident Response**: Have plan for security incidents

### Maintenance

1. **Tool Updates**: Keep security tools updated
2. **Configuration Review**: Regular review of security configs
3. **Training**: Regular security training for team
4. **Audit**: Regular security audits

## 📞 Support

### Security Issues

For security issues or questions:

- **GitHub Issues**: Create issue with `security` label
- **Email**: Contact via [LinkedIn](https://www.linkedin.com/in/bvlucena)
- **Documentation**: Check this guide and vulnerability report

### Tool Issues

For tool-specific issues:

- **Trivy**: [GitHub Repository](https://github.com/aquasecurity/trivy)
- **GoSec**: [GitHub Repository](https://github.com/securecodewarrior/gosec)
- **GitLeaks**: [GitHub Repository](https://github.com/gitleaks/gitleaks)
- **TruffleHog**: [GitHub Repository](https://github.com/trufflesecurity/trufflehog)

## 📚 Additional Resources

- [Vulnerability Analysis Report](../VULNERABILITY_ANALYSIS_REPORT.md)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Security Headers](https://securityheaders.com/)
- [Mozilla Security Guidelines](https://infosec.mozilla.org/guidelines/)

---

*This guide should be updated regularly as security tools and practices evolve.*
