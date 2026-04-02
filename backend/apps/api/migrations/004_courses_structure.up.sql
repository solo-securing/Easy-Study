CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_courses_tenant_status ON courses (tenant_id, status);

CREATE TABLE IF NOT EXISTS course_sections (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    sort_order INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_course_sections_course ON course_sections (course_id, sort_order);

CREATE TABLE IF NOT EXISTS course_subsections (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    section_id UUID NOT NULL REFERENCES course_sections(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    sort_order INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_course_subsections_section ON course_subsections (section_id, sort_order);

CREATE TABLE IF NOT EXISTS course_units (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    subsection_id UUID NOT NULL REFERENCES course_subsections(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    unit_type TEXT NOT NULL,
    sort_order INT NOT NULL,
    content_ref TEXT,
    quiz_id UUID
);

CREATE INDEX IF NOT EXISTS idx_course_units_subsection ON course_units (subsection_id, sort_order);
