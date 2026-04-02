package models

import "time"

type QuizScoringPolicy string

const (
	QuizScoringPolicyHighestValidAttempt QuizScoringPolicy = "highest_valid_attempt"
)

type QuizAttemptStatus string

const (
	QuizAttemptStatusInProgress   QuizAttemptStatus = "in_progress"
	QuizAttemptStatusSubmitted    QuizAttemptStatus = "submitted"
	QuizAttemptStatusAutoSubmitted QuizAttemptStatus = "auto_submitted"
	QuizAttemptStatusExpired      QuizAttemptStatus = "expired"
)

type Quiz struct {
	QuizID            string            `json:"quizId"`
	TenantID          string            `json:"tenantId"`
	CourseID          string            `json:"courseId"`
	Title             string            `json:"title"`
	MaxScore          float64           `json:"maxScore"`
	PassingScore      float64           `json:"passingScore"`
	TimeLimitMinutes  int               `json:"timeLimitMinutes"`
	AttemptLimit      int               `json:"attemptLimit"`
	ShuffleQuestions  bool              `json:"shuffleQuestions"`
	ScoringPolicy     QuizScoringPolicy `json:"scoringPolicy"`
}

type QuizAnswer struct {
	QuestionID        string   `json:"questionId"`
	SelectedOptionIDs []string `json:"selectedOptionIds,omitempty"`
	TextAnswer        string   `json:"textAnswer,omitempty"`
}

type QuizAttempt struct {
	AttemptID  string            `json:"attemptId"`
	TenantID   string            `json:"tenantId"`
	QuizID     string            `json:"quizId"`
	UserID     string            `json:"userId"`
	AttemptNo  int               `json:"attemptNo"`
	Status     QuizAttemptStatus `json:"status"`
	StartedAt  time.Time         `json:"startedAt"`
	ExpiresAt  time.Time         `json:"expiresAt"`
	SubmittedAt *time.Time       `json:"submittedAt,omitempty"`
	Answers    []QuizAnswer      `json:"answers,omitempty"`
	Score      float64           `json:"score,omitempty"`
}

type QuizAttemptStartInput struct {
	TenantID        string `json:"tenantId"`
	QuizID          string `json:"quizId"`
	UserID          string `json:"userId"`
	ResumeAttemptID string `json:"resumeAttemptId,omitempty"`
}

type QuizAutosaveInput struct {
	TenantID      string       `json:"tenantId"`
	QuizID        string       `json:"quizId"`
	AttemptID     string       `json:"attemptId"`
	UserID        string       `json:"userId"`
	Answers       []QuizAnswer `json:"answers"`
	ClientSavedAt time.Time    `json:"clientSavedAt"`
}

type QuizSubmitResult struct {
	AttemptID                 string            `json:"attemptId"`
	Status                    QuizAttemptStatus `json:"status"`
	Score                     float64           `json:"score"`
	OfficialScorePolicy       QuizScoringPolicy `json:"officialScorePolicy"`
	CountedTowardAttemptLimit bool              `json:"countedTowardAttemptLimit"`
}
