package entity

import "time"

// URL сущность
type URL struct {
	ID          int64     `db:"id"`
	OriginalURL string    `db:"original_url"`
	Alias       string    `db:"alias"`
	CreatedAt   time.Time `db:"created_at"`
}

// ClickLogs сущность аналитической записи
type ClickLogs struct {
	ID        int64     `db:"id"`     // ID аналитической записи
	URLID     int64     `db:"url_id"` // id of URL struct
	UserAgent string    `db:"user_agent"`
	IPAddress string    `db:"ip_address"`
	ClickedAt time.Time `db:"clicked_at"`
}

// StatDay - кол-во переходов за день.
type StatDay struct {
	Day   time.Time `json:"day"`
	Count int64     `json:"count"`
}

// StatMonth - кол-во переходов за месяц.
type StatMonth struct {
	Month time.Time `json:"month"`
	Count int64     `json:"count"`
}

// StatAgent - кол-во переходов у конкретного user agent.
type StatAgent struct {
	UserAgent string `json:"user_agent"`
	Count     int64  `json:"count"`
}
