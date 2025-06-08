-- Add new varchar column
ALTER TABLE staff ADD COLUMN password_old VARCHAR(40);

-- Convert bytea back to hex string
UPDATE staff 
SET password_old = encode(password, 'hex')
WHERE password IS NOT NULL;

-- Drop the bytea column
ALTER TABLE staff DROP COLUMN password;

-- Rename the varchar column to the original name
ALTER TABLE staff RENAME COLUMN password_old TO password; 