CREATE TABLE users (
    id VARCHAR(26) PRIMARY KEY,
    email TEXT,
    name TEXT,
    password TEXT,
    isSystem BOOL DEFAULT false,
    isActive BOOL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Buat trigger untuk auto-update kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE roles (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    isSystem BOOL DEFAULT false,
    permissions JSON NOT NULL
);

CREATE TABLE user_roles (
    user_id VARCHAR(26) NOT NULL,
    role_id VARCHAR(26) NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
