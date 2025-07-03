-- Migration to add test_case_count column to builds table
-- Run this against your existing database before reseeding

ALTER TABLE builds ADD COLUMN test_case_count INTEGER;

-- Optionally, you can set a default value for existing records
-- UPDATE builds SET test_case_count = 0 WHERE test_case_count IS NULL;
