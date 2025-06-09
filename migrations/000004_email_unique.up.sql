-- First, let's identify and handle duplicate emails
WITH duplicates AS (
    SELECT email, COUNT(*) as count
    FROM customer
    GROUP BY email
    HAVING COUNT(*) > 1
)
UPDATE customer c
SET email = email || '_' || user_id
WHERE email IN (SELECT email FROM duplicates);

-- Now we can safely add the unique constraint
ALTER TABLE customer ADD CONSTRAINT customer_email_unique UNIQUE (email); 