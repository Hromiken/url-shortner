package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/wb-go/wbf/dbpg"

	"shortner/internal/entity"
)

var (
	ErrUrlNotFound   = errors.New("url not found")
	ErrAliasExists   = errors.New("alias already exists")
	ErrAliasTooLarge = errors.New("alias too large")
)

type URLRepo struct {
	db *dbpg.DB
}

func NewURLRepo(db *dbpg.DB) *URLRepo {
	return &URLRepo{db: db}
}

// SaveURL - запись alias и получение его id.
func (p *URLRepo) SaveURL(ctx context.Context, url entity.URL) (int64, error) {
	if len(url.Alias) > 64 {
		return 0, fmt.Errorf("%s: alias too long (max 64 characters)", ErrAliasTooLarge)
	}

	query := `INSERT INTO public.short_urls(original_url, alias) VALUES($1, $2) RETURNING id`
	row := p.db.Master.QueryRowContext(ctx, query, url.OriginalURL, url.Alias)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf(" %w", ErrAliasExists)
			}
		}
		return 0, fmt.Errorf("%w", err)
	}

	return id, nil
}

// GetURL - получение URL entity по его alias.
func (p *URLRepo) GetURL(ctx context.Context, url entity.URL) (*entity.URL, error) {
	query := `SELECT id, original_url, alias, created_at FROM short_urls WHERE alias = $1`
	row := p.db.QueryRowContext(ctx, query, url.Alias)

	var resultURL entity.URL
	err := row.Scan(&resultURL.ID, &resultURL.OriginalURL, &resultURL.Alias, &resultURL.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUrlNotFound
		}
		return nil, err
	}

	return &resultURL, nil
}

// SaveUserClick - запись переходов, ip user, user agent.
func (p *URLRepo) SaveUserClick(ctx context.Context, click entity.ClickLogs) error {
	query := `INSERT INTO click_logs (url_id, user_agent, ip_address)VALUES ($1, $2, $3)`

	_, err := p.db.Master.ExecContext(ctx, query, click.URLID, click.UserAgent, click.IPAddress)
	return err
}

// GetAnalytics - получение аналитики по alias
func (p *URLRepo) GetAnalytics(ctx context.Context, url entity.URL) ([]entity.ClickLogs, error) {
	query := `
SELECT a.id, a.url_id, a.user_agent, a.ip_address, a.clicked_at
FROM click_logs a
JOIN short_urls s ON s.id = a.url_id
WHERE s.alias = $1
ORDER BY a.clicked_at DESC`

	rows, err := p.db.QueryContext(ctx, query, url.Alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []entity.ClickLogs
	for rows.Next() {
		var r entity.ClickLogs
		err = rows.Scan(&r.ID, &r.URLID, &r.UserAgent, &r.IPAddress, &r.ClickedAt)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}

// CountClicks - подсчет кол-ва переходов.
func (p *URLRepo) CountClicks(ctx context.Context, url entity.URL) (int64, error) {
	query := `SELECT COUNT(*) FROM click_logs a JOIN short_urls s ON s.id = a.url_id WHERE s.alias = $1`
	row := p.db.QueryRowContext(ctx, query, url.Alias)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// StatisticDay - аналитика по дню.
func (p *URLRepo) StatisticDay(ctx context.Context, url entity.URL) ([]entity.StatDay, error) {
	query := `
		SELECT date_trunc('day', a.clicked_at) AS day,
		       COUNT(*) AS cnt
		FROM click_logs a
		JOIN short_urls s ON s.id = a.url_id
		WHERE s.alias = $1
		GROUP BY day
		ORDER BY day DESC
	`

	rows, err := p.db.QueryContext(ctx, query, url.Alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []entity.StatDay

	for rows.Next() {
		var s entity.StatDay
		err = rows.Scan(&s.Day, &s.Count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

// StatisticMonth - аналитика за месяц.
func (p *URLRepo) StatisticMonth(ctx context.Context, url entity.URL) ([]entity.StatMonth, error) {
	query := `
		SELECT date_trunc('month', a.clicked_at) AS month,
		       COUNT(*) AS cnt
		FROM click_logs a
		JOIN short_urls s ON s.id = a.url_id
		WHERE s.alias = $1
		GROUP BY month
		ORDER BY month DESC
	`

	rows, err := p.db.QueryContext(ctx, query, url.Alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []entity.StatMonth

	for rows.Next() {
		var s entity.StatMonth
		err = rows.Scan(&s.Month, &s.Count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

// StatisticAgent - аналитика по user agent.
func (p *URLRepo) StatisticAgent(ctx context.Context, url entity.URL) ([]entity.StatAgent, error) {
	query := `
		SELECT a.user_agent AS agent,
		       COUNT(*) AS cnt
		FROM click_logs a
		JOIN short_urls s ON s.id = a.url_id
		WHERE s.alias = $1
		GROUP BY user_agent
		ORDER BY user_agent DESC
	`

	rows, err := p.db.QueryContext(ctx, query, url.Alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []entity.StatAgent

	for rows.Next() {
		var s entity.StatAgent
		err = rows.Scan(&s.UserAgent, &s.Count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}
