#!/bin/bash

# 🔒 Bruno Site Security Scanning Script
# This script provides comprehensive security scanning for the Bruno Site project

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCAN_TYPE="${1:-all}"
OUTPUT_DIR="${2:-./security-reports}"
VERBOSE="${3:-false}"
FAIL_ON_CRITICAL="${4:-true}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if required tools are installed
check_tools() {
    local tools=("docker" "trivy" "gosec" "gitleaks" "trufflehog")
    local missing_tools=()
    
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_info "Run 'make security-install-tools' to install missing tools"
        exit 1
    fi
    
    log_success "All required security tools are available"
}

# Container vulnerability scanning
scan_containers() {
    log_info "Scanning Docker containers for vulnerabilities..."
    
    local containers=("api" "frontend")
    local critical_vulns=0
    
    for container in "${containers[@]}"; do
        log_info "Building and scanning $container container..."
        
        # Build container
        docker build -t "bruno-site-$container:scan" "./$container" || {
            log_error "Failed to build $container container"
            continue
        }
        
        # Scan with Trivy
        local output_file="$OUTPUT_DIR/trivy-$container-results.json"
        trivy image --severity CRITICAL,HIGH,MEDIUM --format json --output "$output_file" "bruno-site-$container:scan" || {
            log_warning "Trivy scan failed for $container"
            continue
        }
        
        # Count critical vulnerabilities
        local critical_count=$(jq -r '.Results[]?.Vulnerabilities[]? | select(.Severity == "CRITICAL") | .VulnerabilityID' "$output_file" | wc -l)
        critical_vulns=$((critical_vulns + critical_count))
        
        log_success "Container scan completed for $container"
        log_info "Results saved to: $output_file"
    done
    
    if [ "$critical_vulns" -gt 0 ]; then
        log_warning "Found $critical_vulns critical vulnerabilities in containers"
        if [ "$FAIL_ON_CRITICAL" = "true" ]; then
            log_error "Critical vulnerabilities detected. Failing scan."
            exit 1
        fi
    fi
}

# Dependency vulnerability scanning
scan_dependencies() {
    log_info "Scanning dependencies for vulnerabilities..."
    
    # Go dependencies
    if [ -d "api" ]; then
        log_info "Scanning Go dependencies..."
        local go_output="$OUTPUT_DIR/trivy-go-results.json"
        cd api && trivy fs --severity CRITICAL,HIGH,MEDIUM --format json --output "../$go_output" . || {
            log_warning "Go dependency scan failed"
        }
        cd ..
        log_success "Go dependency scan completed"
    fi
    
    # Node.js dependencies
    if [ -d "frontend" ]; then
        log_info "Scanning Node.js dependencies..."
        local npm_output="$OUTPUT_DIR/npm-audit-results.json"
        cd frontend && npm audit --audit-level=high --json > "../$npm_output" 2>/dev/null || {
            log_warning "npm audit failed"
        }
        cd ..
        log_success "Node.js dependency scan completed"
    fi
}

# Code security analysis
scan_code() {
    log_info "Scanning code for security issues..."
    
    # Go code security
    if [ -d "api" ]; then
        log_info "Scanning Go code with GoSec..."
        local gosec_output="$OUTPUT_DIR/gosec-results.json"
        cd api && gosec -fmt=json -out="../$gosec_output" ./... || {
            log_warning "GoSec scan failed"
        }
        cd ..
        log_success "Go code security scan completed"
    fi
    
    # JavaScript/TypeScript code security
    if [ -d "frontend" ]; then
        log_info "Scanning JavaScript/TypeScript code..."
        cd frontend
        
        # ESLint security rules
        if [ -f "package.json" ] && grep -q "eslint" package.json; then
            local eslint_output="../$OUTPUT_DIR/eslint-security-results.json"
            npx eslint . --ext .js,.jsx,.ts,.tsx --format json --output-file "$eslint_output" || {
                log_warning "ESLint security scan failed"
            }
        fi
        
        # npm audit for security issues
        npm audit --audit-level=high || {
            log_warning "npm audit failed"
        }
        
        cd ..
        log_success "JavaScript/TypeScript code security scan completed"
    fi
}

# Secrets detection
scan_secrets() {
    log_info "Scanning for secrets and sensitive information..."
    
    # GitLeaks scan
    log_info "Running GitLeaks scan..."
    local gitleaks_output="$OUTPUT_DIR/gitleaks-results.json"
    gitleaks detect --source . --report-format json --report-path "$gitleaks_output" || {
        log_warning "GitLeaks scan failed"
    }
    
    # TruffleHog scan
    log_info "Running TruffleHog scan..."
    local trufflehog_output="$OUTPUT_DIR/trufflehog-results.json"
    trufflehog --only-verified --format json . > "$trufflehog_output" || {
        log_warning "TruffleHog scan failed"
    }
    
    log_success "Secrets detection scan completed"
}

# Generate comprehensive security report
generate_report() {
    log_info "Generating comprehensive security report..."
    
    local report_file="$OUTPUT_DIR/security-scan-report.md"
    
    cat > "$report_file" << EOF
# 🔒 Bruno Site Security Scan Report

**Scan Date:** $(date -u +"%Y-%m-%d %H:%M:%S UTC")
**Repository:** Bruno Site
**Branch:** $(git branch --show-current 2>/dev/null || echo "unknown")
**Commit:** $(git rev-parse HEAD 2>/dev/null || echo "unknown")
**Scan Type:** $SCAN_TYPE
**Output Directory:** $OUTPUT_DIR

## 📊 Scan Summary

### Container Security
- **Status:** $(if [ -f "$OUTPUT_DIR/trivy-api-results.json" ] || [ -f "$OUTPUT_DIR/trivy-frontend-results.json" ]; then echo "✅ Completed"; else echo "❌ Failed"; fi)
- **Critical Vulnerabilities:** $(find "$OUTPUT_DIR" -name "trivy-*-results.json" -exec jq -r '.Results[]?.Vulnerabilities[]? | select(.Severity == "CRITICAL") | .VulnerabilityID' {} \; 2>/dev/null | wc -l)
- **High Vulnerabilities:** $(find "$OUTPUT_DIR" -name "trivy-*-results.json" -exec jq -r '.Results[]?.Vulnerabilities[]? | select(.Severity == "HIGH") | .VulnerabilityID' {} \; 2>/dev/null | wc -l)

### Dependency Security
- **Status:** $(if [ -f "$OUTPUT_DIR/trivy-go-results.json" ] || [ -f "$OUTPUT_DIR/npm-audit-results.json" ]; then echo "✅ Completed"; else echo "❌ Failed"; fi)
- **Go Dependencies:** $(if [ -f "$OUTPUT_DIR/trivy-go-results.json" ]; then jq -r '.Results[]?.Vulnerabilities[]? | select(.Severity == "CRITICAL" or .Severity == "HIGH") | .VulnerabilityID' "$OUTPUT_DIR/trivy-go-results.json" 2>/dev/null | wc -l; else echo "0"; fi)
- **Node.js Dependencies:** $(if [ -f "$OUTPUT_DIR/npm-audit-results.json" ]; then jq -r '.vulnerabilities | to_entries[] | select(.value.severity == "high" or .value.severity == "critical") | .key' "$OUTPUT_DIR/npm-audit-results.json" 2>/dev/null | wc -l; else echo "0"; fi)

### Code Security
- **Status:** $(if [ -f "$OUTPUT_DIR/gosec-results.json" ] || [ -f "$OUTPUT_DIR/eslint-security-results.json" ]; then echo "✅ Completed"; else echo "❌ Failed"; fi)
- **Go Security Issues:** $(if [ -f "$OUTPUT_DIR/gosec-results.json" ]; then jq -r '.Issues | length' "$OUTPUT_DIR/gosec-results.json" 2>/dev/null || echo "0"; else echo "0"; fi)
- **JavaScript Security Issues:** $(if [ -f "$OUTPUT_DIR/eslint-security-results.json" ]; then jq -r 'length' "$OUTPUT_DIR/eslint-security-results.json" 2>/dev/null || echo "0"; else echo "0"; fi)

### Secrets Detection
- **Status:** $(if [ -f "$OUTPUT_DIR/gitleaks-results.json" ] || [ -f "$OUTPUT_DIR/trufflehog-results.json" ]; then echo "✅ Completed"; else echo "❌ Failed"; fi)
- **GitLeaks Findings:** $(if [ -f "$OUTPUT_DIR/gitleaks-results.json" ]; then jq -r 'length' "$OUTPUT_DIR/gitleaks-results.json" 2>/dev/null || echo "0"; else echo "0"; fi)
- **TruffleHog Findings:** $(if [ -f "$OUTPUT_DIR/trufflehog-results.json" ]; then jq -r 'length' "$OUTPUT_DIR/trufflehog-results.json" 2>/dev/null || echo "0"; else echo "0"; fi)

## 📁 Generated Files

$(find "$OUTPUT_DIR" -name "*.json" -o -name "*.md" | sort | while read -r file; do
    echo "- $(basename "$file"): $(basename "$(dirname "$file")")/$(basename "$file")"
done)

## 🚨 Critical Issues Summary

$(if [ -f "$OUTPUT_DIR/trivy-api-results.json" ] || [ -f "$OUTPUT_DIR/trivy-frontend-results.json" ]; then
    echo "### Container Vulnerabilities"
    find "$OUTPUT_DIR" -name "trivy-*-results.json" -exec jq -r '.Results[]?.Vulnerabilities[]? | select(.Severity == "CRITICAL") | "- " + .VulnerabilityID + ": " + .Title' {} \; 2>/dev/null | head -10
    echo ""
fi)

$(if [ -f "$OUTPUT_DIR/gitleaks-results.json" ]; then
    echo "### Exposed Secrets"
    jq -r '.[] | "- " + .RuleID + " in " + .File' "$OUTPUT_DIR/gitleaks-results.json" 2>/dev/null | head -5
    echo ""
fi)

## 🔧 Recommendations

1. **Immediate Actions:**
   - Address all critical vulnerabilities
   - Remove any exposed secrets
   - Update dependencies with known vulnerabilities

2. **Short-term Actions:**
   - Implement security headers
   - Add rate limiting
   - Enable authentication for API endpoints

3. **Long-term Actions:**
   - Set up automated security scanning in CI/CD
   - Implement security monitoring
   - Regular security assessments

## 📞 Security Contact

**Security Team:** Bruno Lucena  
**Email:** [Contact via LinkedIn](https://www.linkedin.com/in/bvlucena)  
**GitHub:** [brunovlucena](https://github.com/brunovlucena)

---

*This report was generated automatically by the Bruno Site security scanning script.*
EOF
    
    log_success "Security report generated: $report_file"
}

# Main execution
main() {
    log_info "Starting Bruno Site security scan..."
    log_info "Scan type: $SCAN_TYPE"
    log_info "Output directory: $OUTPUT_DIR"
    log_info "Fail on critical: $FAIL_ON_CRITICAL"
    
    # Check tools
    check_tools
    
    # Run scans based on type
    case "$SCAN_TYPE" in
        "all")
            scan_containers
            scan_dependencies
            scan_code
            scan_secrets
            ;;
        "containers")
            scan_containers
            ;;
        "dependencies")
            scan_dependencies
            ;;
        "code")
            scan_code
            ;;
        "secrets")
            scan_secrets
            ;;
        *)
            log_error "Invalid scan type: $SCAN_TYPE"
            log_info "Valid options: all, containers, dependencies, code, secrets"
            exit 1
            ;;
    esac
    
    # Generate report
    generate_report
    
    log_success "Security scan completed successfully!"
    log_info "Check the output directory for detailed results: $OUTPUT_DIR"
}

# Run main function
main "$@"
