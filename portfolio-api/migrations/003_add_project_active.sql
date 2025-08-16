-- Migration: Add active column to projects table
-- This allows controlling which projects are visible on the portfolio

-- Add active column with default value true (all existing projects remain active)
ALTER TABLE projects ADD COLUMN active BOOLEAN DEFAULT true NOT NULL;

-- Add index for better query performance when filtering by active status
CREATE INDEX idx_projects_active ON projects(active);

-- Add comment to document the column purpose
COMMENT ON COLUMN projects.active IS 'Controls whether the project is visible on the portfolio (true = visible, false = hidden)'; 