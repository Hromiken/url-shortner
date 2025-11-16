-- Таблица сокращённых ссылок
CREATE TABLE IF NOT EXISTS public.short_urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    alias VARCHAR(64) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- Статистика переходов
CREATE TABLE IF NOT EXISTS click_logs (
    id SERIAL PRIMARY KEY,
    url_id INT NOT NULL REFERENCES short_urls(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address INET,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для ускорения аналитики по alias
CREATE INDEX IF NOT EXISTS idx_short_urls_alias ON short_urls(alias);

-- Индекс для аналитики по short_url_id
CREATE INDEX IF NOT EXISTS idx_analytics_short_url_id ON click_logs(url_id);
