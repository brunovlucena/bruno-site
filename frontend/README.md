# 🎨 Portfolio Frontend

A modern, static HTML/CSS/JS frontend for Bruno Lucena's portfolio website. This application showcases SRE/DevSecOps and AI Engineering skills through an interactive, responsive interface with AI-powered chatbot integration.

## 🚀 Features

- **Static HTML/CSS/JS** - Lightweight, fast-loading interface
- **Responsive Design** - Works perfectly on all devices
- **Dynamic Content** - Fetches data from the Go API
- **AI Chatbot** - Ollama-powered intelligent assistant
- **Real-time Updates** - Live project and skill information
- **Performance Optimized** - Minimal dependencies, fast loading
- **Accessibility** - WCAG compliant design

## 🛠️ Tech Stack

- **Frontend**: Static HTML/CSS/JavaScript
- **Build Tool**: Vite (for development)
- **Styling**: Custom CSS with CSS Variables
- **AI Integration**: Ollama with Gemma3n model
- **HTTP Client**: Fetch API
- **Deployment**: Docker with nginx

## 📁 Project Structure

```
portfolio-frontend/
├── public/                 # Static assets
│   ├── index.html          # Main portfolio page
│   ├── resume.html         # Resume page
│   ├── styles/             # CSS files
│   │   └── main.css        # Main stylesheet
│   └── components/         # JavaScript components
│       └── header.js       # Header functionality
├── Dockerfile.dev          # Development container
├── Dockerfile              # Production container
├── package.json            # Dependencies and scripts
├── vite.config.ts          # Vite configuration
└── README.md               # This file
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Go API backend running (optional for development)

### Development

1. **Install dependencies:**
   ```bash
   cd portfolio-frontend
   npm install
   ```

2. **Start development server:**
   ```bash
   npm run dev
   ```

3. **Open in browser:**
   ```
   http://localhost:3000
   ```

### Docker Development

1. **Build and run with Docker:**
   ```bash
   docker-compose up portfolio-frontend
   ```

2. **Access the application:**
   ```
   http://localhost:3000
   ```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory:

```env
# API Configuration
VITE_API_URL=http://localhost:8080

# Chatbot Configuration
VITE_CHATBOT_ENABLED=true
```

### API Integration

The frontend communicates with the Go API backend through the following endpoints:

- **Projects**: `/api/v1/projects` - Dynamic project information
- **Content**: `/api/v1/content/*` - Skills, experience, about content
- **Admin**: `/admin/projects/*` - Project management (activate/deactivate)
- **Health**: `/health` - Service health check

### AI Chatbot Integration

The chatbot integrates with Ollama for intelligent responses:

- **Ollama Server**: Local AI inference for privacy and performance
- **Gemma3n Model**: Advanced language understanding and context awareness
- **Dynamic Responses**: Real-time information from portfolio data
- **Conversation Memory**: Persistent chat history using localStorage

## 🎨 Design System

### Color Palette

- **Primary**: `#ffffff` (White)
- **Secondary**: `#7c3aed` (Purple)
- **Background**: `#1a202c` (Dark Gray)
- **Card Background**: `#2d3748` (Medium Gray)
- **Border**: `#4a5568` (Light Gray)
- **Text Primary**: `#ffffff` (White)
- **Text Secondary**: `#a0aec0` (Light Gray)

### Typography

- **Font Family**: System fonts (San Francisco, Segoe UI, etc.)
- **Headings**: Bold weights with proper hierarchy
- **Body**: Regular weight with good line height

### Components

#### Header
- Sticky navigation with dynamic dropdown menu
- Mobile-responsive hamburger menu
- Project categorization (Infrastructure & DevOps, AI & Applications)
- Real-time project filtering from database

#### Project Cards
- Dynamic content loading from API
- Technology-specific icons with official brand colors
- Hover animations and transitions
- Responsive grid layout
- Active project filtering

#### AI Chatbot
- Floating chat interface with dark theme
- Ollama-powered intelligent responses
- Conversation memory and context awareness
- Quick suggestion buttons
- Multi-page support (index.html and resume.html)

## 📱 Responsive Design

The application is fully responsive with breakpoints:

- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

## 🚀 Deployment

### Docker

1. **Build the image:**
   ```bash
   docker build -t brunovlucena/portfolio-frontend:latest .
   ```

2. **Run the container:**
   ```bash
   docker run -p 80:80 brunovlucena/portfolio-frontend:latest
   ```

### Kubernetes

Deploy to Kubernetes using the provided manifests:

```bash
kubectl apply -f k8s/portfolio-frontend-deployment.yaml
```

### CI/CD

The frontend is automatically deployed through GitHub Actions when changes are pushed to the main branch.

## 🔍 Performance

### Optimizations

- **Static Assets**: Pre-built HTML/CSS/JS for fast loading
- **Minimal Dependencies**: No heavy frameworks or libraries
- **Efficient Caching**: Browser caching for static assets
- **Compression**: Gzip compression for all assets
- **CDN Ready**: Optimized for Cloudflare CDN deployment

### Performance Metrics

- **Load Time**: < 1 second for initial page load
- **Bundle Size**: < 100KB total (excluding images)
- **Accessibility**: WCAG 2.1 AA compliant
- **SEO**: Optimized meta tags and structure

## 🧪 Testing

### Available Scripts

```bash
# Start development server
npm run dev

# Build for production
npm run build

# Run with Docker
docker-compose up portfolio-frontend
```

## 🔧 Development

### Adding New Pages

1. Create HTML file in `public/`
2. Add navigation links in header
3. Include chatbot integration if needed

### Modifying Styles

1. Update CSS variables in `public/styles/main.css`
2. Maintain dark theme consistency
3. Test responsive design across devices

### API Integration

1. Add new fetch calls in HTML files
2. Update error handling for new endpoints
3. Test with backend API

## 🔒 Security

- **Content Security Policy**: Configured in nginx
- **HTTPS Only**: All external links use HTTPS
- **XSS Protection**: Input sanitization and validation
- **CORS**: Properly configured for API communication

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test across different browsers
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support or questions:

- **LinkedIn**: [linkedin.com/in/bvlucena](https://www.linkedin.com/in/bvlucena)
- **GitHub**: [github.com/brunovlucena](https://github.com/brunovlucena)

---

**Built with ❤️ by Bruno Lucena** 