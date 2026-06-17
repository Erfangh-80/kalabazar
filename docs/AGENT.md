# General Agent Rules for the Project

## 🎯 Main Mission

You are an **Expert Software Engineer in the Domain Layer**.

Your responsibility is limited to implementing:

- Business Logic
- Domain Entities
- Domain Events
- Repository Interfaces
- Domain Errors

inside the `internal/domain` package.

---

## 🚫 Strict Restrictions

### Never work outside the Domain Layer

You must never:

- Create or modify UseCases
- Create or modify Controllers
- Create or modify Handlers
- Create or modify Services outside Domain
- Create or modify Repository Implementations
- Create or modify Database code
- Create or modify Infrastructure code
- Create or modify API contracts
- Create or modify files outside the allowed Domain package

### Forbidden Dependencies

Do not use:

- Gin
- Echo
- GORM
- SQL Drivers
- External Frameworks
- Infrastructure packages

Domain must remain completely independent.

---

## ✅ Architectural Requirements

- All implementation must remain inside `internal/domain`
- Domain must have zero dependency on upper or lower layers
- Follow DDD principles
- Follow Clean Architecture principles

---

## 📂 Allowed Directory Structure

```text
internal/domain/
├── entity/          # Domain entities (with behavior)
├── event/           # Domain events
├── errors.go        # Domain errors
└── repository/      # Repository interfaces only
```

---

## 📖 Documentation Reading Order

Before starting any task, always read:

```text
AGENT.md
```

Then read the Scenario file related to the requested entity.

Example:

```text
docs/domain/store/Scenario.md
```

If the task is related to the Store entity, the above file is the source of truth and must be read before implementation.

---

## ⚠️ Error Handling Rules

All domain errors must be defined or reused from:

```text
internal/domain/errors.go
```

Do not create duplicate error definitions elsewhere.

---

## ⚠️ Repository Rules

All repository interfaces must be placed inside:

```text
internal/domain/repository/
```

Only repository interfaces are allowed.

Repository implementations are strictly forbidden.

---

## 🧪 Expected Output

### Required

- Clean code
- Readable code
- Meaningful naming
- Godoc comments for exported types, functions, methods, and interfaces
- Business rules implemented inside the Domain Layer
- Behavior-rich entities

### Forbidden

Debug or experimental code:

```go
fmt.Println(...)
```

```go
log.Println(...)
```

```go
panic(...)
```

unless explicitly required by domain rules.

---

## ✅ Domain Purity Rules

Entities must contain business behavior and invariants.

Entities must not depend on:

- DTOs
- HTTP models
- Request/Response objects
- ORM models
- Database schemas
- Infrastructure services

Domain objects must remain framework-independent and persistence-independent.

---

## ✅ Task Completion Checklist

Before considering a task complete, ensure that:

- AGENT.md has been read.
- The relevant `Scenario.md` has been read.
- Domain implementation has been completed.
- Domain errors are defined in `internal/domain/errors.go`.
- Repository interfaces are placed in `internal/domain/repository/`.
- No code outside the Domain Layer has been created or modified.
- All exported code includes proper Godoc comments.

```

```
