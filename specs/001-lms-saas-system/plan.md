# Kế Hoạch Triển Khai: Nền Tảng LMS SaaS Multi-Tenant

**Branch**: `001-lms-saas-system` | **Date**: 2026-04-01 | **Spec**: `/specs/001-lms-saas-system/spec.md`  
**Input**: Đặc tả tính năng từ `/specs/001-lms-saas-system/spec.md`

## Tóm Tắt

Xây dựng nền tảng LMS SaaS multi-tenant cho 4 vai trò: Super Admin, Tenant Admin, Instructor và Student trên hạ tầng dùng chung nhưng cách ly tenant nghiêm ngặt. Việc triển khai tuân theo API-first governance, lấy OpenAPI 3.1 làm source of truth, đồng bộ type generation cho backend/frontend, kiểm soát quota-entitlement, hỗ trợ quiz autosave/recovery, đảm bảo độ tin cậy khi gửi notification, và đầy đủ auditability.

Baseline thiết kế hiện tại đã được rà soát lại và bao gồm:

- Mô hình shared-schema multi-tenancy với cách ly bằng `tenant_id`.
- Frontend dùng Next.js 16 + TypeScript trong Turborepo.
- Backend dùng Go 1.26 + Gin với generated interfaces từ OpenAPI.
- PostgreSQL 18 + Redis + object storage tương thích S3 (Rustfs).
- Hợp đồng API v1 đã mở rộng (54 paths / 65 operations) và đồng bộ với data model.

## Bối Cảnh Kỹ Thuật

**Language/Version**: Go 1.26 (backend); TypeScript 5.x trên Next.js 16 (frontend)

**Primary Dependencies**:

- Backend: Gin, oapi-codegen, golang-migrate/migrate, JWT middleware, thư viện OAuth2
- Frontend: Next.js 16, Turborepo, shadcn/ui, openapi-typescript
- API contract: OpenAPI 3.1, swagger-cli validation, Spectral (khi có ruleset), Schemathesis

**Storage**:

- PostgreSQL 18 (dữ liệu quan hệ chính)
- Redis (session/cache/rate-limit counters)
- Rustfs S3-compatible object storage (luồng metadata cho documents/submissions/media)

**Testing**:

- Backend unit/integration: `go test ./...`
- Frontend unit/component và E2E: dùng test runners theo package trong Turborepo
- Contract/fuzz testing: Schemathesis chạy theo OpenAPI endpoints
- Contract structure validation: swagger-cli

**Target Platform**:

- Linux containers (Docker) để đồng nhất local, staging, production
- Trình duyệt hiện đại (Chromium/Firefox/Safari/Edge)

**Project Type**: Web application (monorepo gồm frontend + backend)

**Performance Goals**:

- API p95 latency <= 200ms (read), <= 500ms (write)
- Dashboard load <= 3s với tenant cỡ tới 10k users
- In-app notifications được phân phối trong vòng 30s

**Constraints**:

- Bắt buộc tenant isolation cho mọi read/write path
- Quy trình API-first phải đi trước khi hiện thực business logic
- JWT expiry 1 giờ, RBAC tại API layer, tuân thủ OWASP/GDPR
- Áp dụng account lock + rate limit theo FR-042

**Scale/Scope**:

- Benchmark mục tiêu: 100 active tenants, tối đa 1,000 concurrent users mỗi tenant
- Phạm vi v1 bao gồm tenant lifecycle, course lifecycle, quiz attempts, reporting, notifications, billing entitlements

## Kiểm Tra Constitution

*GATE: Bắt buộc pass trước Phase 0 research. Rà soát lại sau Phase 1 design.*

### Rà Soát Gate Trước Phase 0

- Principle I (Code Quality): PASS. Kế hoạch duy trì boundary rõ ràng và typed interfaces.
- Principle II (Testing): PASS. Đã nêu rõ chiến lược unit/integration/E2E/contract.
- Principle III (UX Consistency): PASS. Hành vi giao diện tenant được thống nhất qua shared frontend stack.
- Principle IV (Performance): PASS. Có ngân sách hiệu năng định lượng và gắn với success criteria.
- Principle V (Security): PASS. Đã mô hình hóa isolation, RBAC, lockout, rate-limit, SSO mapping và audit logging.
- Principle VI (Simplicity): PASS. Shared schema + generated contracts giúp giảm độ phức tạp không cần thiết.
- Principle VII (API-First): PASS. Quy trình Design -> Validate -> Generate -> Test được giữ xuyên suốt.

### Rà Soát Lại Sau Phase 1

- Principle I: PASS. Data model đã phủ policy/config, delivery outbox và auth session state.
- Principle II: PASS. Quickstart có luồng validation thực thi được và các test gates.
- Principle III: PASS. Contract thể hiện nhất quán hành vi theo role và theo tenant.
- Principle IV: PASS. Các entity usage/quota/report hỗ trợ chiến lược dashboard theo mục tiêu hiệu năng.
- Principle V: PASS. Invariants bảo mật cốt lõi (lockout, rate limit, activation, SSO linkage, audit) đã rõ ràng.
- Principle VI: PASS. Mọi phần bổ sung đều bám requirement, không tạo tầng thừa.
- Principle VII: PASS. OpenAPI và data model đã đồng bộ và được kiểm chứng.

## Cấu Trúc Dự Án

### Tài Liệu (feature này)

```text
specs/001-lms-saas-system/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── openapi-lms-v1.yaml
│   └── domain-events.md
└── tasks.md
```

### Source Code (repository root)

```text
backend/
├── apps/
│   └── api/
│       ├── cmd/api/main.go
│       └── internal/
│           ├── database/
│           └── server/
└── packages/
    ├── configloader/
    └── database/

frontend/
├── apps/
│   ├── web/
│   └── docs/
└── packages/
    ├── ui/
    ├── eslint-config/
    └── typescript-config/

dev/
└── docker-compose.yaml
```

**Structure Decision**: Giữ cấu trúc web-application hiện có (`backend` + `frontend`) và triển khai contract-first trong `specs/001-lms-saas-system/contracts` trước khi code service/UI.

## Kế Hoạch Theo Phase

### Phase 0 - Research

- Chốt các quyết định cho tenancy model, auth/session model, quiz autosave policy, quota enforcement, notification reliability và ranh giới SSO/JIT.
- Tổng hợp alternatives và rationale vào `research.md`.

### Phase 1 - Design và Contracts

- Duy trì entity model và invariants trong `data-model.md`.
- Duy trì REST contract ở `contracts/openapi-lms-v1.yaml` và async event contracts ở `contracts/domain-events.md`.
- Duy trì onboarding/validation flow có thể chạy được trong `quickstart.md`.

### Phase 2 - Chuẩn Bị Lập Task

- Dùng baseline plan/design đã ổn định làm đầu vào cho `/speckit.tasks` để sinh task theo dependency order.

## Trạng Thái Validation Hiện Tại

- Plan, research, contract và data model đã được rà soát end-to-end.
- OpenAPI contract đã validate bằng swagger-cli.
- Data model và API đã khớp ở mức endpoint, schema và entity-invariant.

## Theo Dõi Độ Phức Tạp

Không có vi phạm constitution nào cần exception ở giai đoạn planning.