# General Agent Rules for the Project

## 🎯 Main Mission

You are an Expert Software Engineer working on Domain and Application Layers.

---

## 🧱 Domain Layer Responsibilities

- Business Logic
- Domain Entities
- Domain Events
- Domain Errors
- Repository Interfaces

---

## ⚙️ Application Layer Responsibilities

- Use Cases
- DTOs
- Command / Query Models
- Application Validation
- Entity Orchestration

---

# 🧪 TDD Requirement

- Write tests first
- Tests must describe business behavior
- All UseCase tests must be placed in:
  test/unit/usecase/\*\_test.go
- Implement minimal code to pass tests
- Refactor with tests green

---

## 🚫 Forbidden in TDD

- Production code before tests
- Skipping tests
- Tests after implementation
- Implementation-focused tests

---

# 🏗 Architecture Rules

- Clean Architecture enforced
- Domain owns business rules
- Application owns orchestration
- Infrastructure is external

---

## 🧱 Domain Rules

- Entities contain all business rules
- Entities control state transitions
- Entities emit domain events
- Entities define repository interfaces
- Must be framework independent

---

## ⚙️ Use Case Rules

- Orchestration only
- May call repositories
- May call entities
- May coordinate multiple entities
- May return DTOs

---

## 🚫 Use Case Forbidden

- Business logic
- Domain rules
- Infrastructure access
- ORM usage

---

# 🚫 Global Restrictions

- No Controllers
- No Handlers
- No Routers
- No Middleware
- No Infrastructure code
- No DB implementation
- No API contracts outside layers

---

# 🚫 Forbidden Dependencies

- Gin
- Echo
- Fiber
- GORM
- SQL Drivers
- Any framework or infra package

---

# 📂 Structure

Domain:
internal/domain/entity/<entity>.go

Application:
internal/application/usecase/<usecase_name>.go

---

# 📖 Required Reading

- Agent.md
- docs/domain/<entity>/SCENARIO.md
- docs/application/usecase/<usecase_name>/SCENARIO.md

---

# ⚠️ Errors

- Must be inside entity file
- No shared errors
- No cross-entity reuse

---

# ⚠️ Repository Rules

- Interfaces inside Domain
- No implementations in Domain

---

# 🧪 Output Rules

- Clean code
- Strong naming
- GoDoc required
- Thin UseCases
- Rich Entities
- Full TDD compliance

---

# 🚫 Debug Code

- No fmt.Println
- No log.Println
- No panic (unless domain requires)

---

# 🧠 Mental Model

- Entity = rules
- UseCase = flow
- Infrastructure = external
