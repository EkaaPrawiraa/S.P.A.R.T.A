Build a minimalist, responsive gym dashboard web app with “Greek / classical gym” vibes.

Tech:

- Next.js (App Router) + TypeScript
- Tailwind CSS
- shadcn/ui components
- Framer Motion for smooth, subtle animations (page transitions, section reveals, hover)
- Recharts for charts
- Mobile-first responsive layout, also great on laptop

Theme (required):

- Implement a real light/dark mode toggle using `next-themes` + shadcn/ui theming.
- Default to system theme.
- When **dark mode** is enabled, add a subtle “glow” accent on primary CTAs and key surfaces (hero CTA, selected nav item, primary cards, chart container).
  - Use only theme tokens (`primary`, `ring`, `border`, `background`, `foreground`) and Tailwind utilities (e.g. `ring-1 ring-primary/30`, `shadow-lg`, `shadow-black/20`, `backdrop-blur`) — do not introduce custom hard-coded colors.
  - Glow should be tasteful: readable, not neon.

Visual design:

- Minimalist, “Greek gym” vibe: clean typography, generous whitespace, subtle stone/marble-inspired background texture (very light), strong contrast, modern.
- No clutter. Keep UI understandable for non-technical users.
- Smooth animations, but not flashy.

Core UX requirement:

- The first screen is a Dashboard “Hero” that shows:
  - Today’s motivation/quote (from backend)
  - Date
  - A single primary CTA like “Start” or “Scroll to Dashboard”
- After that, the user can scroll/slide up to reveal the rest of the app. Implement as a vertical two-section layout:
  - Section 1: Motivation Hero (full height)
  - Section 2: Main App Content
    Use CSS scroll-snap OR a clearly polished scroll transition. Add a subtle “scroll indicator”.

Backend integration (REST):

- Base URL via env var: NEXT_PUBLIC_API_BASE_URL
- All API responses use an envelope:
  - success: { "status": "success", "data": ... }
  - error: { "status": "error", "message": "..." }
- Protected endpoints require Authorization: Bearer <JWT>
- JWT includes claim `user_id` (string UUID). Many routes also require user_id in path/body and must match token.

Backend auth (required — implement this FIRST):

- Add an auth flow to the Go (Gin) backend so the frontend does NOT rely on “paste JWT” as the primary path.
- Use free, common libraries:
  - JWT: `github.com/golang-jwt/jwt/v5`
  - Password hashing: `golang.org/x/crypto/bcrypt`
- Endpoints (enveloped responses, consistent errors):
  - `POST /api/v1/auth/register`
    - Body: `{ "email": string, "password": string, "name": string }`
    - Behavior: create user (unique email), store bcrypt hash, return `{ user_id, token }`.
  - `POST /api/v1/auth/login`
    - Body: `{ "email": string, "password": string }`
    - Behavior: verify password, return `{ user_id, token }`.
- Token:
  - HMAC SHA256 signing method, secret from env `JWT_SECRET`.
  - Include claims: `user_id`, `exp`, `iat`.
- Persistence:
  - If the DB already has a users table, use it; otherwise add a minimal migration that supports login (id UUID, email unique, name, password_hash, created_at).
- Update/confirm the existing JWT middleware accepts tokens from the new auth endpoints.

Endpoints to support (must implement screens + data fetching for these):
Public:

- GET /api/v1/health
- GET /api/v1/metrics

Auth handling (updated — backend DOES provide auth now):

- Provide dedicated pages:
  - `/login` (email + password)
  - `/register` (name + email + password)
- Store `{ token, userId }` in localStorage after successful auth.
- Keep the “Session” panel in Settings as an advanced fallback for developers (manual paste/override), but not the primary flow.
- Create a fetch client that attaches the Bearer token automatically.
- Route protection:
  - If not authenticated, redirect to `/login`.
  - Add a logout action that clears localStorage.

AI + Motivation:

- GET /api/v1/ai/motivation (show on Hero; also show in dashboard feed)
- GET /api/v1/ai/coaching (show as simple bullet cards)
- POST /api/v1/ai/explain-workout (a small form: split_day_name, fatigue, exercises[]; render returned summary + per-exercise notes)
- POST /api/v1/ai/workout (form: split_day_id, fatigue; render plan)
- POST /api/v1/ai/overload (form: exercise_id; render action/message)
- POST /api/v1/ai/generate-split (form: days_per_week, focus_muscle; render generated template)

Exercises:

- POST /api/v1/exercises (create exercise)
- POST /api/v1/exercises/:id/media (attach media url)
- GET /api/v1/exercises (list)
- GET /api/v1/exercises/:id (detail)
  Design: exercise list is clean cards; detail page shows media list (image/video URLs) simply.

Splits:

- POST /api/v1/splits (create template)
- GET /api/v1/splits/:id (detail)
- PUT /api/v1/splits/:id (update)
- POST /api/v1/splits/:id/activate (activate)
- GET /api/v1/splits/user/:user_id (list)
  Design: split templates list + detail with days/exercises; keep editing UI simple (basic forms, no complex builder).

Workouts:

- POST /api/v1/workouts (create session)
- GET /api/v1/workouts/:id (detail)
- GET /api/v1/workouts/user/:user_id (list)
  Design: workout list shows date, duration, quick summary; detail shows exercises and sets.

Nutrition:

- POST /api/v1/nutrition (upsert daily nutrition)
- GET /api/v1/nutrition/user/:user_id?date=YYYY-MM-DD
  Design: simple daily protein progress + weekly trend chart (client-computed by fetching multiple dates is OK if no endpoint exists; otherwise show only selected day).

Planner:

- POST /api/v1/planner/generate/:user_id (generates AI coaching and saves it)
- GET /api/v1/planner/user/:user_id
  Design: recommendations list (most recent first), show as bullet content.

Analytics requirement (graphs “journey” by muscle group):

- Compute client-side analytics by combining:
  - workouts list (contains exercise_id + sets with reps/weight/rpe)
  - exercises list (map exercise_id -> primary_muscle)
- Provide an “Analytics” view:
  - Muscle group selector (from known primary_muscle values)
  - Show 3 charts:
    1. Volume trend over time (sum reps\*weight) per muscle group
    2. Average working weight trend
    3. Total sets trend
- Keep chart UI minimal and readable on mobile (tabs to switch metrics on small screens).

Navigation:

- After the hero section, show main app navigation:
  - Desktop: left sidebar
  - Mobile: bottom tab bar
    Tabs: Dashboard, Workouts, Splits, Exercises, Nutrition, Analytics, AI Tools, Planner, Settings

Dashboard (Section 2 top area):

- Show quick tiles: “This week workouts count”, “Today protein”, “Latest recommendation”
- Show coaching suggestions preview
- Keep it simple and not overwhelming.

Quality:

- Good loading/error states, no blank screens.
- Use toasts for success/error.
- Use React Server Components where appropriate, but API calls requiring JWT should be client-side with a clean API wrapper.
- Include example mock data fallback only if JWT is missing, but clearly label “Not connected”.

Non-negotiables:

- Must use Next.js + shadcn/ui (do not swap UI kits).
- Must include light/dark toggle; dark mode has subtle token-based glow.
- Must include real backend auth endpoints (`/auth/register`, `/auth/login`) and wire frontend to them.

Page-by-page requirements (DETAILED — generate ALL pages below)

Routing/layout architecture (App Router):

- Use route groups:
  - `app/(auth)/login/page.tsx`
  - `app/(auth)/register/page.tsx`
  - `app/(app)/layout.tsx` (protected shell: sidebar on desktop, bottom tabs on mobile)
  - `app/(app)/page.tsx` (Dashboard with Hero + Main Content sections)
  - `app/(app)/workouts/page.tsx`
  - `app/(app)/workouts/[id]/page.tsx`
  - `app/(app)/splits/page.tsx`
  - `app/(app)/splits/[id]/page.tsx`
  - `app/(app)/exercises/page.tsx`
  - `app/(app)/exercises/[id]/page.tsx`
  - `app/(app)/nutrition/page.tsx`
  - `app/(app)/analytics/page.tsx`
  - `app/(app)/ai-tools/page.tsx`
  - `app/(app)/planner/page.tsx`
  - `app/(app)/settings/page.tsx`
- Protected routing rule:
  - Any route in `(app)` requires auth.
  - If `{token,userId}` missing or invalid → redirect to `/login`.
- Motion:
  - Use Framer Motion for subtle route transitions and section reveals.
  - Keep transitions quick and minimal (opacity/translate only).

Global UI shell (for all `(app)` routes):

- Desktop:
  - Left sidebar navigation (shadcn `Button`/`NavigationMenu` style; keep minimal).
  - Main content area with max width and generous padding.
- Mobile:
  - Bottom tab bar with the same tabs.
- Tabs (exactly these; no extra): Dashboard, Workouts, Splits, Exercises, Nutrition, Analytics, AI Tools, Planner, Settings
- Add a small top-right user menu (no modal) with:
  - Theme toggle
  - Link to Settings
  - Logout

Page: `/login`

- Purpose: authenticate the user.
- UI:
  - Minimal centered card with: Email, Password, “Login” button.
  - Secondary link to `/register`.
  - Respect Greek/minimal design; dark-mode glow on primary button.
- Behavior:
  - Call `POST /api/v1/auth/login`.
  - On success: store `{token,userId}` in localStorage; toast success; redirect to `/`.
  - On error: show toast + inline message.
- States:
  - Disable button + show loading spinner while submitting.

Page: `/register`

- Purpose: create a new account.
- UI:
  - Centered card with: Name, Email, Password, “Create account”.
  - Secondary link to `/login`.
- Behavior:
  - Call `POST /api/v1/auth/register`.
  - On success: store `{token,userId}`; redirect to `/`.
- Validation:
  - Basic client validation: required fields; email format; password min length.

Page: `/` (Dashboard: Hero + Main Content)

- Must be a 2-section vertical layout:
  - Section 1: Motivation Hero (100vh)
  - Section 2: Main App Content
- Implement scroll-snap OR a polished scroll transition.
- Hero section:
  - Fetch `GET /api/v1/ai/motivation` on mount (client-side due to auth).
  - Show: quote text, optional author/source if provided, today’s date.
  - Primary CTA: “Scroll to Dashboard” that smoothly scrolls to Section 2.
  - Add subtle scroll indicator.
- Main content section (top area):
  - Quick tiles (3 cards):
    1. “This week workouts”
       - Fetch workouts list `GET /api/v1/workouts/user/:user_id` and count workouts whose date is within the last 7 days.
    2. “Today protein”
       - Fetch `GET /api/v1/nutrition/user/:user_id?date=YYYY-MM-DD` for today.
       - Show protein grams and a simple progress bar (if goal exists; otherwise just show grams).
    3. “Latest recommendation”
       - Fetch planner list `GET /api/v1/planner/user/:user_id`, show the newest item preview.
  - Coaching preview:
    - Fetch `GET /api/v1/ai/coaching` and render as a short list of bullet cards.
  - Motivation feed:
    - Repeat the motivation quote in a smaller “Today” card (no extra feed endpoints).
- System status (required to cover public endpoints, but keep minimal):
  - Show a small “System” card that calls:
    - `GET /api/v1/health` (show OK / error)
    - `GET /api/v1/metrics` (render a tiny JSON preview or key counters if available)
- States:
  - If auth missing: show “Not connected” placeholders (no hard crash) + link to Settings Session panel.
  - Always show skeletons while loading.

Page: `/workouts`

- Purpose: list and create workout sessions.
- Data:
  - Fetch `GET /api/v1/workouts/user/:user_id`.
- UI:
  - Header with title + a compact inline “New workout” form (no modal):
    - Minimum fields needed by backend (match your API DTO); keep it simple.
  - List workouts as clean cards:
    - Show date, duration (if available), and a 1-line summary (e.g. number of exercises / sets).
    - Clicking a card routes to `/workouts/[id]`.
- Actions:
  - Create via `POST /api/v1/workouts` then refresh list and toast success.

Page: `/workouts/[id]`

- Purpose: workout detail view.
- Data:
  - Fetch `GET /api/v1/workouts/:id`.
- UI:
  - Summary header: date, duration, notes if present.
  - Exercises table/list:
    - For each exercise: show name (if present), and its sets with reps/weight/RPE.
  - Provide a small “Explain this workout” CTA that links to AI Tools → Explain Workout tab and pre-fills data if possible (optional; if not, just provide a link).

Page: `/splits`

- Purpose: list split templates and create a new template.
- Data:
  - Fetch `GET /api/v1/splits/user/:user_id`.
- UI:
  - Inline “Create split template” form:
    - Minimal fields required by backend (template name, days, focus).
  - List split templates as cards with:
    - Name, days per week, focus muscle.
    - Status badge if active.
- Actions:
  - Create via `POST /api/v1/splits`.
  - Navigate to detail on click.

Page: `/splits/[id]`

- Purpose: view + edit + activate split template.
- Data:
  - Fetch `GET /api/v1/splits/:id`.
- UI:
  - Show split days and exercises in a readable list.
  - Simple editing UI (basic forms only, no complex drag/drop builder):
    - Edit template name/focus/days and save.
- Actions:
  - Save updates via `PUT /api/v1/splits/:id`.
  - Activate via `POST /api/v1/splits/:id/activate`.
  - Toast on success/failure.

Page: `/exercises`

- Purpose: manage exercises.
- Data:
  - Fetch `GET /api/v1/exercises`.
- UI:
  - Inline “Create exercise” form (minimal fields required).
  - Exercise list as clean cards:
    - Name + primary muscle.
    - Click navigates to `/exercises/[id]`.
- Actions:
  - Create via `POST /api/v1/exercises` then refresh.

Page: `/exercises/[id]`

- Purpose: exercise detail and media attachment.
- Data:
  - Fetch `GET /api/v1/exercises/:id`.
- UI:
  - Show exercise metadata.
  - Media section:
    - List media URLs plainly (image/video URL text + optional preview if easy).
  - Inline “Attach media” form:
    - URL input → submit.
- Actions:
  - Call `POST /api/v1/exercises/:id/media`.

Page: `/nutrition`

- Purpose: daily nutrition view + upsert.
- Data:
  - Date picker (simple input type=date).
  - Fetch `GET /api/v1/nutrition/user/:user_id?date=YYYY-MM-DD` for selected date.
- UI:
  - Daily summary card:
    - Protein (required), plus calories/carbs/fat if supported by backend.
  - Upsert form (same fields as backend expects).
  - Weekly trend chart:
    - Compute client-side by fetching the last 7 days (7 requests is OK).
    - Render a minimal Recharts line chart for protein.
- Actions:
  - Save via `POST /api/v1/nutrition` (upsert) then refetch selected day and weekly series.

Page: `/analytics`

- Purpose: muscle-group journey charts.
- Data:
  - Fetch exercises `GET /api/v1/exercises` (to map `exercise_id -> primary_muscle`).
  - Fetch workouts list `GET /api/v1/workouts/user/:user_id`.
- Compute:
  - Filter sets by selected muscle group.
  - For each workout date compute:
    - Volume = sum(reps \* weight)
    - Avg working weight = average(weight across working sets)
    - Total sets = count(sets)
- UI:
  - Muscle group selector.
  - Charts:
    - Desktop: show all 3 charts stacked.
    - Mobile: use shadcn `Tabs` to switch between the 3 metrics.
  - Empty states:
    - If no workouts/exercises: show clear guidance.

Page: `/ai-tools`

- Purpose: single page with tabs for all AI endpoints.
- UI:
  - Use shadcn `Tabs` with exactly these tools:
    1. Coaching (GET)
    2. Explain Workout (POST)
    3. Workout Plan (POST)
    4. Overload (POST)
    5. Generate Split (POST)
- Tool: Coaching
  - Call `GET /api/v1/ai/coaching` and render bullet cards.
- Tool: Explain Workout
  - Form fields: `split_day_name` (text), `fatigue` (select/number), and `exercises[]` (simple repeatable inputs: name, sets summary).
  - Submit to `POST /api/v1/ai/explain-workout` and render returned summary + per-exercise notes.
- Tool: Workout Plan
  - Form: `split_day_id`, `fatigue` → submit `POST /api/v1/ai/workout`.
  - Render plan as clean sections (warmup/main/accessories).
- Tool: Overload
  - Form: `exercise_id` → submit `POST /api/v1/ai/overload`.
  - Render action/message in a card.
- Tool: Generate Split
  - Form: `days_per_week`, `focus_muscle` → submit `POST /api/v1/ai/generate-split`.
  - Render generated template in readable format.
- States:
  - Each tool has its own loading + error UI; no blank tab content.

Page: `/planner`

- Purpose: generate and view saved planner recommendations.
- Data:
  - Fetch `GET /api/v1/planner/user/:user_id`.
- UI:
  - “Generate recommendation” button at top.
  - List recommendations (most recent first) as cards; content is bullet text.
- Actions:
  - Generate via `POST /api/v1/planner/generate/:user_id` then refresh list.

Page: `/settings`

- Purpose: session + preferences.
- Sections (single page, no modals):
  1. Appearance
     - Light/dark/system toggle (next-themes)
  2. Session (developer fallback)
     - Inputs to paste/override token + userId
     - Validate button that calls `GET /api/v1/health` and shows auth status
  3. System
     - Show `GET /api/v1/health` result
     - Show `GET /api/v1/metrics` output (minimal rendering)
  4. Logout
     - Clears localStorage and redirects to `/login`

Deliverables:

- Full Next.js app scaffold with pages/components
- Reusable api client: `apiFetch(path, options)` that unwraps the envelope and throws on error
- localStorage-backed auth token/userId handling (from login/register)
- Responsive layout + animations + charts
- Minimalist Greek/gym styling

Now generate the code.

---
