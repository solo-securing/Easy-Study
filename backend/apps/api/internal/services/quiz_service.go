package services

import (
	"context"
	"errors"

	"api/internal/database"
	"api/internal/models"
)

type QuizService struct {
	repo  database.QuizRepository
	audit *AuditService
}

func NewQuizService(repo database.QuizRepository, audit *AuditService) *QuizService {
	return &QuizService{
		repo:  repo,
		audit: audit,
	}
}

func (s *QuizService) GetQuiz(ctx context.Context, tenantID, quizID string) (models.Quiz, error) {
	return s.repo.GetQuiz(ctx, tenantID, quizID)
}

func (s *QuizService) StartAttempt(ctx context.Context, actorID string, input models.QuizAttemptStartInput) (models.QuizAttempt, error) {
	quiz, err := s.repo.GetQuiz(ctx, input.TenantID, input.QuizID)
	if err != nil {
		return models.QuizAttempt{}, err
	}

	attempt, err := s.repo.StartOrResumeAttempt(ctx, input)
	if err != nil {
		return models.QuizAttempt{}, err
	}

	if attempt.AttemptNo > quiz.AttemptLimit {
		return models.QuizAttempt{}, errors.New("attempt limit reached")
	}
	if attempt.ExpiresAt.Before(attempt.StartedAt) {
		return models.QuizAttempt{}, errors.New("invalid attempt expiration")
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: input.TenantID,
			ActorID:  actorID,
			Action:   "quiz.attempt_start",
			Resource: "quiz_attempt",
			Result:   "success",
		})
	}
	return attempt, nil
}

func (s *QuizService) GetAttempt(ctx context.Context, tenantID, quizID, attemptID string) (models.QuizAttempt, error) {
	return s.repo.GetAttempt(ctx, tenantID, quizID, attemptID)
}

func (s *QuizService) Autosave(ctx context.Context, actorID string, input models.QuizAutosaveInput) (models.QuizAttempt, error) {
	attempt, err := s.repo.SaveAttemptAnswers(ctx, input)
	if err != nil {
		return models.QuizAttempt{}, err
	}
	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: input.TenantID,
			ActorID:  actorID,
			Action:   "quiz.autosave",
			Resource: "quiz_attempt",
			Result:   "success",
		})
	}
	return attempt, nil
}

func (s *QuizService) Submit(ctx context.Context, actorID, tenantID, quizID, attemptID string, auto bool) (models.QuizSubmitResult, error) {
	attempt, err := s.repo.SubmitAttempt(ctx, tenantID, quizID, attemptID, auto)
	if err != nil {
		return models.QuizSubmitResult{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "quiz.submit",
			Resource: "quiz_attempt",
			Result:   "success",
		})
	}

	return models.QuizSubmitResult{
		AttemptID:                 attempt.AttemptID,
		Status:                    attempt.Status,
		Score:                     attempt.Score,
		OfficialScorePolicy:       models.QuizScoringPolicyHighestValidAttempt,
		CountedTowardAttemptLimit: true,
	}, nil
}

func (s *QuizService) SubmitExpiredAttempt(ctx context.Context, tenantID, quizID, attemptID string) (models.QuizSubmitResult, error) {
	return s.Submit(ctx, "system-timeout", tenantID, quizID, attemptID, true)
}

func (s *QuizService) SeedQuiz(ctx context.Context, quiz models.Quiz) error {
	return s.repo.UpsertQuiz(ctx, quiz)
}
