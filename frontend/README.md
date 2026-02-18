# φ Gym Dashboard

A minimalist, Greek-inspired fitness tracking dashboard built with Next.js 16, TypeScript, and Tailwind CSS. Track your workouts, nutrition, splits, and get AI-powered recommendations to optimize your training.

## Features

### Core Features
- **Authentication**: Secure login/register with JWT tokens stored in localStorage
- **Dashboard**: Hero section with daily motivation quotes and quick statistics
- **Workouts**: Log and track your training sessions with duration and notes
- **Splits**: Create and manage workout splits (Push/Pull/Legs, Upper/Lower, etc.)
- **Exercises**: Build your exercise database with primary muscle group tracking
- **Nutrition**: Track daily protein intake and visualize weekly trends
- **Analytics**: Visualize volume, weight, and set progression across muscle groups
- **AI Tools**: 
  - Daily coaching tips
  - Workout explanation
  - AI-generated workout plans
  - Progressive overload strategies
  - Split template generation
- **Planner**: Get personalized training recommendations
- **Settings**: Manage appearance, session, and system health

### Design
- **Minimalist Aesthetic**: Clean, generous whitespace with strong contrast
- **Greek/Classical Vibes**: Elegant typography and proportions
- **Stone/Marble Texture**: Subtle background texture for depth
- **Dark Mode Support**: Full dark mode with thoughtful color palette
- **Responsive Design**: Mobile-first approach with seamless desktop experience
- **Smooth Animations**: Framer Motion for page transitions and reveals

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS 4.1 + shadcn/ui components
- **State Management**: Client-side with localStorage + SWR for data sync
- **Charts**: Recharts for data visualization
- **Animations**: Framer Motion
- **Forms**: React Hook Form + Zod validation
- **Theme**: next-themes for light/dark mode
- **Notifications**: sonner for toast messages
- **Icons**: Lucide React

## Project Structure

```
app/
├── layout.tsx                 # Root layout with theme provider
├── page.tsx                   # Root redirect to app or login
├── globals.css                # Global styles and design tokens
├── (auth)/                    # Auth route group
│   ├── layout.tsx             # Auth layout
│   ├── login/page.tsx         # Login page
│   └── register/page.tsx      # Register page
└── (app)/                     # Protected app route group
    ├── layout.tsx             # Main app layout with navigation
    ├── page.tsx               # Dashboard with hero section
    ├── workouts/
    │   ├── page.tsx           # Workouts list
    │   └── [id]/page.tsx      # Workout detail
    ├── splits/
    │   ├── page.tsx           # Splits list
    │   └── [id]/page.tsx      # Split detail
    ├── exercises/
    │   ├── page.tsx           # Exercises list
    │   └── [id]/page.tsx      # Exercise detail
    ├── nutrition/page.tsx      # Nutrition tracking
    ├── analytics/page.tsx      # Training analytics
    ├── ai-tools/page.tsx       # AI features
    ├── planner/page.tsx        # Recommendations
    └── settings/page.tsx       # User settings

components/
├── auth-form.tsx              # Login/register form
├── app-header.tsx             # Top app header with theme toggle
├── app-nav.tsx                # Navigation component
├── dashboard/
│   ├── hero-section.tsx       # Hero with motivation quotes
│   └── quick-tiles.tsx        # Quick stats cards
├── workouts/
│   └── workout-form.tsx       # Create workout form
└── ui/                        # shadcn components

lib/
├── api.ts                     # Centralized API client with Bearer token handling
├── auth.ts                    # Auth state management helpers
└── utils.ts                   # Utility functions (cn, etc.)
```

## Getting Started

### Prerequisites
- Node.js 18+ and pnpm
- A backend API server running (see API Specification below)
- Optional: Vercel account for deployment

### Installation

1. Clone or download the project
2. Install dependencies:
   ```bash
   pnpm install
   ```

3. Set up environment variables in `.env.local`:
   ```
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
   ```

4. Start the development server:
   ```bash
   pnpm dev
   ```

5. Open [http://localhost:3000](http://localhost:3000) in your browser

### First Time Setup

1. **Create an account**: Go to `/register` and create a new account
2. **Log in**: Use your credentials at `/login`
3. **Explore the dashboard**: Check out the dashboard, create a split, and log your first workout
4. **Add nutrition data**: Track your daily protein intake
5. **Check analytics**: View your training progress

## API Specification

The dashboard expects a backend API at `NEXT_PUBLIC_API_BASE_URL` with the following endpoints:

### Authentication
- `POST /api/v1/auth/login` - Login with email/password
- `POST /api/v1/auth/register` - Register new account

### Workouts
- `GET /api/v1/workouts/user/:userId` - List user's workouts
- `GET /api/v1/workouts/:id` - Get workout details
- `POST /api/v1/workouts` - Create workout

### Splits
- `GET /api/v1/splits/user/:userId` - List user's splits
- `GET /api/v1/splits/:id` - Get split details
- `POST /api/v1/splits` - Create split
- `PUT /api/v1/splits/:id` - Update split
- `POST /api/v1/splits/:id/activate` - Activate a split

### Exercises
- `GET /api/v1/exercises` - List all exercises
- `GET /api/v1/exercises/:id` - Get exercise details
- `POST /api/v1/exercises` - Create exercise
- `POST /api/v1/exercises/:id/media` - Attach media URL

### Nutrition
- `GET /api/v1/nutrition/user/:userId?date=YYYY-MM-DD` - Get nutrition for date
- `POST /api/v1/nutrition` - Create/update nutrition

### AI Tools
- `GET /api/v1/ai/coaching` - Get daily coaching tips
- `POST /api/v1/ai/explain-workout` - Explain a workout
- `POST /api/v1/ai/workout` - Generate workout plan
- `POST /api/v1/ai/overload` - Get progressive overload strategy
- `POST /api/v1/ai/generate-split` - Generate split template

### Planner
- `GET /api/v1/planner/user/:userId` - List recommendations
- `POST /api/v1/planner/generate/:userId` - Generate recommendation

### System
- `GET /api/v1/health` - Health check

### Response Format

All API responses should follow this envelope format:
```json
{
  "status": "success" | "error",
  "data": { /* response data */ },
  "message": "optional error message"
}
```

### Authentication

Include the JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

The token is automatically injected by the API client (`lib/api.ts`).

## Configuration

### Colors & Theme

The app uses a stone/marble color palette designed for Greek/classical aesthetic:

**Light Mode**:
- Background: #faf9f7 (off-white)
- Foreground: #1a1a1a (near-black)
- Primary: #1a1a1a (stone black)
- Muted: #e8e3db (light stone)

**Dark Mode**:
- Background: #0f0f0f (deep black)
- Foreground: #faf9f7 (off-white)
- Primary: #e8e3db (light stone)
- Muted: #3a3a3a (dark gray)

Modify colors in `app/globals.css` CSS variables.

### Fonts

The app uses Geist for sans-serif and Geist Mono for monospace. Change fonts in `app/layout.tsx`.

### API Base URL

Update `NEXT_PUBLIC_API_BASE_URL` in `.env.local` to point to your backend.

## Development

### File Structure Guidelines
- Components in `/components` - reusable UI components
- Pages in `/app` - route handlers
- Utilities in `/lib` - helper functions and API client
- Styles in `app/globals.css` - global styles and design tokens

### Adding New Features

1. **Create the component**: Add to `/components` (or `/components/feature-name`)
2. **Create the page**: Add to `/app/(app)/feature-name/page.tsx`
3. **Connect the API**: Use `api.get()`, `api.post()`, etc. from `lib/api.ts`
4. **Handle loading/errors**: Use `toast.error()` from sonner for notifications
5. **Style with Tailwind**: Use design tokens (`bg-primary`, `text-muted-foreground`, etc.)
6. **Add navigation**: Update navigation items in `components/app-nav.tsx`

### Common Patterns

**Data Fetching**:
```tsx
const [data, setData] = useState(null)
const [isLoading, setIsLoading] = useState(true)

useEffect(() => {
  api.get('/api/v1/endpoint')
    .then(setData)
    .catch(err => toast.error(err.message))
    .finally(() => setIsLoading(false))
}, [])
```

**Forms with Validation**:
```tsx
const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault()
  if (!email) {
    toast.error('Email is required')
    return
  }
  
  try {
    await api.post('/api/v1/endpoint', { email })
    toast.success('Success!')
  } catch (err) {
    toast.error(err.message)
  }
}
```

**Protected Routes**:
Routes in `/app/(app)` are automatically protected by middleware that checks for auth tokens.

## Deployment

### Deploy to Vercel

1. Push your code to GitHub
2. Import the repository on [Vercel](https://vercel.com)
3. Set environment variables:
   - `NEXT_PUBLIC_API_BASE_URL` - Your backend API URL
4. Deploy

The app will automatically optimize with Turbopack (Next.js 16 default) and deploy.

### Deploy Elsewhere

The app is a standard Next.js 16 app and can be deployed anywhere:

```bash
pnpm build
pnpm start
```

## Troubleshooting

### API Connection Issues
- Check `NEXT_PUBLIC_API_BASE_URL` in `.env.local`
- Ensure backend is running and accessible
- Check network tab in browser DevTools for failed requests
- Verify API response format matches the envelope structure

### Auth Issues
- Clear browser localStorage if token becomes invalid
- Check that backend returns `{ status: 'success', data: { token, userId } }`
- Verify Bearer token is being sent in Authorization header

### Styling Issues
- Ensure Tailwind CSS is properly configured
- Check that design tokens are set in `app/globals.css`
- Clear `.next` cache: `rm -rf .next && pnpm dev`

### Performance
- Charts may be slow with large datasets - consider pagination
- Use React DevTools Profiler to identify slow components
- Ensure animations aren't disabled with `prefers-reduced-motion`

## Customization

### Add New API Endpoints

1. Create a new function in `lib/api.ts` (or just use `api.get()`, `api.post()`, etc.)
2. Call from your component
3. Handle loading and errors with try/catch + toast notifications

### Change Color Scheme

Edit CSS variables in `app/globals.css`:
```css
:root {
  --primary: #yourcolor;
  --foreground: #yourcolor;
  /* ... other colors */
}

.dark {
  --primary: #yourcolor;
  /* ... */
}
```

### Add New Pages

1. Create folder structure: `/app/(app)/feature-name/page.tsx`
2. Add navigation item in `components/app-nav.tsx`
3. The page is automatically protected by middleware

## License

This project is provided as-is for personal use.

## Support

For issues, questions, or feature requests, please refer to the Vercel documentation or create an issue in your repository.

---

Built with [v0](https://v0.app) • Powered by [Next.js](https://nextjs.org) • Styled with [Tailwind CSS](https://tailwindcss.com)
