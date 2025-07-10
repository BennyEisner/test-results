# Hexagonal Architecture (Ports and Adapters) Explained Simply

## What is Hexagonal Architecture?

Hexagonal Architecture, also known as Ports and Adapters, is a way to design software so that the core logic (the "heart" of your app) is completely separated from outside concerns like databases, user interfaces, or external services. This makes your code easier to test, maintain, and adapt to new requirements.

## Why Use It?
- **Separation of concerns:** Keeps business logic isolated from technical details.
- **Testability:** You can test your core logic without needing a database or web server.
- **Flexibility:** Easily swap out databases, UIs, or APIs without changing your core logic.

## The Main Ideas

### 1. The Hexagon (the Core)
- Imagine your application as a hexagon (the shape is just for illustration).
- Inside the hexagon is your business logic: the rules and processes that make your app unique.

### 2. Ports
- **Ports** are interfaces (contracts) that define how the outside world can interact with your app.
- Think of a port like a USB port on your computer: it doesn't care if you plug in a mouse, keyboard, or printer, as long as the device fits the port.
- There are two types of ports:
  - **Driving Ports:** For things that call into your app (like a web controller or CLI command).
  - **Driven Ports:** For things your app calls out to (like a database or email service).

### 3. Adapters
- **Adapters** are the "plugs" that connect the outside world to your ports.
- For example, a REST API controller is an adapter that lets HTTP requests reach your app through a port.
- A database repository is an adapter that lets your app save or load data through a port.
- You can have many adapters for a single port (e.g., both a web UI and a CLI can use the same business logic).

## How Does It Work?

- The core app defines ports (interfaces) for what it needs (e.g., "save a user", "get orders").
- Adapters implement these ports using specific technologies (e.g., Postgres, MongoDB, HTTP, CLI).
- The core app never knows or cares about the details of the adapters.

## Visual Example

```
[ Web UI ]         [ CLI ]
     |                |
     v                v
 [ Driving Adapters (Controllers) ]
     |                |
     v                v
 [ Driving Ports (Interfaces) ]
     |                |
     v                v
   [  Hexagon: Core Business Logic  ]
     ^                ^
     |                |
 [ Driven Ports (Interfaces) ]
     ^                ^
     |                |
 [ Driven Adapters (DB, Email, APIs) ]
     ^                ^
[ Database ]      [ Email Service ]
```

## Real-World Analogy
- Think of a power strip (the hexagon) with different types of sockets (ports).
- You can plug in a lamp, a phone charger, or a fan (adapters) as long as they fit the socket.
- The power strip doesn't care what you plug in; it just provides power through the port.

## Key Benefits
- **No business logic leaks into the UI or database code.**
- **Easy to test:** You can plug in "fake" adapters for testing.
- **Easy to change:** Want to switch from a SQL database to NoSQL? Just write a new adapter.

## Conclusion
Hexagonal Architecture (Ports and Adapters) helps you build software that is clean, testable, and adaptable. By keeping your core logic separate from technical details, you can grow and change your app with confidence.

---

*Inspired by [this article on Medium](https://medium.com/ssense-tech/hexagonal-architecture-there-are-always-two-sides-to-every-story-bc0780ed7d9c).* 