-- Populate database with initial data including Bruno Site
-- Migration: 002_populate_data.sql

-- Clear existing data
DELETE FROM projects;
DELETE FROM experience;

-- Insert only Bruno Site and Knative Lambda projects
INSERT INTO projects (title, description, type, github_url, live_url, technologies, featured, "order") VALUES
(
    'Bruno Site',
    'Personal portfolio and homelab showcase website built with React, TypeScript, Go, and modern cloud-native technologies. Features real-time project updates, interactive chatbot, and comprehensive skill showcase.',
    'Portfolio Website',
    'https://github.com/brunovlucena/bruno-site',
    'https://www.youtube.com/watch?v=lkkGlVWvkLk',
    ARRAY['React', 'TypeScript', 'Go', 'PostgreSQL', 'Redis', 'Docker', 'Kubernetes', 'Nginx'],
    TRUE,
    1
),
(
    'Knative Lambda',
    'Serverless functions and cloud-native development platform using Knative for scalable, event-driven applications with Kubernetes.',
    'Serverless',
    'https://github.com/brunovlucena/knative-lambda',
    'https://www.youtube.com/watch?v=lkkGlVWvkLk',
    ARRAY['Knative', 'Kubernetes', 'Serverless', 'CloudEvents', 'Go'],
    TRUE,
    2
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
('Tempo', 'Observability', 4, '⏱️', 16),
('OpenTelemetry', 'Observability', 4, '👁️', 17),

-- Infrastructure
('Terraform', 'Infrastructure', 5, '🏗️', 18),
('Pulumi', 'Infrastructure', 4, '☁️', 19),
('Docker', 'Infrastructure', 5, '🐳', 20),
('Flux', 'Infrastructure', 4, '⚡', 21),
('Helm', 'Infrastructure', 4, '⚓', 22),

-- Programming Languages
('Go', 'Programming', 5, '🐹', 23),
('Python', 'Programming', 4, '🐍', 24),
('TypeScript', 'Programming', 4, '📘', 25),
('JavaScript', 'Programming', 4, '📗', 26),
('Bash', 'Programming', 4, '💻', 27),

-- Databases & Messaging
('PostgreSQL', 'Database', 5, '🐘', 28),
('Redis', 'Database', 4, '🔴', 29),
('RabbitMQ', 'Messaging', 4, '🐰', 30),
('MongoDB', 'Database', 3, '🍃', 31),

-- AI/ML
('Machine Learning', 'AI/ML', 4, '🤖', 32),
('TensorFlow', 'AI/ML', 4, '📊', 33),
('Natural Language Processing', 'AI/ML', 4, '💬', 34),
('Computer Vision', 'AI/ML', 3, '👁️', 35),

-- DevOps & SRE
('Site Reliability Engineering', 'DevOps', 5, '⚙️', 36),
('DevSecOps', 'DevOps', 5, '🔒', 37),
('CI/CD', 'DevOps', 5, '🔄', 38),
('GitOps', 'DevOps', 4, '📦', 39),
('Infrastructure as Code', 'DevOps', 5, '🏗️', 40),

-- Monitoring & Alerting
('Monitoring', 'Monitoring', 5, '📊', 41),
('Alerting', 'Monitoring', 5, '🚨', 42),
('Logging', 'Monitoring', 5, '📝', 43),
('Tracing', 'Monitoring', 4, '🔍', 44),
('Metrics', 'Monitoring', 5, '📈', 45),

-- Cloud Platforms
('AWS', 'Cloud', 5, '☁️', 46),
('Google Cloud Platform', 'Cloud', 4, '☁️', 47),
('Azure', 'Cloud', 3, '☁️', 48),
('Multi-cloud', 'Cloud', 4, '☁️', 49),

-- Networking & Security
('Network Security', 'Security', 4, '🛡️', 50),
('Load Balancing', 'Networking', 4, '⚖️', 51),
('API Gateway', 'Networking', 4, '🚪', 52),
('Service Mesh', 'Networking', 4, '🕸️', 53),
('VPN', 'Security', 4, '🔐', 54),

-- Tools & Platforms
('GitHub', 'Tools', 5, '🐙', 55),
('GitLab', 'Tools', 4, '🦊', 56),
('Jenkins', 'Tools', 4, '🤖', 57),
('ArgoCD', 'Tools', 4, '🚀', 58),
('Knative', 'Platforms', 4, '☸️', 59),
('Serverless', 'Platforms', 4, '⚡', 60);

-- Insert experience data in chronological order (oldest to newest)
INSERT INTO experience (title, company, start_date, end_date, current, description, technologies, "order", active) VALUES
(
    'Operations Engineer',
    'Crealytics',
    '2017-08-01',
    '2018-03-31',
    FALSE,
    'Key Responsibilities:

- Cloud Operations: Managed and maintained complex cloud infrastructure on AWS and GCP.
Automation: Implemented automation tools (Saltstack) to streamline operations and reduce manual effort.

- Monitoring and Logging: Deployed and configured monitoring and logging solutions (Prometheus, ELK) to ensure system health and performance.

- Distributed Systems: Worked with distributed systems technologies like Mesos, Consul, Kafka, and Linkerd to build scalable and resilient applications.',
    ARRAY['AWS', 'GCP', 'Saltstack', 'Prometheus', 'ELK', 'Mesos', 'Consul', 'Kafka', 'Linkerd', 'Distributed Systems', 'Cloud Operations', 'Automation', 'Monitoring', 'Logging']::TEXT[],
    1,
    TRUE
),
(
    'DevOps Engineer',
    'Lesara',
    '2018-04-01',
    '2018-12-31',
    FALSE,
    'Key Responsibilities:

- Cloud-Native Infrastructure: Designed and implemented a Kubernetes cluster on bare-metal to modernize the infrastructure.

- Automation and CI/CD: Automated infrastructure provisioning and configuration management using Saltstack and Chef.

- Monitoring and Logging: Deployed and configured monitoring and logging solutions (Prometheus, ELK) to gain visibility into system health and performance.

- Collaboration: Worked closely with development teams to improve deployment processes and reduce downtime.',
    ARRAY['Kubernetes', 'Bare-metal', 'Saltstack', 'Chef', 'Prometheus', 'ELK', 'Automation', 'CI/CD', 'Monitoring', 'Logging', 'Infrastructure', 'Collaboration']::TEXT[],
    2,
    TRUE
),
(
    'Cloud Consultant',
    'Namecheap, Inc',
    '2019-03-01',
    '2019-08-31',
    FALSE,
    'Key Responsibilities:

- Cloud Migration and Modernization: Led the migration of legacy infrastructure from VMware ESXi to a Kubernetes-based platform on OpenStack.

- Infrastructure as Code: Implemented infrastructure as code practices using Terraform to automate provisioning and configuration management.
 
- Automation and CI/CD: Developed and maintained automation scripts (Bash, Golang, Ansible, Helm) to streamline operations and improve efficiency.',
    ARRAY['Cloud Migration', 'VMware ESXi', 'Kubernetes', 'OpenStack', 'Terraform', 'Bash', 'Golang', 'Ansible', 'Helm', 'Infrastructure as Code', 'Automation', 'CI/CD']::TEXT[],
    3,
    TRUE
),
(
    'Senior Infrastructure Engineer',
    'Mobimeo',
    '2020-02-01',
    '2023-03-31',
    FALSE,
    'Key Responsibilities:

- Cloud Native Infrastructure: Designed, implemented, and maintained a robust cloud-native infrastructure on AWS, leveraging services like EKS, Kops, and Kubernetes.

- Automation and CI/CD: Automated infrastructure provisioning, deployment, and configuration management using Terraform and GitHub Actions/GitLab CI/CD.

- Observability: Implemented and optimized monitoring, logging, and tracing solutions (Prometheus, Loki, Grafana, Thanos, EFK) to gain deep insights into system performance and behavior.

- Problem-Solving and Troubleshooting: Quickly identified and resolved complex infrastructure issues, minimizing downtime and service disruptions.',
    ARRAY['AWS', 'EKS', 'Kops', 'Kubernetes', 'Terraform', 'GitHub Actions', 'GitLab CI/CD', 'Prometheus', 'Loki', 'Grafana', 'Thanos', 'EFK', 'Infrastructure', 'Automation', 'CI/CD', 'Observability', 'Troubleshooting']::TEXT[],
    4,
    TRUE
),
(
    'SRE Chapter Lead',
    'Mobimeo',
    '2021-12-01',
    '2023-03-31',
    FALSE,
    'The SRE chapter lead is the line manager for the chapter members, responsible for developing people and the things happening in the SRE chapter but still is a member of the infrastructure & Operations Team and does day-to-day work.',
    ARRAY['SRE', 'Team Leadership', 'People Management', 'Infrastructure', 'Operations']::TEXT[],
    5,
    TRUE
),
(
    'SRE/DevOps',
    'Notifi',
    '2023-06-01',
    NULL,
    TRUE,
    'Key Responsibilities:

- Cloud Native Infrastructure: Architect, build, and maintain highly available, scalable, and resilient cloud-native infrastructure using Kubernetes, AWS, GCP, Pulumi, and many others

- Observability: Implement and optimize monitoring, logging, and tracing solutions (Prometheus, Loki, Tempo, Grafana, OpenTelemetry) to gain deep insights into system performance and behavior.

- Chatbot for SRE: RAG, Vertex AI

- Automation and CI/CD: Automate infrastructure provisioning, deployment, and configuration management using Terraform, Atmos, and GitHub Actions to accelerate development and reduce errors.

- Serverless and Function-as-a-Service: Develop and deploy serverless applications on AWS Lambda 

- Serverless on K8s: Knative (CloudEvents, RabbitMQ), Golang 

- Security and Compliance: Ensure the security and compliance of systems and applications by implementing best practices and leveraging security tools.',
    ARRAY['Kubernetes', 'AWS', 'GCP', 'Pulumi', 'Prometheus', 'Loki', 'Tempo', 'Grafana', 'OpenTelemetry', 'RAG', 'Vertex AI', 'Terraform', 'Atmos', 'GitHub Actions', 'AWS Lambda', 'Knative', 'CloudEvents', 'RabbitMQ', 'Golang', 'Security', 'Compliance']::TEXT[],
    6,
    TRUE
);

-- Insert content data
INSERT INTO content (key, value) VALUES
(
    'about',
    '{"description": "Senior Cloud Native Infrastructure Engineer with extensive experience in designing, implementing, and maintaining scalable, resilient cloud-native infrastructure. Passionate about automation, observability, and modern DevOps practices."}'
),
(
    'contact',
    '{"email": "bruno@lucena.cloud", "location": "Brazil", "linkedin": "https://www.linkedin.com/in/bvlucena", "github": "https://github.com/brunovlucena", "availability": "Open to new opportunities"}'
);

-- Verify all data
SELECT 'Projects' as table_name, COUNT(*) as count FROM projects
UNION ALL
SELECT 'Skills' as table_name, COUNT(*) as count FROM skills
UNION ALL
SELECT 'Experience' as table_name, COUNT(*) as count FROM experience
UNION ALL
SELECT 'Content' as table_name, COUNT(*) as count FROM content; 