package database

import (
	"context"
	"sync"
	"time"

	"api/internal/models"
)

type UsageMetricsRepository interface {
	UpsertDaily(ctx context.Context, metric models.UsageMetricDaily) error
	LatestTenantDashboard(ctx context.Context, tenantID string) (models.TenantDashboard, error)
	ListAdminUsage(ctx context.Context, page, size int) (models.AdminUsageReport, error)
	GetCourseReport(ctx context.Context, tenantID, courseID string) (models.CourseReport, error)
}

type InMemoryUsageMetricsRepository struct {
	mu       sync.RWMutex
	byTenant map[string]models.UsageMetricDaily
}

func NewInMemoryUsageMetricsRepository() *InMemoryUsageMetricsRepository {
	return &InMemoryUsageMetricsRepository{
		byTenant: map[string]models.UsageMetricDaily{},
	}
}

func (r *InMemoryUsageMetricsRepository) UpsertDaily(_ context.Context, metric models.UsageMetricDaily) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if metric.MetricDate.IsZero() {
		metric.MetricDate = time.Now().UTC()
	}
	r.byTenant[metric.TenantID] = metric
	return nil
}

func (r *InMemoryUsageMetricsRepository) LatestTenantDashboard(_ context.Context, tenantID string) (models.TenantDashboard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metric, ok := r.byTenant[tenantID]
	if !ok {
		return models.TenantDashboard{
			TenantID:          tenantID,
			ActiveUserCount:   0,
			ActiveCourseCount: 0,
			CompletionRate:    0,
		}, nil
	}
	return models.TenantDashboard{
		TenantID:          tenantID,
		ActiveUserCount:   metric.ActiveUserCount,
		ActiveCourseCount: metric.ActiveCourseCount,
		CompletionRate:    metric.CompletionRate,
	}, nil
}

func (r *InMemoryUsageMetricsRepository) ListAdminUsage(_ context.Context, page, size int) (models.AdminUsageReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	rows := make([]models.AdminUsageRow, 0, len(r.byTenant))
	for tenantID, metric := range r.byTenant {
		rows = append(rows, models.AdminUsageRow{
			TenantID:        tenantID,
			TenantName:      tenantID,
			UserCount:       metric.ActiveUserCount,
			CourseCount:     metric.ActiveCourseCount,
			LoginCount:      metric.LoginCount,
			LearningMinutes: metric.LearningMinutes,
		})
	}

	start := (page - 1) * size
	if start > len(rows) {
		start = len(rows)
	}
	end := start + size
	if end > len(rows) {
		end = len(rows)
	}
	return models.AdminUsageReport{
		Items: rows[start:end],
		Meta: models.PageMeta{
			Page:  page,
			Size:  size,
			Total: len(rows),
		},
	}, nil
}

func (r *InMemoryUsageMetricsRepository) GetCourseReport(_ context.Context, tenantID, courseID string) (models.CourseReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metric, ok := r.byTenant[tenantID]
	if !ok {
		metric = models.UsageMetricDaily{TenantID: tenantID}
	}
	return models.CourseReport{
		CourseID:       courseID,
		LearnerCount:   metric.ActiveUserCount,
		CompletionRate: metric.CompletionRate,
		Items: []models.StudentProgressItem{
			{
				UserID:             "student-001",
				CompletionPercent:  metric.CompletionRate,
				CompletedUnitCount: 0,
				OfficialQuizScore:  0,
			},
		},
	}, nil
}
