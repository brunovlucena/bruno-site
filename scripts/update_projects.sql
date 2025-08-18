-- Update projects in Bruno site database
-- This script will:
-- 1. Delete "DevOps CLI" project
-- 2. Deactivate all remaining projects except "Knative Lambda" and "Monitoring Platform"

-- First, let's see the current projects
SELECT id, title, active FROM projects ORDER BY id;

-- Delete "DevOps CLI" project
DELETE FROM projects WHERE title = 'DevOps CLI';

-- Deactivate all remaining projects except "Knative Lambda" and "Monitoring Platform"
UPDATE projects 
SET active = false 
WHERE title NOT IN ('Knative Lambda', 'Monitoring Platform');

-- Verify the changes
SELECT id, title, active FROM projects ORDER BY id;
