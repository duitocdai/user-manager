-- Drop trigger
DROP TRIGGER IF EXISTS set_user_updated_at ON users;
-- Drop trigger function
DROP FUNCTION   IF EXISTS update_user_updated_at_column ON users;
-- Drop table
DROP TABLE IF EXISTS users;