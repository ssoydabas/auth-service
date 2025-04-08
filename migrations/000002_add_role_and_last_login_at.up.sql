ALTER TABLE accounts
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'common' CHECK (role IN ('common', 'admin', 'manager', 'teacher', 'student')),
ADD COLUMN last_login_at TIMESTAMP WITH TIME ZONE; 