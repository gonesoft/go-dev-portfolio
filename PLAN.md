# Go + Postgres + AWS Roadmap (20 Weeks)

## Purpose

A structured plan to master backend development with **Go**, **PostgreSQL**, and **AWS**.  
This document serves as a living roadmap, updated as I progress.

## Day 1 – Project Setup

- ✅ Defined stack: Go + PostgreSQL + AWS
- ✅ Created `README.md` with project goals
- ✅ Created `PLAN.md` with roadmap
- ✅ Set up GitHub repo `gonesoft/go-dev-portfolio`
- ✅ Added `docker-compose.yml` for PostgreSQL
- ✅ Created `init.sql` to seed `users` table
- ✅ First commit & push to GitHub

## Day 2 – Database Connection

- ✅ Created `.env` for DB credentials
- ✅ Installed and used `godotenv` to load `.env`
- ✅ Connected Go to PostgreSQL using `database/sql` and `lib/pq`
- ✅ Queried and printed users from DB
- ✅ Added `Makefile` with commands (`db-up`, `db-down`, `run`)
- ✅ Committed and pushed Day 2 changes

## Day 3 – Basic HTTP API (No Framework)

- ✅ Implemented `/users` endpoint using Go’s `net/http`
- ✅ Refactored to standard Go folder structure:
  ```
  cmd/api/main.go
  internal/db/db.go
  internal/users/model.go
  internal/users/handler.go
  ```
- ✅ Tested:
  ```bash
  curl http://localhost:8082/users
  ```

## Day 4 – Single User & Create User Endpoints

- ✅ Implemented `GET /users/{id}` to fetch a single user by ID
- ✅ Implemented `POST /users` to create a new user with JSON input
- ✅ Added route handling for GET and POST requests
- ✅ Tested:
  ```bash
  curl http://localhost:8080/users
  curl http://localhost:8080/users/1
  curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"name":"Jane Doe","email":"jane@example.com"}'
  ```

## Day 5 – Update & Delete Users

- ✅ Implemented `PUT /users/{id}` to update a user's name and email
- ✅ Implemented `DELETE /users/{id}` to remove a user by ID
- ✅ Added route handling for PUT and DELETE requests
- ✅ Tested:

  ```bash
  # Update user
  curl -X PUT http://localhost:8080/users/1 \
    -H "Content-Type: application/json" \
    -d '{"name":"Updated Name","email":"updated@example.com"}'

  # Delete user
  curl -X DELETE http://localhost:8080/users/1
  ```

## Day 6 - 8 – Validation, Partial Updates, Soft Delete & Pagination

- ✅ Added request validation for POST and PUT.
- ✅ Created unified JSON response helper for success/errors.
- ✅ Updated PUT /users/{id} to allow partial updates.
- ✅ Modified DELETE /users/{id} to perform soft deletes (added deleted_at column).
- ✅ Implemented pagination for GET /users with page & limit query params.
- ✅ Updated queries to ignore deleted users.
- ✅ Tested endpoints with curl for valid/invalid inputs.

## Day 9 – Search, Sorting, and Total Count for Users

- ✅ Enhanced `GET /users` to support:
  - Search by name or email using `?search=` query parameter (case-insensitive).
  - Sorting by `name` or `created_at` with `?sort=` and order control via `?order=asc|desc`.
  - Pagination with `?page=` and `?limit=`.
  - Total record count and total pages returned in the JSON response.
- ✅ Ensured soft-deleted users (`deleted_at` IS NOT NULL) are excluded from results.
- ✅ Tested with curl:
  ```bash
  curl "http://localhost:8083/users?search=john&sort=created_at&order=desc&page=1&limit=5"
  curl "http://localhost:8083/users?search=doe&page=2&limit=2"
  ```

## Day 10 – Unit Tests for Database Layer

- ✅ Set up a dedicated test database using Docker.
- ✅ Installed `testify` for assertions.
- ✅ Created unit tests for:
  - Creating a user.
  - Fetching a user by ID.
  - Listing users with pagination.
- ✅ Ensured tests run in isolation by cleaning data between tests.
- ✅ Added `TestMain` in repository test to ensure DB connection before tests.
- ✅ Cleared `users` table before running tests to avoid data collisions.
- ✅ Created `TestCreateUserAndFetch` to insert and fetch users from DB.
- ✅ Confirmed tests run successfully against `postgres-test` container.
- ✅ Validated DB retry logic in `Connect()` resolves timing issues.
- ✅ Ran tests successfully:
  ```bash
  make test
  ```

## Day 11 – Update & Delete Operations in Repository

• ✅ Implemented `UpdateUser` repository function to modify name/email by ID.
• ✅ Implemented `DeleteUser` repository function to soft-delete a user by setting `deleted_at` column.
• ✅ Added unit tests for `UpdateUser` and `DeleteUser`:

- Verified update modifies correct record and persists changes.
- Verified delete marks record as deleted without removing it from DB.
  • ✅ Confirmed tests run against `postgres-test` and pass consistently.
  • ✅ Maintained DB cleanup before tests to ensure deterministic results.

## Day 12 – Data Integrity & End-to-End (E2E) Tests

• ✅ Added `deleted_at IS NULL` filter to all SELECT queries in repository layer to exclude soft-deleted users from results.
• ✅ Implemented unique email validation in `CreateUser` repository method.
• ✅ Updated HTTP handlers to:

- Return 409 Conflict when attempting to create a user with an existing email.
- Return 404 Not Found when fetching a deleted or non-existent user.
  • ✅ Created integration tests using `httptest.NewServer` to:

1. Create a user via POST /users
2. Fetch created user via GET /users/{id}
3. Update user via PUT /users/{id}
4. Soft delete user via DELETE /users/{id}
5. Attempt to fetch deleted user (expect 404)
   • ✅ Verified all E2E tests pass against `postgres-test`.

## Day 13 – Robust Update & Delete Operations

• ✅ Added email uniqueness validation to `UpdateUser` repository method to prevent duplicate emails when updating.
• ✅ Modified `UpdateUser` and `DeleteUser` to return `sql.ErrNoRows` when no matching row exists or the user is soft-deleted.
• ✅ Updated HTTP handlers to:

- Return 404 Not Found when attempting to update or delete a non-existent/deleted user.
- Return 409 Conflict when updating to an email that already exists for another user.
  • ✅ Added E2E tests covering:

1. Updating a user with an existing email → expect 409 Conflict.
2. Updating a non-existent user → expect 404 Not Found.
3. Deleting a non-existent user → expect 404 Not Found.
   • ✅ Verified all repository and E2E tests pass for update/delete scenarios.

## Day 14 – Pagination & Sorting for GET /users

• ✅ Added pagination support to GET /users using query parameters:

- `limit` (number of records to return, default: 10)
- `offset` (number of records to skip, default: 0)
  • ✅ Added sorting support to GET /users:
- `sort` (column name to sort by, allowed: id, name, email, created_at)
- `order` (sort direction, allowed: ASC, DESC)
  • ✅ Validated query parameters to prevent SQL injection.
  • ✅ Updated repository function to accept pagination and sorting parameters and construct safe SQL queries using placeholders for limit/offset.
  • ✅ Updated HTTP handler to parse query params, apply defaults, and pass them to the repository layer.
  • ✅ Added unit tests for repository function covering:
- Default parameters
- Custom limit/offset
- Sorting by different fields in ascending and descending order
  • ✅ Added E2E tests for GET /users:

1. Fetch with default params.
2. Fetch with limit & offset.
3. Fetch with sorting by name ASC.
4. Fetch with sorting by name DESC.
5. Invalid sort or order param returns HTTP 400.

## Weeks 1–2 — Fundamentals

**Goal:** Build a basic REST API in Go with CRUD against Postgres, containerized with Docker.

**Key Deliverables:**

- CRUD for `users` table.
- Unit tests for DB layer.
- Integration tests with docker-compose.
- JWT auth skeleton.
- Deployment to AWS (Elastic Beanstalk/ECS or Lambda).
- Structured logging + background worker.

---

## Weeks 3–4 — API & DB Consolidation

- Add migrations (`golang-migrate`).
- Apply indexes, EXPLAIN queries, JSONB fields.
- Secure endpoints.
- Refactor to follow Go best practices.

---

## Weeks 5–8 — Patterns & Concurrency

- Read & apply **100 Go Mistakes**.
- Implement Go design patterns (from _Go Design Patterns_).
- Add worker queues, caching (Redis), rate limiting.
- Prepare microservice-ready architecture.

---

## Weeks 9–12 — AWS & Microservices

- Split into multiple services (e.g., users, auth).
- Use gRPC for service-to-service communication.
- Deploy to AWS ECS/Fargate or Lambda+API Gateway.
- Setup RDS Postgres.
- Implement CI/CD with GitHub Actions.

---

## Weeks 13–15 — System Design

- Study **System Design Interview** (Alex Xu) + **Grokking the SDI**.
- Create architecture diagrams for 3 projects.
- Document trade-offs and scalability considerations.

---

## Weeks 16–18 — Algorithms & Interviews

- Practice **Cracking the Coding Interview** + LeetCode (Go).
- Track 100+ solved problems.
- Mock interviews (technical + system design).

---

## Weeks 19–20 — Final Prep

- Polish portfolio projects, README files, and LinkedIn.
- Apply to jobs (10+/week).
- Reach out to recruiters.
- Conduct final interview rehearsals.

---

## Check-in Format

_TBD_
