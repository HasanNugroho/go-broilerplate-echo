DROP TABLE users;

CREATE TABLE "users" (
    "id" UUID PRIMARY KEY,
    "email" TEXT,
    "name" TEXT,
    "password" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) engine = innodb;

CREATE INDEX idx_users_id ON users (id);