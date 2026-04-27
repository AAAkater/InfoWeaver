# InfoWeaver AI Agent Instructions

## Project Overview

InfoWeaver is a multi-modal RAG (Retrieval-Augmented Generation) system built with Go backend and Vue frontend. The backend uses Echo framework with a clean architecture separating concerns into API, Service, Models, and Database layers.

## Architecture & Key Components

### Core Structure

- **API Layer** (`/api/v1/`): REST endpoints with Swagger documentation
- **Service Layer** (`/service/`): Business logic implementation
- **Models Layer** (`/models/`): Data structures and database entities
- **Database Layer** (`/db/`): Multi-database integration (PostgreSQL, Redis, MinIO, RabbitMQ, Milvus)
- **Middleware** (`/middleware/`): JWT authentication, CORS, logging, error handling
- **Config** (`/config/`): Viper-based configuration management

### Key Services

- **File Service**: Handles file uploads to MinIO object storage with PostgreSQL metadata
- **Dataset Service**: Manages datasets that files belong to (ownership-based access control)
- **User Service**: Authentication and user management with JWT tokens

### Database Integration Pattern

Each service uses direct database access through global DB instances:

- `db.PgSqlDB` for PostgreSQL operations
- `db.MinioClient` for MinIO object storage
- `db.RedisClient` for Redis caching
- Services use GORM for ORM operations with context-aware queries

## Development Workflow

### Local Development

```bash
# Backend setup
cd server
cp .env.example .env.local
go run main.go

# Frontend setup
cd web
pnpm i
pnpm dev
```

### Configuration

- Environment variables defined in `.env.example`
- Copy to `.env.local` for local development
- Key services: PostgreSQL (5432), Redis (6379), MinIO (9000)

### Testing

- Limited test coverage currently exists only in `/tests/`
- Run tests with: `go test -v ./tests`
- Focus on integration tests for database operations

### API Documentation

- Swagger docs auto-generated at `/swagger/*`
- API routes defined in `/api/v1/` with GoDoc comments
- Base path: `/api/v1`

## Coding Patterns & Conventions

### Error Handling

- Use predefined service errors: `ErrUnknown`, `ErrNotFound`, `ErrDuplicatedKey`
- API responses use `response.ResponseBase[T]` wrapper
- HTTP status codes: 400 (BadRequest), 401 (NoAuth), 403 (Forbidden)

### Authentication

- JWT middleware validates tokens on protected routes
- Current user extracted via `utils.GetCurrentUser(ctx)`
- All file/dataset operations require ownership validation

### File Upload Pattern

```go
// Example from fileApi.uploadFile
currentUser := utils.GetCurrentUser(ctx)
datasetID := ctx.FormValue("id")
// Validate dataset ownership
datasetService.GetDatasetInfoByID(ctx, datasetID, currentUser.ID)
// Upload to MinIO + create DB record
fileService.CreateFileInfo(ctx, currentUser.ID, datasetID, filename, fileType, fileReader, fileSize)
```

### Database Operations (GORM Patterns)

The service layer uses **two GORM APIs**: the modern `gorm.G[T]` generic API (preferred for most operations) and the legacy `db.PgSqlDB.Model()` API (used when Joins or projection-struct scanning is needed).

#### gorm.G[T] Generic API (preferred)

Use for: Create / First / Updates / Delete / Count. Return types are explicit — `(T, error)` or `(int64, error)`.

```go
// Create — insert a single record
dbRecord := models.Dataset{Name: "abc", OwnerID: ownerID}
return gorm.G[models.Dataset](db.PgSqlDB).Create(ctx, &dbRecord)

// First — query a single record, returns gorm.ErrRecordNotFound if not found
dbUser, err := gorm.G[models.User](db.PgSqlDB).
    Where("email = ?", email).
    First(ctx)
return &dbUser, err

// First with Select — query specific columns only
dbFile, err := gorm.G[models.File](db.PgSqlDB).
    Select("minio_path").
    Where("id = ? AND user_id = ?", fileID, ownerID).
    First(ctx)

// Updates — update records, returns (rowsAffected, error)
rowsAffected, err := gorm.G[models.Dataset](db.PgSqlDB).
    Where("id = ? AND owner_id = ?", id, ownerID).
    Updates(ctx, newDatasetInfo)
if rowsAffected == 0 {
    return gorm.ErrRecordNotFound  // zero affected rows = record not found
}

// Delete — delete records, returns (rowsAffected, error)
rowsAffected, err := gorm.G[models.Provider](db.PgSqlDB).
    Where("id = ? AND owner_id = ?", providerID, ownerID).
    Delete(ctx)
if rowsAffected == 0 {
    return gorm.ErrRecordNotFound
}

// Count — existence / count check
cnt, err := gorm.G[models.Dataset](db.PgSqlDB).
    Where("name = ? AND owner_id = ?", name, ownerID).
    Count(ctx, "*")
return cnt > 0, err

// Update single column
_, err := gorm.G[models.User](db.PgSqlDB).
    Where("id = ?", userID).
    Update(ctx, "password", hashedPassword)
```

#### db.PgSqlDB.Model() Legacy API

Use for: JOIN queries, scanning into projection structs (e.g. `ChunkInfo`).

```go
// Joins + Smart Select — Model declares the source table; Find's target struct
// automatically determines the SELECT columns
err = db.PgSqlDB.Model(&models.Chunk{}).
    Joins("JOIN files ON files.id = chunks.file_id").
    Where("files.dataset_id = ? AND files.user_id = ?", datasetID, ownerID).
    Order("chunks.id DESC").
    Offset((page - 1) * pageSize).
    Limit(pageSize).
    Find(&chunks).Error  // chunks is []ChunkInfo — GORM auto-selects only matching columns

// First with projection struct
result := db.PgSqlDB.Model(&models.Dataset{}).
    Where("id = ? AND owner_id = ?", id, ownerID).
    First(&dbDataset)  // dbDataset is *DatasetInfo — only selects matching fields

// Count with Joins
countResult := db.PgSqlDB.Model(&models.Chunk{}).
    Joins("JOIN files ON files.id = chunks.file_id").
    Where("files.dataset_id = ? AND files.user_id = ?", datasetID, ownerID).
    Count(&total)
if countResult.Error != nil { ... }
```

#### General Conventions

| Convention                | Description                                                                                                            |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| **Ownership validation**  | All queries must include `user_id = ?` / `owner_id = ?` in WHERE                                                       |
| **Pagination**            | `Offset((page-1)*pageSize).Limit(pageSize)`, paired with `Order("id DESC")`                                            |
| **Zero-row detection**    | Check `rowsAffected == 0` after Updates/Delete → `ErrNotFound`                                                         |
| **Existence check**       | Use `Count(ctx, "*") > 0`, not `First` + ignoring `ErrRecordNotFound`                                                  |
| **Context passing**       | Generic API's first argument is ctx; legacy API uses `ctx.Request().Context()` in handlers                             |
| **Smart field selection** | `Model(&Source{}).Find(&Projection{})` auto-selects only fields present in the target struct — no manual Select needed |
| **Return value handling** | Generic API returns `(result, error)` directly; Model API uses `result.Error`                                          |

### Service Initialization

- Global service instances: `var FileServiceApp = new(FileService)`
- No dependency injection - services access global DB instances directly

## Key Files to Reference

- **Main entry**: `main.go` - Shows initialization sequence
- **API routing**: `api/v1/enter.go` - All route registration
- **File service**: `service/file.go` - Complete file upload/download pattern
- **Database init**: `db/enter.go` - Multi-database connection setup
- **Configuration**: `config/enter.go` - All config structure definitions
- **Error responses**: `models/response/` - Standard response format

## External Dependencies

- **PostgreSQL**: Primary data store for metadata
- **MinIO**: Object storage for uploaded files
- **Redis**: Caching and session management
- **RabbitMQ**: Message queue (currently unused in main flow)
- **Milvus**: Vector database for embeddings (commented out in init)

## Important Notes

- **Ownership Security**: Every database operation must validate user ownership
- **Context Propagation**: Always pass context from API handlers through service layers
- **File Paths**: MinIO paths use format `{userID}/{filename}` for organization
- **Error Logging**: Use `utils.Logger` for structured logging with Zap
- **Swagger Comments**: Maintain GoDoc comments for API documentation generation
