# Nghiên Cứu: LMS SaaS Multi-Tenant (Phase 0)

Date: 2026-04-01  
Feature: 001-lms-saas-system

Toàn bộ các điểm cần làm rõ trong technical context đã được giải quyết cho phase planning này. Hiện không còn mục "NEEDS CLARIFICATION" đang mở.

## Quyết Định 1: Shared-schema tenancy với tenant scoping nghiêm ngặt

Quyết định: Dùng shared PostgreSQL schema và bắt buộc có `tenant_id` trên mọi bảng thuộc tenant, được cưỡng chế bằng tenant filter ở repository layer và DB policy tùy chọn.

Lý do:
- Phù hợp mục tiêu benchmark v1 (100 active tenants) với chi phí vận hành thấp hơn schema-per-tenant hoặc database-per-tenant.
- Đơn giản hóa migration và quy trình deploy cho nhóm SaaS quy mô nhỏ/trung bình.
- Vẫn hỗ trợ platform analytics xuyên tenant cho Super Admin nhưng giữ ownership theo từng row.

Phương án thay thế đã cân nhắc:
- Database-per-tenant: cách ly mạnh nhất nhưng ops phức tạp và chi phí cao.
- Schema-per-tenant: cách ly vừa phải nhưng migration fan-out phức tạp, CI pipeline khó hơn.

## Quyết Định 2: API-first contract governance với OpenAPI 3.1

Quyết định: Dùng OpenAPI 3.1 làm source of truth và bắt buộc vòng đời contract trong CI (Spectral lint -> codegen -> contract tests).

Lý do:
- Phù hợp Constitution Principle VII (API-First Design).
- Tránh frontend/backend drift nhờ generate artifacts cho cả Go và TypeScript.
- Cho phép tự động hóa negative/fuzz coverage qua Schemathesis.

Phương án thay thế đã cân nhắc:
- Code-first API framework: coding nhanh lúc đầu nhưng yếu về kiểm soát hợp đồng giữa các nhóm.
- Viết tay SDK types: rủi ro drift cao và trùng lặp công sức.

## Quyết Định 3: Mô hình xác thực và SSO theo entitlement

Quyết định: Dùng JWT access token (1 giờ) + rotating refresh token cho local auth; chỉ bật OAuth2/SSO khi tenant có entitlement SSO.

Lý do:
- Đáp ứng yêu cầu bảo mật và hành vi entitlement ở FR-043/FR-044.
- Tách rõ local auth khỏi enterprise SSO features.
- Hỗ trợ JIT provisioning với mapping key xác định duy nhất `(tenant_id, idp_issuer, idp_subject)`.

Phương án thay thế đã cân nhắc:
- Session-only auth: dễ revoke hơn nhưng kém phù hợp cho distributed API clients.
- Luôn bật SSO cho mọi plan: không đúng yêu cầu đóng gói sản phẩm.

## Quyết Định 4: Vòng đời activation link cho invitation và self-registration

Quyết định: Dùng single-use activation token, hết hạn sau 24 giờ; khi cấp token mới sẽ vô hiệu toàn bộ token cũ chưa dùng của cùng user.

Lý do:
- Đáp ứng trực tiếp FR-016.
- Giảm rủi ro takeover từ link cũ còn hiệu lực.
- Hành vi resend rõ ràng, dễ hỗ trợ vận hành.

Phương án thay thế đã cân nhắc:
- Multi-use activation link: triển khai đơn giản hơn nhưng bảo mật yếu.
- Tăng thời gian hết hạn: tiện hơn nhưng mở rộng cửa sổ tấn công.

## Quyết Định 5: Hành vi quiz attempt, autosave và timeout

Quyết định: Lưu quiz attempt snapshots theo chu kỳ cố định (ví dụ 15 giây) và tại các mốc quan trọng; khi hết giờ thì auto-submit đáp án đã lưu gần nhất và đánh dấu attempt là counted.

Lý do:
- Đáp ứng ngữ nghĩa FR-022 và FR-037 (HighestValidAttempt, autosave continuity, timeout auto-submit).
- Hạn chế mất bài làm khi gián đoạn mạng ngắn hạn.
- Giữ nhất quán việc chấm điểm qua các trạng thái attempt tường minh (`InProgress`, `Submitted`, `AutoSubmitted`, `Expired`).

Phương án thay thế đã cân nhắc:
- Chỉ lưu khi bấm submit: không đủ khả năng chịu lỗi kết nối.
- Lưu draft phía client: không đáp ứng độ tin cậy và khả năng audit.

## Quyết Định 6: Enforce quota và xử lý downgrade plan

Quyết định: Thực thi quota checks tại transactional service layer ở thời điểm ghi dữ liệu; cho phép downgrade nhưng chặn tạo mới khi đang over quota.

Lý do:
- Đáp ứng FR-034/FR-035/FR-038 mà không cần xóa dữ liệu hiện hữu.
- Tránh race conditions trong CSV import và các thao tác bulk.
- Hỗ trợ cảnh báo chính xác ở ngưỡng 80% và 90%.

Phương án thay thế đã cân nhắc:
- Chỉ kiểm tra quota theo lô (batch): phản hồi chậm và UX kém.
- Hard delete khi downgrade: vi phạm yêu cầu retention và continuity.

## Quyết Định 7: Độ tin cậy notification bằng outbox + worker retries

Quyết định: Lưu jobs gửi notification/email vào outbox table và xử lý async với tối đa 3 lần retry, sau đó chuyển trạng thái terminal failed.

Lý do:
- Đáp ứng yêu cầu reliability và observability của FR-039.
- Có delivery status rõ ràng để hiển thị trên dashboard quản trị.
- Tách độ trễ giao dịch người dùng khỏi sự không ổn định của nhà cung cấp email.

Phương án thay thế đã cân nhắc:
- Gửi email đồng bộ trong request: dễ vỡ và tăng latency.
- Fire-and-forget queue không có durable outbox: có thể mất event khi lỗi.

## Quyết Định 8: Mô hình reporting dựa trên bảng tổng hợp usage

Quyết định: Ghi nhận immutable usage events và duy trì bảng aggregate theo tenant/ngày, kết hợp Redis caching cho truy vấn dashboard.

Lý do:
- Hỗ trợ SC-006 về latency dashboard và phạm vi FR-027/FR-028.
- Tránh real-time joins nặng trên bảng giao dịch lớn.
- Giữ logic billing/analytics minh bạch và dễ kiểm thử.

Phương án thay thế đã cân nhắc:
- Tính dashboard trực tiếp từ transactional tables: khó scale ổn định.
- Dùng data warehouse đầy đủ ngay từ v1: quá phức tạp cho giai đoạn đầu.

## Quyết Định 9: Luồng upload file bằng presigned URL và xử lý async

Quyết định: API Go tạo presigned upload URL, client upload trực tiếp lên Rustfs, sau đó xử lý storage events bất đồng bộ và lưu metadata vào PostgreSQL.

Lý do:
- Phù hợp kiến trúc đã cung cấp và tránh nghẽn băng thông qua API.
- Tăng khả năng scale cho luồng upload media/document.
- Cho phép gắn pipeline malware scanning/transcoding mà không block request người dùng.

Phương án thay thế đã cân nhắc:
- Proxy toàn bộ upload qua API: semantics đơn giản hơn nhưng kém scale.
- Cho client ghi trực tiếp không qua presigned control: mô hình bảo mật yếu hơn.

## Quyết Định 10: Kiểm soát abuse ở mức API và account

Quyết định: Dùng Redis-backed IP rate limiting (100 req/min/IP cho public endpoints) và account lockout (15 phút sau 5 lần đăng nhập sai).

Lý do:
- Đáp ứng trực tiếp FR-042.
- Enforcement độ trễ thấp, phù hợp triển khai phân tán.
- Hỗ trợ observability và alerting qua counters/logs tập trung.

Phương án thay thế đã cân nhắc:
- In-memory counters theo từng node: hỏng khi scale ngang.
- Chỉ dựa vào WAF: không đủ cho lock semantics theo account.