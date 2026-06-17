## Scenario: Warehouse Management for Sellers

### 🧭 Overview

A seller needs a way to store and manage physical goods.  
To do this, the seller can register one or more warehouses.

Each warehouse represents a real physical location where goods are stored and managed.

---

### 👤 Actor

- Seller: a user who owns goods and manages storage locations

---

### 🏢 Core Concept: Warehouse

A warehouse is a physical storage location used by a seller.

Each warehouse has:

- A name to identify it
- A real-world address
- A physical location
- Storage capacity limits
- A current usage level

---

### 📌 Key Rules

- A seller can register more than one warehouse
- A warehouse cannot exist without being assigned to a seller
- Each warehouse represents a single physical place
- The storage capacity of a warehouse is limited
- A warehouse cannot exceed its maximum storage capacity
- A warehouse must have a valid physical address

---

### ⚙️ Main Capabilities (System Behavior)

The system should allow:

- Creating a warehouse for a seller
- Viewing all warehouses of a seller
- Updating warehouse information
- Activating or deactivating a warehouse
- Tracking how much storage space is used in a warehouse

---

### 🚫 Important Constraints

- A warehouse cannot be used if it is inactive
- A warehouse cannot exceed its capacity
- A warehouse must always belong to one seller
- A warehouse must represent a real physical location

---

### 📈 Business Understanding Notes

- Warehouses are treated as real-world assets
- Sellers manage their own storage spaces independently
- Each warehouse is controlled by a single seller, not shared between multiple sellers
