---
name: hexagonal-architecture
description: Guide LLMs in implementing hexagonal architecture (Ports & Adapters) for backend applications using Go. Covers feature-based organization, domain/application/infrastructure/interfaces layers, dependency injection, ports and adapters patterns, cross-feature communication, and LLM-friendly coding conventions.
---

# Hexagonal Architecture Skill - Backend (Go)

## Purpose
Guide LLMs in implementing hexagonal architecture (Ports & Adapters) for backend applications using Go.

## Core Principles

### 1. Feature-Based Organization
- Each feature is a self-contained vertical slice
- All code for a feature lives in one directory
- Minimal dependencies between features
- Optimized for LLM context windows

### 2. Dependency Rule
- Domain depends on nothing
- Application depends on domain only
- Infrastructure depends on domain + application
- Interfaces depend on application

### 3. Ports (Interfaces) vs Adapters (Implementations)
- Ports = interfaces defined by domain/application needs
- Adapters = technical implementations of ports
- Domain defines what it needs, infrastructure provides it

---

## Directory Structure

```
internal/
├── catalog/                    # Feature: Catalog management
│   ├── domain/                 # Domain layer (pure business logic)
│   │   ├── catalog.go         # Aggregate root
│   │   ├── endpoint.go        # Entity
│   │   ├── repository.go      # Port (interface)
│   │   └── events.go          # Domain events
│   ├── application/            # Application layer (use cases)
│   │   ├── create_catalog.go
│   │   ├── import_swagger.go
│   │   └── catalog_service.go
│   ├── infrastructure/         # Adapters (technical implementations)
│   │   ├── file_repository.go
│   │   └── sqlite_repository.go (if needed)
│   └── interfaces/             # External interfaces (HTTP, CLI)
│       ├── http_handlers.go
│       └── dto.go
│
├── request/                    # Feature: API request execution
├── invocation/                 # Feature: Flight recorder
├── sandbox/                    # Feature: Sandbox APIs
├── importexport/               # Feature: Import/Export formats
├── tab/                        # Feature: Tab management
├── settings/                   # Feature: Configuration
│
└── shared/                     # Truly shared code ONLY
    ├── events/
    │   └── bus.go
    ├── types/
    │   └── id.go
    └── errors/
        └── errors.go
```

---

## Layer Responsibilities

### Domain Layer (`domain/`)

**Rules:**
- NO external dependencies (no database, HTTP, JSON tags)
- Pure Go types and business logic
- Defines ports (interfaces) it needs
- Emits domain events

**What goes here:**
- Aggregates and Entities
- Value Objects
- Repository interfaces (ports)
- Domain Services (if needed)
- Domain Events
- Business validation rules

**Example:**
```go
// domain/catalog.go
package domain

import "time"

// Aggregate root
type Catalog struct {
    ID          string
    Name        string
    Description string
    Endpoints   []Endpoint
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Metadata    CatalogMetadata
}

// Business validation
func (c *Catalog) Validate() error {
    if c.Name == "" {
        return ErrCatalogNameRequired
    }
    if len(c.Name) > 255 {
        return ErrCatalogNameTooLong
    }
    return nil
}

// Business operation
func (c *Catalog) AddEndpoint(endpoint Endpoint) error {
    if err := endpoint.Validate(); err != nil {
        return err
    }
    c.Endpoints = append(c.Endpoints, endpoint)
    c.UpdatedAt = time.Now()
    return nil
}

// domain/repository.go
package domain

// Port - interface defined by domain needs
type CatalogRepository interface {
    Save(catalog *Catalog) error
    FindByID(id string) (*Catalog, error)
    FindAll() ([]*Catalog, error)
    Delete(id string) error
}
```

---

### Application Layer (`application/`)

**Rules:**
- Orchestrates domain objects
- Implements use cases
- Depends on domain layer only
- Uses ports (interfaces), not concrete implementations

**What goes here:**
- Use cases (one per file)
- Application services (orchestration)
- DTOs (if needed for use case input/output)

**Example:**
```go
// application/create_catalog.go
package application

import (
    "myapp/internal/catalog/domain"
    "myapp/internal/shared/events"
)

type CreateCatalogUseCase struct {
    repo     domain.CatalogRepository  // Port interface
    eventBus events.EventBus            // Port interface
}

func NewCreateCatalogUseCase(
    repo domain.CatalogRepository,
    eventBus events.EventBus,
) *CreateCatalogUseCase {
    return &CreateCatalogUseCase{
        repo:     repo,
        eventBus: eventBus,
    }
}

type CreateCatalogInput struct {
    Name        string
    Description string
}

func (uc *CreateCatalogUseCase) Execute(input CreateCatalogInput) (*domain.Catalog, error) {
    // Create domain object
    catalog := &domain.Catalog{
        ID:          generateID(),
        Name:        input.Name,
        Description: input.Description,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // Domain validation
    if err := catalog.Validate(); err != nil {
        return nil, err
    }
    
    // Persist via port
    if err := uc.repo.Save(catalog); err != nil {
        return nil, err
    }
    
    // Emit domain event
    uc.eventBus.Publish(domain.CatalogCreatedEvent{
        CatalogID: catalog.ID,
        Name:      catalog.Name,
    })
    
    return catalog, nil
}
```

---

### Infrastructure Layer (`infrastructure/`)

**Rules:**
- Implements ports defined by domain/application
- Contains all technical details (DB, HTTP, files)
- Depends on domain + application layers

**What goes here:**
- Repository implementations (SQLite, file, etc.)
- HTTP clients
- External API adapters
- Event bus implementations

**Example:**
```go
// infrastructure/file_repository.go
package infrastructure

import (
    "encoding/json"
    "myapp/internal/catalog/domain"
    "os"
    "path/filepath"
)

// Adapter - implements domain.CatalogRepository port
type FileCatalogRepository struct {
    dataDir string
}

func NewFileCatalogRepository(dataDir string) *FileCatalogRepository {
    return &FileCatalogRepository{dataDir: dataDir}
}

func (r *FileCatalogRepository) Save(catalog *domain.Catalog) error {
    // Technical implementation details
    filename := filepath.Join(r.dataDir, catalog.ID+".json")
    data, err := json.Marshal(catalog)
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

func (r *FileCatalogRepository) FindByID(id string) (*domain.Catalog, error) {
    filename := filepath.Join(r.dataDir, id+".json")
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var catalog domain.Catalog
    err = json.Unmarshal(data, &catalog)
    return &catalog, err
}

// ... implement other interface methods
```

---

### Interfaces Layer (`interfaces/`)

**Rules:**
- External interfaces (HTTP, CLI, WebSocket)
- Converts external requests to use case inputs
- Depends on application layer

**What goes here:**
- HTTP handlers
- DTOs for API requests/responses
- WebSocket handlers
- CLI commands

**Example:**
```go
// interfaces/http_handlers.go
package interfaces

import (
    "encoding/json"
    "myapp/internal/catalog/application"
    "net/http"
)

type CatalogHandlers struct {
    createUseCase *application.CreateCatalogUseCase
}

func NewCatalogHandlers(createUseCase *application.CreateCatalogUseCase) *CatalogHandlers {
    return &CatalogHandlers{createUseCase: createUseCase}
}

type CreateCatalogRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

func (h *CatalogHandlers) CreateCatalog(w http.ResponseWriter, r *http.Request) {
    // Parse HTTP request
    var req CreateCatalogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Convert to use case input
    input := application.CreateCatalogInput{
        Name:        req.Name,
        Description: req.Description,
    }
    
    // Execute use case
    catalog, err := h.createUseCase.Execute(input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return HTTP response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(catalog)
}
```

---

## Dependency Injection

All dependencies are injected at startup. Use constructor injection.

**Example:**
```go
// cmd/webserver/main.go
package main

func main() {
    // Infrastructure adapters
    catalogRepo := infrastructure.NewFileCatalogRepository("./data/catalogs")
    eventBus := events.NewInMemoryEventBus()
    
    // Application use cases
    createCatalog := application.NewCreateCatalogUseCase(catalogRepo, eventBus)
    
    // Interface handlers
    handlers := interfaces.NewCatalogHandlers(createCatalog)
    
    // Setup HTTP routes
    http.HandleFunc("/api/catalogs", handlers.CreateCatalog)
    
    http.ListenAndServe(":8080", nil)
}
```

---

## Cross-Feature Communication

### Option 1: Application Services (Synchronous)
```go
// application/request_orchestrator.go
package application

type RequestOrchestrator struct {
    executeRequest    *ExecuteRequestUseCase
    captureInvocation *invocation.CaptureInvocationUseCase
}

func (o *RequestOrchestrator) ExecuteAndCapture(input ExecuteRequestInput) error {
    // Execute request
    result, err := o.executeRequest.Execute(input)
    
    // Capture invocation (different feature)
    o.captureInvocation.Execute(invocation.CaptureInput{
        URL:      input.URL,
        Status:   result.Status,
        Response: result.Body,
    })
    
    return err
}
```

### Option 2: Domain Events (Asynchronous)
```go
// Emit event
eventBus.Publish(RequestExecutedEvent{
    RequestID: req.ID,
    Status:    200,
    Duration:  123,
})

// Different feature subscribes
// invocation/application/capture_listener.go
func (l *CaptureListener) OnRequestExecuted(event RequestExecutedEvent) {
    l.captureUseCase.Execute(...)
}
```

---

## LLM-Friendly Patterns

### 1. One File Per Concept
- One aggregate per file
- One use case per file
- One handler per file
- Makes search and context loading efficient

### 2. Explicit Names
- `CreateCatalogUseCase` not `Create`
- `FileCatalogRepository` not `FileRepo`
- `CatalogCreatedEvent` not `Created`
- LLMs can understand from name alone

### 3. Minimal Interfaces
```go
// Good - single responsibility
type CatalogRepository interface {
    Save(catalog *Catalog) error
    FindByID(id string) (*Catalog, error)
}

// Bad - too many responsibilities
type Repository interface {
    SaveCatalog(*Catalog) error
    SaveRequest(*Request) error
    SaveInvocation(*Invocation) error
}
```

### 4. Self-Documenting Code
```go
// Good - clear intent
func (c *Catalog) AddEndpoint(endpoint Endpoint) error {
    if err := endpoint.Validate(); err != nil {
        return fmt.Errorf("invalid endpoint: %w", err)
    }
    c.Endpoints = append(c.Endpoints, endpoint)
    return nil
}

// Bad - requires comments to understand
func (c *Catalog) Add(e Endpoint) error {
    // validate first
    if err := e.Validate(); err != nil {
        return err
    }
    // then add
    c.Endpoints = append(c.Endpoints, e)
    return nil
}
```

---

## Testing Strategy

### Domain Tests (No Dependencies)
```go
// domain/catalog_test.go
func TestCatalog_AddEndpoint(t *testing.T) {
    catalog := &Catalog{ID: "1", Name: "Test"}
    endpoint := Endpoint{ID: "e1", Name: "GET /users"}
    
    err := catalog.AddEndpoint(endpoint)
    
    assert.NoError(t, err)
    assert.Len(t, catalog.Endpoints, 1)
}
```

### Application Tests (Mock Ports)
```go
// application/create_catalog_test.go
func TestCreateCatalogUseCase(t *testing.T) {
    mockRepo := &MockCatalogRepository{}
    mockBus := &MockEventBus{}
    useCase := NewCreateCatalogUseCase(mockRepo, mockBus)
    
    catalog, err := useCase.Execute(CreateCatalogInput{
        Name: "My API",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "My API", catalog.Name)
    assert.True(t, mockRepo.SaveCalled)
}
```

### Infrastructure Tests (Real Dependencies)
```go
// infrastructure/file_repository_test.go
func TestFileCatalogRepository_Save(t *testing.T) {
    tmpDir := t.TempDir()
    repo := NewFileCatalogRepository(tmpDir)
    
    catalog := &domain.Catalog{ID: "1", Name: "Test"}
    err := repo.Save(catalog)
    
    assert.NoError(t, err)
    // Verify file exists
    assert.FileExists(t, filepath.Join(tmpDir, "1.json"))
}
```

---

## Common Mistakes to Avoid

### ❌ Domain depending on infrastructure
```go
// BAD - domain imports database package
package domain

import "database/sql"

type Catalog struct {
    DB *sql.DB  // NO!
}
```

### ❌ Leaking infrastructure details into domain
```go
// BAD - JSON tags in domain model
type Catalog struct {
    ID   string `json:"id" db:"id"`  // NO!
    Name string `json:"name"`
}
```

### ❌ Use case knowing about HTTP
```go
// BAD - use case returns HTTP status
func (uc *CreateCatalogUseCase) Execute(input Input) (int, error) {
    // ...
    return http.StatusCreated, nil  // NO!
}
```

### ✅ Correct layering
```go
// GOOD - domain is pure
package domain

type Catalog struct {
    ID   string
    Name string
}

// GOOD - infrastructure handles persistence details
package infrastructure

type catalogDTO struct {
    ID   string `json:"id" db:"id"`
    Name string `json:"name"`
}

func (r *Repository) Save(catalog *domain.Catalog) error {
    dto := catalogDTO{
        ID:   catalog.ID,
        Name: catalog.Name,
    }
    // save dto...
}
```

---

## Quick Reference

### Where does this code go?

| Code Type | Layer | Example |
|-----------|-------|---------|
| Business rules | Domain | `catalog.Validate()` |
| Entities/Aggregates | Domain | `type Catalog struct` |
| Interfaces (ports) | Domain | `type CatalogRepository interface` |
| Use cases | Application | `CreateCatalogUseCase` |
| Orchestration | Application | `RequestOrchestrator` |
| Database code | Infrastructure | `FileCatalogRepository` |
| HTTP handlers | Interfaces | `CreateCatalogHandler` |
| DTOs for API | Interfaces | `CreateCatalogRequest` |

### Dependency Flow
```
Interfaces → Application → Domain
                ↓
          Infrastructure
```

### File Naming
- Domain: `catalog.go`, `endpoint.go`, `repository.go`
- Application: `create_catalog.go`, `import_swagger.go`
- Infrastructure: `file_repository.go`, `sqlite_repository.go`
- Interfaces: `http_handlers.go`, `dto.go`

---

## Summary

1. **Feature-based** organization for LLM context efficiency
2. **Dependency rule** always flows inward (toward domain)
3. **Ports** defined by domain/application needs
4. **Adapters** provide technical implementations
5. **Minimal shared code** - duplicate when in doubt
6. **Explicit naming** for LLM discoverability
7. **One concept per file** for clean context loading