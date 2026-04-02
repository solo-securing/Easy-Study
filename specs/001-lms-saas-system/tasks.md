# Tasks: LMS SaaS Multi-Tenant Platform

**Input**: Design documents from `/specs/001-lms-saas-system/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Bao gom test tasks vi plan/spec yeu cau quality gates va contract validation.  
**Organization**: Tasks duoc nhom theo user story de moi story co the implement va test doc lap.

## Format: `[ID] [P?] [Story?] Description with file path`

- **[P]**: Co the chay song song (khac file, khong phu thuoc task chua xong)
- **[Story]**: Nhan user story (`[US1]`, `[US2]`, ...)
- Moi task phai ghi ro file path

## Path Conventions

- Backend: `backend/apps/api/`
- Frontend: `frontend/apps/web/`
- Shared packages: `backend/packages/`, `frontend/packages/`
- Specs and contracts: `specs/001-lms-saas-system/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Khoi tao khung du an va toolchain theo plan

- [X] T001 Cap nhat command codegen OpenAPI trong `backend/apps/api/Makefile`
- [X] T002 Them script generate API types trong `frontend/package.json`
- [X] T003 [P] Khoi tao cau truc thu muc backend handlers/services/repositories trong `backend/apps/api/internal/`
- [X] T004 [P] Khoi tao route groups cho Super Admin/Tenant Admin/Instructor/Student trong `frontend/apps/web/src/app/`
- [X] T005 [P] Tao file environment mau cho API trong `backend/apps/api/.env.example`
- [X] T006 [P] Chuan hoa local infra services trong `dev/docker-compose.yaml`
- [X] T007 [P] Tao workflow validate OpenAPI contract trong `.github/workflows/contract-validation.yml`
- [X] T008 Khoi tao thu muc migrations va baseline placeholder trong `backend/apps/api/migrations/.keep`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Nen tang bat buoc truoc khi implement bat ky user story nao

**CRITICAL**: Khong bat dau user story khi phase nay chua xong

- [X] T009 Tao migration runner bootstrap trong `backend/apps/api/internal/database/migrate.go`
- [X] T010 [P] Implement tenant context middleware trong `backend/apps/api/internal/middleware/tenant_context.go`
- [X] T011 [P] Implement JWT authentication middleware trong `backend/apps/api/internal/middleware/auth_jwt.go`
- [X] T012 [P] Implement RBAC authorization middleware trong `backend/apps/api/internal/middleware/rbac.go`
- [X] T013 [P] Implement request logging va trace middleware trong `backend/apps/api/internal/middleware/request_logger.go`
- [X] T014 [P] Implement problem-details error mapper trong `backend/apps/api/internal/server/problem.go`
- [X] T015 [P] Implement audit log writer service trong `backend/apps/api/internal/services/audit_service.go`
- [X] T016 [P] Implement Redis rate limiter cho public APIs trong `backend/apps/api/internal/middleware/rate_limit.go`
- [X] T017 [P] Implement account lockout service (5 lan sai, khoa 15 phut) trong `backend/apps/api/internal/services/lockout_service.go`
- [X] T018 [P] Khoi tao worker entrypoint cho async jobs trong `backend/apps/api/cmd/workers/main.go`
- [X] T019 [P] Wire generated OpenAPI router vao server trong `backend/apps/api/internal/server/routes.go`
- [X] T020 Tao integration test bootstrap trong `backend/apps/api/tests/integration/test_main.go`

**Checkpoint**: Foundation san sang - co the bat dau user stories

---

## Phase 3: User Story 1 - Super Admin Tao va Quan Ly Tenant Moi (Priority: P1) MVP

**Goal**: Tao tenant moi, quan ly lifecycle (active/suspended/deactivated), va bao dam subdomain uniqueness

**Independent Test**: Super Admin tao tenant moi, thay tenant trong danh sach, truy cap duoc subdomain, suspend/reactivate thanh cong

- [X] T021 [P] [US1] Tao tenant core migration (tenants, plans, tenant_plan_assignments) trong `backend/apps/api/migrations/001_tenant_core.up.sql`
- [X] T022 [P] [US1] Tao tenant models trong `backend/apps/api/internal/models/tenant.go`
- [X] T023 [P] [US1] Implement tenant repository trong `backend/apps/api/internal/database/tenant_repository.go`
- [X] T024 [US1] Implement tenant lifecycle service trong `backend/apps/api/internal/services/tenant_service.go`
- [X] T025 [US1] Implement reserved subdomain validation service trong `backend/apps/api/internal/services/subdomain_service.go`
- [X] T026 [US1] Implement tenant handlers (list/create/get/status) trong `backend/apps/api/internal/handlers/tenant_handler.go`
- [X] T027 [US1] Dang ky tenant routes va RBAC policies trong `backend/apps/api/internal/server/routes.go`
- [X] T028 [US1] Tao trang danh sach tenant cho Super Admin trong `frontend/apps/web/src/app/(super-admin)/tenants/page.tsx`
- [X] T029 [P] [US1] Tao tenant create form component trong `frontend/apps/web/src/components/tenants/TenantCreateForm.tsx`
- [X] T030 [P] [US1] Tao tenant status actions component trong `frontend/apps/web/src/components/tenants/TenantStatusActions.tsx`
- [X] T031 [US1] Viet integration test cho tenant lifecycle flow trong `backend/apps/api/tests/integration/tenant_lifecycle_test.go`

**Checkpoint**: US1 hoan chinh va test doc lap

---

## Phase 4: User Story 2 - Tenant Admin Cau Hinh va Quan Ly Nguoi Dung (Priority: P1)

**Goal**: Quan ly users/groups, branding/settings, invitation/activation, CSV import

**Independent Test**: Tenant Admin cap nhat branding, import users tu CSV, gan role, tao group, user kich hoat account qua activation link

- [X] T032 [P] [US2] Tao users va activation migrations trong `backend/apps/api/migrations/002_users_and_auth.up.sql`
- [X] T033 [P] [US2] Tao groups va import jobs migrations trong `backend/apps/api/migrations/003_groups_and_imports.up.sql`
- [X] T034 [P] [US2] Tao user/auth models trong `backend/apps/api/internal/models/user.go`
- [X] T035 [P] [US2] Implement user repository trong `backend/apps/api/internal/database/user_repository.go`
- [X] T036 [P] [US2] Implement group repository trong `backend/apps/api/internal/database/group_repository.go`
- [X] T037 [US2] Implement activation token service (single-use, 24h, invalidate old) trong `backend/apps/api/internal/services/activation_service.go`
- [X] T038 [US2] Implement user service (create/invite/role assignment) trong `backend/apps/api/internal/services/user_service.go`
- [X] T039 [US2] Implement CSV import processor va row-level error reporting trong `backend/apps/api/internal/services/user_import_service.go`
- [X] T040 [US2] Implement tenant branding/settings service trong `backend/apps/api/internal/services/tenant_settings_service.go`
- [X] T041 [US2] Implement user/group/activation handlers trong `backend/apps/api/internal/handlers/user_handler.go`
- [X] T042 [US2] Implement branding/settings handlers trong `backend/apps/api/internal/handlers/tenant_settings_handler.go`
- [X] T043 [US2] Tao tenant admin users page trong `frontend/apps/web/src/app/(tenant-admin)/users/page.tsx`
- [X] T044 [P] [US2] Tao CSV import va invitation panel components trong `frontend/apps/web/src/components/users/UserImportPanel.tsx`
- [X] T045 [US2] Viet integration tests cho activation va CSV import flows trong `backend/apps/api/tests/integration/user_activation_import_test.go`

**Checkpoint**: US2 hoan chinh va test doc lap

---

## Phase 5: User Story 3 - Instructor Tao va Quan Ly Khoa Hoc (Priority: P2)

**Goal**: Tao course hierarchy, publish validation, assignment theo user/group/role

**Independent Test**: Instructor tao course day du section/subsection/unit, publish thanh cong, assign nhom hoc vien, theo doi progress

- [X] T046 [P] [US3] Tao course structure migration trong `backend/apps/api/migrations/004_courses_structure.up.sql`
- [X] T047 [P] [US3] Tao enrollment va progress migration trong `backend/apps/api/migrations/005_enrollment_progress.up.sql`
- [X] T048 [P] [US3] Tao course hierarchy models trong `backend/apps/api/internal/models/course.go`
- [X] T049 [P] [US3] Tao course enrollment model trong `backend/apps/api/internal/models/course_enrollment.go`
- [X] T050 [US3] Implement course repository trong `backend/apps/api/internal/database/course_repository.go`
- [X] T051 [US3] Implement enrollment repository trong `backend/apps/api/internal/database/enrollment_repository.go`
- [X] T052 [US3] Implement course service (draft/publish/archive + hierarchy validation) trong `backend/apps/api/internal/services/course_service.go`
- [X] T053 [US3] Implement enrollment assignment service theo user/group/role trong `backend/apps/api/internal/services/enrollment_service.go`
- [X] T054 [US3] Implement course/enrollment handlers trong `backend/apps/api/internal/handlers/course_handler.go`
- [X] T055 [US3] Tao instructor course editor page trong `frontend/apps/web/src/app/(instructor)/courses/[courseId]/edit/page.tsx`
- [X] T056 [P] [US3] Tao publish validation panel component trong `frontend/apps/web/src/components/courses/PublishValidationPanel.tsx`
- [X] T057 [US3] Viet integration test cho publish validation va enrollment assignment trong `backend/apps/api/tests/integration/course_publish_enroll_test.go`

**Checkpoint**: US3 hoan chinh va test doc lap

---

## Phase 6: User Story 4 - Student Tham Gia va Hoan Thanh Khoa Hoc (Priority: P2)

**Goal**: Hoc vien hoc course, lam quiz, autosave/recovery, timeout auto-submit, cap nhat progress

**Independent Test**: Student vao khoa hoc, hoan thanh unit, lam quiz voi autosave, mat ket noi va quay lai van tiep tuc duoc, diem duoc ghi nhan dung

- [X] T058 [P] [US4] Tao quiz attempts migration trong `backend/apps/api/migrations/006_quiz_attempts.up.sql`
- [X] T059 [P] [US4] Tao quiz models trong `backend/apps/api/internal/models/quiz.go`
- [X] T060 [US4] Implement quiz repository trong `backend/apps/api/internal/database/quiz_repository.go`
- [X] T061 [US4] Implement quiz service (start/autosave/submit/timeout/highest-valid-attempt) trong `backend/apps/api/internal/services/quiz_service.go`
- [X] T062 [US4] Implement progress aggregation service cho student trong `backend/apps/api/internal/services/progress_service.go`
- [X] T063 [US4] Implement quiz handlers trong `backend/apps/api/internal/handlers/quiz_handler.go`
- [X] T064 [US4] Implement student me handlers (my courses/progress) trong `backend/apps/api/internal/handlers/student_handler.go`
- [X] T065 [US4] Tao student my-courses page trong `frontend/apps/web/src/app/(student)/courses/page.tsx`
- [X] T066 [US4] Tao quiz attempt page co autosave timer trong `frontend/apps/web/src/app/(student)/quizzes/[quizId]/attempt/page.tsx`
- [X] T067 [P] [US4] Tao progress tracker component trong `frontend/apps/web/src/components/student/ProgressTracker.tsx`
- [X] T068 [US4] Viet integration tests cho autosave/recovery/timeout auto-submit trong `backend/apps/api/tests/integration/quiz_recovery_test.go`
- [X] T069 [US4] Viet frontend e2e test cho student learning flow trong `frontend/apps/web/tests/e2e/student_learning_flow.spec.ts`

**Checkpoint**: US4 hoan chinh va test doc lap

---

## Phase 7: User Story 5 - Super Admin Impersonate Tenant de Ho Tro (Priority: P2)

**Goal**: Ho tro tenant qua impersonation mode, co audit day du, chan destructive actions

**Independent Test**: Super Admin bat impersonation, thuc hien thao tac read/config, khong the xoa du lieu, audit log luu du thong tin phien va hanh dong

- [ ] T070 [P] [US5] Tao impersonation migrations trong `backend/apps/api/migrations/007_impersonation_audit.up.sql`
- [ ] T071 [P] [US5] Tao impersonation models trong `backend/apps/api/internal/models/impersonation_session.go`
- [ ] T072 [US5] Implement impersonation repository trong `backend/apps/api/internal/database/impersonation_repository.go`
- [ ] T073 [US5] Implement impersonation service (start/end/action tracking) trong `backend/apps/api/internal/services/impersonation_service.go`
- [ ] T074 [US5] Implement destructive action guard middleware trong `backend/apps/api/internal/middleware/impersonation_guard.go`
- [ ] T075 [US5] Implement impersonation handlers va audit endpoint trong `backend/apps/api/internal/handlers/impersonation_handler.go`
- [ ] T076 [US5] Tao impersonation launcher component cho Super Admin trong `frontend/apps/web/src/components/impersonation/ImpersonationLauncher.tsx`
- [ ] T077 [US5] Tao tenant admin impersonation audit page trong `frontend/apps/web/src/app/(tenant-admin)/audit/impersonation/page.tsx`
- [ ] T078 [US5] Viet integration test cho destructive action blocking trong `backend/apps/api/tests/integration/impersonation_guard_test.go`
- [ ] T079 [US5] Viet integration test dam bao Super Admin khong xuat hien trong tenant learner projections trong `backend/apps/api/tests/integration/impersonation_visibility_test.go`

**Checkpoint**: US5 hoan chinh va test doc lap

---

## Phase 8: User Story 6 - Tenant Admin Xem Bao Cao va Dashboard (Priority: P3)

**Goal**: Dashboard cho Tenant Admin va report theo course cho Instructor

**Independent Test**: Dashboard hien thi dung active users/active courses/completion; Instructor xem report khoa hoc khop du lieu thuc te

- [ ] T080 [P] [US6] Tao usage metrics migration trong `backend/apps/api/migrations/008_usage_metrics_daily.up.sql`
- [ ] T081 [P] [US6] Tao usage metric model trong `backend/apps/api/internal/models/usage_metric_daily.go`
- [ ] T082 [US6] Implement usage metrics repository trong `backend/apps/api/internal/database/usage_metrics_repository.go`
- [ ] T083 [US6] Implement reporting service cho tenant dashboard va course report trong `backend/apps/api/internal/services/reporting_service.go`
- [ ] T084 [US6] Implement reporting handlers trong `backend/apps/api/internal/handlers/report_handler.go`
- [ ] T085 [US6] Tao tenant admin dashboard page trong `frontend/apps/web/src/app/(tenant-admin)/dashboard/page.tsx`
- [ ] T086 [US6] Tao instructor course report page trong `frontend/apps/web/src/app/(instructor)/reports/courses/[courseId]/page.tsx`
- [ ] T087 [US6] Viet integration test cho reporting accuracy voi seeded data trong `backend/apps/api/tests/integration/reporting_accuracy_test.go`
- [ ] T088 [US6] Viet performance test dashboard latency cho 10k users trong `backend/apps/api/tests/performance/dashboard_p95_test.go`

**Checkpoint**: US6 hoan chinh va test doc lap

---

## Phase 9: User Story 7 - Super Admin Theo Doi Usage va Quan Ly Billing (Priority: P3)

**Goal**: Quan ly plans/custom plans, enforce quota, canh bao nguong 80/90, cau hinh SSO entitlement, JIT provisioning

**Independent Test**: Super Admin tao custom plan, gan tenant, verify quota va entitlement co hieu luc ngay; SSO JIT tao/linked user dung quy tac uniqueness

- [ ] T089 [P] [US7] Tao plan/quota/sso core migration trong `backend/apps/api/migrations/009_plan_quota_sso.up.sql`
- [ ] T090 [P] [US7] Tao user identity links migration trong `backend/apps/api/migrations/010_user_identity_links.up.sql`
- [ ] T091 [P] [US7] Tao billing va SSO models trong `backend/apps/api/internal/models/tenant_sso_config.go`
- [ ] T092 [US7] Implement plan repository trong `backend/apps/api/internal/database/plan_repository.go`
- [ ] T093 [US7] Implement plan service (custom plan + assignment) trong `backend/apps/api/internal/services/plan_service.go`
- [ ] T094 [US7] Implement quota service (enforce/warn 80-90/downgrade behavior) trong `backend/apps/api/internal/services/quota_service.go`
- [ ] T095 [US7] Implement SSO entitlement/config service trong `backend/apps/api/internal/services/sso_service.go`
- [ ] T096 [US7] Implement SSO JIT provisioning va linked-existing logic trong `backend/apps/api/internal/services/sso_jit_service.go`
- [ ] T097 [US7] Implement billing handlers (plans/quota/sso) trong `backend/apps/api/internal/handlers/billing_handler.go`
- [ ] T098 [US7] Tao super admin billing dashboard page trong `frontend/apps/web/src/app/(super-admin)/billing/page.tsx`
- [ ] T099 [P] [US7] Tao custom plan form component trong `frontend/apps/web/src/components/plans/CustomPlanForm.tsx`
- [ ] T100 [P] [US7] Tao tenant SSO settings page trong `frontend/apps/web/src/app/(tenant-admin)/settings/sso/page.tsx`
- [ ] T101 [US7] Viet integration tests cho quota threshold va SSO mapping uniqueness trong `backend/apps/api/tests/integration/billing_sso_test.go`

**Checkpoint**: US7 hoan chinh va test doc lap

---

## Phase 10: User Story 8 - Thong Bao Multi-Tenant (Priority: P3)

**Goal**: Thong bao in-tenant va global (in-app/email), projection trang thai delivery failures cho admin

**Independent Test**: Assign course sinh thong bao cho student; Super Admin gui global banner cho tat ca tenants; email failures hien thi dung sau 3 retries

- [ ] T102 [P] [US8] Tao notifications migrations trong `backend/apps/api/migrations/011_notifications.up.sql`
- [ ] T103 [P] [US8] Tao notification models trong `backend/apps/api/internal/models/notification.go`
- [ ] T104 [US8] Implement notification repository trong `backend/apps/api/internal/database/notification_repository.go`
- [ ] T105 [US8] Implement notification service voi audience expansion trong `backend/apps/api/internal/services/notification_service.go`
- [ ] T106 [US8] Implement notification delivery worker va retry policy trong `backend/apps/api/cmd/workers/notification_worker.go`
- [ ] T107 [US8] Implement notification handlers (tenant/global/list/read) trong `backend/apps/api/internal/handlers/notification_handler.go`
- [ ] T108 [US8] Integrate course-assigned va quiz-deadline triggers vao notifier trong `backend/apps/api/internal/services/notification_triggers.go`
- [ ] T109 [US8] Tao notification center component trong `frontend/apps/web/src/components/notifications/NotificationCenter.tsx`
- [ ] T110 [P] [US8] Tao global banner component trong `frontend/apps/web/src/components/notifications/GlobalBanner.tsx`
- [ ] T111 [US8] Tao tenant notification compose page trong `frontend/apps/web/src/app/(tenant-admin)/notifications/page.tsx`
- [ ] T112 [US8] Viet integration tests cho notification scope va delivery failure projection trong `backend/apps/api/tests/integration/notification_scope_delivery_test.go`

**Checkpoint**: Tat ca user stories hoat dong doc lap

---

## Phase 11: Polish & Cross-Cutting Concerns

**Purpose**: Cac cong viec lien story de harden va release

- [ ] T113 [P] Chay full OpenAPI validation va regenerate server/client types tu `specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml`
- [ ] T114 [P] Them Schemathesis regression config trong `backend/apps/api/tests/contract/schemathesis.yaml`
- [ ] T115 [P] Them tenant isolation matrix regression tests trong `backend/apps/api/tests/integration/tenant_isolation_matrix_test.go`
- [ ] T116 [P] Implement retention purge worker cho deactivated tenants trong `backend/apps/api/cmd/workers/retention_purge_worker.go`
- [ ] T117 [P] Implement GDPR data export service trong `backend/apps/api/internal/services/gdpr_export_service.go`
- [ ] T118 [P] Them monitoring alerts config cho quota/rate-limit/worker failures trong `dev/monitoring/lms-alerts.yaml`
- [ ] T119 [P] Harden security headers va anti-abuse response policies trong `backend/apps/api/internal/middleware/security_headers.go`
- [ ] T120 [P] Them frontend accessibility smoke tests cho core pages trong `frontend/apps/web/tests/e2e/accessibility_smoke.spec.ts`
- [ ] T121 [P] Cap nhat quickstart verification checklist trong `specs/001-lms-saas-system/quickstart.md`
- [ ] T122 Chay full smoke checklist va ghi ket qua trong `specs/001-lms-saas-system/checklists/requirements.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: khong phu thuoc, bat dau ngay
- **Phase 2 (Foundational)**: phu thuoc Phase 1, block toan bo user stories
- **Phase 3-10 (User Stories)**: deu phu thuoc Phase 2
- **Phase 11 (Polish)**: phu thuoc stories can release da hoan tat

### User Story Dependencies

- **US1 (P1)**: bat dau ngay sau Foundational
- **US2 (P1)**: bat dau ngay sau Foundational
- **US3 (P2)**: can US2 du user/group cho enrollment test
- **US4 (P2)**: can US3 du course hierarchy de hoc/quiz
- **US5 (P2)**: co the doc lap sau Foundational
- **US6 (P3)**: can data tu US2-US4 de report co y nghia
- **US7 (P3)**: co the bat dau sau Foundational, tich hop voi US1/US2
- **US8 (P3)**: tich hop su kien tu US2-US4

### Suggested Story Completion Order

1. US1 + US2 (song song)
2. US3 + US5 + US7 (song song co dieu kien)
3. US4 + US6 + US8
4. Polish phase

---

## Parallel Execution Examples Per Story

### US1

- T022 + T023 + T029 + T030 co the chay song song

### US2

- T034 + T035 + T036 chay song song
- T043 + T044 chay song song sau khi APIs on dinh

### US3

- T048 + T049 chay song song
- T055 + T056 chay song song

### US4

- T059 + T060 chay song song
- T065 + T067 chay song song

### US5

- T071 + T072 chay song song
- T076 + T077 chay song song

### US6

- T081 + T082 chay song song
- T085 + T086 chay song song

### US7

- T091 + T092 chay song song
- T099 + T100 chay song song

### US8

- T103 + T104 chay song song
- T109 + T110 chay song song

---

## Implementation Strategy

### MVP First

1. Hoan tat Phase 1 + Phase 2
2. Hoan tat US1
3. Hoan tat US2
4. Validate onboarding tenant va user activation end-to-end

### Incremental Delivery

1. Release increment A: US1 + US2
2. Release increment B: US3 + US4
3. Release increment C: US5 + US6
4. Release increment D: US7 + US8
5. Release hardening: Phase 11

### Quality Gates

- Moi story phai pass integration tests truoc khi merge
- Contract changes phai re-validate bang swagger-cli + Schemathesis
- Khong merge neu vi pham tenant isolation, lockout/rate-limit, hoac audit requirements
