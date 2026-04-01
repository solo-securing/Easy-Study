# Quickstart: LMS SaaS Multi-Tenant

Tài liệu này hướng dẫn cách bootstrap môi trường local và validate các planning contracts cho feature `001-lms-saas-system`.

## 1. Điều kiện tiên quyết

- Môi trường phát triển Linux/macOS
- Docker và Docker Compose
- Go 1.26+
- Node.js 22+
- pnpm 9+

## 2. Clone project và cài dependencies

```bash
cd /home/dev/Projects/tmp/Easy-Study

# Backend dependencies
cd backend/apps/api
go mod download

# Frontend dependencies
cd /home/dev/Projects/tmp/Easy-Study/frontend
pnpm install
```

## 3. Khởi động hạ tầng dịch vụ

Dùng file compose cấp project cho các dịch vụ dùng chung ở local:

```bash
cd /home/dev/Projects/tmp/Easy-Study/dev
docker compose up -d
```

Nếu chỉ cần DB stack cho API local:

```bash
cd /home/dev/Projects/tmp/Easy-Study/backend/apps/api
make docker-run
```

## 4. Chạy backend API

```bash
cd /home/dev/Projects/tmp/Easy-Study/backend/apps/api
make run
```

Tùy chọn live reload:

```bash
make watch
```

## 5. Chạy frontend apps

```bash
cd /home/dev/Projects/tmp/Easy-Study/frontend
pnpm turbo dev --filter=web --filter=docs
```

## 6. Quy trình contract-first

### 6.1 Validate OpenAPI contract

```bash
cd /home/dev/Projects/tmp/Easy-Study
pnpm --dir frontend dlx @apidevtools/swagger-cli validate specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml
```

Tùy chọn Spectral lint (cần có Spectral ruleset trong repository):

```bash
cd /home/dev/Projects/tmp/Easy-Study
pnpm --dir frontend dlx @stoplight/spectral-cli lint specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml
```

### 6.2 Generate backend interfaces từ OpenAPI

```bash
cd /home/dev/Projects/tmp/Easy-Study/backend/apps/api
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
  -generate types,gin,strict-server \
  -package api \
  /home/dev/Projects/tmp/Easy-Study/specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml \
  > internal/server/generated_api.go
```

### 6.3 Generate frontend TypeScript types

```bash
cd /home/dev/Projects/tmp/Easy-Study/frontend
pnpm dlx openapi-typescript \
  /home/dev/Projects/tmp/Easy-Study/specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml \
  -o packages/ui/src/api-types.ts
```

## 7. Test và quality gates

### Backend tests

```bash
cd /home/dev/Projects/tmp/Easy-Study/backend/apps/api
make test
make itest
```

### Frontend quality checks

```bash
cd /home/dev/Projects/tmp/Easy-Study/frontend
pnpm turbo run lint typecheck test
```

### Contract testing (sau khi API đã chạy)

```bash
schemathesis run \
  --checks all \
  /home/dev/Projects/tmp/Easy-Study/specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml \
  --base-url http://localhost:8080
```

## 8. Smoke checks cho phạm vi feature

- Super Admin tạo/suspend/reactivate tenant thành công.
- Tenant Admin import users bằng CSV và xem được lỗi theo từng dòng.
- Instructor chỉ publish được course hierarchy hợp lệ.
- Student quiz autosave và timeout auto-submit được lưu đầy đủ.
- Super Admin tạo custom plan và impersonate có audit logs.
- Global và tenant notifications được queue và quan sát được delivery failures.