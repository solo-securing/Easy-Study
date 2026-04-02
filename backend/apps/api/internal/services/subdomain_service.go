package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"
)

var (
	ErrInvalidSubdomain  = errors.New("invalid subdomain")
	ErrReservedSubdomain = errors.New("subdomain is reserved")
	subdomainPattern     = regexp.MustCompile(`^[a-z0-9-]{3,63}$`)
)

type SubdomainService struct {
	mu       sync.RWMutex
	reserved map[string]struct{}
}

func NewSubdomainService(initialReserved []string) *SubdomainService {
	set := make(map[string]struct{}, len(initialReserved))
	for _, value := range initialReserved {
		set[strings.ToLower(strings.TrimSpace(value))] = struct{}{}
	}
	return &SubdomainService{reserved: set}
}

func (s *SubdomainService) Validate(_ context.Context, subdomain string) error {
	normalized := strings.ToLower(strings.TrimSpace(subdomain))
	if !subdomainPattern.MatchString(normalized) {
		return ErrInvalidSubdomain
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, found := s.reserved[normalized]; found {
		return ErrReservedSubdomain
	}

	return nil
}
