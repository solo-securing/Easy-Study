# Feature Specification: Hệ Thống LMS SaaS Multi-Tenant

**Feature Branch**: `001-lms-saas-system`  
**Created**: 2026-03-26  
**Status**: Draft  
**Input**: User description: "Hệ thống LMS SaaS multi-tenant phục vụ nhiều tổ chức khách hàng trên cùng hạ tầng, với quản lý tenant, phân quyền, khóa học, báo cáo, thông báo và billing."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Super Admin Tạo và Quản Lý Tenant Mới (Priority: P1)

Super Admin đăng nhập vào trang quản trị hệ thống SaaS, tạo một tenant mới bằng cách cung cấp thông tin: tên tổ chức, subdomain, gói dịch vụ (plan), giới hạn user/course. Hệ thống tự động khởi tạo không gian riêng cho tenant, gán Tenant Admin đầu tiên (owner), và kích hoạt tenant. Super Admin cũng có thể tạm khóa hoặc deactivate tenant khi cần.

**Why this priority**: Đây là nền tảng cốt lõi của mô hình SaaS — nếu không có khả năng tạo và quản lý tenant, toàn bộ hệ thống không thể vận hành. Mọi tính năng khác đều phụ thuộc vào việc tenant tồn tại.

**Independent Test**: Có thể kiểm tra đầy đủ bằng cách Super Admin tạo tenant mới, xác nhận tenant xuất hiện trong danh sách, truy cập subdomain của tenant, và thực hiện kích hoạt/tạm khóa. Cung cấp giá trị ngay: onboarding khách hàng mới.

**Acceptance Scenarios**:

1. **Given** Super Admin đã đăng nhập, **When** tạo tenant mới với thông tin hợp lệ (tên, subdomain, plan, owner email), **Then** hệ thống tạo tenant thành công, gửi invitation đến Tenant Admin, và tenant có thể truy cập qua subdomain.
2. **Given** tenant đang active, **When** Super Admin chọn tạm khóa (suspend) tenant, **Then** toàn bộ user trong tenant không thể đăng nhập, dữ liệu được bảo toàn, và tenant hiển thị trạng thái "Suspended".
3. **Given** tenant đang suspended, **When** Super Admin kích hoạt lại tenant, **Then** toàn bộ user có thể đăng nhập lại và dữ liệu không thay đổi.
4. **Given** Super Admin tạo tenant với subdomain đã tồn tại, **When** submit form, **Then** hệ thống hiển thị lỗi trùng subdomain và không tạo tenant.

---

### User Story 2 - Tenant Admin Cấu Hình và Quản Lý Người Dùng (Priority: P1)

Tenant Admin đăng nhập vào không gian LMS của tổ chức mình, cấu hình thương hiệu (logo, màu sắc, tên hiển thị), và quản lý người dùng: tạo tài khoản thủ công, upload CSV, invite qua email, gán role (Student, Instructor, Tenant Admin phụ), và nhóm user theo phòng ban/team.

**Why this priority**: Tenant Admin là người vận hành hàng ngày. Nếu không quản lý được user và cấu hình không gian, tenant không hoạt động — cùng mức ưu tiên với tạo tenant.

**Independent Test**: Tenant Admin đăng nhập, cấu hình branding, tạo user bằng CSV, gán role, tạo nhóm. Xác nhận user nhận invitation, đăng nhập thành công, và thấy branding đúng.

**Acceptance Scenarios**:

1. **Given** Tenant Admin đã đăng nhập, **When** cập nhật logo và màu sắc chủ đạo, **Then** tất cả user trong tenant nhìn thấy giao diện mới ngay lập tức.
2. **Given** Tenant Admin upload file CSV chứa 100 user hợp lệ, **When** hệ thống xử lý xong, **Then** 100 tài khoản được tạo, email invitation được gửi, và báo cáo import hiển thị kết quả (thành công/lỗi từng dòng).
3. **Given** Tenant Admin tạo phòng ban "Kỹ Thuật" và gán 10 user, **When** gán khóa học cho phòng ban, **Then** cả 10 user đều thấy khóa học trong danh sách của mình.
4. **Given** user thuộc tenant A, **When** user cố truy cập dữ liệu của tenant B, **Then** hệ thống từ chối truy cập và ghi log bảo mật.
5. **Given** Tenant Admin tạo user và hệ thống gửi invitation email, **When** email gửi thất bại sau 3 lần retry, **Then** hệ thống ghi log lỗi và hiển thị trạng thái "Gửi invitation thất bại" cho Tenant Admin.
6. **Given** một user mới nhận được invitation email, **When** click vào "Magic Link" trong email, **Then** hệ thống yêu cầu thiết lập mật khẩu, ghi nhận tài khoản kích hoạt thành công và tự động login vào tenant.

---

### User Story 3 - Instructor Tạo và Quản Lý Khóa Học (Priority: P2)

Instructor đăng nhập vào LMS của tenant, tạo khóa học mới với cấu trúc: section → subsection → unit (video, tài liệu, quiz). Instructor publish khóa học, quy định ai được học (theo user, group, role), và theo dõi tiến độ học viên.

**Why this priority**: Khóa học là giá trị cốt lõi mà LMS cung cấp. Tuy nhiên, cần có tenant và user trước (P1), nên đây là P2 — bước tiếp theo tự nhiên sau khi tenant được thiết lập.

**Independent Test**: Instructor tạo khóa học đầy đủ (section, subsection, unit), publish, gán cho nhóm học viên, và xem dashboard theo dõi. Xác nhận học viên thấy khóa học và tiến độ được ghi nhận.

**Acceptance Scenarios**:

1. **Given** Instructor đã đăng nhập, **When** tạo khóa học với 3 section, mỗi section có subsection, mỗi subsection có unit (video, tài liệu, quiz), **Then** khóa học được lưu ở trạng thái Draft và chỉ Instructor thấy.
2. **Given** khóa học ở trạng thái Draft, **When** Instructor publish khóa học và gán cho phòng ban "Kỹ Thuật", **Then** toàn bộ thành viên phòng ban thấy khóa học trong danh sách và có thể bắt đầu học.
3. **Given** Instructor thuộc tenant A, **When** tạo khóa học, **Then** khóa học chỉ hiển thị trong tenant A, không xuất hiện ở tenant B.
4. **Given** khóa học có section trống hoặc unit thiếu nội dung, **When** Instructor cố publish, **Then** hệ thống chặn publish và hiển thị danh sách lỗi chi tiết chỉ rõ từng section/subsection/unit bị thiếu nội dung.

---

### User Story 4 - Student Tham Gia và Hoàn Thành Khóa Học (Priority: P2)

Student đăng nhập vào LMS của tenant, xem danh sách khóa học được gán, bắt đầu học (xem video, đọc tài liệu, làm quiz), và theo dõi tiến độ cá nhân. Student nhận thông báo khi được assign khóa mới hoặc khi deadline quiz sắp đến.

**Why this priority**: Trải nghiệm học viên là lý do tồn tại của LMS. Cùng mức với tạo khóa học vì cả hai tạo thành luồng sử dụng hoàn chỉnh.

**Independent Test**: Student đăng nhập, mở khóa học, xem video, làm quiz, kiểm tra tiến độ. Xác nhận tiến độ được cập nhật chính xác và thông báo hoạt động.

**Acceptance Scenarios**:

1. **Given** Student được gán khóa học, **When** đăng nhập, **Then** thấy khóa học trong danh sách "Khóa của tôi" với tiến độ 0%.
2. **Given** Student đang học khóa, **When** hoàn thành 1 unit, **Then** tiến độ được cập nhật tự động và Instructor thấy sự thay đổi.
3. **Given** Student hoàn thành quiz, **When** submit câu trả lời, **Then** nhận kết quả ngay và điểm được ghi nhận.
4. **Given** Student thuộc tenant A, **When** tìm kiếm khóa học, **Then** chỉ thấy khóa học của tenant A.
5. **Given** Student đang làm quiz, **When** mất kết nối hoặc đóng trình duyệt, **Then** hệ thống đã auto-save câu trả lời và student có thể tiếp tục làm khi quay lại (trong thời hạn quiz).
6. **Given** Student mất kết nối khi đang làm quiz, **When** hết thời hạn quiz mà chưa quay lại, **Then** hệ thống tự động submit các câu đã auto-save và ghi nhận kết quả.

---

### User Story 5 - Super Admin Impersonate Tenant để Hỗ Trợ (Priority: P2)

Super Admin cần hỗ trợ một tenant gặp vấn đề. Super Admin sử dụng chức năng "impersonate" để vào không gian của tenant, xem và thao tác như Tenant Admin nhưng không xuất hiện trong dữ liệu học tập của tenant. Mọi hành động được ghi audit log.

**Why this priority**: Giảm gánh nặng hỗ trợ khách hàng, nhưng vẫn đảm bảo bảo mật — cần thiết cho vận hành nhưng không block các luồng chính.

**Independent Test**: Super Admin impersonate vào tenant, thực hiện thao tác (xem user, xem khóa, chỉnh cấu hình), thoát impersonate. Xác nhận audit log ghi đủ và Super Admin không xuất hiện trong dữ liệu học tập.

**Acceptance Scenarios**:

1. **Given** Super Admin chọn impersonate tenant A, **When** thực hiện impersonate, **Then** giao diện hiển thị banner cảnh báo "Đang impersonate tenant A" và Super Admin thấy mọi thứ như Tenant Admin.
2. **Given** Super Admin đang impersonate, **When** xem danh sách học viên, **Then** không thấy chính mình trong danh sách user hoặc dữ liệu học tập.
3. **Given** Super Admin thực hiện thao tác trong impersonate mode, **When** kiểm tra audit log, **Then** mọi hành động đều được ghi với thông tin "Super Admin X impersonating Tenant A".
4. **Given** Super Admin đang impersonate tenant, **When** cố thực hiện destructive action (xóa user, xóa khóa học, xóa dữ liệu), **Then** hệ thống chặn hành động và hiển thị thông báo "Không được phép thực hiện hành động này trong chế độ impersonate".
5. **Given** Super Admin đã impersonate tenant (plan Enterprise), **When** Tenant Admin mở impersonate audit log, **Then** thấy bản ghi các phiên impersonate với thông tin Super Admin, thời gian bắt đầu/kết thúc, và các hành động đã thực hiện.

---

### User Story 6 - Tenant Admin Xem Báo Cáo và Dashboard (Priority: P3)

Tenant Admin truy cập dashboard để xem tổng quan: số user active, số khóa học active, tỷ lệ hoàn thành khóa. Instructor xem báo cáo chi tiết theo khóa: danh sách học viên, tiến độ, kết quả quiz.

**Why this priority**: Báo cáo quan trọng cho đánh giá hiệu quả training, nhưng hệ thống vẫn hoạt động được nếu chưa có báo cáo. Phụ thuộc vào dữ liệu từ khóa học và học viên (P1, P2).

**Independent Test**: Tenant Admin mở dashboard, xác nhận số liệu khớp với dữ liệu thực tế. Instructor mở báo cáo khóa học, xác nhận danh sách học viên và tiến độ chính xác.

**Acceptance Scenarios**:

1. **Given** tenant có 50 user active và 10 khóa học, **When** Tenant Admin mở dashboard, **Then** hiển thị số liệu tổng quan chính xác.
2. **Given** khóa học có 20 học viên đã hoàn thành, **When** Instructor xem báo cáo khóa, **Then** thấy danh sách học viên với tiến độ và kết quả quiz của từng người.

---

### User Story 7 - Super Admin Theo Dõi Usage và Quản Lý Billing (Priority: P3)

Super Admin xem dashboard toàn hệ thống: thống kê per-tenant (số user, số khóa, usage), so sánh tenants. Super Admin gán plan cho tenant, theo dõi quota, và hệ thống cảnh báo khi tenant gần chạm giới hạn.

**Why this priority**: Billing và usage tracking quan trọng cho mô hình kinh doanh SaaS, nhưng hệ thống có thể hoạt động cơ bản mà chưa cần billing automation. Đây là lớp thương mại bổ sung.

**Independent Test**: Super Admin mở dashboard hệ thống, xem usage per-tenant, gán plan cho tenant, xác nhận cảnh báo khi tenant gần hết quota.

**Acceptance Scenarios**:

1. **Given** hệ thống có 5 tenants, **When** Super Admin mở dashboard, **Then** thấy bảng tổng hợp usage (user count, course count, login count, learning hours) cho từng tenant.
2. **Given** tenant sử dụng 90% quota user, **When** Tenant Admin thêm user mới, **Then** hệ thống hiển thị cảnh báo "Gần chạm giới hạn plan" cho cả Tenant Admin và Super Admin.
3. **Given** Super Admin thay đổi plan của tenant từ Free sang Pro, **When** cập nhật xong, **Then** giới hạn user/course mới được áp dụng ngay và Tenant Admin nhận thông báo.
4. **Given** tenant có 50 user (plan Pro), **When** Super Admin downgrade sang Free (giới hạn 10 user), **Then** hệ thống cho phép downgrade, 50 user hiện tại vẫn hoạt động bình thường, nhưng chặn tạo thêm user/course mới cho đến khi giảm xuống dưới quota.
5. **Given** tenant đã bị downgrade và vượt quota, **When** Tenant Admin cố thêm user mới, **Then** hệ thống chặn và hiển thị thông báo "Đã vượt giới hạn plan, vui lòng nâng cấp hoặc giảm số lượng hiện tại".

---

### User Story 8 - Thông Báo Multi-Tenant (Priority: P3)

Hệ thống gửi thông báo trong từng tenant (assign khóa mới, deadline quiz, thông báo chung từ Tenant Admin) và từ SaaS provider đến tất cả tenants (maintenance, update, tính năng mới) qua in-app notification và/hoặc email.

**Why this priority**: Thông báo giữ người dùng engaged, nhưng hệ thống vẫn hoạt động cơ bản không có thông báo tự động. Phụ thuộc vào user, khóa học đã tồn tại.

**Independent Test**: Gán khóa cho student, xác nhận student nhận thông báo. Super Admin gửi thông báo hệ thống, xác nhận hiển thị cho tất cả tenant.

**Acceptance Scenarios**:

1. **Given** Instructor gán khóa mới cho student, **When** student đăng nhập, **Then** thấy thông báo "Bạn được assign khóa mới: [tên khóa]" trong notification center.
2. **Given** Tenant Admin tạo thông báo chung, **When** publish, **Then** toàn bộ user trong tenant nhận thông báo.
3. **Given** Super Admin tạo thông báo bảo trì hệ thống, **When** publish, **Then** tất cả user trên tất cả tenant thấy global banner với nội dung thông báo.
4. **Given** hệ thống gửi email thông báo cho student, **When** email gửi thất bại sau 3 lần retry, **Then** hệ thống ghi log lỗi và hiển thị trạng thái "Gửi thất bại" trong notification management cho admin.

---

### Edge Cases

- Điều gì xảy ra khi tenant bị suspend nhưng có student đang làm quiz? → Hệ thống cho phép hoàn thành quiz hiện tại, sau đó chặn truy cập.
- Điều gì xảy ra khi Tenant Admin bị xóa và là owner duy nhất? → Hệ thống yêu cầu gán owner mới trước khi cho phép xóa.
- Điều gì xảy ra khi tenant vượt quá quota trong khi user đang bulk import? → Hệ thống import đến khi đạt giới hạn, dừng lại, và báo cáo số user đã import thành công vs thất bại.
- Điều gì xảy ra khi subdomain tenant trùng với route hệ thống? → Hệ thống có danh sách reserved subdomain và từ chối khi tạo.
- Điều gì xảy ra khi Super Admin impersonate tenant và tenant bị suspend cùng lúc? → Phiên impersonate kết thúc, Super Admin nhận thông báo tenant đã bị suspend.
- Điều gì xảy ra khi hai Tenant Admin cùng sửa cấu hình branding đồng thời? → Hệ thống áp dụng cơ chế "last write wins" kèm thông báo cho bên còn lại.

## Requirements *(mandatory)*

### Functional Requirements

**Quản lý Tenant (Tenant Management)**

- **FR-001**: Hệ thống PHẢI cho phép Super Admin tạo tenant mới với thông tin: tên tổ chức, subdomain, plan, giới hạn user, giới hạn course, và email của Tenant Admin đầu tiên.
- **FR-002**: Hệ thống PHẢI hỗ trợ các trạng thái tenant: Active, Suspended, Deactivated. Super Admin PHẢI có thể chuyển đổi giữa các trạng thái.
- **FR-003**: Hệ thống PHẢI ngăn chặn tạo tenant với subdomain trùng lặp hoặc nằm trong danh sách reserved.
- **FR-004**: Hệ thống PHẢI tự động khởi tạo không gian riêng cho tenant mới với cấu hình mặc định.

**Cấu hình Tenant (Tenant Configuration)**

- **FR-005**: Tenant Admin PHẢI có thể cấu hình branding: logo, màu chủ đạo, favicon, tên hiển thị.
- **FR-006**: Tenant Admin PHẢI có thể thiết lập ngôn ngữ mặc định và timezone cho tenant.
- **FR-007**: Tenant Admin PHẢI có thể cấu hình policy: password policy, thời gian session, quy tắc privacy.
- **FR-008**: Hệ thống PHẢI cho phép tenant truy cập qua subdomain riêng (ví dụ: companyA.platform.com).

**Phân quyền Đa Tenant (Multi-tenant Authorization)**

- **FR-009**: Mỗi user PHẢI được gắn với một tenant chính duy nhất. Role và quyền PHẢI được scope theo tenant.
- **FR-010**: Hệ thống PHẢI đảm bảo data isolation hoàn toàn: user của tenant A KHÔNG ĐƯỢC truy cập bất kỳ dữ liệu nào của tenant B.
- **FR-011**: Super Admin PHẢI có khả năng impersonate vào tenant để hỗ trợ. Chế độ impersonate CHỈ cho phép read và cấu hình (branding, policy), KHÔNG ĐƯỢC thực hiện destructive actions (xóa user, xóa khóa học, xóa dữ liệu). Phiên impersonate KHÔNG có giới hạn thời gian, kéo dài cho đến khi Super Admin chủ động thoát. Mọi hành động impersonate PHẢI được ghi audit log.
- **FR-012**: Trong chế độ impersonate, Super Admin KHÔNG ĐƯỢC xuất hiện trong dữ liệu học tập hoặc danh sách user của tenant.

**Quản lý Người Dùng (User Management)**

- **FR-013**: Tenant Admin PHẢI có thể tạo user thủ công, upload CSV (bulk import), và invite qua email.
- **FR-014**: Tenant Admin PHẢI có thể gán role: Student, Instructor, Tenant Admin phụ cho từng user.
- **FR-015**: Hệ thống PHẢI hỗ trợ nhóm user theo phòng ban/team để gán khóa học theo nhóm.
- **FR-016**: Hệ thống PHẢI hỗ trợ student tự đăng ký qua subdomain của tenant (nếu tenant cho phép) hoặc qua invitation. Tài khoản tự đăng ký PHẢI được tự động duyệt nhưng BẮT BUỘC phải xác thực email (Email Verification) trước khi được login.
- **FR-017**: Instructor CHỈ ĐƯỢC tham gia tenant qua invitation từ Tenant Admin.

**Quản lý Khóa Học (Course Management)**

- **FR-018**: Instructor/Tenant Admin PHẢI có thể tạo khóa học với cấu trúc: section → subsection → unit (video, tài liệu, quiz).
- **FR-019**: Khóa học PHẢI có trạng thái: Draft, Published, Archived. Chỉ khóa Published mới hiển thị cho student.
- **FR-020**: Hệ thống PHẢI cho phép quy định quyền truy cập khóa học theo user, group, hoặc role.
- **FR-021**: Khóa học của tenant A KHÔNG ĐƯỢC hiển thị hoặc truy cập được từ tenant B.

**Theo Dõi Tiến Độ (Progress Tracking)**

- **FR-022**: Hệ thống PHẢI tự động ghi nhận tiến độ học viên: unit đã xem, quiz đã làm, điểm số.
- **FR-023**: Student PHẢI có thể xem tiến độ cá nhân cho mỗi khóa học (phần trăm hoàn thành, kết quả quiz).
- **FR-024**: Instructor PHẢI có thể xem tiến độ của tất cả học viên trong khóa học mình quản lý.

**Báo Cáo (Reporting)**

- **FR-025**: Tenant Admin PHẢI có dashboard hiển thị: số user active, số khóa active, tỷ lệ hoàn thành khóa.
- **FR-026**: Instructor PHẢI có báo cáo theo khóa: danh sách học viên, tiến độ, kết quả quiz.
- **FR-027**: Super Admin PHẢI có dashboard toàn hệ thống: thống kê per-tenant bao gồm số user, số khóa, usage (giờ học, lượt login).
- **FR-028**: Dữ liệu usage per-tenant chỉ Super Admin xem được, KHÔNG ĐƯỢC hiển thị chéo giữa các tenant.

**Thông Báo (Notifications)**

- **FR-029**: Hệ thống PHẢI gửi thông báo in-app khi student được assign khóa mới, khi deadline quiz đến gần, và khi có thông báo chung từ Tenant Admin.
- **FR-030**: Tenant Admin PHẢI có thể gửi thông báo chung đến toàn bộ user trong tenant.
- **FR-031**: Super Admin PHẢI có thể gửi thông báo toàn hệ thống (maintenance, update) hiển thị dưới dạng global banner.

**Billing & Plan**

- **FR-032**: Hệ thống PHẢI hỗ trợ các plan cơ bản từ mã nguồn (Free/Trial, Pro, Enterprise). Enterprise có user lớn, SLA riêng (quản lý offline qua hợp đồng, không chia luồng traffic ưu tiên ở mức hệ thống trong v1), SSO, impersonate audit log. Đối với tenant có cấu hình SSO, hệ thống PHẢI tự động tạo user (Just-In-Time provisioning) với role mặc định là Student khi user login thành công lần đầu.
- **FR-033**: Super Admin PHẢI có thể gán, thay đổi plan cho tenant, và tạo "Custom Plan" để override các limit (số user, số khóa học) riêng biệt cho tenant đó. Thay đổi plan PHẢI có hiệu lực ngay.
- **FR-034**: Hệ thống PHẢI theo dõi usage so với quota (số user, số khóa, dung lượng) và cảnh báo khi tenant đạt 80% và 90% giới hạn.
- **FR-035**: Hệ thống PHẢI ngăn chặn vượt quota: không cho thêm user/course khi đạt giới hạn plan.
- **FR-036**: Hệ thống PHẢI chặn publish khóa học nếu có section/subsection trống hoặc unit thiếu nội dung, và PHẢI hiển thị danh sách lỗi chi tiết chỉ rõ từng vị trí bị thiếu.
- **FR-037**: Hệ thống PHẢI tự động lưu câu trả lời quiz định kỳ. Khi student mất kết nối, PHẢI cho phép tiếp tục làm khi quay lại trong thời hạn quiz. Khi hết thời hạn mà student chưa quay lại, PHẢI tự động submit các câu đã lưu.
- **FR-038**: Khi downgrade plan và tenant đã vượt quota của plan mới, hệ thống PHẢI cho phép downgrade, PHẢI giữ nguyên dữ liệu hiện tại, và PHẢI chặn tạo thêm user/course mới cho đến khi tenant giảm xuống dưới quota.
- **FR-039**: Khi gửi email thất bại (invitation, notification, cảnh báo quota), hệ thống PHẢI retry tự động tối đa 3 lần. Nếu vẫn thất bại, PHẢI ghi log lỗi và PHẢI hiển thị trạng thái "gửi thất bại" cho admin tương ứng.
- **FR-040**: Tenant Admin của plan Enterprise PHẢI có thể xem impersonate audit log trong tenant của mình, bao gồm: thông tin Super Admin, thời gian bắt đầu/kết thúc phiên, và danh sách hành động đã thực hiện. Các plan khác KHÔNG hiển thị impersonate audit log cho Tenant Admin.
- **FR-041**: Khi tenant chuyển sang trạng thái Deactivated, hệ thống PHẢI giữ lại dữ liệu trong 90 ngày (grace period). Sau 90 ngày, dữ liệu của tenant PHẢI được xóa vĩnh viễn (hard delete).
- **FR-042**: Hệ thống PHẢI tạm khóa tài khoản 15 phút sau 5 lần đăng nhập sai liên tiếp, và API public PHẢI được rate limit ở mức 100 requests/phút/IP để chống brute force và spam.

### Key Entities

- **Tenant**: Đại diện cho một tổ chức khách hàng. Thuộc tính chính: tên, subdomain, trạng thái (Active/Suspended/Deactivated), plan, quota, cấu hình branding, policy. Quan hệ: chứa nhiều User, Course, Group.
- **User**: Người dùng trong hệ thống, gắn với 1 tenant. Thuộc tính: thông tin cá nhân, role (Super Admin, Tenant Admin, Instructor, Student), trạng thái. Quan hệ: thuộc Tenant, thuộc Group, tham gia Course.
- **Course**: Khóa học do Instructor/Tenant Admin tạo. Thuộc tính: tên, mô tả, trạng thái (Draft/Published/Archived), cấu trúc nội dung. Quan hệ: thuộc Tenant, chứa nhiều Section, gán cho User/Group.
- **Section**: Phần của khóa học. Chứa nhiều Subsection. Quan hệ: thuộc Course.
- **Subsection**: Phần của section. Chứa nhiều Unit. Quan hệ: thuộc Section.
- **Unit**: Bài học trong subsection. Loại: video, tài liệu, quiz. Quan hệ: thuộc Subsection.
- **Quiz**: Bài kiểm tra gắn với unit. Chứa câu hỏi, đáp án, điểm. Quan hệ: thuộc Unit.
- **Group**: Nhóm user theo phòng ban/team trong tenant. Quan hệ: thuộc Tenant, chứa User, gán Course.
- **Plan**: Gói dịch vụ (Free/Trial, Pro, Enterprise, Custom). Thuộc tính: giới hạn user, course, dung lượng, danh sách feature. Quan hệ: gán cho Tenant.
- **Progress**: Tiến độ học của student. Thuộc tính: khóa đã hoàn thành (%), unit đã xem, kết quả quiz. Quan hệ: liên kết User với Course.
- **Notification**: Thông báo trong hệ thống. Thuộc tính: nội dung, loại (in-tenant, global), trạng thái đọc. Quan hệ: gắn Tenant hoặc hệ thống, nhận bởi User.
- **AuditLog**: Bản ghi hành động quan trọng. Thuộc tính: actor, action, timestamp, context (impersonate, role change, data export). Quan hệ: gắn User, Tenant.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Super Admin có thể tạo và kích hoạt tenant mới sẵn sàng sử dụng trong vòng 5 phút.
- **SC-002**: Tenant Admin có thể bulk import 500 user qua CSV trong vòng 3 phút, với tỷ lệ import thành công ≥98% cho dữ liệu hợp lệ.
- **SC-003**: 100% trường hợp truy cập dữ liệu cross-tenant bị chặn — không có data leakage giữa các tenant.
- **SC-004**: Hệ thống hỗ trợ ít nhất 100 tenants active (đang có user đăng nhập) cùng lúc, mỗi tenant tối đa 1.000 concurrent users, mà không suy giảm hiệu năng ngoài SLA. Đo bằng load testing và APM monitoring.
- **SC-005**: Student có thể bắt đầu học khóa mới trong vòng 2 phút kể từ lúc đăng nhập.
- **SC-006**: Dashboard báo cáo tải trong vòng 3 giây cho tenant có đến 10.000 user.
- **SC-007**: Thông báo in-app được gửi đến người nhận trong vòng 30 giây sau sự kiện kích hoạt.
- **SC-008**: Cảnh báo quota được gửi chính xác 100% khi tenant đạt ngưỡng 80% và 90% giới hạn plan.

## Non-Functional Requirements

> Các chỉ số dưới đây bổ sung cho constitution (`.specify/memory/constitution.md`). Những giá trị đã định nghĩa trong constitution được reference, không lặp lại.

### Performance (reference constitution)

- API p95 latency: ≤200ms (read), ≤500ms (write) — *định nghĩa tại Constitution §IV*
- Page load LCP: ≤2.5s (4G) — *Constitution §IV*
- Client-side navigation: ≤500ms — *Constitution §IV*

### Availability (reference constitution)

- Uptime SLA: ≥99.9% hàng tháng — *Constitution Performance Standards*
- Error rate: <0.1% request — *Constitution Performance Standards*

### Security (reference constitution + bổ sung)

- JWT token expiry: 1 giờ — *Constitution §V*
- OWASP Top 10 compliance — *Constitution §V*
- RBAC tại API layer — *Constitution §V*
- GDPR data privacy: hỗ trợ xuất và xóa dữ liệu theo yêu cầu — *Constitution §V*
- **Data retention**: giữ dữ liệu 90 ngày sau khi tenant bị deactivate, sau đó xóa vĩnh viễn ("hard delete").
- **Brute force & Rate limit**: Khóa tài khoản 15 phút sau 5 lần đăng nhập sai; giới hạn API public 100 req/phút/IP.

### Scalability

- **Horizontal scale**: Không có hard limit kiến trúc về tổng số tenant. Thiết kế cho phép scale out theo traffic.
- **v1 Benchmark Target**: Hệ thống chịu tải 100 tenants active cùng lúc (đang có user đăng nhập), mỗi tenant tối đa 1.000 concurrent users — *SC-004*
- 1.000 concurrent users/tenant limit — *Constitution §IV*

## Assumptions

- Người dùng truy cập hệ thống qua web browser hiện đại (Chrome, Firefox, Safari, Edge phiên bản mới nhất hoặc N-1). Ứng dụng mobile native nằm ngoài phạm vi phiên bản đầu tiên.
- Mỗi user thuộc duy nhất một tenant. Trường hợp user cần truy cập nhiều tenant sẽ sử dụng nhiều tài khoản riêng biệt.
- Hệ thống email gửi invitation/notification đã có sẵn hoặc sử dụng dịch vụ bên ngoài (transactional email service).
- Thanh toán (payment processing) nằm ngoài phạm vi hệ thống này — billing module chỉ quản lý plan assignment và usage tracking, không xử lý giao dịch tài chính.
- Video hosting sử dụng dịch vụ bên ngoài (streaming service) — hệ thống quản lý metadata và liên kết, không lưu trữ/transcode video trực tiếp.
- Ngôn ngữ giao diện ban đầu hỗ trợ tiếng Việt và tiếng Anh. Thêm ngôn ngữ khác qua i18n framework.
- SSO integration (SAML, OIDC) chỉ khả dụng cho plan Enterprise — các plan khác sử dụng email/password authentication.
- Constitution đã xác lập: hệ thống tuân thủ OWASP Top 10, GDPR data privacy, RBAC tại API layer, audit logging cho hành động quan trọng.
