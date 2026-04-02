package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrQuizNotFound         = errors.New("quiz not found")
	ErrQuizAttemptNotFound  = errors.New("quiz attempt not found")
	ErrQuizAttemptConflict  = errors.New("quiz attempt conflict")
)

type QuizRepository interface {
	GetQuiz(ctx context.Context, tenantID, quizID string) (models.Quiz, error)
	UpsertQuiz(ctx context.Context, quiz models.Quiz) error
	StartOrResumeAttempt(ctx context.Context, input models.QuizAttemptStartInput) (models.QuizAttempt, error)
	GetAttempt(ctx context.Context, tenantID, quizID, attemptID string) (models.QuizAttempt, error)
	SaveAttemptAnswers(ctx context.Context, input models.QuizAutosaveInput) (models.QuizAttempt, error)
	SubmitAttempt(ctx context.Context, tenantID, quizID, attemptID string, auto bool) (models.QuizAttempt, error)
}

type InMemoryQuizRepository struct {
	mu          sync.RWMutex
	quizzes     map[string]models.Quiz
	attempts    map[string]models.QuizAttempt
	userAttempt map[string][]string
}

func NewInMemoryQuizRepository() *InMemoryQuizRepository {
	return &InMemoryQuizRepository{
		quizzes:     map[string]models.Quiz{},
		attempts:    map[string]models.QuizAttempt{},
		userAttempt: map[string][]string{},
	}
}

func quizKey(tenantID, quizID string) string {
	return tenantID + ":" + quizID
}

func userQuizKey(tenantID, quizID, userID string) string {
	return tenantID + ":" + quizID + ":" + userID
}

func (r *InMemoryQuizRepository) UpsertQuiz(_ context.Context, quiz models.Quiz) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.quizzes[quizKey(quiz.TenantID, quiz.QuizID)] = quiz
	return nil
}

func (r *InMemoryQuizRepository) GetQuiz(_ context.Context, tenantID, quizID string) (models.Quiz, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	quiz, ok := r.quizzes[quizKey(tenantID, quizID)]
	if !ok {
		return models.Quiz{}, ErrQuizNotFound
	}
	return quiz, nil
}

func (r *InMemoryQuizRepository) StartOrResumeAttempt(_ context.Context, input models.QuizAttemptStartInput) (models.QuizAttempt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	quiz, ok := r.quizzes[quizKey(input.TenantID, input.QuizID)]
	if !ok {
		return models.QuizAttempt{}, ErrQuizNotFound
	}

	if input.ResumeAttemptID != "" {
		attempt, ok := r.attempts[input.ResumeAttemptID]
		if !ok || attempt.TenantID != input.TenantID || attempt.QuizID != input.QuizID || attempt.UserID != input.UserID {
			return models.QuizAttempt{}, ErrQuizAttemptNotFound
		}
		if attempt.Status != models.QuizAttemptStatusInProgress {
			return models.QuizAttempt{}, ErrQuizAttemptConflict
		}
		return attempt, nil
	}

	userKey := userQuizKey(input.TenantID, input.QuizID, input.UserID)
	for _, attemptID := range r.userAttempt[userKey] {
		attempt := r.attempts[attemptID]
		if attempt.Status == models.QuizAttemptStatusInProgress && time.Now().UTC().Before(attempt.ExpiresAt) {
			return attempt, nil
		}
	}

	attemptNo := len(r.userAttempt[userKey]) + 1
	now := time.Now().UTC()
	expiry := now.Add(time.Duration(quiz.TimeLimitMinutes) * time.Minute)
	if quiz.TimeLimitMinutes <= 0 {
		expiry = now.Add(1 * time.Second)
	}
	attempt := models.QuizAttempt{
		AttemptID: uuid.NewString(),
		TenantID:  input.TenantID,
		QuizID:    input.QuizID,
		UserID:    input.UserID,
		AttemptNo: attemptNo,
		Status:    models.QuizAttemptStatusInProgress,
		StartedAt: now,
		ExpiresAt: expiry,
		Answers:   []models.QuizAnswer{},
	}

	r.attempts[attempt.AttemptID] = attempt
	r.userAttempt[userKey] = append(r.userAttempt[userKey], attempt.AttemptID)
	return attempt, nil
}

func (r *InMemoryQuizRepository) GetAttempt(_ context.Context, tenantID, quizID, attemptID string) (models.QuizAttempt, error) {
	r.mu.RLock()
	attempt, ok := r.attempts[attemptID]
	r.mu.RUnlock()
	if !ok || attempt.TenantID != tenantID || attempt.QuizID != quizID {
		return models.QuizAttempt{}, ErrQuizAttemptNotFound
	}
	if attempt.Status == models.QuizAttemptStatusInProgress && time.Now().UTC().After(attempt.ExpiresAt) {
		r.mu.Lock()
		now := time.Now().UTC()
		attempt.Status = models.QuizAttemptStatusAutoSubmitted
		attempt.SubmittedAt = &now
		attempt.Score = float64(len(attempt.Answers)) * 10
		r.attempts[attemptID] = attempt
		r.mu.Unlock()
	}
	return attempt, nil
}

func (r *InMemoryQuizRepository) SaveAttemptAnswers(_ context.Context, input models.QuizAutosaveInput) (models.QuizAttempt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	attempt, ok := r.attempts[input.AttemptID]
	if !ok || attempt.TenantID != input.TenantID || attempt.QuizID != input.QuizID || attempt.UserID != input.UserID {
		return models.QuizAttempt{}, ErrQuizAttemptNotFound
	}
	if attempt.Status != models.QuizAttemptStatusInProgress {
		return models.QuizAttempt{}, ErrQuizAttemptConflict
	}
	if time.Now().UTC().After(attempt.ExpiresAt) {
		now := time.Now().UTC()
		attempt.Status = models.QuizAttemptStatusAutoSubmitted
		attempt.SubmittedAt = &now
		attempt.Score = float64(len(attempt.Answers)) * 10
		r.attempts[input.AttemptID] = attempt
		return models.QuizAttempt{}, ErrQuizAttemptConflict
	}

	attempt.Answers = input.Answers
	r.attempts[input.AttemptID] = attempt
	return attempt, nil
}

func (r *InMemoryQuizRepository) SubmitAttempt(_ context.Context, tenantID, quizID, attemptID string, auto bool) (models.QuizAttempt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	attempt, ok := r.attempts[attemptID]
	if !ok || attempt.TenantID != tenantID || attempt.QuizID != quizID {
		return models.QuizAttempt{}, ErrQuizAttemptNotFound
	}
	if attempt.Status != models.QuizAttemptStatusInProgress {
		return models.QuizAttempt{}, ErrQuizAttemptConflict
	}

	now := time.Now().UTC()
	attempt.SubmittedAt = &now
	if auto {
		attempt.Status = models.QuizAttemptStatusAutoSubmitted
	} else {
		attempt.Status = models.QuizAttemptStatusSubmitted
	}
	attempt.Score = float64(len(attempt.Answers)) * 10
	r.attempts[attemptID] = attempt
	return attempt, nil
}
