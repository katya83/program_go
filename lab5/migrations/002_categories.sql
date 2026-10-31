CREATE TABLE IF NOT EXISTS categories (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);
 
-- Базовые категории.
INSERT INTO categories (name) VALUES
    ('Еда'),
    ('Транспорт'),
    ('Жильё'),
    ('Развлечения'),
    ('Здоровье'),
    ('Другое');
 
--  Новое поле в expenses: внешний ключ на categories.
ALTER TABLE expenses ADD COLUMN category_id INTEGER REFERENCES categories(id);
 
-- Уже существующим тратам ставим категорию "Другое".
UPDATE expenses SET category_id = (SELECT id FROM categories WHERE name = 'Другое');
