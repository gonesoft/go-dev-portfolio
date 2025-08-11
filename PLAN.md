# Go + Postgres + AWS Roadmap (20 Weeks)

## Purpose

A structured plan to master backend development with **Go**, **PostgreSQL**, and **AWS**.  
This document serves as a living roadmap, updated as I progress.

---

Day 1 – Project Setup
• ✅ Defined stack: Go + PostgreSQL + AWS
• ✅ Created README.md with project goals
• ✅ Created PLAN.md with roadmap
• ✅ Set up GitHub repo gonesoft/go-dev-portfolio
• ✅ Added docker-compose.yml for PostgreSQL
• ✅ Created init.sql to seed users table
• ✅ First commit & push to GitHub

⸻

Day 2 – Database Connection
• ✅ Created .env for DB credentials
• ✅ Installed and used godotenv to load .env
• ✅ Connected Go to PostgreSQL using database/sql and lib/pq
• ✅ Queried and printed users from DB
• ✅ Added Makefile with commands (db-up, db-down, run)
• ✅ Committed and pushed Day 2 changes

⸻

Day 3 – Basic HTTP API (No Framework)
• ✅ Implemented /users endpoint using Go’s net/http
• ✅ Refactored to standard Go folder structure:
cmd/api/main.go
internal/db/db.go
internal/users/model.go
internal/users/handler.go
curl http://localhost:8082/users

Day 4 – Single User & Create User Endpoints
• ✅ Implemented GET /users/{id} to fetch a single user by ID
• ✅ Implemented POST /users to create a new user with JSON input
• ✅ Added route handling for GET and POST requests
• ✅ Tested with curl:

curl http://localhost:8080/users
curl http://localhost:8080/users/1
curl -X POST http://localhost:8080/users \
 -H "Content-Type: application/json" \
 -d '{"name":"Jane Doe","email":"jane@example.com"}'

---

## Week 0 (Day 1 — Today)

- Install/update Go, Docker, psql.
- Initialize GitHub repo with README.md & PLAN.md.
- Create local Postgres DB via Docker (`users` table).
- Complete first modules of Go Tour.
- Commit initial code (`hello.go`, SQL script).

---

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
