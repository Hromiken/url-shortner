package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
	"shortner/internal/entity"
	"shortner/internal/handler/dto"
	"time"

	"github.com/wb-go/wbf/redis"
)

var (
	ErrAliasIncorrect = errors.New("alias must be alphanumeric, -, _")
)

// URLRepo - описывает методы URLRepo.
type URLRepo interface {
	SaveURL(ctx context.Context, url entity.URL) (int64, error)
	GetURL(ctx context.Context, url entity.URL) (*entity.URL, error)
	SaveUserClick(ctx context.Context, click entity.ClickLogs) error
	GetAnalytics(ctx context.Context, url entity.URL) ([]entity.ClickLogs, error)
	StatisticDay(ctx context.Context, url entity.URL) ([]entity.StatDay, error)
	StatisticMonth(ctx context.Context, url entity.URL) ([]entity.StatMonth, error)
	StatisticAgent(ctx context.Context, url entity.URL) ([]entity.StatAgent, error)
	CountClicks(ctx context.Context, url entity.URL) (int64, error)
}

// ShortenerService - shorter service.
type ShortenerService struct {
	URLRepo URLRepo
	Cache   *redis.Client
}

// NewShorterService - return instance of shorter service.
func NewShorterService(repo URLRepo) *ShortenerService {
	return &ShortenerService{URLRepo: repo}
}

// CreateShortURL - вызов SaveURL
func (s *ShortenerService) CreateShortURL(ctx context.Context, request dto.CreateURLShorterRequest) (string, error) {
	if request.URL == "" {
		return "", errors.New("url is empty")
	}

	alias := request.CustomAlias
	if alias == "" {
		alias = generateAlias(8)
	}

	if !isValidAlias(alias) {
		return "", ErrAliasIncorrect
	}

	_, err := s.URLRepo.SaveURL(ctx, entity.URL{Alias: alias, OriginalURL: request.URL})
	if err != nil {
		return "", err
	}
	return alias, nil
}

// GetURL - получить оригинальный URL по alias.
func (s *ShortenerService) GetURL(ctx context.Context, url entity.URL) (*entity.URL, error) {
	cacheKey := "short:" + url.Alias

	if s.Cache != nil {
		cached, err := s.Cache.Get(ctx, cacheKey)
		if err == nil {
			return &entity.URL{
				Alias:       url.Alias,
				OriginalURL: cached,
			}, nil
		}
	}

	return s.URLRepo.GetURL(ctx, url)
}

// GetAnalytics - Получить аналитику по ссылке.
func (s *ShortenerService) GetAnalytics(ctx context.Context, url entity.URL) ([]entity.ClickLogs, error) {
	return s.URLRepo.GetAnalytics(ctx, url)
}

// SaveUserClick - Сохраняем данные о переходе.
func (s *ShortenerService) SaveUserClick(ctx context.Context, click entity.ClickLogs) error {
	return s.URLRepo.SaveUserClick(ctx, click)
}

// CountClicks - Подсчет кол-ва переходов.
func (s *ShortenerService) CountClicks(ctx context.Context, url entity.URL) (int64, error) {
	return s.URLRepo.CountClicks(ctx, url)
}

func (s *ShortenerService) StatisticDay(ctx context.Context, url entity.URL) ([]entity.StatDay, error) {
	return s.URLRepo.StatisticDay(ctx, url)
}

func (s *ShortenerService) StatisticMonth(ctx context.Context, url entity.URL) ([]entity.StatMonth, error) {
	return s.URLRepo.StatisticMonth(ctx, url)
}

func (s *ShortenerService) StatisticAgent(ctx context.Context, url entity.URL) ([]entity.StatAgent, error) {
	return s.URLRepo.StatisticAgent(ctx, url)
}

func generateAlias(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return time.Now().Format("150405")
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func isValidAlias(a string) bool {
	var aliasRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	return aliasRe.MatchString(a)
}
