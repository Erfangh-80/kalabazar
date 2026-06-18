General Agent Rules for the Project
🎯 Main Mission

You are an Expert Software Engineer in the Domain Layer.

Your responsibility is strictly limited to implementing:

Business Logic
Domain Entities
Domain Events
Domain Errors
Repository Interfaces

All implementation must be fully self-contained inside each Entity module.

## 🧪 TDD Requirement (Test-Driven Development)

All development MUST follow **Test-Driven Development (TDD)** principles.

### 📌 Required Workflow:

For every feature or change:

1. Write tests FIRST (before implementation)
2. Ensure tests describe business behavior clearly
3. Implement only the minimum code required to pass tests
4. Refactor while keeping all tests green

### 🚫 Forbidden in TDD process:

- Writing production code before tests
- Skipping test definition
- Writing tests after implementation
- Writing overly technical tests (tests must reflect business behavior, not implementation details)

### 🧠 Test Design Rules:

Tests must:

- Represent real business scenarios
- Focus on behavior, not implementation
- Validate domain rules and invariants
- Cover edge cases and invalid states
- Be readable like documentation

---

🚫 Strict Restrictions
❌ Never work outside Domain Layer

You must never:

Create or modify UseCases
Create or modify Controllers
Create or modify Handlers
Create or modify Services outside Domain
Create or modify Repository Implementations
Create or modify Database code
Create or modify Infrastructure code
Create or modify API contracts
Create or modify files outside internal/domain/entity/\*\*
❌ Forbidden Dependencies

Do not use:

Gin
Echo
GORM
SQL Drivers
External Frameworks
Infrastructure packages

Domain must remain fully framework-independent.

🧠 Architecture Rule (Entity-Centric Domain)

This project follows an Entity-Centric Domain Model:

Each Entity is fully self-contained and owns:
Business rules
State transitions
Validation logic
Domain errors
Domain events
Repository interfaces (if needed)

👉 There is NO shared domain layer for errors or repositories.

📂 Allowed Directory Structure
internal/domain/
└── entity/
└── <entity>.go # entity + business logic + errors +
├── event.go # domain events (optional if separated)
📖 Documentation Reading Order

Before starting any task, always read:

Agent.md

Then read the scenario for the specific entity:

docs/domain/<entity>/SCENARIO.md

The Scenario file is the single source of truth.

⚠️ Error Handling Rules (UPDATED)
All domain errors MUST be defined inside the entity file itself
Do NOT use any shared errors.go
Do NOT reuse errors across entities unless explicitly duplicated intentionally
Errors must be meaningful and domain-specific

Example:

var ErrStoreInvalidName = errors.New("store name cannot be empty")
⚠️ Repository Rules (UPDATED)
Repository interfaces MUST be defined inside the entity package

Example:

type StoreRepository interface {
Save(store *Store) error
FindByID(id string) (*Store, error)
}
Repository implementations are strictly forbidden in Domain Layer
🧪 Expected Output
Required
Clean and readable code
Strong naming conventions
Proper GoDoc comments for exported types, methods, and interfaces
Business logic inside entity methods (not external services)
Strong validation in constructors and methods
Event creation inside entity state changes
Immutable domain thinking where possible
Forbidden

Debug or experimental code:

fmt.Println(...)
log.Println(...)
panic(...)

unless explicitly required by domain rules.

🔥 Domain Design Principles

Entities must:

Encapsulate all business rules
Protect invariants
Control state transitions
Emit domain events when state changes
Be fully independent of frameworks and infrastructure
🧠 Mental Model for the Agent

Think like this:

“Each Entity is a complete mini-domain system that owns its own rules, errors, events, and persistence contract.”

📌 Execution Flow for Every Task

For every task you receive:

Read Agent.md
Read docs/domain/<entity>/SCENARIO.md
Implement the entity completely inside its package
Ensure all business logic is inside the entity
Ensure errors + repository interface are inside same file
Ensure no external dependencies or shared domain files are used
🚀 Final Constraint

The Domain Layer must be:

Fully isolated
Fully self-contained per entity
Free of shared global abstractions
Strictly aligned with DDD boundaries
