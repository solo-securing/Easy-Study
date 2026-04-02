package services

import "time"

type LockoutService struct {
	MaxFailedAttempts int
	LockDuration      time.Duration
}

func NewLockoutService() *LockoutService {
	return &LockoutService{
		MaxFailedAttempts: 5,
		LockDuration:      15 * time.Minute,
	}
}

func (s *LockoutService) ShouldLock(failedAttempts int) bool {
	return failedAttempts >= s.MaxFailedAttempts
}

func (s *LockoutService) LockedUntil(now time.Time) time.Time {
	return now.Add(s.LockDuration)
}
