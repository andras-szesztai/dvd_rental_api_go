-- First, update any NULL address_id values with a default address
-- Replace 1 with a valid address_id from your address table
UPDATE customer SET address_id = 1 WHERE address_id IS NULL;

-- Then set the NOT NULL constraint
ALTER TABLE customer ALTER COLUMN address_id SET NOT NULL;