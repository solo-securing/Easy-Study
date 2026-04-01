# Mô Hình Dữ Liệu: LMS SaaS Multi-Tenant

## Quy Tắc Mô Hình Hóa

- Mọi bảng thuộc tenant phải có `tenant_id UUID NOT NULL`.
- Mọi unique constraint trong phạm vi tenant phải là composite với `tenant_id`.
- Cross-tenant reads/writes bị cấm ở service boundary và được xác minh bằng integration tests.
- Soft-delete áp dụng cho thực thể có thể phục hồi ở tầng người dùng; hard-delete áp dụng cho job xóa theo retention.
- Mọi thay đổi cần audit phải ghi immutable entries vào `audit_logs`.
- Time-window counters cho abuse prevention được lưu trên Redis dưới dạng dữ liệu vận hành tạm thời (ephemeral).

## Thực Thể Cốt Lõi

### Tenant

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| name | TEXT | NOT NULL |
| subdomain | TEXT | NOT NULL, globally unique, not reserved |
| status | ENUM(active,suspended,deactivated,purged) | NOT NULL, default `active` |
| default_locale | TEXT | NOT NULL, default `en` |
| timezone | TEXT | NOT NULL |
| plan_id | UUID | FK -> plans.id (denormalized active plan pointer) |
| deactivated_at | TIMESTAMPTZ | NULL |
| purge_due_at | TIMESTAMPTZ | NULL (deactivated_at + 90 days) |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- `subdomain` phải khớp DNS-safe pattern và không nằm trong reserved list.
- `purge_due_at` là bắt buộc khi status là `deactivated`.
- `plan_id` luôn phải khớp với `tenant_plan_assignments.plan_id` đang active.

Chuyển trạng thái:
- `active -> suspended`
- `suspended -> active`
- `active|suspended -> deactivated`
- `deactivated -> purged` (background retention job sau 90 ngày)

### ReservedSubdomain

| Field | Type | Constraints |
|---|---|---|
| subdomain | TEXT | PK |
| reason | TEXT | NOT NULL |
| created_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- `tenants.subdomain` không được tồn tại trong `reserved_subdomains`.

### Plan

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| code | TEXT | UNIQUE (`free`,`trial`,`pro`,`enterprise`,`custom-*`) |
| name | TEXT | NOT NULL |
| user_limit | INTEGER | >= 1 |
| course_limit | INTEGER | >= 1 |
| storage_limit_gb | INTEGER | >= 1 |
| entitlements | JSONB | NOT NULL (ví dụ `{"sso": true}`) |
| is_custom | BOOLEAN | NOT NULL |
| created_by | UUID | nullable (super admin actor) |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- `entitlements.sso` phải là kiểu boolean.
- Custom plans phải có immutable audit trail.

### TenantPlanAssignment

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| plan_id | UUID | FK -> plans.id |
| effective_from | TIMESTAMPTZ | NOT NULL |
| changed_by | UUID | FK -> users.id (super admin) |
| reason | TEXT | NULL |

Quy tắc validation:
- Mỗi tenant chỉ có đúng một active assignment tại một thời điểm.

### TenantBranding

| Field | Type | Constraints |
|---|---|---|
| tenant_id | UUID | PK, FK -> tenants.id |
| display_name | TEXT | NOT NULL |
| logo_asset_key | TEXT | NULL |
| favicon_asset_key | TEXT | NULL |
| primary_color_hex | TEXT | NOT NULL |
| updated_by | UUID | FK -> users.id |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- `primary_color_hex` phải là mã màu hex hợp lệ.

### TenantSettings

| Field | Type | Constraints |
|---|---|---|
| tenant_id | UUID | PK, FK -> tenants.id |
| default_locale | TEXT | NOT NULL |
| timezone | TEXT | NOT NULL |
| session_ttl_minutes | INTEGER | NOT NULL, range 15..1440 |
| password_policy_json | JSONB | NOT NULL |
| privacy_policy_url | TEXT | NULL |
| allow_student_self_registration | BOOLEAN | NOT NULL default false |
| updated_by | UUID | FK -> users.id |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Phải được khởi tạo bằng safe defaults ngay khi tenant được tạo.

### TenantSsoConfig

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | UNIQUE, FK -> tenants.id |
| provider | ENUM(google_oidc,oidc,saml) | NOT NULL |
| issuer | TEXT | NOT NULL |
| client_id | TEXT | NOT NULL |
| client_secret_ref | TEXT | NOT NULL |
| redirect_uri | TEXT | NOT NULL |
| enabled | BOOLEAN | NOT NULL default false |
| configured_by | UUID | FK -> users.id |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- SSO config chỉ được enabled khi tenant có entitlement `sso=true`.

### User

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id, NULL for platform super admins only |
| email | CITEXT | NOT NULL |
| full_name | TEXT | NOT NULL |
| role | ENUM(super_admin,tenant_admin,instructor,student) | NOT NULL |
| status | ENUM(invited,active,suspended,deactivated) | NOT NULL |
| password_hash | TEXT | NULL for SSO-only accounts |
| failed_login_count | INTEGER | NOT NULL default 0 |
| locked_until | TIMESTAMPTZ | NULL |
| email_verified_at | TIMESTAMPTZ | NULL |
| last_login_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Unique email scope: `(tenant_id, email)`.
- Instructor không được self-register (chỉ qua invitation).

### ActivationToken

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| user_id | UUID | FK -> users.id |
| token_hash | TEXT | NOT NULL |
| token_version | INTEGER | NOT NULL |
| expires_at | TIMESTAMPTZ | NOT NULL |
| consumed_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Token là single-use.
- Cấp token mới phải tăng `token_version` và vô hiệu hóa mọi token cũ.
- Thời gian hết hạn cố định 24 giờ.

### RefreshTokenSession

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| user_id | UUID | FK -> users.id |
| token_hash | TEXT | NOT NULL |
| issued_at | TIMESTAMPTZ | NOT NULL |
| expires_at | TIMESTAMPTZ | NOT NULL |
| revoked_at | TIMESTAMPTZ | NULL |
| replaced_by_session_id | UUID | FK -> refresh_token_sessions.id, NULL |
| user_agent | TEXT | NULL |
| ip_address | INET | NULL |

Quy tắc validation:
- Token hash phải unique trong tập non-revoked sessions.
- Refresh rotation phải tạo row mới và revoke row cũ theo cơ chế atomic.

### OutboundMessage

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id, NULL for platform messages |
| message_type | ENUM(invitation,notification,quota_warning) | NOT NULL |
| recipient | TEXT | NOT NULL |
| payload_json | JSONB | NOT NULL |
| status | ENUM(pending,sent,failed) | NOT NULL default `pending` |
| attempts | INTEGER | NOT NULL default 0 |
| next_retry_at | TIMESTAMPTZ | NULL |
| last_error_code | TEXT | NULL |
| sent_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Giới hạn retry tối đa 3 lần; sau đó status bắt buộc là `failed`.

### UserIdentityLink

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| user_id | UUID | FK -> users.id |
| idp_issuer | TEXT | NOT NULL |
| idp_subject | TEXT | NOT NULL |
| linked_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Unique `(tenant_id, idp_issuer, idp_subject)`.
- Mapping đã tồn tại là nguồn chân lý ngay cả khi email claim thay đổi.

### Group

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| name | TEXT | NOT NULL |
| description | TEXT | NULL |
| created_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Unique `(tenant_id, name)`.

### GroupMember

| Field | Type | Constraints |
|---|---|---|
| group_id | UUID | FK -> groups.id |
| user_id | UUID | FK -> users.id |
| tenant_id | UUID | FK -> tenants.id |
| created_at | TIMESTAMPTZ | NOT NULL |

Khóa chính tổng hợp (Composite PK): `(group_id, user_id)`.

### UserImportJob

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| created_by | UUID | FK -> users.id |
| file_asset_ref | TEXT | NOT NULL |
| default_role | ENUM(student,instructor) | NOT NULL |
| status | ENUM(queued,processing,completed,failed) | NOT NULL |
| total_rows | INTEGER | NOT NULL default 0 |
| success_count | INTEGER | NOT NULL default 0 |
| failed_count | INTEGER | NOT NULL default 0 |
| started_at | TIMESTAMPTZ | NULL |
| finished_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL |

### UserImportJobError

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| job_id | UUID | FK -> user_import_jobs.id |
| row_number | INTEGER | NOT NULL |
| error_code | TEXT | NOT NULL |
| error_message | TEXT | NOT NULL |

Quy tắc validation:
- Unique `(job_id, row_number)`.

### Course

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| created_by | UUID | FK -> users.id |
| title | TEXT | NOT NULL |
| description | TEXT | NULL |
| status | ENUM(draft,published,archived) | NOT NULL |
| published_at | TIMESTAMPTZ | NULL |
| archived_at | TIMESTAMPTZ | NULL |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Publish bị chặn nếu bất kỳ section/subsection/unit nào chưa đầy đủ.

Chuyển trạng thái:
- `draft -> published`
- `published -> archived`
- `archived -> draft` (luồng admin tùy chọn, có audit)

### Section

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| course_id | UUID | FK -> courses.id |
| title | TEXT | NOT NULL |
| sort_order | INTEGER | >= 0 |

### Subsection

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| section_id | UUID | FK -> sections.id |
| title | TEXT | NOT NULL |
| sort_order | INTEGER | >= 0 |

### Unit

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| subsection_id | UUID | FK -> subsections.id |
| type | ENUM(video,document,quiz) | NOT NULL |
| title | TEXT | NOT NULL |
| content_ref | TEXT | NULL |
| sort_order | INTEGER | >= 0 |

Quy tắc validation:
- `content_ref` là bắt buộc cho unit type `video` và `document`.
- Unit type `quiz` phải có đúng một row liên kết trong `quizzes` qua `quizzes.unit_id`.

### Quiz

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| unit_id | UUID | FK -> units.id UNIQUE |
| title | TEXT | NOT NULL |
| max_score | NUMERIC(6,2) | >= 0 |
| passing_score | NUMERIC(6,2) | >= 0 and <= max_score |
| time_limit_minutes | INTEGER | >= 1 |
| attempt_limit | INTEGER | >= 1 |
| shuffle_questions | BOOLEAN | NOT NULL default false |
| open_at | TIMESTAMPTZ | NULL |
| close_at | TIMESTAMPTZ | NULL |
| scoring_policy | ENUM(highest_valid_attempt) | NOT NULL |

### QuizAttempt

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| quiz_id | UUID | FK -> quizzes.id |
| student_id | UUID | FK -> users.id |
| attempt_no | INTEGER | >= 1 |
| status | ENUM(in_progress,submitted,auto_submitted,expired) | NOT NULL |
| started_at | TIMESTAMPTZ | NOT NULL |
| submitted_at | TIMESTAMPTZ | NULL |
| score | NUMERIC(6,2) | NULL |
| answers_snapshot_ref | TEXT | NULL |
| counted_toward_limit | BOOLEAN | NOT NULL default true |
| created_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Unique `(tenant_id, quiz_id, student_id, attempt_no)`.
- `auto_submitted` attempts được tính vào attempt limit.

Chuyển trạng thái:
- `in_progress -> submitted`
- `in_progress -> auto_submitted`
- `in_progress -> expired` (khi không thể submit theo policy)

### QuizAttemptSnapshot

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| quiz_attempt_id | UUID | FK -> quiz_attempts.id |
| snapshot_no | INTEGER | NOT NULL |
| answers_payload_ref | TEXT | NOT NULL |
| save_source | ENUM(auto,manual,submit,timeout) | NOT NULL |
| saved_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Unique `(quiz_attempt_id, snapshot_no)` và `snapshot_no` tăng đơn điệu.

### CourseEnrollment

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| course_id | UUID | FK -> courses.id |
| user_id | UUID | FK -> users.id |
| source_type | ENUM(direct,group,role) | NOT NULL |
| assigned_by | UUID | FK -> users.id |
| assigned_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Unique `(tenant_id, course_id, user_id)`.

### ProgressRecord

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| course_id | UUID | FK -> courses.id |
| student_id | UUID | FK -> users.id |
| completion_percent | NUMERIC(5,2) | >= 0 and <= 100 |
| completed_unit_count | INTEGER | >= 0 |
| official_quiz_score | NUMERIC(6,2) | NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Official quiz score được suy ra từ highest valid attempt.

### Notification

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id, NULL for global |
| type | ENUM(course_assigned,deadline_reminder,tenant_announcement,global_banner,quota_warning) | NOT NULL |
| channel | ENUM(in_app,email,both) | NOT NULL |
| title | TEXT | NOT NULL |
| body | TEXT | NOT NULL |
| published_by | UUID | FK -> users.id |
| published_at | TIMESTAMPTZ | NOT NULL |

### NotificationAudience

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id, NULL for global |
| notification_id | UUID | FK -> notifications.id |
| audience_type | ENUM(all_users,user_ids,group_ids,roles) | NOT NULL |
| audience_payload_json | JSONB | NOT NULL |

Quy tắc validation:
- Với `all_users`, payload là object rỗng.
- Với `user_ids`, `group_ids`, hoặc `roles`, payload phải chứa target arrays không rỗng.

### NotificationRecipient

| Field | Type | Constraints |
|---|---|---|
| notification_id | UUID | FK -> notifications.id |
| tenant_id | UUID | FK -> tenants.id |
| user_id | UUID | FK -> users.id |
| read_at | TIMESTAMPTZ | NULL |
| email_status | ENUM(pending,sent,failed) | NOT NULL default `pending` |
| email_attempts | INTEGER | NOT NULL default 0 |

Khóa chính tổng hợp (Composite PK): `(notification_id, user_id)`.

### QuotaThresholdAlert

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| metric | ENUM(users,courses,storage_bytes) | NOT NULL |
| threshold_percent | ENUM(80,90) | NOT NULL |
| plan_assignment_id | UUID | FK -> tenant_plan_assignments.id |
| triggered_at | TIMESTAMPTZ | NOT NULL |
| notification_id | UUID | FK -> notifications.id |

Quy tắc validation:
- Unique `(tenant_id, metric, threshold_percent, plan_assignment_id)`.

### ImpersonationSession

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| super_admin_id | UUID | FK -> users.id |
| started_at | TIMESTAMPTZ | NOT NULL |
| ended_at | TIMESTAMPTZ | NULL |
| status | ENUM(active,ended,force_ended) | NOT NULL |

Quy tắc validation:
- Mỗi cặp `(tenant_id, super_admin_id)` chỉ có đúng một active impersonation session.
- Các destructive actions phải bị chặn khi session đang active.

### ImpersonationAction

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| session_id | UUID | FK -> impersonation_sessions.id |
| action | TEXT | NOT NULL |
| resource_type | TEXT | NOT NULL |
| resource_id | UUID | NULL |
| happened_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Chỉ cho phép các action type không phá hủy dữ liệu.

### AuditLog

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id, NULL for global actions |
| actor_id | UUID | FK -> users.id |
| actor_role | TEXT | NOT NULL |
| action | TEXT | NOT NULL |
| context | JSONB | NOT NULL |
| happened_at | TIMESTAMPTZ | NOT NULL |
| immutable_hash | TEXT | NOT NULL |

Quy tắc validation:
- Áp dụng append-only write pattern.

### UsageMetricDaily

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK -> tenants.id |
| day | DATE | NOT NULL |
| active_user_count | INTEGER | >= 0 |
| course_count | INTEGER | >= 0 |
| login_count | INTEGER | >= 0 |
| learning_minutes | INTEGER | >= 0 |
| storage_used_bytes | BIGINT | >= 0 |

Quy tắc validation:
- Unique `(tenant_id, day)`.

### PublicRateLimitCounter (Redis logical model)

| Field | Type | Constraints |
|---|---|---|
| bucket_key | TEXT | PK logical key |
| window_started_at | TIMESTAMPTZ | NOT NULL |
| request_count | INTEGER | NOT NULL |
| expires_at | TIMESTAMPTZ | NOT NULL |

Quy tắc validation:
- Key pattern nên bao gồm route scope + client IP.
- Cửa sổ counter thực thi giới hạn 100 requests/phút/IP cho public APIs.

## Tóm Tắt Quan Hệ (Relationships Summary)

- Tenant 1..N Users, Groups, Courses, Notifications, UsageMetricDaily, UserImportJobs, OutboundMessages, RefreshTokenSessions.
- Tenant 1..1 TenantBranding, TenantSettings, optional TenantSsoConfig.
- Plan 1..N TenantPlanAssignment; Tenant có một active plan assignment.
- Course 1..N Sections; Section 1..N Subsections; Subsection 1..N Units.
- Unit (quiz type) 1..1 Quiz; Quiz 1..N QuizAttempts.
- QuizAttempt 1..N QuizAttemptSnapshots.
- User N..N Group qua GroupMember.
- Course N..N User qua CourseEnrollment.
- Notification N..N User qua NotificationRecipient.
- Notification 1..N NotificationAudience targets cho campaign recipient expansion.
- OutboundMessage là delivery outbox cho các luồng email invitation/notification/quota-warning.
- QuotaThresholdAlert tham chiếu notification records cho vòng đời cảnh báo 80/90%.
- ImpersonationSession 1..N ImpersonationActions; mọi privileged mutations đều ghi AuditLog.

## Bất Biến Liên Thực Thể (Cross-Entity Invariants)

- Mọi truy vấn trong phạm vi tenant bắt buộc có predicate `tenant_id`.
- User role checks là bắt buộc tại API boundary trước khi chạy business logic.
- Tenant subdomain phải globally unique và không được tồn tại trong `reserved_subdomains`.
- Course publish yêu cầu hierarchy đầy đủ và phải pass validation.
- Quota checks chạy trước khi tạo user/course và trong suốt quá trình bulk import.
- Cảnh báo quota chỉ được trigger tại đúng các ngưỡng 80% và 90%, đồng thời lưu qua `quota_threshold_alerts`.
- Email invitation, notification và quota-warning phải đi qua `outbound_messages` với tối đa 3 retries.
- Refresh token rotation phải revoke phiên cũ và duy trì replacement chain có thể audit.
- `notification_recipients.email_status` là read projection suy ra từ delivery state của `outbound_messages`.
- Ở impersonation mode, actor không được materialize thành learner/enrollment/progress subject trong tenant data.
- SSO JIT mapping phải luôn unique theo `(tenant_id, idp_issuer, idp_subject)` và update user hiện có khi email claim thay đổi.
- Public API abuse prevention phải enforce cả account lockout lẫn Redis-backed IP rate limits.
- Dữ liệu của deactivated tenant ở chế độ read-only cho tới khi purge job đạt `purge_due_at`.