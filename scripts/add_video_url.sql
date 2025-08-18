-- Migration: Add video_url column to projects table
-- Purpose: allow embedding YouTube videos in Homelab cards

-- Add nullable video_url column
ALTER TABLE projects ADD COLUMN IF NOT EXISTS video_url VARCHAR(500);

-- Seed example video for active projects (adjust as needed)
UPDATE projects
SET video_url = CASE
  WHEN title = 'Knative Lambda' THEN 'https://www.youtube.com/embed/ZToicYcHIOU'
  WHEN title = 'Monitoring Platform' THEN 'https://www.youtube.com/embed/ZToicYcHIOU'
  ELSE video_url
END;
