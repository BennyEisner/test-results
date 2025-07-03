-- Migration SQL to add the duration column
ALTER TABLE builds ADD COLUMN duration DOUBLE PRECISION;