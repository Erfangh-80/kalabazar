## Mission

Build software using:

- TDD (Test Driven Development)
- Clean Architecture
- Domain First Design

The goal is correctness, simplicity, maintainability, and testability.

---

## Development Order

Layers must be implemented strictly from inside to outside.

Allowed orders:

``text
Domain
→ Application
→ Infrastructure
→ Presentation

```

The agent MUST NOT start a higher layer until the current layer is fully completed and approved.

---

## TDD Rules

For each feature:

``text
1. Write Test
2. Run Test (Fail)
3. Implement Minimum Code
4. Run Test (Pass)
5. Refactor
```

Never write production code before a failing test exists.

---

## Domain Layer Rules

Responsibilities:

- Entities
- Value Objects
- Domain Events
- Domain Errors
- Repository Interfaces
- Domain Services

Forbidden:

- Database
- HTTP
- gRPC
- Frameworks
- ORM
- External Libraries
- Infrastructure Concerns

Domain must be pure business logic.

---

## Application Layer Rules

Responsibilities:

- Use Cases
- DTOs
- Commands
- Queries
- Application Validation
- Use Case Orchestration

Forbidden:

- Database Implementations
- HTTP Handlers
- Framework Logic
- Infrastructure Logic

Application depends only on Domain.

---

## Code Quality Rules

Code must be:

- Readable
- Simple
- Small
- Explicit
- Testable

Preference:

- Clear names
- Small functions
- Small files
- Composition over complexity

Avoid:

- Clever code
- Unnecessary abstractions
- Premature optimization
- Dead code

---

## Dependency Rule

``text
Application
↓
Domain

```

Never violate this rule.

---

## Output Rules

Before writing code:

1. Explain what will be implemented.
2. Explain why it belongs to the current layer.
3. Generate tests first.
4. Generate production code after tests.
5. Stop when the requested layer is completed.

```

````

```md
## Scope Rule

Implement only the requested feature.

Do not generate future features unless explicitly requested.
````

```md
## Architecture Rule

If a responsibility belongs to another service or bounded context,
do not implement it.

Consume events instead of owning external business processes.
```

## You should create the folder structure like this:

.
├── cmd/
│ └── main.go
│
├── docs/
│ └── AGENT.md
│
├── internal/
│
│ ├── domain/
│ │
│ │ ├── seller/
│ │ │ ├── seller.go
│ │ │ ├── errors.go
│ │ │ ├── repository.go
│ │ │ ├── validation.go
│ │ │ └── events.go
│ │ │
│ │ ├── store/
│ │ ├── warehouse/
│ │ ├── product/
│ │ ├── inventory/
│ │ ├── campaign/
│ │ ├── commission/
│ │ ├── settlement/
│ │ └── payout/
│ │
│ ├── application/
│ │ ├── seller/
│ │ │ ├── create_seller.go
│ │ │ ├── verify_seller.go
│ │ │ ├── dto.go
│ │ │ └── create_seller_test.go
│ │ ├── store/
│ │ ├── warehouse/
│ │ ├── product/
│ │ ├── inventory/
│ │ ├── campaign/
│ │ ├── commission/
│ │ ├── settlement/
│ │ └── payout/
│  
│
├── test/
│
│ ├── unit/
│ │
│ │ ├── domain/
│ │ └── application/
│
└── go.mod
