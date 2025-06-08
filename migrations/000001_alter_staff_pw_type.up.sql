-- Add new bytea column
ALTER TABLE staff ADD COLUMN password_new BYTEA;

-- Convert existing passwords to bytea using decode
UPDATE staff 
SET password_new = decode(password, 'hex')
WHERE password IS NOT NULL;

-- Drop the old column
ALTER TABLE staff DROP COLUMN password;

-- Rename the new column to the original name
ALTER TABLE staff RENAME COLUMN password_new TO password; 