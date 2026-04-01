# Hợp Đồng Domain Events

Tài liệu này mô tả các hợp đồng tích hợp bất đồng bộ (asynchronous integration contracts) bổ trợ cho REST API contract.

## Transport và Event Envelope

- Transport: queue/stream bền vững (triển khai có thể là Redis Streams, RabbitMQ hoặc Kafka).
- Delivery: at-least-once.
- Ordering: đảm bảo theo aggregate key khi phù hợp (`tenant_id`, `course_id`, `quiz_attempt_id`).
- Envelope fields (bắt buộc cho mọi event):

| Field | Type | Description |
|---|---|---|
| event_id | UUID | Định danh duy nhất của event |
| event_type | string | Tên event có version |
| occurred_at | RFC3339 timestamp | Thời điểm event được tạo |
| tenant_id | UUID hoặc null | Phạm vi tenant, null cho event cấp hệ thống |
| producer | string | Tên service phát event |
| payload | object | Dữ liệu nghiệp vụ của event |

## Event: tenant.created.v1

Producer: tenant service  
Consumers: notification service, billing service, audit service

Payload:
- tenant_id (UUID)
- subdomain (string)
- plan_code (string)
- owner_user_id (UUID)

## Event: tenant.status_changed.v1

Producer: tenant service  
Consumers: auth service, reporting service, audit service

Payload:
- tenant_id (UUID)
- previous_status (`active|suspended|deactivated`)
- new_status (`active|suspended|deactivated`)
- changed_by (UUID)
- reason (string, optional)

Ghi chú:
- Nếu `new_status` là `deactivated`, producer phải gửi thêm `purge_due_at`.

## Event: users.import_job_completed.v1

Producer: user service  
Consumers: notification service, reporting service

Payload:
- tenant_id (UUID)
- job_id (UUID)
- total_rows (integer)
- success_count (integer)
- failed_count (integer)
- failure_report_ref (string, optional)

## Event: course.published.v1

Producer: course service  
Consumers: enrollment service, notification service, reporting service

Payload:
- tenant_id (UUID)
- course_id (UUID)
- published_by (UUID)
- assignee_scope (object: direct/group/role targets)

## Event: quiz.autosaved.v1

Producer: quiz service  
Consumers: analytics service

Payload:
- tenant_id (UUID)
- quiz_attempt_id (UUID)
- student_id (UUID)
- saved_at (RFC3339)
- answer_count (integer)

## Event: quiz.auto_submitted.v1

Producer: quiz service  
Consumers: grading service, notification service, analytics service

Payload:
- tenant_id (UUID)
- quiz_attempt_id (UUID)
- student_id (UUID)
- submitted_at (RFC3339)
- counted_toward_attempt_limit (boolean, bắt buộc true)

## Event: quota.threshold_reached.v1

Producer: billing/usage service  
Consumers: notification service

Payload:
- tenant_id (UUID)
- metric (`users|courses|storage_bytes`)
- threshold_percent (integer, cho phép: 80, 90)
- used_value (number)
- limit_value (number)

## Event: email.delivery_failed.v1

Producer: notification service  
Consumers: admin UI projection, audit service

Payload:
- tenant_id (UUID)
- message_type (`invitation|notification|quota_warning`)
- recipient (string email)
- attempts (integer, bắt buộc bằng 3)
- terminal_status (`failed`)
- error_code (string)

## Event: file.uploaded.v1

Producer: storage ingestion service  
Consumers: media worker, malware scan worker

Payload:
- tenant_id (UUID)
- object_key (string)
- bucket (string)
- size_bytes (integer)
- content_type (string)
- uploaded_by (UUID)

## Event: sso.user_provisioned.v1

Producer: auth service  
Consumers: audit service, notification service

Payload:
- tenant_id (UUID)
- user_id (UUID)
- idp_issuer (string)
- idp_subject (string)
- provision_mode (`jit_created|linked_existing`)

## Event: tenant.data_purged.v1

Producer: retention worker  
Consumers: audit service, billing service

Payload:
- tenant_id (UUID)
- purged_at (RFC3339)
- purge_window_days (integer, kỳ vọng 90)

## Quy Tắc Tương Thích (Compatibility Rules)

- Version event dùng hậu tố `.vN`.
- Mọi thay đổi breaking payload bắt buộc tạo event name/version mới.
- Consumer phải bỏ qua unknown fields để bảo đảm forward compatibility.
- Producer phải publish idempotency key = `event_id`; các lần retry không được làm thay đổi payload.