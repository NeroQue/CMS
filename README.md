# Course Management System

A web app for managing and taking courses. Built with a Go backend, React frontend, and PostgreSQL database.

## What it does

- Browse and view courses, modules, and content
- Track user progress through courses
- Gamification features (XP, streaks)
- Admin tools for scanning and managing course content
- User profiles and sessions

## Tech stack

- **Backend:** Go, PostgreSQL, Goose (migrations), SQLC (query generation)
- **Frontend:** React, TypeScript, Vite
- **Infrastructure:** Docker, Docker Compose

## Prerequisites

- Docker and Docker Compose

## Getting started

1. Clone the repo

2. Create a `.env` file inside the `docker/` folder with the following variables:
   ```
   POSTGRES_USER=your_db_user
   POSTGRES_PASSWORD=your_db_password
   POSTGRES_DB=course_management
   BACKEND_PORT=8080
   FRONTEND_PORT=3000
   COURSES_BASE_DIR=/path/to/your/courses
   ```

3. Create a `.env` file inside the `frontend/` folder (you can copy `frontend/.env.example`):
   ```
   VITE_BASE_URL=http://localhost:8080
   ```
   Make sure the port matches the `BACKEND_PORT` you set in step 2.

4. Start everything with Docker Compose:
   ```
   cd docker
   docker compose up --build
   ```

5. Open the frontend at `http://localhost:3000` (or whatever `FRONTEND_PORT` you set).

## Project structure

```
backend/          Go API server
  cmd/api/        Entry point
  internal/       Core logic (handlers, services, database)
  pkg/            Shared packages (gamification, parser, etc.)
  sql/            SQL migrations and queries
frontend/         React app (Vite + TypeScript)
docker/           Dockerfiles and docker-compose config
```

## API docs

Swagger docs are available at `http://localhost:8080/swagger/` when the backend is running.

