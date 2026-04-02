CREATE TABLE IF NOT EXISTS usage_metrics_daily (
    tenant_id TEXT NOT NULL,
    metric_date DATE NOT NULL,
    active_user_count INT NOT NULL DEFAULT 0,
    active_course_count INT NOT NULL DEFAULT 0,
    completion_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    login_count INT NOT NULL DEFAULT 0,
    learning_minutes INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, metric_date)
);

CREATE INDEX IF NOT EXISTS idx_usage_metrics_daily_tenant_date ON usage_metrics_daily (tenant_id, metric_date DESC);
