package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aanjardev/waste-api/internal/dto"
)

type ReportRepository interface {
	WasteSummary(ctx context.Context) ([]dto.WasteSummaryItem, error)
	PaymentSummary(ctx context.Context) (*dto.PaymentSummaryResponse, error)
}

type reportRepository struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) ReportRepository {
	return &reportRepository{pool: pool}
}

func (r *reportRepository) WasteSummary(ctx context.Context) ([]dto.WasteSummaryItem, error) {
	query := `
		SELECT type::text, status::text, COUNT(*) AS count
		FROM waste_pickups
		GROUP BY type, status
		ORDER BY type, status
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query waste summary: %w", err)
	}
	defer rows.Close()

	var items []dto.WasteSummaryItem
	for rows.Next() {
		var item dto.WasteSummaryItem
		if err := rows.Scan(&item.Type, &item.Status, &item.Count); err != nil {
			return nil, fmt.Errorf("scan waste summary: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return items, nil
}

func (r *reportRepository) PaymentSummary(ctx context.Context) (*dto.PaymentSummaryResponse, error) {
	query := `
		SELECT status::text, COUNT(*) AS count, COALESCE(SUM(amount), 0) AS total
		FROM payments
		GROUP BY status
		ORDER BY status
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query payment summary: %w", err)
	}
	defer rows.Close()

	resp := &dto.PaymentSummaryResponse{
		ByStatus: []dto.PaymentSummaryItem{},
	}

	for rows.Next() {
		var item dto.PaymentSummaryItem
		if err := rows.Scan(&item.Status, &item.Count, &item.Total); err != nil {
			return nil, fmt.Errorf("scan payment summary: %w", err)
		}
		resp.ByStatus = append(resp.ByStatus, item)

		// Total revenue = hanya yang paid
		if item.Status == "paid" {
			resp.TotalRevenue = item.Total
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return resp, nil
}