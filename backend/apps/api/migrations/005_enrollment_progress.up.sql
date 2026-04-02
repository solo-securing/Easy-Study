CREATE TABLE IF NOT EXISTS course_enrollments (
    tenant_id TEXT NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    source_type TEXT NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (course_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_course_enrollments_tenant_course ON course_enrollments (tenant_id, course_id);

CREATE TABLE IF NOT EXISTS course_progress_snapshots (
    tenant_id TEXT NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    completion_percent NUMERIC(5,2) NOT NULL DEFAULT 0,
    completed_unit_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (course_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_course_progress_tenant_course ON course_progress_snapshots (tenant_id, course_id);
