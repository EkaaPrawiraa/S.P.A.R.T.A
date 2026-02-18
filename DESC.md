---

# 🧠 PROJECT CONTEXT BLUEPRINT

## Project Name: AI Gym OS (Personal AI Training System)

---

# 1️⃣ PROJECT VISION

AI Gym OS is a **personal AI-powered training decision system**, not just a workout tracker.

The system combines:

* Deterministic training logic
* Rule-based progression engine
* AI reasoning layer
* Data-driven analytics dashboard

It is:

* Personally useful
* Engineer-designed
* Portfolio-level project

Core philosophy:

> Deterministic logic first. AI as reasoning layer on top.

---

# 2️⃣ CORE FEATURES

## A. Workout Tracking System

* Log workout sessions
* Log sets (exercise, reps, weight, RPE, rest)
* Auto-calculate:

  * Volume
  * Estimated 1RM
  * PR detection
  * Weekly muscle distribution

---

## B. Split Template System

User can:

* Create and save split templates:

  * Arnold Split
  * Upper / Lower
  * Push Pull Legs
  * Custom split

Each template contains:

* Days
* Exercises per day
* Default sets/reps

User can:

* Import split into calendar
* Modify exercises per day
* Convert workout → template

---

## C. Manual Workout Builder

User can:

* Manually choose exercises
* Add sets
* Save session
* Optionally save as reusable template

Burnfit-style flexibility.

---

## D. Daily Protein Tracker (Simple)

User logs:

* protein grams
* protein target

System shows:

* % completion
* daily progress
* weekly protein trend

No complex calorie tracking (v1 intentionally simple).

---

## E. Progressive Overload Engine (Rule-Based)

Deterministic logic:

Examples:

* If RPE ≤ 7 → increase weight next session
* If RPE ≥ 9 → maintain or reduce
* If stagnation ≥ 3 weeks → trigger plateau flag

Outputs:

* Suggested next weight
* Rep scheme change suggestion
* Deload suggestion

This is NOT AI.
This is algorithmic logic.

---

## F. AI Cognitive Layer (OpenAI Integration)

AI is used for interpretation, not math.

### AI Use Cases:

1. AI Coach Chat

   * Context-aware responses
   * Reads recent sessions + trends

2. Weekly Insight Generator

   * Explains performance trends

3. Adaptive Planning Assistant

   * Suggest weekly adjustments

4. Personalized Motivation Quote

   * Based on recent performance

AI never replaces deterministic calculations.
AI reasons on structured outputs.

---

# 3️⃣ SYSTEM ARCHITECTURE

## Backend

Language: Go
Framework: Gin
Database: PostgreSQL
Cache: Redis (optional for quotes + session caching)
AI: OpenAI API

Architecture Style:
→ Modular Monolith (NOT microservices initially)

Modules:

* auth
* workouts
* splits
* analytics
* planner
* nutrition
* ai-coach

Reason:
Cleaner dev process, easier debugging, still architecturally clean.

---

## Frontend

Framework: Vue (Dashboard heavy)
Purpose:

* Data visualization
* Workout logging
* Split management
* Planner display
* AI chat UI

Frontend communicates only via REST API.

No business logic inside frontend.

---

## Database Design Philosophy

Relational and structured.

Core Tables:

users
exercises
split_templates
split_days
split_day_exercises
workout_sessions
workout_sets
nutrition_daily

Analytics derived via queries (not stored blindly).

Important:
Workout data must be clean and normalized for AI reasoning later.

---

# 4️⃣ INTELLIGENCE LAYERS

Layer 1: Data Layer

* Raw logs
* Protein intake
* Templates

Layer 2: Deterministic Engine

* Volume
* PR
* Overload suggestion
* Plateau detection

Layer 3: AI Reasoning Layer

* Interpretation
* Advice
* Motivation
* Strategic adjustment

---

# 5️⃣ DEVELOPMENT ORDER

Phase 1:

* Database schema
* Workout logging
* Split templates
* Protein tracker

Phase 2:

* Progressive overload engine
* Analytics dashboard

Phase 3:

* AI integration (coach + insights)

---

# 6️⃣ DESIGN PRINCIPLES

* Backend-first development
* Deterministic before AI
* AI only where interpretation needed
* Structured data → AI prompt
* Clean architecture (domain separation)

---

# 7️⃣ LONG-TERM VISION

System evolves into:

> Hybrid Intelligence Training Engine
> (Algorithmic + AI reasoning combined)

Potential future upgrades:

* Recovery scoring
* Vector memory for long-term trend awareness
* Wearable integration
* Public SaaS version

---

# 8️⃣ PROJECT TYPE

Personal tool first.
Portfolio killer second.
Scalable architecture possible.


---

# Master Strategic Order (Portfolio-Killer Roadmap)

## Phase 1 — Production-Grade Backend Foundation (Stability First)

Before anything AI-heavy.

1. DTO Layer

   * `delivery/http/dto`
   * Request / Response models
   * Mapping DTO ↔ Domain

2. Validation System

   * Struct validation (go-playground/validator)
   * Domain invariant checks

3. Error Architecture

   * Domain errors
   * HTTP error mapper
   * Standard API response format

4. Auth Integration

   * JWT middleware
   * User context injection
   * Auth library experimentation (your goal)

5. Transaction / UnitOfWork

   * Multi-repo transaction control
   * Required for planner + analytics flows

---

## Phase 2 — Performance & Scalability Core

Now backend becomes **serious infrastructure**.

6. Redis caching

   * Exercise library cache
   * Favorite exercises cache

7. Background Worker System

   * Async planner generation
   * Async AI insights
   * Queue: Asynq / Redis

8. Observability

   * Request ID middleware
   * Structured logs
   * Health endpoint `/healthz`

---

## Phase 3 — AI Intelligence Layer (The Differentiator)

Now we make this **AI-first**, not AI-addon.

9. AI Orchestrator Service

   * Prompt templates
   * Structured JSON outputs
   * Model switching
   * Retry / timeout handling

10. AI Features

* Workout planner generator
* Progressive overload analyzer
* Smart split generator (Arnold / PPL / custom)
* Daily coaching recommendations
* Nutrition assistant (protein suggestions)

---

## Phase 4 — Product Layer (User-Facing Impact)

11. Dashboard analytics endpoints

* Strength progression
* Volume trends
* Consistency metrics

12. Motivation Engine

* AI generated daily quotes
* Cached daily message

---

## Phase 5 — Frontend Integration

13. Frontend (Vue / React)

* Dashboard
* Workout tracker
* Split builder
* AI coach panel

---

## Phase 6 — Portfolio-Killer Finish

14. OpenAPI / Swagger docs
15. Unit tests + integration tests
16. CI pipeline
17. Docker deployment
18. Cloud deployment (Fly.io / AWS / GCP)

---



## Muscle
Muscle = {upper chest, middlechest, lower chest, upper back, middle back,lower back, front delt shoulder, side delt shoulder, back delt shoulder, traps, lats, abs, quads, glute, calf, forearm,long head tricep, short head tricep, brachialis, short head bicep, long head bicep, etc}