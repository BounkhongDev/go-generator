# Go Generator

[![Go Reference](https://pkg.go.dev/badge/github.com/BounkhongDev/go-generator.svg)](https://pkg.go.dev/github.com/BounkhongDev/go-generator)

A simple and opinionated project generator for Go (Golang), designed to help you quickly scaffold RESTful APIs with a clean layered architecture (controllers, services, repositories).

**Documentation:** [pkg.go.dev/github.com/BounkhongDev/go-generator](https://pkg.go.dev/github.com/BounkhongDev/go-generator) · [Library (pkg/generator)](https://pkg.go.dev/github.com/BounkhongDev/go-generator/pkg/generator)

---

## ✅ Features

- Generate a new Go project structure
- Easily create new modules with full CRUD templates
- RESTful API with standard response format (success, pagination, validation error, general error)
- ECS-formatted JSON logs for Elasticsearch and Kibana
- Supports `macOS`, `Linux`, and `Windows`
- Automatically adds the generator to your system path

---

## 🧪 Testing

Run tests and coverage locally:

```bash
go test ./...
go test -cover ./...
```

With coverage profile and report:

```bash
make test      # run tests
make cover     # tests with coverage summary
make cover-profile   # coverage profile + func report
make cover-html      # generate coverage.html (open in browser)
```

CI runs tests and coverage on push/PR to `main` or `master` (see [.github/workflows/test.yml](.github/workflows/test.yml)).

---

## 📦 Requirements

- [Go (Golang)](https://golang.org/dl/) installed
- One of the following operating systems:
  - macOS
  - Linux
  - Windows

---

## 🚀 Installation

### All platforms (recommended)

If you have Go installed, the same steps work on **Windows, macOS, and Linux**:

```bash
go install github.com/BounkhongDev/go-generator/cmd/go-gen-r@latest
```

Ensure the Go bin directory is in your `PATH`:

- **macOS / Linux:** `$HOME/go/bin` (or `$GOPATH/bin` if set). Add to `PATH` in `~/.bashrc`, `~/.zshrc`, or `~/.profile` if needed.
- **Windows:** `%USERPROFILE%\go\bin` (e.g. `C:\Users\YourName\go\bin`). Add it to your user PATH in System Properties → Environment Variables → Path.

Then run from any terminal:

```bash
go-gen-r init
```

---

### Build from source

**1. Clone the repository**

```bash
git clone https://github.com/BounkhongDev/go-generator.git
cd go-generator
```

**2. Build the binary**

- **macOS / Linux:**
  ```bash
  go build -o go-gen-r ./cmd/go-gen-r
  ```
  Binary: `go-gen-r`

- **Windows (Command Prompt or PowerShell):**
  ```cmd
  go build -o go-gen-r.exe ./cmd/go-gen-r
  ```
  Binary: `go-gen-r.exe`

**3. Install / add to PATH**

| OS        | Method |
|-----------|--------|
| **macOS / Linux** | Run `./install.sh` to copy the binary to `/usr/local/bin` (requires write access). Or copy `go-gen-r` to a directory already in your `PATH`. |
| **Windows**       | In PowerShell (from the repo folder): `.\install.ps1`. If you get an execution policy error, run `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned` once, then run `.\install.ps1` again. Or manually add the folder containing `go-gen-r.exe` to your [Environment Variables → Path](https://learn.microsoft.com/en-us/windows/win32/procthread/environment-variables). |

### Windows: manual PATH setup

If you built the binary and did not use `install.ps1`:

1. Note the folder that contains `go-gen-r.exe` (e.g. `C:\go-generator`).
2. Open **Settings** → **System** → **About** → **Advanced system settings** (or search "environment variables").
3. Click **Environment Variables**. Under **User variables**, select **Path** → **Edit** → **New**.
4. Add the folder path (e.g. `C:\go-generator`). Confirm with **OK**.
5. Open a **new** terminal; run `go-gen-r`.

Now you can use `go-gen-r` from any terminal on all supported OSes.

---

## 🛠️ Usage

### 🔧 Initialize a New Project

```bash
mkdir my-project
cd my-project
go-gen-r init
```

You'll be prompted to enter your project name (module name for Go):

```
Enter Project Name: hrms-service
```

- Project names can include hyphens (e.g. `hrms-service`, `user-api`)
- Project names must not contain spaces

The generator will:

- Run `go mod init <projectName>`
- Install dependencies (Fiber, GORM, Viper, Zap, Validator, etc.)
- Create the project structure (config, database, routes, middleware, etc.)
- Generate an `example` module with full CRUD

---

### 🧱 Generate a New Module

Inside your project directory:

```bash
go-gen-r user
go-gen-r category
go-gen-r product
```

- Module names must not contain hyphens or spaces (use underscores: `user_account`)

This will generate:

- `internal/models/<module>.go`
- `internal/requests/<module>_request.go`
- `internal/responses/<module>_response.go`
- `internal/repositories/<module>_repository.go`
- `internal/services/<module>_service.go`
- `internal/controllers/<module>_controller.go`
- Test stubs and migrations

**Important:** You must manually register new modules in `main.go` and `routes/fiber_routes.go` if you add modules after the initial `init`.

---

### 📚 Use as a library

You can use the generator programmatically from your own Go code:

```go
import "github.com/BounkhongDev/go-generator/pkg/generator"

func main() {
    // Initialize a new project (run from the project directory)
    if err := generator.Init("my-api"); err != nil {
        log.Fatal(err)
    }

    // Generate a new module (run from an existing project root)
    if err := generator.GenerateModule("user"); err != nil {
        log.Fatal(err)
    }
}
```

- `generator.Init(projectName)` runs `go mod init`, installs dependencies, and creates the full project structure including an example module.
- `generator.GenerateModule(moduleName)` creates a new module (model, repository, service, controller, etc.); the project name is read from `go.mod` in the current directory.

---

## 📐 RESTful API Conventions

The generator produces APIs that follow RESTful conventions with a standard response format.

### HTTP Methods & Status Codes

| Action   | Method | Path              | Status Code     |
|----------|--------|-------------------|-----------------|
| List all | GET    | `/api/v1/:resource` | 200 OK          |
| Get by ID| GET    | `/api/v1/:resource/:id` | 200 OK       |
| Create   | POST   | `/api/v1/:resource` | 201 Created     |
| Update   | PUT    | `/api/v1/:resource/:id` | 200 OK      |
| Delete   | DELETE | `/api/v1/:resource/:id` | 204 No Content |

### API Versioning

All routes are prefixed with `/api/v1`:

```
GET    /api/v1/users
GET    /api/v1/users/1
POST   /api/v1/users
PUT    /api/v1/users/1
DELETE /api/v1/users/1
```

### Resource Naming

- Use **plural nouns** for resources: `/users`, `/categories`, `/products`
- Use **lowercase** with hyphens for multi-word resources: `/user-profiles`

### Response Format

#### ✅ 1. Success Response

Used when the request is successfully processed.

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { "id": 1, "name": "John" },
  "errors": null
}
```

| Field    | Description                                              |
|----------|----------------------------------------------------------|
| success  | Always `true` for successful responses                   |
| message  | Human-readable summary                                   |
| data     | Result (object, array, etc.)                             |
| errors   | Always `null` on success                                 |

#### ✅ 2. Pagination Response

Used when returning paginated data. Pagination values are always integers.

**Parameters:** `page` (current page, starts from 1), `limit` (items per page)

```json
{
  "success": true,
  "message": "Operation successful",
  "pagination": {
    "total_items": 0,
    "items_per_page": 10,
    "current_page": 1,
    "total_pages": 0,
    "next_page": 0,
    "previous_page": null
  },
  "data": [],
  "errors": null
}
```

| Field                     | Description                                                         |
|---------------------------|---------------------------------------------------------------------|
| pagination.total_items    | Total items matching the query                                      |
| pagination.items_per_page | Number of items per page (limit)                                    |
| pagination.current_page   | Current page number                                                 |
| pagination.total_pages    | `ceil(total_items / items_per_page)`                                |
| pagination.next_page      | Next page number; `0` when on last page                             |
| pagination.previous_page  | Previous page number; `null` on first page                          |

**Empty cases:** Use `[]` for empty list, `null` for no single object found.

#### ❌ 3. Validation Error Response

Used when the request has invalid or missing inputs (HTTP 422).

```json
{
  "success": false,
  "message": "Validation failed",
  "data": null,
  "errors": [
    { "field": "email", "message": "Email is required" },
    { "field": "password", "message": "Password must be at least 8 characters" }
  ]
}
```

| Field   | Description                                   |
|---------|-----------------------------------------------|
| success | Always `false`                                |
| message | Short explanation of the failure              |
| data    | Always `null`                                 |
| errors  | Array of `{ field, message }`                 |

#### ❌ 4. General Error Response

Used when the request fails for reasons other than validation (4xx/5xx).

```json
{
  "success": false,
  "message": "Something went wrong",
  "data": null,
  "errors": null
}
```

| Field   | Description                                   |
|---------|-----------------------------------------------|
| success | Always `false`                                |
| message | Short description of the error                |
| data    | Always `null`                                 |
| errors  | Always `null`                                 |

### What You Must Do

1. **Use plural resource paths** – e.g. `/users` not `/user`
2. **Return correct HTTP status codes** – 200, 201, 204 for success; 4xx/5xx for errors
3. **Use JSON for request/response bodies** with `Content-Type: application/json`
4. **Validate input** before processing (the generator adds validation tags)
5. **Handle errors consistently** using `responses.NewErrorResponse`, `responses.NewValidationError`
6. **Use idempotent methods** – GET, PUT, DELETE should be idempotent

---

## 📋 Log Output Samples (ECS format for Elastic/Kibana)

All logs use [Elastic Common Schema (ECS)](https://www.elastic.co/guide/en/ecs/current/index.html) for Elasticsearch and Kibana.

### HTTP Request Log (Fiber middleware – each request)

**Success (200):**
```json
{"@timestamp":"2026-02-02T10:30:00.000Z","log.level":"info","message":"http request","http.request.method":"GET","http.response.status_code":200,"event.duration":"12.5ms","client.ip":"192.168.1.1","url.path":"/api/v1/users","trace.id":"abc-123","service.name":"hrms-service"}
```

**Not Found (404):**
```json
{"@timestamp":"2026-02-02T10:30:01.000Z","log.level":"info","message":"http request","http.request.method":"GET","http.response.status_code":404,"event.duration":"2ms","client.ip":"192.168.1.1","url.path":"/api/v1/users/999","trace.id":"abc-124","service.name":"hrms-service"}
```

**Server Error (500):**
```json
{"@timestamp":"2026-02-02T10:30:03.000Z","log.level":"info","message":"http request","http.request.method":"POST","http.response.status_code":500,"event.duration":"150ms","client.ip":"192.168.1.1","url.path":"/api/v1/users","trace.id":"abc-126","service.name":"hrms-service"}
```

### Application Logs (Zap – `logs` package)

**Info:**
```json
{"@timestamp":"2026-02-02T10:30:00.123Z","log.level":"info","message":"User created","service.name":"hrms-service","trace.id":"abc-123","caller":"service.go:50","user_id":1}
```

**Error:**
```json
{"@timestamp":"2026-02-02T10:30:01.456Z","log.level":"error","message":"database connection failed","service.name":"hrms-service","trace.id":"abc-124","error.message":"connection refused","caller":"repository.go:30"}
```

**Debug:**
```json
{"@timestamp":"2026-02-02T10:30:00.999Z","log.level":"debug","message":"processing request","service.name":"hrms-service","caller":"handler.go:100"}
```

### ECS fields per log type

| Field | HTTP log | App log |
|-------|:--------:|:-------:|
| `@timestamp` | ✓ | ✓ |
| `log.level` | ✓ | ✓ |
| `message` | ✓ | ✓ |
| `trace.id` | ✓ | ✓ |
| `service.name` | ✓ | ✓ |
| `error.message` | — | ✓ (on error) |
| `http.request.method` | ✓ | — |
| `http.response.status_code` | ✓ | — |
| `event.duration` | ✓ | — |
| `client.ip` | ✓ | — |
| `url.path` | ✓ | — |

**Note:** Ensure `requestid` middleware runs before your handlers so `trace.id` is set. Use `logs.Info(msg, ctx.Context(), ...)` from handlers so the context carries `trace.id`.

---

## 📁 Generated Project Structure

```
.
├── main.go
├── go.mod
├── go.sum
├── example.config.yaml
├── config/          # Environment & config
├── database/        # PostgreSQL connection
├── errs/            # Application errors
├── logs/            # Logging
├── middleware/      # HTTP middleware
├── migrations/      # Database migrations
├── paginates/       # Pagination helpers
├── responses/       # API response helpers
├── routes/          # Route registration
├── validation/      # Request validation
└── internal/
    ├── controllers/ # HTTP handlers
    ├── models/      # GORM models
    ├── repositories/# Data access
    ├── requests/    # Request DTOs
    ├── responses/   # Response DTOs (internal)
    └── services/    # Business logic
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Bounkhong CHUANGTHEVY** – Backend Developer

---

## Donate

If this project helps you, you can support it via Lightning or BCEL One:

**Lightning**

- Wallet: `spoiledhate85@walletofsatoshi.com`

<img src="assets/lightning-qr.jpg" alt="Lightning QR" width="200">

**BCEL One**

- Account: `050120001275407001`

<img src="assets/bcel-qr.png" alt="BCEL QR" width="200">
