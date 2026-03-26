# Easy-Study Constitution

## Core Principles

### I. Multi-Tenant SaaS Architecture
- Mỗi tenant (trường/tổ chức) MUST isolate hoàn toàn ở tầng ứng dụng và cơ sở dữ liệu
- Single codebase, shared infrastructure, per-tenant configuration
- Quy trình onboarding tenant MUST tự động hóa (self-service registration)
- MUST support custom domain/subdomain cho mỗi tenant

### II. API-First Design
- Mọi feature MUST expose qua RESTful API trước; UI consume API
- MUST apply API versioning (ví dụ: `/api/v1/`)
- Xác thực bằng JWT; phân quyền bằng RBAC (Role-Based Access Control)
- Tài liệu API (OpenAPI/Swagger) MUST đồng bộ với implementation

### III. Modular & Scalable
- Kiến trúc module-based: mỗi business domain MUST là module độc lập
- Các module MUST giao tiếp qua well-defined interfaces; MUST NOT phụ thuộc chéo trực tiếp
- Horizontal scaling: stateless services, externalized sessions/cache
- Database migrations MUST backward-compatible

### IV. Security-First
- MUST tuân thủ OWASP Top 10
- Dữ liệu MUST mã hóa at rest và in transit (TLS 1.2+)
- Input validation & sanitization MUST apply ở mọi entry point
- Audit logging MUST ghi lại mọi thao tác nhạy cảm (enrollment, grading, user management)
- GDPR/data privacy: MUST support xuất và xóa dữ liệu theo yêu cầu người dùng

### V. Test-Driven & Quality Gates
- Unit tests MUST cover business logic (tối thiểu 80% coverage)
- Integration tests MUST cover API endpoints và cross-module interactions
- CI pipeline MUST pass trước khi merge: lint → test → build
- MUST NOT push trực tiếp vào `main`; MUST change qua Pull Request + review

### VI. Accessibility (a11y)
- MUST đạt chuẩn WCAG 2.1 AA cho mọi UI component
- MUST sử dụng semantic HTML, ARIA attributes, keyboard navigation
- Color contrast ratio MUST đạt chuẩn; MUST support screen reader
- MUST support light/dark mode
- MUST responsive cho Desktop và Mobile
- Kiểm thử accessibility MUST nằm trong quy trình QA

### VII. Internationalization (i18n)
- MUST support tối thiểu 2 ngôn ngữ: Tiếng Việt (VI) + English (EN)
- Mọi text hiển thị cho người dùng MUST use translation keys; MUST NOT hardcode chuỗi
- Định dạng theo locale cho ngày tháng, số, tiền tệ
- Ngôn ngữ mặc định: VI

### VIII. Observability
- Mọi service MUST use structured logging (định dạng JSON)
- Mỗi deployable unit MUST expose health check endpoint (`/healthz`)
- Error response MUST tuân theo format nhất quán (error code, message, trace ID)
- SHOULD triển khai request tracing xuyên suốt các service boundaries

### IX. Backup & Disaster Recovery
- Database MUST have lịch backup tự động
- MUST định nghĩa RTO (Recovery Time Objective) và RPO (Recovery Point Objective) trước khi lên production
- Quy trình khôi phục backup MUST kiểm thử định kỳ
- SHOULD have tài liệu disaster recovery runbook

## Technology Stack

- **Frontend**: Next.js (React) + TypeScript, Turborepo (monorepo build system), shadcn/ui (base UI) cho component library
- **Backend**: Golang
- **Database**: PostgreSQL (primary), Redis (cache/session)
- **Storage**: S3-compatible object storage (tài liệu khóa học, bài nộp)
- **Auth**: JWT + Refresh Token, OAuth2 (Google) integration
- **Infrastructure**: Docker containers, CI/CD qua GitHub Actions
- **Monitoring**: Structured logging (JSON), health check endpoints

## Development Workflow

1. **Branch Strategy**: `main` (production) ← `develop` (staging) ← `feature/*`, `fix/*`
2. **Commit Convention**: Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`)
3. **PR Process**: Feature branch → PR → Code review (≥1 approval) → CI pass → Merge
4. **Database Changes**: Chỉ dùng migration files, không DDL thủ công; migrations MUST có thể rollback
5. **Environment Parity**: Local dev, staging, production MUST dùng cùng Docker config base

## Governance

- Constitution này là tài liệu tối cao; mọi quyết định kỹ thuật MUST tuân theo các principles trên
- Sửa đổi yêu cầu: documented rationale, team review, migration plan nếu có breaking changes
- Versioning policy: MAJOR (bỏ/tái định nghĩa principle), MINOR (thêm principle/section), PATCH (clarification/typo)
- Khi conflict giữa tốc độ và chất lượng → chất lượng thắng (trừ khi có ngoại lệ được document rõ ràng)
- Compliance review: mọi PR/review MUST xác minh sự tuân thủ với các principles trong constitution

**Version**: 1.0.0 | **Ratified**: 2026-03-26 | **Last Amended**: 2026-03-26
