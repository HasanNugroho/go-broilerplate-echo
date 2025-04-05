-- Menghapus tabel user_roles 
DROP TABLE IF EXISTS user_roles;

-- Menghapus trigger dan function yang terkait dengan users
DROP TRIGGER IF EXISTS trigger_update_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Menghapus tabel users
DROP TABLE IF EXISTS users;

-- Menghapus tabel roles setelah tabel user_roles
DROP TABLE IF EXISTS roles;
