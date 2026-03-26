# Easy-Study Constitution

## Core Principles

### I. Code Quality & Maintainability

Toàn bộ mã nguồn PHẢI tuân thủ các tiêu chuẩn sau:

- **Single Responsibility**: Mỗi module, class và function PHẢI có một mục đích được định nghĩa rõ ràng. Các file vượt quá 300 dòng PHẢI được xem xét để tách nhỏ.
- **Consistent Style**: Tất cả mã nguồn PHẢI vượt qua kiểm tra linting và formatting (ESLint, Prettier, hoặc công cụ tương đương) với không cảnh báo nào trước khi merge.
- **Meaningful Naming**: Biến, function và module PHẢI sử dụng tên mô tả rõ ý nghĩa. Viết tắt bị cấm ngoại trừ các thuật ngữ phổ biến toàn cầu (ví dụ: `id`, `url`, `api`).
- **Type Safety**: TypeScript strict mode PHẢI được bật cho toàn bộ frontend. Các backend service PHẢI sử dụng typed interface cho tất cả public API boundary.
- **Documentation**: Mỗi public API, service interface và thuật toán phức tạp PHẢI có tài liệu inline giải thích *tại sao*, không chỉ *cái gì*.
- **No Dead Code**: Mã không thể truy cập, import không sử dụng và các đoạn code bị comment PHẢI được xóa trước khi merge.

### II. Testing Standards — NON-NEGOTIABLE

Kiểm thử là bắt buộc và PHẢI tuân theo cách tiếp cận phân tầng có cấu trúc:

- **Unit Test**: Tất cả business logic function PHẢI có unit test với ≥80% branch coverage. Test PHẢI được đặt cùng vị trí với source hoặc trong thư mục `tests/` tương ứng.
- **Integration Test**: Mỗi API endpoint và database interaction PHẢI có integration test xác minh hành vi đúng qua các component boundary.
- **E2E Test**: Các user journey quan trọng (đăng ký khóa học, hoàn thành khóa học, nộp bài đánh giá, xem điểm) PHẢI có end-to-end test chạy trên môi trường staging.
- **Test-First cho Bug Fix**: Mỗi bug được xác nhận PHẢI có failing test được viết *trước* khi fix được triển khai.
- **Test Naming**: Test PHẢI tuân theo mẫu `should_[expected]_when_[condition]` để dễ đọc.
- **Không có Skipped Test**: Các test được đánh dấu `skip` hoặc `pending` KHÔNG ĐƯỢC tồn tại trong codebase quá một sprint mà không có issue liên kết và kế hoạch khắc phục.

### III. User Experience Consistency

Nền tảng PHẢI cung cấp trải nghiệm thống nhất, accessible trên mọi bề mặt:

- **Design System**: Tất cả UI component PHẢI được lấy từ design system chung. Styling tùy tiện KHÔNG ĐƯỢC thêm vào mà không có sự phê duyệt của design review.
- **Responsive Design**: Mỗi trang PHẢI hiển thị chính xác trên viewport từ 320px (mobile) đến 2560px (ultrawide). Hành vi breakpoint PHẢI được kiểm thử.
- **Accessibility**: Tất cả interactive element PHẢI đạt chuẩn WCAG 2.1 AA. Bao gồm keyboard navigation, screen reader support, tỷ lệ color contrast ≥4.5:1 và ARIA label trên các element không phải text.
- **Loading State**: Mỗi asynchronous operation PHẢI hiển thị loading indicator trong vòng 100ms. Skeleton screen được ưu tiên hơn spinner cho vùng nội dung.
- **Error State**: Tất cả error condition PHẢI hiển thị thông báo thân thiện với người dùng kèm hướng dẫn hành động. Raw error code hoặc stack trace KHÔNG ĐƯỢC hiển thị.
- **Internationalization**: Tất cả chuỗi hiển thị cho người dùng PHẢI được tách ra thành locale file ngay từ đầu. Hardcode display text bị cấm.

### IV. Performance Requirements

Nền tảng PHẢI đạt các chỉ tiêu hiệu năng đo lường được sau:

- **Page Load**: Initial page load (LCP) PHẢI hoàn thành trong 2.5 giây trên kết nối 4G. Client-side navigation PHẢI hoàn thành trong 500ms.
- **API Response**: 95th percentile API response time PHẢI ≤200ms cho read operation và ≤500ms cho write operation dưới tải bình thường.
- **Database Query**: Không query nào PHẢI vượt quá 100ms thời gian thực thi. Các query chạy trên 50ms PHẢI được ghi log để review.
- **Bundle Size**: Frontend JavaScript bundle KHÔNG ĐƯỢC vượt quá 250KB gzipped cho lần tải đầu tiên. Code splitting PHẢI được áp dụng theo route.
- **Concurrent User**: Hệ thống PHẢI hỗ trợ 1.000 concurrent user mỗi tenant mà không suy giảm hiệu năng vượt quá SLA đã công bố.
- **Memory**: Các backend service KHÔNG ĐƯỢC vượt quá 512MB resident memory mỗi instance dưới tải bình thường. Memory leak PHẢI kích hoạt cảnh báo.

### V. Security & Data Protection

- **OWASP Top 10**: PHẢI đảm bảo tuân thủ OWASP Top 10.
- **Authentication**: Tất cả endpoint PHẢI yêu cầu authentication ngoại trừ các public route được whitelist rõ ràng. JWT token PHẢI hết hạn trong 1 giờ.
- **Authorization**: Role-based access control (RBAC) PHẢI được thực thi tại API layer. Frontend route guard đơn thuần KHÔNG đủ.
- **Data Encryption**: Tất cả dữ liệu truyền tải PHẢI sử dụng TLS 1.2+. Dữ liệu nhạy cảm lưu trữ (PII, điểm số, câu trả lời bài đánh giá) PHẢI được mã hóa.
- **Input Validation**: Tất cả user input PHẢI được validate và sanitize tại API boundary. Bảo vệ chống SQL injection, XSS và CSRF PHẢI được xác minh qua automated security scan.
- **Audit Logging**: Tất cả thay đổi điểm số, thay đổi role và export dữ liệu PHẢI được ghi vào immutable audit log với thông tin actor, timestamp và action.
- **Dependency Security**: Tất cả third-party dependency PHẢI được scan tìm vulnerability đã biết trên mỗi build.
- **GDPR/data privacy**: PHẢI hỗ trợ xuất và xóa dữ liệu theo yêu cầu người dùng.

### VI. Simplicity & Pragmatism

Độ phức tạp PHẢI được biện minh. Câu trả lời mặc định là cách tiếp cận đơn giản hơn:

- **YAGNI**: Feature, abstraction và infrastructure KHÔNG ĐƯỢC thêm vào một cách suy đoán. Xây dựng cho yêu cầu hiện tại với các extension point rõ ràng.
- **Standard Tooling**: Ưu tiên thư viện được sử dụng rộng rãi, có tài liệu tốt hơn so với custom implementation. Giải pháp tùy chỉnh PHẢI có biện minh bằng văn bản giải thích tại sao không có thư viện hiện có nào đáp ứng.
- **Flat Architecture**: Ưu tiên ít layer trừu tượng hơn. Mỗi architectural layer PHẢI biện minh sự tồn tại của nó bằng cách giải quyết một vấn đề cụ thể, được ghi nhận.

## Performance Standards & SLAs

| Chỉ Số | Mục Tiêu | Phương Pháp Đo Lường |
|--------|----------|---------------------|
| LCP (Largest Contentful Paint) | ≤2.5s (4G) | Lighthouse CI trên mỗi deploy |
| FID (First Input Delay) | ≤100ms | Real User Monitoring (RUM) |
| CLS (Cumulative Layout Shift) | ≤0.1 | Lighthouse CI trên mỗi deploy |
| API p95 Latency (read) | ≤200ms | APM dashboard (Datadog/New Relic) |
| API p95 Latency (write) | ≤500ms | APM dashboard |
| Uptime | ≥99.9% hàng tháng | Status page monitoring |
| Error Rate | <0.1% request | APM dashboard |
| Time to Interactive | ≤3.5s (4G) | Lighthouse CI |

Performance budget PHẢI được thực thi trong CI. Bất kỳ build nào làm giảm Core Web Vital vượt quá mục tiêu PHẢI fail pipeline và yêu cầu sự phê duyệt rõ ràng từ team lead để override.

## Development Workflow & Quality Gates

### Worlkflow

1. **Branch Strategy**: `main` (production) ← `develop` (staging) ← `feature/*`, `fix/*`
2. **Commit Convention**: Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`)
3. **PR Process**: Feature branch → PR → Code review (≥1 approval) → CI pass → Merge
4. **Database Changes**: Chỉ dùng migration files, không DDL thủ công; migrations MUST có thể rollback
5. **Environment Parity**: Local dev, staging, production MUST dùng cùng Docker config base

### Yêu Cầu Pull Request

1. **Automated Check** (PHẢI vượt qua tất cả):
   - Linting và formatting (không cảnh báo nào)
   - Toàn bộ unit test suite (duy trì ≥80% branch coverage)
   - Integration test suite
   - Type checking với strict mode
   - Kiểm tra bundle size budget
   - Security dependency scan

2. **Human Review** (PHẢI hoàn thành):
   - Xác minh tuân thủ Constitution (reviewer PHẢI xác nhận sự phù hợp với các nguyên tắc liên quan)

3. **Deployment Gate**:
   - Staging deployment PHẢI thành công và E2E test PHẢI vượt qua trước khi
     promote lên production
   - Feature flag PHẢI được sử dụng cho các thay đổi ảnh hưởng >1% người dùng
   - Rollback plan PHẢI được ghi nhận cho các thay đổi infrastructure

### Code Review Checklist

Reviewer PHẢI xác minh:

- [ ] Principle I: Code dễ đọc, đặt tên tốt và typed đúng cách
- [ ] Principle II: Test có mặt, có ý nghĩa và passing
- [ ] Principle III: Thay đổi UI tuân theo design system và accessibility standard
- [ ] Principle IV: Không có performance regression
- [ ] Principle V: Các vấn đề security được xử lý (auth, validation, logging)
- [ ] Principle VI: Giải pháp đủ đơn giản cho vấn đề

## Governance

- Constitution này là tài liệu tối cao; mọi quyết định kỹ thuật MUST tuân theo các principles trên
- Sửa đổi yêu cầu: documented rationale, team review, migration plan nếu có breaking changes
- Versioning policy: MAJOR (bỏ/tái định nghĩa principle), MINOR (thêm principle/section), PATCH (clarification/typo)
- Khi conflict giữa tốc độ và chất lượng → chất lượng thắng (trừ khi có ngoại lệ được document rõ ràng)
- Compliance review: mọi PR/review MUST xác minh sự tuân thủ với các principles trong constitution
- Khi đưa ra quyết định kỹ thuật, thành viên đội PHẢI áp dụng các nguyên tắc theo thứ tự ưu tiên:

  1. **Security & Data Protection** (Principle V)
  2. **Testing Standards** (Principle II)
  3. **Performance Requirements** (Principle IV)
  4. **Code Quality** (Principle I)
  5. **UX Consistency** (Principle III)
  6. **Simplicity** (Principle VI)

**Version**: 1.0.0 | **Ratified**: 2026-03-26 | **Last Amended**: 2026-03-26
