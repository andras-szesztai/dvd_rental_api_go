-- First remove the unique constraint
ALTER TABLE customer DROP CONSTRAINT customer_email_unique;

-- Restore original email addresses by removing the appended IDs
UPDATE customer
SET email = regexp_replace(email, '_\d+$', '')
WHERE email ~ '_\d+$';