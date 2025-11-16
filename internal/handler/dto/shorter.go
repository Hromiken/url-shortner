package dto

import "shortner/internal/entity"

type ShortenResponseAnalytics struct {
	Records     []entity.ClickLogs `json:"records,omitempty"`
	Alias       string             `json:"alias"`
	ByDay       []entity.StatDay   `json:"by_day"`
	ByMonth     []entity.StatMonth `json:"by_month"`
	ByAgent     []entity.StatAgent `json:"by_agent"`
	ClicksTotal int64              `json:"clicks_total"`
}

type CreateURLShorterRequest struct {
	URL         string `json:"url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

type ShortenResponse struct {
	Alias string `json:"alias"`
	Short string `json:"short"`
}
