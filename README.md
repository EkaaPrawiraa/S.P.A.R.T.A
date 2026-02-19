<div align="center">
	<img
		src="assets/spartan-helmet-with-plume.svg"
		alt="S.P.A.R.T.A logo"
		width="160"
		height="160"
	/>
	<h1>S.P.A.R.T.A</h1>
	<p><b>Strength. Planning. Analytics. Recovery. Training. Assistant.</b></p>
	<p>A full-stack fitness tracking dashboard with workouts, splits, nutrition, analytics, and AI coaching.</p>
</div>

---

## Monorepo

- `frontend/` — Next.js (App Router) web app
- `backend/` — Go (Gin) API + Postgres + Redis
- `postman/` — Postman collection + environment
- `output/` — product screenshots used in this README

## Screenshots

<div align="center">
	<img src="output/Dashboard_Light.png" alt="Dashboard (Light)" width="92%" />
</div>

<table>
	<tr>
		<td width="50%" align="center">
			<img src="output/Dashboard_Dark.png" alt="Dashboard (Dark)" width="100%" />
		</td>
		<td width="50%" align="center">
			<img src="output/Workouts_Main.png" alt="Workouts" width="100%" />
		</td>
	</tr>
	<tr>
		<td width="50%" align="center">
			<img src="output/Splits_Main.png" alt="Splits" width="100%" />
		</td>
		<td width="50%" align="center">
			<img src="output/Analytics.png" alt="Analytics" width="100%" />
		</td>
	</tr>
	<tr>
		<td width="50%" align="center">
			<img src="output/AI_Coaching.png" alt="AI Coaching" width="100%" />
		</td>
		<td width="50%" align="center">
			<img src="output/Planner_Main.png" alt="Planner" width="100%" />
		</td>
	</tr>
</table>

## Features

- Auth (register/login) with JWT
- Workouts: create sessions, view details
- Splits: multi-day split templates, activate/deactivate
- Exercises: searchable exercise library + media URLs
- Nutrition: daily tracking + trends
- Analytics: training volume and progression insights
- AI tools: split generation, overload suggestions, workout plans, explanations, coaching + daily motivation
- Planner: personalized recommendations

## Tech Stack

**Frontend**

- Next.js (App Router), TypeScript
- Tailwind CSS + shadcn/ui
- next-themes (light/dark mode)

**Backend**

- Go + Gin
- Postgres (persistence)
- Redis (caching + daily motivation)
- OpenAI-compatible API client (optional; only needed for AI features)

## Quick Start (Local)

### Prerequisites

- Go (for the API)
- Node.js 18+ (for the web app)
- Postgres
- Redis

### 1) Backend (API)

From the repo root:

```bash
cd backend
go mod download
```

#### Database setup

1. Create a database in Postgres (name is up to you).

2. Run migrations (SQL files are in `backend/migrations/`). Example using `psql`:

```bash
psql -U postgres -h localhost -p 5432 -d <your_db> -f backend/migrations/001_init.sql
psql -U postgres -h localhost -p 5432 -d <your_db> -f backend/migrations/002_admin_auth.sql
```

Alternatively, you can use `golang-migrate` (see `backend/migrations/README.md`).

#### Environment variables

The API best-effort loads `backend/.env` automatically on startup.

Minimal example (do not commit real secrets):

```env
# backend/.env
APP_ENV=local
APP_PORT=8080

DATABASE_URL=postgres://postgres:postgres@localhost:5432/<your_db>?sslmode=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=change_me

# Optional (only required for AI features)
OPENAI_API_KEY=
OPENAI_MODEL=gpt-4o-mini
OPENAI_BASE_URL=https://api.openai.com
```

#### Run the API

```bash
cd backend
go run ./cmd/api
```

Health check:

```bash
curl -sS http://localhost:8080/api/v1/health
```

### 2) Frontend (Web)

```bash
cd frontend
npm install
```

Create `frontend/.env.local`:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Run the dev server:

```bash
cd frontend
npm run dev
```

Open http://localhost:3000

## API Notes

- Base path: `/api/v1`
- CORS is configured for local dev (`http://localhost:3000`)
- Protected routes require: `Authorization: Bearer <token>`

Useful endpoints:

- `GET /api/v1/health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/exercises`
- `POST /api/v1/splits`
- `POST /api/v1/ai/motivation/reset`

## Postman

Import these files into Postman:

- `postman/SPARTA.postman_collection.json`
- `postman/SPARTA.local.postman_environment.json`

## Attribution

The Spartan helmet logo SVG in `assets/spartan-helmet-with-plume.svg` is derived from an SVG published on svgstud.io and is licensed under **CC BY-SA 4.0**. See the SVG metadata and the source file in `assets/spartan_helmet_with_plume_2.svg` for attribution details.
