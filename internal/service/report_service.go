package service

import (
	"context"
	"fmt"

	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/repository"
)

type ReportService interface {
	WasteSummary(ctx context.Context) ([]dto.WasteSummaryItem, error)
	PaymentSummary(ctx context.Context) (*dto.PaymentSummaryResponse, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) WasteSummary(ctx context.Context) ([]dto.WasteSummaryItem, error) {
	items, err := s.repo.WasteSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("waste summary: %w", err)
	}
	return items, nil
}

func (s *reportService) PaymentSummary(ctx context.Context) (*dto.PaymentSummaryResponse, error) {
	resp, err := s.repo.PaymentSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("payment summary: %w", err)
	}
	return resp, nil
}