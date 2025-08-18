-- Migration: Add short_description column to projects table
-- This allows projects to have a brief description for the header dropdown

-- Add short_description column with default value
ALTER TABLE projects ADD COLUMN short_description VARCHAR(200);

-- Update existing projects with short descriptions
UPDATE projects SET short_description = 
  CASE 
    WHEN title = 'Monitoring Platform' THEN 'Prometheus, Grafana, Loki, OpenTelemetry & Tempo'
    WHEN title = 'Knative Lambda' THEN 'Serverless functions and cloud-native development'
    WHEN title = 'Doctor Chatbot' THEN 'AI-powered medical assistance and health guidance'
    WHEN title = 'SRE Agent on K8s' THEN 'Intelligent SRE agent for automated monitoring'
    WHEN title = 'DJ Double' THEN 'Advanced music mixing and DJ application'
    WHEN title = 'Analista Financeiro' THEN 'Financial analysis platform with automated reporting'
    ELSE LEFT(description, 100)
  END;

-- Make the column NOT NULL after populating data
ALTER TABLE projects ALTER COLUMN short_description SET NOT NULL;

-- Add index for better query performance
CREATE INDEX idx_projects_short_description ON projects(short_description);

-- Add comment to document the column purpose
COMMENT ON COLUMN projects.short_description IS 'Brief description for header display (max 200 characters)';
