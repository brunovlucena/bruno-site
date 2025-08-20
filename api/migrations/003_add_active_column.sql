-- Add active column to experience table
-- Migration: 003_add_active_column.sql

-- Add active column to experience table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'experience' 
        AND column_name = 'active'
    ) THEN
        ALTER TABLE experience ADD COLUMN active BOOLEAN DEFAULT TRUE;
    END IF;
END $$;

-- Add active column to projects table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'projects' 
        AND column_name = 'active'
    ) THEN
        ALTER TABLE projects ADD COLUMN active BOOLEAN DEFAULT TRUE;
    END IF;
END $$;

-- Add active column to skills table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'skills' 
        AND column_name = 'active'
    ) THEN
        ALTER TABLE skills ADD COLUMN active BOOLEAN DEFAULT TRUE;
    END IF;
END $$;

-- Create index for active column on experience table
CREATE INDEX IF NOT EXISTS idx_experience_active ON experience(active);

-- Create index for active column on projects table
CREATE INDEX IF NOT EXISTS idx_projects_active ON projects(active);

-- Create index for active column on skills table
CREATE INDEX IF NOT EXISTS idx_skills_active ON skills(active);
