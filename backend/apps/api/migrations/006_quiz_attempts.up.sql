CREATE TABLE IF NOT EXISTS quizzes (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    course_id UUID NOT NULL,
    title TEXT NOT NULL,
    max_score NUMERIC(10,2) NOT NULL,
    passing_score NUMERIC(10,2) NOT NULL,
    time_limit_minutes INT NOT NULL,
    attempt_limit INT NOT NULL,
    scoring_policy TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_quizzes_tenant_course ON quizzes (tenant_id, course_id);

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    attempt_no INT NOT NULL,
    status TEXT NOT NULL,
    answers JSONB NOT NULL DEFAULT '[]'::jsonb,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    submitted_at TIMESTAMPTZ,
    score NUMERIC(10,2),
    UNIQUE (quiz_id, user_id, attempt_no)
);

CREATE INDEX IF NOT EXISTS idx_quiz_attempts_tenant_quiz_user ON quiz_attempts (tenant_id, quiz_id, user_id);
