# Test Results

## Overview

This project is a full-stack application designed to:

* Accept test results in **JUnit** and **ReadyAPI** formats from various CI systems (GitHub Actions, Travis, Jenkins) via a **Go-based CLI**
* Post the test results to a **Go-based REST API**
* Store and manage test results in a **PostgreSQL** database
* Display test results over time using a **React frontend dashboard**

The project provides experience in:

* Backend development with Go
* Frontend development with React
* Working with RESTful APIs
* Using relational databases
* Unit, integration, and end-to-end testing

---

## System Architecture Diagram

```
+------------------------+        +----------------+        +---------------------+
|   CLI / Build Agent    +------->+     REST API    +------->+  PostgreSQL DB       |
| (junit/readyapi input) |        | (Go-based)      |        | (test results store) |
+------------------------+        +--------+-------+        +-----------+---------+
                                            |
                                            v
                                    +---------------+
                                    | React Frontend|
                                    +---------------+
```

---

## Getting Started Guide

### Install Required Tools

* Go
* Node.js & npm
* Docker

### Clone the Repository

```bash
git clone https://github.com/your-org/fullstack-test-tracker.git
cd fullstack-test-tracker
```

### Start the Project

```bash
docker-compose up -d
cd api && go run main.go
cd ../frontend && npm install && npm start
```

### Seed Database

Insert test data manually or using a seed script.

### Run CLI

Point the CLI to sample JUnit/ReadyAPI result files and post them to the API.

---

## Contribution Guide

* Use feature branches
* Format code with `gofmt`, Prettier, etc.
* Write tests for new code
* Use pre-commit hooks if available

---

## Project Layout (Monorepo)

```text
fullstack-test-tracker/
├── api/               # Go REST API
│   ├── main.go
│   ├── handlers/
│   ├── models/
│   └── tests/
├── cli/               # Go CLI
│   ├── main.go
│   ├── parsers/
│   └── uploader/
├── frontend/          # React app
│   ├── public/
│   ├── src/
│   └── tests/
├── db/                # SQL migrations
│   └── schema.sql
├── docker-compose.yml
├── Makefile
├── README.md
└── docs/
```

---

## Milestone Plan

### Week 1: Setup & Orientation

* Install dependencies
* Run services locally
* Understand architecture and DB schema

### Week 2: Backend API

* Create endpoints for tests, suites, runs
* Connect to PostgreSQL
* Unit test handlers and models

### Week 3: CLI Tool

* Parse JUnit & ReadyAPI XML
* Post data to API
* Include metadata flags (CI type, run ID, etc.)

### Week 4: Frontend App

* Display list of test runs
* Filter by CI system and test status
* Drill-down into run details

### Week 5+: Testing & CI

* Add integration and E2E tests
* Set up CI pipelines (GitHub Actions)
* Optional: Add authentication

---

## Story Descriptions
* [ ] Setup local dev environment and verify Docker + Go + Node installation
* [ ] Run Postgres with Docker and connect via psql or GUI
* [ ] Create initial DB schema for test results and test suites
* [ ] Build a basic Go server with a health check endpoint
* [ ] Add API route to submit test run metadata (test name, suite, CI system, timestamp)
* [ ] Implement unit tests for API handlers
* [ ] Write CLI parser for JUnit XML files
* [ ] Send parsed JUnit data to the API using HTTP POST
* [ ] Expand CLI to support ReadyAPI XML format
* [ ] Add metadata flags to CLI for CI source, run ID, and suite ID
* [ ] Build frontend layout with a basic test run list view
* [ ] Fetch and render test data from the API using React Query
* [ ] Add filter UI for CI system and test status
* [ ] Build detail page for a single test run
* [ ] Implement integration tests for API + DB flow
* [ ] Add E2E test
* [ ] Set up GitHub Actions CI pipeline for Go tests and React build

---

## Key Concepts to Learn

* REST API design
* Git workflows
* Relational modeling
* Unit, integration, and E2E testing
* Containerized dev environments with Docker

---

## Optional

We can generate a GitHub repo scaffold if needed to jumpstart development.
