package services

import (
	"context"
	"time"

	"api/internal/database"
	"api/internal/models"
)

type ReportingService struct {
	repo database.UsageMetricsRepository
}

func NewReportingService(repo database.UsageMetricsRepository) *ReportingService {
	return &ReportingService{repo: repo}
}

func (s *ReportingService) SeedTenantMetric(ctx context.Context, metric models.UsageMetricDaily) error {
	if metric.MetricDate.IsZero() {
		metric.MetricDate = time.Now().UTC()
	}
	return s.repo.UpsertDaily(ctx, metric)
}

func (s *ReportingService) GetTenantDashboard(ctx context.Context, tenantID string) (models.TenantDashboard, error) {
	return s.repo.LatestTenantDashboard(ctx, tenantID)
}

func (s *ReportingService) GetCourseReport(ctx context.Context, tenantID, courseID string) (models.CourseReport, error) {
	return s.repo.GetCourseReport(ctx, tenantID, courseID)
}

func (s *ReportingService) GetAdminUsage(ctx context.Context, page, size int) (models.AdminUsageReport, error) {
	return s.repo.ListAdminUsage(ctx, page, size)
}
