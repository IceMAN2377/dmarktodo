CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       status TEXT CHECK (status IN ('active', 'done')) DEFAULT 'active',
                       priority TEXT CHECK (priority IN ('low', 'medium', 'high')) DEFAULT 'medium',
                       completed_at TIMESTAMP NULL,
                       created_at TIMESTAMP DEFAULT now(),
                       updated_at TIMESTAMP DEFAULT now()
);
