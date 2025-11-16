package handler

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"shortner/internal/entity"
	"shortner/internal/handler/dto"
	"shortner/internal/service"
	"strings"
)

type ShortenerService interface {
	CreateShortURL(ctx context.Context, request dto.CreateURLShorterRequest) (string, error)
	GetURL(ctx context.Context, url entity.URL) (*entity.URL, error)
	SaveUserClick(ctx context.Context, click entity.ClickLogs) error
	GetAnalytics(ctx context.Context, url entity.URL) ([]entity.ClickLogs, error)
	CountClicks(ctx context.Context, url entity.URL) (int64, error)
	StatisticDay(ctx context.Context, url entity.URL) ([]entity.StatDay, error)
	StatisticMonth(ctx context.Context, url entity.URL) ([]entity.StatMonth, error)
	StatisticAgent(ctx context.Context, url entity.URL) ([]entity.StatAgent, error)
}

// Handler - handler + validator.
type Handler struct {
	svc *service.ShortenerService
}

// NewHandler - return instance of handler.
func NewHandler(svc *service.ShortenerService) *Handler {
	return &Handler{
		svc: svc,
	}
}

// SaveURL - POST /shorten.
func (h *Handler) SaveURL(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateURLShorterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		writeError(w, "url is required", http.StatusBadRequest)
		return
	}
	if !isValidURL(req.URL) {
		writeError(w, "invalid url format", http.StatusBadRequest)
		return
	}

	if req.CustomAlias != "" && !isValidAlias(req.CustomAlias) {
		writeError(w, "invalid custom alias", http.StatusBadRequest)
		return
	}

	createdAlias, err := h.svc.CreateShortURL(r.Context(), req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	shortURL := "http://" + host + "/s/" + createdAlias

	resp := dto.ShortenResponse{
		Alias: createdAlias,
		Short: shortURL,
	}

	writeJSON(w, resp, http.StatusCreated)
}

// Redirect - переход по ссылке.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := strings.TrimPrefix(r.URL.Path, "/s/")
	alias = strings.TrimSpace(alias)
	if alias == "" {
		writeError(w, "alias is required", http.StatusBadRequest)
		return
	}

	urlEntity, err := h.svc.GetURL(r.Context(), entity.URL{Alias: alias})
	if err != nil || urlEntity == nil {
		writeError(w, "url not found", http.StatusNotFound)
		return
	}

	var c entity.ClickLogs
	c.UserAgent = r.UserAgent()
	c.IPAddress = getIP(r)
	c.URLID = urlEntity.ID

	go func() {
		_ = h.svc.SaveUserClick(context.Background(), c)
	}()

	http.Redirect(w, r, urlEntity.OriginalURL, http.StatusFound)
}

// GetAnalytics - (GET /analytics) получение аналитики по URL.
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	alias := strings.TrimPrefix(r.URL.Path, "/a/")
	alias = strings.TrimSpace(alias)

	if alias == "" {
		writeError(w, "alias is required", http.StatusBadRequest)
		return
	}

	urlEntity := entity.URL{Alias: alias}

	total, err := h.svc.CountClicks(r.Context(), urlEntity)
	if err != nil {
		writeError(w, "failed to count clicks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	records, err := h.svc.GetAnalytics(r.Context(), urlEntity)
	if err != nil {
		writeError(w, "failed to get analytics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	byDay, _ := h.svc.StatisticDay(r.Context(), urlEntity)
	byMonth, _ := h.svc.StatisticMonth(r.Context(), urlEntity)
	byAgent, _ := h.svc.StatisticAgent(r.Context(), urlEntity)

	// 3. Формируем DTO-ответ
	resp := dto.ShortenResponseAnalytics{
		Alias:       alias,
		ClicksTotal: total,
		Records:     records,
		ByDay:       byDay,
		ByAgent:     byAgent,
		ByMonth:     byMonth,
	}

	writeJSON(w, resp, http.StatusOK)
}

func isValidAlias(a string) bool {
	for _, r := range a {
		if !(r >= 'a' && r <= 'z') &&
			!(r >= 'A' && r <= 'Z') &&
			!(r >= '0' && r <= '9') &&
			r != '-' && r != '_' {
			return false
		}
	}
	return true
}

func isValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	writeJSON(w, map[string]string{"error": msg}, code)
}
