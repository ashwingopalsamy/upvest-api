# upvest-api

A dockerized Go web service that implement [**Upvest - Public APIs**](https://docs.upvest.co/api) with Postgres and Kafka in a simple event-driven architecture.

> **Note**: This repository is an independent **technical demonstration** showcasing my abilities in **API product design**, **event-driven architecture** and **distributed systems**.
> 
> It is neither affiliated with; nor endorsed by **Upvest** and does not attempt to replicate their official APIs or services. This project is strictly serves to demonstrate my **technical expertise**.

---

## Table of Contents

1. [**Implementation Architecture**](#1-implementation-architecture)
2. [**Implementation Highlights**](#2-implementation-highlights)
3. [**Repository Structure**](#3-repository-structure)
4. [**Endpoints**](#4-endpoints)
5. [**Getting Started**](#5-getting-started)
6. [**Design Highlights**](#6-design-highlights)
7. [**Future Enhancements**](#7-future-enhancements)

---

## 1. Implementation Architecture

- **Domain-Driven Design (DDD):** Enabled separation of concerns across domains like users, accounts, etc.
- **Event-Driven Architecture:** Kafka-powered event streams to handle asynchronous processes.
- **Scalability & Maintainability:** Modular repository structure and clean abstractions.

---

## 2. Implementation Highlights

- Postgres with Goose support for migrations to manage schema evolution.
- Kafka Message Broker for publishing and consuming events.
- Docker Compose for seamless containerized local development.
- Middleware support for Paging, Sorting, Offset, Limit for APIs.
- Comprehensive unit and integration tests using `sqlmock` and `testify`.

---

## 3. Repository Structure

```
upvest-api/
├── cmd/                   # Entrypoints for services
│   ├── upvest-api-publisher/
│   ├── upvest-api-subscriber/
├── internal/              # Core application logic
│   ├── domain/            # Domain models and validation
│   ├── pkg/               # Handlers and repository code
│   ├── event/             # Kafka publishers/subscribers
│   ├── middleware/        # API middleware
│   └── util/              # Utilities and mocks
├── schema/
│   ├── migrations/        # Database schema migrations
├── Makefile               # Helper commands for building/testing
├── Dockerfile             # Helper commands for building/testing
└── README.md              # Documentation
```

---

## 4. Endpoints

Each API adheres to principles of validation, structured request, responses and a clear error handling.

- **POST** `/users` – Create a user
- **GET** `/users` – Retrieve a paginated list of users
- **GET** `/users/{user_id}` – Fetch a specific user by ID
- **DELETE** `/users/{user_id}` – Offboard a user

---

## 5. Getting Started

### Setup Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/ashwingopalsamy/upvest-api.git
   cd upvest-api
   ```

2. Build and run the application, execute the endpoints:
   ```bash
   make up
   ```

3. Run tests:
   ```bash
   make test
   ```

---

## 6. Design Highlights

### Event-Driven Architecture
- **Publisher:** Trigger Kafka events on user creation, deletion, and data changes.
- **Subscriber:** Asynchronously process user-related events, e.g., offboarding workflows.

### Database Management
- **Postgres:** Schema migrations are tracked using the `goose` migration tool.
- **Domain Modeling:** User-centric tables with JSONB fields for flexible data storage.

### Middleware
- Centralized paging and sorting logic applied required list APIs.
- Structured error handling ensures consistent client responses.

---

## 7. Future Enhancements

- **OAuth2 Authentication:** Integrate secure API access control.
- **Accounts Domain:** Extend functionality to manage user accounts.
- **Cloud Deployment:** Add deployment scripts for GCP with Kubernetes.
- **Observability:** Attempt to integrate Prometheus for metrics and dashboards or OpenTelemetry.