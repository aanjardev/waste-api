package dto

type WasteSummaryItem struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type PaymentSummaryItem struct {
	Status string  `json:"status"`
	Count  int     `json:"count"`
	Total  float64 `json:"total"`
}

type PaymentSummaryResponse struct {
	ByStatus     []PaymentSummaryItem `json:"by_status"`
	TotalRevenue float64              `json:"total_revenue"`
}