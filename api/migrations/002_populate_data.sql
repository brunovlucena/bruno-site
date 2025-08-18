-- Populate database with initial data from HTML files
-- Migration: 002_populate_data.sql

-- Insert projects from index.html
INSERT INTO projects (title, description, type, modules, github_url, technologies, featured, "order") VALUES
(
    'DevOps CLI',
    'A powerful command-line interface for DevSecOps/SRE automation, built with TypeScript and featuring AI-powered code generation and project management tools.',
    'CLI Tool',
    4,
    'https://github.com/brunovlucena/bruno-cli',
    ARRAY['TypeScript', 'Node.js', 'AI/ML', 'CLI', 'DevSecOps'],
    TRUE,
    1
),
(
    'Monitoring Platform',
    'Complete observability solution with Prometheus, Grafana, Loki, OpenTelemetry & Tempo for comprehensive application monitoring and tracing.',
    'Observability Stack',
    5,
    'https://github.com/brunovlucena/observability-stack',
    ARRAY['Prometheus', 'Grafana', 'Loki', 'OpenTelemetry', 'Tempo', 'Kubernetes'],
    TRUE,
    2
),
(
    'Knative Lambda',
    'Serverless functions and cloud-native development platform using Knative for scalable, event-driven applications with Kubernetes.',
    'Serverless',
    3,
    'https://github.com/brunovlucena/knative-lambda',
    ARRAY['Knative', 'Kubernetes', 'Serverless', 'CloudEvents', 'Go'],
    FALSE,
    3
),
(
    'Doctor Chatbot',
    'AI-powered medical assistance and health guidance system built with advanced natural language processing and medical knowledge integration.',
    'AI Application',
    3,
    'https://github.com/brunovlucena/doctor-companion',
    ARRAY['AI/ML', 'NLP', 'Healthcare', 'Python', 'TensorFlow'],
    TRUE,
    4
),
(
    'SRE Agent on K8s',
    'Intelligent SRE agent deployed on Kubernetes for automated monitoring, incident response, and infrastructure optimization using AI/ML capabilities.',
    'Kubernetes Agent',
    5,
    'https://github.com/brunovlucena/sre-agent-k8s',
    ARRAY['Kubernetes', 'AI/ML', 'SRE', 'Monitoring', 'Automation'],
    TRUE,
    5
),
(
    'DJ Double',
    'Advanced music mixing and DJ application with real-time audio processing, beat matching, and professional sound engineering capabilities.',
    'Music Application',
    4,
    'https://github.com/brunovlucena/dj-double',
    ARRAY['Audio Processing', 'Real-time', 'Music', 'C++', 'DSP'],
    FALSE,
    6
),
(
    'Analista Financeiro',
    'Comprehensive financial analysis platform with automated reporting, risk assessment, and data visualization for investment decision-making.',
    'Financial Analysis',
    6,
    'https://github.com/brunovlucena/financial-analyst',
    ARRAY['Financial Analysis', 'Data Visualization', 'Risk Assessment', 'Python', 'Pandas'],
    FALSE,
    7
);

-- Insert skills from about section
INSERT INTO skills (name, category, proficiency, icon, "order") VALUES
-- IT Security
('IT Security', 'Security', 5, '🔒', 1),
('Vulnerability Assessment', 'Security', 5, '🔍', 2),
('Nessus', 'Security', 4, '🛡️', 3),
('Security Auditing', 'Security', 4, '📋', 4),

-- Project Management
('Project Management', 'Management', 4, '📊', 5),
('Team Leadership', 'Management', 4, '👥', 6),
('Agile/Scrum', 'Management', 4, '🔄', 7),

-- Kubernetes & Cloud
('Kubernetes', 'Cloud', 5, '☸️', 8),
('AWS EKS', 'Cloud', 5, '☁️', 9),
('GCP', 'Cloud', 4, '☁️', 10),
('AWS Lambda', 'Cloud', 4, '⚡', 11),
('OpenStack', 'Cloud', 3, '☁️', 12),

-- Observability
('Prometheus', 'Observability', 5, '📊', 13),
('Grafana', 'Observability', 5, '📈', 14),
('Loki', 'Observability', 4, '📝', 15),
('OpenTelemetry', 'Observability', 4, '🔍', 16),
('Thanos', 'Observability', 3, '📊', 17),

-- AI/ML
('AI/LLMOps', 'AI', 4, '🤖', 18),
('Vertex AI', 'AI', 4, '🧠', 19),
('RAG', 'AI', 4, '🔍', 20),
('Machine Learning', 'AI', 3, '🤖', 21),

-- Automation & CI/CD
('Terraform', 'Automation', 5, '🏗️', 22),
('Pulumi', 'Automation', 4, '⚙️', 23),
('GitHub Actions', 'Automation', 5, '🔄', 24),
('GitLab CI', 'Automation', 4, '🔄', 25),
('Atmos', 'Automation', 4, '🌪️', 26),
('Ansible', 'Automation', 3, '🤖', 27),
('Saltstack', 'Automation', 3, '🧂', 28),
('Helm', 'Automation', 4, '⛵', 29),

-- Programming
('Golang', 'Programming', 4, '🐹', 30),
('Python', 'Programming', 4, '🐍', 31),
('Bash', 'Programming', 5, '💻', 32),
('Ruby', 'Programming', 3, '💎', 33),

-- Distributed Systems
('RabbitMQ', 'Distributed', 4, '🐰', 34),
('Kafka', 'Distributed', 3, '📨', 35),
('Consul', 'Distributed', 3, '🏛️', 36),
('CloudEvents', 'Distributed', 4, '☁️', 37);

-- Insert experience from resume.html
INSERT INTO experience (title, company, start_date, end_date, current, description, technologies, "order") VALUES
(
    'SRE/DevOps Engineer',
    'Notifi',
    '2023-06-01',
    NULL,
    TRUE,
    'Architect and maintain highly available, scalable cloud-native infrastructure using Kubernetes, AWS, GCP, and Pulumi. Implement comprehensive observability solutions with Prometheus, Loki, Tempo, Grafana, and OpenTelemetry. Develop RAG-based chatbot for SRE using Vertex AI and advanced AI/ML technologies. Automate infrastructure provisioning and deployment using Terraform, Atmos, and GitHub Actions. Build serverless applications on AWS Lambda and Knative with CloudEvents and RabbitMQ. Lead platform engineering initiatives and mentor junior engineers.',
    ARRAY['Kubernetes', 'AWS', 'GCP', 'Pulumi', 'Prometheus', 'Loki', 'Tempo', 'Grafana', 'OpenTelemetry', 'Vertex AI', 'RAG', 'Terraform', 'Atmos', 'GitHub Actions', 'AWS Lambda', 'Knative', 'CloudEvents', 'RabbitMQ', 'Platform Engineering'],
    1
),
(
    'SRE Chapter Lead & Senior Infrastructure Engineer',
    'Mobimeo',
    '2020-02-01',
    '2023-03-31',
    FALSE,
    'Led SRE chapter as line manager, developing team members and driving infrastructure strategy. Designed and maintained robust cloud-native infrastructure on AWS using EKS, Kops, and Kubernetes. Implemented monitoring and logging solutions with Prometheus, Loki, Grafana, Thanos, and EFK stack. Automated infrastructure using Terraform and CI/CD pipelines with GitHub Actions/GitLab CI. Resolved complex infrastructure issues, minimizing downtime and service disruptions. Established SRE best practices and SLI/SLO frameworks.',
    ARRAY['SRE', 'Leadership', 'AWS EKS', 'Kops', 'Kubernetes', 'Prometheus', 'Loki', 'Grafana', 'Thanos', 'EFK', 'Terraform', 'CI/CD', 'GitHub Actions', 'GitLab CI', 'SLI/SLO'],
    2
),
(
    'Cloud Consultant',
    'Namecheap, Inc',
    '2019-03-01',
    '2019-08-31',
    FALSE,
    'Led migration from VMware ESXi to Kubernetes-based platform on OpenStack. Implemented infrastructure as code practices using Terraform for automation. Developed automation scripts using Bash, Golang, Ansible, and Helm. Designed scalable microservices architecture and implemented CI/CD pipelines. Reduced infrastructure costs by 40% through containerization and automation.',
    ARRAY['VMware', 'Kubernetes', 'OpenStack', 'Terraform', 'Bash', 'Golang', 'Ansible', 'Helm', 'Microservices', 'CI/CD', 'Containerization'],
    3
),
(
    'DevOps Engineer',
    'Lesara',
    '2018-04-01',
    '2018-12-31',
    FALSE,
    'Designed and implemented Kubernetes cluster on bare-metal infrastructure. Automated infrastructure provisioning using Saltstack and Chef. Deployed monitoring and logging solutions with Prometheus and ELK stack. Implemented blue-green deployments and automated rollback mechanisms. Achieved 99.5% uptime through proactive monitoring and incident response.',
    ARRAY['Kubernetes', 'Bare-metal', 'Saltstack', 'Chef', 'Prometheus', 'ELK', 'Blue-Green Deployments', 'Incident Response'],
    4
),
(
    'Operations Engineer',
    'Crealytics',
    '2017-08-01',
    '2018-03-31',
    FALSE,
    'Managed complex cloud infrastructure on AWS and GCP. Implemented automation tools using Saltstack to streamline operations. Deployed monitoring and logging solutions with Prometheus and ELK. Worked with distributed systems including Mesos, Consul, Kafka, and Linkerd. Optimized system performance and reduced response times by 60%.',
    ARRAY['AWS', 'GCP', 'Saltstack', 'Prometheus', 'ELK', 'Mesos', 'Consul', 'Kafka', 'Linkerd', 'Performance Optimization'],
    5
),
(
    'IT Security Analyst',
    'Tempest Security Intelligence',
    '2011-01-01',
    '2013-10-31',
    FALSE,
    'Conducted in-depth vulnerability assessments and security risk analysis. Researched latest security threats and vulnerabilities for impact assessment. Developed automated tools and scripts using Bash and Ruby for vulnerability scanning. Created and customized Nessus Scanner Plugins (NASL) for enhanced detection. Established security frameworks and compliance procedures.',
    ARRAY['Vulnerability Assessment', 'Security Analysis', 'Bash', 'Ruby', 'Nessus', 'NASL', 'Security Frameworks', 'Compliance'],
    6
);

-- Insert content from HTML files
INSERT INTO content (key, value) VALUES
(
    'about',
    '{
        "title": "About Me",
        "description": "Senior Cloud Native Infrastructure Engineer with 12+ years of experience designing, building, and scaling mission-critical cloud-native platforms. Expert in Kubernetes ecosystem, multi-cloud architectures (AWS/GCP), and modern observability stacks. Proven track record in Site Reliability Engineering (SRE), DevSecOps practices, and leading high-performing infrastructure teams. Specialized in AI/ML infrastructure, LLMOps, and building resilient systems that handle millions of requests. Passionate about automation, security-first approaches, and driving innovation in cloud-native technologies. Delivers enterprise-grade solutions with 99.9%+ uptime, optimized performance, and comprehensive security postures.",
        "highlights": [
            {"icon": "☸️", "text": "Kubernetes Expert"},
            {"icon": "☁️", "text": "Multi-Cloud Architect"},
            {"icon": "📊", "text": "Observability & SRE"},
            {"icon": "🤖", "text": "AI/ML Infrastructure"},
            {"icon": "🔒", "text": "DevSecOps"},
            {"icon": "🚀", "text": "Platform Engineering"},
            {"icon": "⚡", "text": "High-Performance Systems"},
            {"icon": "🛡️", "text": "Security-First"}
        ]
    }'::jsonb
),
(
    'contact',
    '{
        "email": "bruno.lucena@example.com",
        "location": "Brazil",
        "linkedin": "https://www.linkedin.com/in/bvlucena",
        "github": "https://github.com/brunovlucena",
        "availability": "Open to new opportunities in SRE, DevSecOps, and AI Engineering roles."
    }'::jsonb
),
(
    'hero',
    '{
        "title": "Senior Cloud Native Infrastructure Engineer",
        "subtitle": "SRE • DevSecOps • AI/ML Infrastructure • Platform Engineering"
    }'::jsonb
); 