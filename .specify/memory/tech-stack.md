# LMS SaaS Multi-Tenant Tech Stack & Architecture

**Status**: Finalized (2026-03-26)
**Scope**: 001-lms-saas-system

## Complete Tech Stack

### Frontend
- **Next.js** (React) + **TypeScript**
- **Turborepo** — monorepo build system for JavaScript and TypeScript codebases
- **shadcn/ui** (Base UI) — base component library
- **openapi-typescript** — generate TypeScript types từ spec

### Backend
- **Golang** + **Gin Web Framework**
- **oapi-codegen** — generate Go types + Gin server interface từ spec

### API Contract
- **OpenAPI 3.1** — source of truth, versioned trong repo
- **Spectral** — spec linting, enforce custom rules
- **Schemathesis** — contract/fuzz testing tự động từ spec

### Database & Storage
- **PostgreSQL** — primary database
- **Redis** — cache + session
- **Rustfs** (S3-compatible) — object storage cho tài liệu, bài nộp
- **Migration DB**: `golang-migrate/migrate`

### Auth
- **JWT** + Refresh Token
- **OAuth2** — Google integration

### Infrastructure & DevOps
- **Docker** containers
- **GitHub Actions** — CI/CD
- **Structured logging** (JSON format)
- Health check endpoints

## Architecture Approach

### Multi-tenant Data Model
- **Sử dụng Shared Schema**: Tất cả các tenant dùng chung database và schema, được phân tách hoàn toàn bằng cột `tenant_id` trên mọi bảng dữ liệu (Pool model). Phù hợp để quản lý vòng đời ứng dụng đơn giản và scale tốt cho mốc 100 benchmark tenants.

### File Upload Flow (Presigned URL Pattern)
1. **Client → Go**: Request lấy presigned URL an toàn.
2. **Client → Rustfs (S3)**: Upload trực tiếp file lên storage, bypass server Go để tránh nghẽn băng thông.
3. **S3 event → Worker**: Object storage trigger event tới async worker (Go/Rust).
4. **Worker xử lý file**: Resize ảnh, scan malware, lấy metadata video, v.v.
5. **Worker → DB**: Lưu metadata chuẩn xác vào PostgreSQL.
