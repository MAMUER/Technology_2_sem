-- Создание таблицы для задач
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date DATE,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    subject VARCHAR(100) NOT NULL
);

-- Индекс для поиска по заголовку
CREATE INDEX IF NOT EXISTS idx_tasks_title ON tasks(title);

-- Индекс для фильтрации по пользователю
CREATE INDEX IF NOT EXISTS idx_tasks_subject ON tasks(subject);

-- Создание таблицы для пользователей
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(50) PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Добавление тестовых пользователей
INSERT INTO users (username, password_hash) VALUES
    ('student', 'student'),
    ('admin', 'admin123')
ON CONFLICT (username) DO NOTHING;