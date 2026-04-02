package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"api/internal/models"
	"github.com/google/uuid"
)

var ErrTenantNotFound = errors.New("tenant not found")
var ErrTenantSubdomainExists = errors.New("tenant subdomain already exists")

type TenantRepository interface {
	List(ctx context.Context, filter models.TenantListFilter) (models.TenantList, error)
	Create(ctx context.Context, input models.TenantCreateInput) (models.Tenant, error)
	GetByID(ctx context.Context, tenantID string) (models.Tenant, error)
	UpdateStatus(ctx context.Context, tenantID string, patch models.TenantStatusPatch) (models.Tenant, error)
}

type InMemoryTenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]models.Tenant
}

func NewInMemoryTenantRepository() *InMemoryTenantRepository {
	return &InMemoryTenantRepository{
		tenants: map[string]models.Tenant{},
	}
}

func (r *InMemoryTenantRepository) List(_ context.Context, filter models.TenantListFilter) (models.TenantList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 20
	}

	items := make([]models.Tenant, 0, len(r.tenants))
	for _, tenant := range r.tenants {
		if filter.Status != "" && tenant.Status != filter.Status {
			continue
		}
		items = append(items, tenant)
	}

	start := (filter.Page - 1) * filter.Size
	if start > len(items) {
		start = len(items)
	}
	end := start + filter.Size
	if end > len(items) {
		end = len(items)
	}

	return models.TenantList{
		Items: items[start:end],
		Meta: models.PageMeta{
			Page:  filter.Page,
			Size:  filter.Size,
			Total: len(items),
		},
	}, nil
}

func (r *InMemoryTenantRepository) Create(_ context.Context, input models.TenantCreateInput) (models.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.tenants {
		if existing.Subdomain == input.Subdomain {
			return models.Tenant{}, ErrTenantSubdomainExists
		}
	}

	now := time.Now().UTC()
	tenant := models.Tenant{
		ID:            uuid.NewString(),
		Name:          input.Name,
		Subdomain:     input.Subdomain,
		Status:        models.TenantStatusActive,
		PlanCode:      input.PlanCode,
		CreatedAt:     now,
		DefaultLocale: "en",
		Timezone:      "UTC",
	}

	r.tenants[tenant.ID] = tenant
	return tenant, nil
}

func (r *InMemoryTenantRepository) GetByID(_ context.Context, tenantID string) (models.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenant, ok := r.tenants[tenantID]
	if !ok {
		return models.Tenant{}, ErrTenantNotFound
	}

	return tenant, nil
}

func (r *InMemoryTenantRepository) UpdateStatus(_ context.Context, tenantID string, patch models.TenantStatusPatch) (models.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenant, ok := r.tenants[tenantID]
	if !ok {
		return models.Tenant{}, ErrTenantNotFound
	}

	tenant.Status = patch.Status
	r.tenants[tenantID] = tenant

	return tenant, nil
}
