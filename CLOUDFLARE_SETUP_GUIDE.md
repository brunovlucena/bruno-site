# 🛡️ Cloudflare Setup Guide for Bruno Site

A comprehensive guide to configure Cloudflare CDN, security, and performance optimization for your Bruno site.

## 📋 Overview

This guide will help you:
- 🚀 Set up Cloudflare CDN for global content delivery
- 🛡️ Configure security features (WAF, DDoS protection)
- ⚡ Optimize performance with caching and compression
- 🔒 Enable SSL/TLS encryption
- 📊 Monitor performance and security

## 🚀 Step 1: Cloudflare Account Setup

### 1.1 Create Cloudflare Account
1. Go to [cloudflare.com](https://cloudflare.com)
2. Click "Sign Up" and create a free account
3. Choose the **Free plan** (perfect for starting out)

### 1.2 Add Your Domain
1. Click "Add a Site"
2. Enter your domain name (e.g., `brunosite.com`)
3. Select the **Free plan**
4. Cloudflare will scan your existing DNS records

## 🔧 Step 2: DNS Configuration

### 2.1 Update Nameservers
1. Copy the Cloudflare nameservers provided
2. Update your domain registrar's nameservers:
   ```
   Nameserver 1: [cloudflare-ns1]
   Nameserver 2: [cloudflare-ns2]
   ```

### 2.2 Configure DNS Records
Add these DNS records in Cloudflare:

```yaml
# Main website
Type: A
Name: @
Content: [YOUR_SERVER_IP]
Proxy: ✅ (Orange cloud)

# API subdomain
Type: A
Name: api
Content: [YOUR_SERVER_IP]
Proxy: ✅ (Orange cloud)

# www subdomain
Type: CNAME
Name: www
Content: [YOUR_DOMAIN]
Proxy: ✅ (Orange cloud)
```

## ⚡ Step 3: Performance Optimization

### 3.1 Caching Configuration
Navigate to **Caching > Configuration**:

```yaml
# Browser Cache TTL
- CSS/JS files: 1 year
- Images: 1 month
- HTML: 4 hours
- API responses: 1 hour

# Edge Cache TTL
- Static assets: 1 year
- API responses: 15 minutes
```

### 3.2 Page Rules Setup
Create these Page Rules for optimal caching:

```yaml
# Rule 1: Cache static assets
URL: brunosite.com/static/*
Settings:
  - Cache Level: Cache Everything
  - Edge Cache TTL: 1 year
  - Browser Cache TTL: 1 year

# Rule 2: API caching
URL: api.brunosite.com/*
Settings:
  - Cache Level: Cache Everything
  - Edge Cache TTL: 15 minutes
  - Browser Cache TTL: 1 hour

# Rule 3: HTML pages
URL: brunosite.com/*
Settings:
  - Cache Level: Standard
  - Edge Cache TTL: 4 hours
  - Browser Cache TTL: 1 hour
```

## 🛡️ Step 4: Security Configuration

### 4.1 SSL/TLS Settings
Navigate to **SSL/TLS > Overview**:

```yaml
Encryption Mode: Full (strict)
Minimum TLS Version: TLS 1.2
Opportunistic Encryption: ✅
TLS 1.3: ✅
Automatic HTTPS Rewrites: ✅
```

### 4.2 Security Headers
Navigate to **Security > Security Headers**:

```yaml
# Content Security Policy
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(), geolocation=()
```

### 4.3 WAF (Web Application Firewall)
Navigate to **Security > WAF**:

```yaml
# Managed Rules
- Cloudflare Managed Rules: ✅
- OWASP Top 10: ✅
- WordPress: ❌ (not applicable)
- Drupal: ❌ (not applicable)

# Custom Rules
- Block suspicious user agents
- Rate limit API requests
- Block known bad IPs
```

### 4.4 Rate Limiting
Create rate limiting rules:

```yaml
# API Rate Limiting
Rule: Rate limit API requests
Expression: (http.request.uri.path contains "/api/")
Rate: 100 requests per minute
Action: Block

# Login Rate Limiting
Rule: Rate limit login attempts
Expression: (http.request.uri.path contains "/login")
Rate: 5 attempts per minute
Action: Block
```

## 🔧 Step 5: Application Configuration

### 5.1 Workers (Optional)
Create a Cloudflare Worker for dynamic content:

```javascript
// workers/dynamic-content.js
addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  // Add custom headers
  const response = await fetch(request)
  const newResponse = new Response(response.body, response)
  
  newResponse.headers.set('X-Custom-Header', 'Bruno-Site')
  newResponse.headers.set('Cache-Control', 'public, max-age=3600')
  
  return newResponse
}
```

### 5.2 Image Optimization
Navigate to **Speed > Optimization**:

```yaml
# Image Optimization
Polish: Lossless
WebP: ✅
AVIF: ✅
Resize Images: ✅
```

## 📊 Step 6: Monitoring & Analytics

### 6.1 Analytics Setup
Navigate to **Analytics > Traffic**:

```yaml
# Enable Analytics
Web Analytics: ✅
Real User Monitoring: ✅
Bot Analytics: ✅
```

### 6.2 Performance Monitoring
Monitor these metrics:

```yaml
# Key Performance Indicators
- Time to First Byte (TTFB)
- First Contentful Paint (FCP)
- Largest Contentful Paint (LCP)
- Cumulative Layout Shift (CLS)
- First Input Delay (FID)
```

## 🚀 Step 7: Deployment Integration

### 7.1 CI/CD Integration
Add Cloudflare deployment to your Makefile:

```makefile
# Add to your existing Makefile
deploy-cloudflare:
	@echo "🚀 Deploying to Cloudflare..."
	# Build frontend
	cd frontend && npm run build
	# Deploy to Cloudflare Pages or Workers
	# Add your deployment commands here

purge-cache:
	@echo "🧹 Purging Cloudflare cache..."
	# Add cache purge commands
```

### 7.2 Environment Variables
Set these in your deployment environment:

```bash
# Cloudflare Configuration
CLOUDFLARE_API_TOKEN=your_api_token
CLOUDFLARE_ZONE_ID=your_zone_id
CLOUDFLARE_ACCOUNT_ID=your_account_id
```

## 🔍 Step 8: Testing & Validation

### 8.1 Performance Testing
Use these tools to validate your setup:

```bash
# Test CDN performance
curl -I https://yourdomain.com
curl -I https://api.yourdomain.com

# Check security headers
curl -I -s https://yourdomain.com | grep -i security

# Test caching
curl -I https://yourdomain.com/static/app.js
```

### 8.2 Security Testing
Validate security configuration:

```bash
# SSL Labs test
https://www.ssllabs.com/ssltest/

# Security Headers test
https://securityheaders.com/

# Observatory test
https://observatory.mozilla.org/
```

## 📈 Step 9: Optimization Tips

### 9.1 Caching Strategy
```yaml
# Static Assets (CSS, JS, Images)
Cache-Control: public, max-age=31536000, immutable

# API Responses
Cache-Control: public, max-age=900, s-maxage=900

# HTML Pages
Cache-Control: public, max-age=14400, s-maxage=14400
```

### 9.2 Compression
Enable these compression options:

```yaml
# Gzip Compression
- CSS: ✅
- JavaScript: ✅
- HTML: ✅
- JSON: ✅
- XML: ✅

# Brotli Compression
- Enable for all text-based content
```

## 🚨 Troubleshooting

### Common Issues and Solutions

```yaml
# Issue: Mixed Content Warnings
Solution: Enable "Always Use HTTPS" in SSL/TLS settings

# Issue: Caching Not Working
Solution: Check Page Rules and Cache-Control headers

# Issue: API Requests Blocked
Solution: Review WAF rules and rate limiting settings

# Issue: Slow Performance
Solution: Check Edge Cache TTL and compression settings
```

## 📞 Support Resources

- **Cloudflare Documentation**: [developers.cloudflare.com](https://developers.cloudflare.com)
- **Community Forum**: [community.cloudflare.com](https://community.cloudflare.com)
- **Status Page**: [cloudflarestatus.com](https://cloudflarestatus.com)

## 🎯 Next Steps

1. **Monitor Performance**: Set up alerts for performance metrics
2. **Security Audits**: Regular security assessments
3. **Optimization**: Continuous performance improvements
4. **Backup Strategy**: Implement backup and recovery procedures

---

*This guide provides a comprehensive Cloudflare setup for your Bruno site. Follow each step carefully and test thoroughly before going live.*
