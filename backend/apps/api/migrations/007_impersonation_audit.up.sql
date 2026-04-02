CREATE TABLE IF NOT EXISTS impersonation_sessions (
    session_id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    super_admin_id TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    mode TEXT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_impersonation_sessions_tenant ON impersonation_sessions (tenant_id, started_at DESC);

CREATE TABLE IF NOT EXISTS impersonation_actions (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES impersonation_sessions(session_id) ON DELETE CASCADE,
    action TEXT NOT NULL,
    at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_impersonation_actions_session ON impersonation_actions (session_id, at DESC);
