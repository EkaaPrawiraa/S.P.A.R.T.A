# API Setup Guide

This guide walks you through setting up and connecting the gym dashboard frontend to your backend API.

## Quick Start

### 1. Set Your API Base URL

In `.env.local`, set your backend API URL:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Replace `http://localhost:8080` with your production API URL when deploying.

### 2. Start Your Backend

Ensure your backend API is running and accessible at the URL specified above.

### 3. Test the Connection

1. Start the frontend: `pnpm dev`
2. Go to `/register` and create an account
3. The frontend will automatically test the connection by calling `POST /api/v1/auth/register`
4. If successful, you'll be logged in and redirected to the dashboard

## API Endpoints Required

Your backend must implement these endpoints. All responses should follow the envelope format below.

### Authentication

#### `POST /api/v1/auth/register`
Register a new user.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response (success)**:
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "userId": "uuid-or-id"
  }
}
```

**Response (error)**:
```json
{
  "status": "error",
  "message": "Email already exists"
}
```

#### `POST /api/v1/auth/login`
Login an existing user.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**: Same as register

---

### Workouts

#### `GET /api/v1/workouts/user/:userId`
Get all workouts for a user.

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "date": "2024-02-15",
      "duration_minutes": 60,
      "exercise_count": 8,
      "notes": "Great session"
    }
  ]
}
```

#### `POST /api/v1/workouts`
Create a new workout.

**Request**:
```json
{
  "split_id": "split-uuid",
  "duration_minutes": 60,
  "notes": "Great session"
}
```

**Response**: Returns created workout object

#### `GET /api/v1/workouts/:id`
Get workout details with exercises.

**Response**:
```json
{
  "status": "success",
  "data": {
    "id": "uuid",
    "date": "2024-02-15",
    "duration_minutes": 60,
    "notes": "Great session",
    "exercises": [
      {
        "id": "ex-uuid",
        "name": "Barbell Bench Press",
        "sets": 4,
        "reps": 8,
        "weight": 225,
        "rpe": 8
      }
    ]
  }
}
```

---

### Splits

#### `GET /api/v1/splits/user/:userId`
Get all workout splits for a user.

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "name": "Push/Pull/Legs",
      "days": 4,
      "focus_muscle": "Legs",
      "is_active": true
    }
  ]
}
```

#### `POST /api/v1/splits`
Create a new split.

**Request**:
```json
{
  "name": "Push/Pull/Legs",
  "days": 4
}
```

#### `GET /api/v1/splits/:id`
Get split details.

#### `PUT /api/v1/splits/:id`
Update a split.

#### `POST /api/v1/splits/:id/activate`
Activate a split (make it the current/active split).

---

### Exercises

#### `GET /api/v1/exercises`
Get all exercises in the database.

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "name": "Barbell Bench Press",
      "primary_muscle": "Chest"
    }
  ]
}
```

#### `POST /api/v1/exercises`
Create a new exercise.

**Request**:
```json
{
  "name": "Barbell Bench Press",
  "primary_muscle": "Chest"
}
```

#### `GET /api/v1/exercises/:id`
Get exercise details with media.

**Response**:
```json
{
  "status": "success",
  "data": {
    "id": "uuid",
    "name": "Barbell Bench Press",
    "primary_muscle": "Chest",
    "media_urls": ["https://example.com/video.mp4"]
  }
}
```

#### `POST /api/v1/exercises/:id/media`
Attach a media URL to an exercise.

**Request**:
```json
{
  "media_url": "https://example.com/video.mp4"
}
```

---

### Nutrition

#### `GET /api/v1/nutrition/user/:userId?date=YYYY-MM-DD`
Get nutrition data for a specific date.

**Response**:
```json
{
  "status": "success",
  "data": {
    "date": "2024-02-15",
    "protein_grams": 150,
    "calories": 2500,
    "carbs": 250,
    "fats": 80
  }
}
```

Return null/empty data if no entry exists for the date (don't error).

#### `POST /api/v1/nutrition`
Create or update nutrition data for a date.

**Request**:
```json
{
  "date": "2024-02-15",
  "protein_grams": 150,
  "calories": 2500
}
```

---

### Analytics

These endpoints are computed client-side by fetching workouts and exercises, but you may want to provide pre-computed endpoints for performance.

The frontend expects:
- `GET /api/v1/workouts/user/:userId` to return workouts with exercises
- `GET /api/v1/exercises` to map exercise IDs to muscle groups

---

### AI Tools

#### `GET /api/v1/ai/coaching`
Get daily coaching tips.

**Response**:
```json
{
  "status": "success",
  "data": {
    "tips": [
      "Focus on progressive overload to build strength",
      "Ensure adequate sleep for recovery",
      "Track your nutrition to optimize performance"
    ]
  }
}
```

#### `POST /api/v1/ai/explain-workout`
Explain the benefits of a workout.

**Request**:
```json
{
  "workout_text": "5x5 barbell squats, 4x8 leg press, 3x10 leg curl"
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "explanation": "This workout targets..."
  }
}
```

#### `POST /api/v1/ai/workout`
Generate a personalized workout plan.

**Response**:
```json
{
  "status": "success",
  "data": {
    "plan": "Monday (Push): 4x5 Bench Press, 4x5 Incline Press, 3x8 Dips..."
  }
}
```

#### `POST /api/v1/ai/overload`
Get progressive overload strategy.

**Response**:
```json
{
  "status": "success",
  "data": {
    "strategy": "Next week, increase weight by 5 lbs on all main compounds..."
  }
}
```

#### `POST /api/v1/ai/generate-split`
Generate a workout split template.

**Response**:
```json
{
  "status": "success",
  "data": {
    "template": "Upper Body: Bench Press, Rows, Shoulder Press..."
  }
}
```

---

### Planner

#### `GET /api/v1/planner/user/:userId`
Get all personalized recommendations.

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "content": "You should deload next week...",
      "created_at": "2024-02-15T10:30:00Z"
    }
  ]
}
```

#### `POST /api/v1/planner/generate/:userId`
Generate a new personalized recommendation.

**Response**: Returns the generated recommendation

---

### Health Check

#### `GET /api/v1/health`
System health check endpoint.

**Response**:
```json
{
  "status": "success",
  "data": {
    "status": "ok",
    "message": "All systems operational"
  }
}
```

---

## Response Format

**All** API responses must follow this envelope format:

```json
{
  "status": "success" | "error",
  "data": { /* response payload */ },
  "message": "optional error message"
}
```

### Rules
1. `status` is always `"success"` or `"error"`
2. `data` contains the actual response payload (only for success)
3. `message` contains error details (only for errors)
4. HTTP status code should be:
   - `200` for successful requests
   - `400` for bad requests
   - `401` for unauthorized
   - `404` for not found
   - `500` for server errors

### Examples

**Success**:
```json
{
  "status": "success",
  "data": { "id": "123", "name": "Workout A" }
}
```

**Error**:
```json
{
  "status": "error",
  "message": "Invalid email or password"
}
```

---

## Authentication

All protected endpoints (except `/api/v1/auth/*`) require the Bearer token in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

The frontend automatically injects this token via `lib/api.ts`. Your backend should:

1. Extract the token from the header
2. Validate it (decode JWT, check expiration, verify signature)
3. Extract the `userId` from the token
4. Use `userId` to scope responses (only return user's own data)

---

## CORS

Ensure your backend allows requests from your frontend domain:

```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

---

## Pagination (Optional)

For production with large datasets, consider adding pagination:

```
GET /api/v1/workouts/user/:userId?page=1&limit=20
```

The frontend can be easily updated to support pagination in `lib/api.ts`.

---

## Error Handling

The frontend expects all errors in the response body with `status: "error"`. Examples:

- Missing required fields: `"Email is required"`
- Invalid email format: `"Invalid email format"`
- Authentication failed: `"Invalid email or password"`
- Not found: `"Workout not found"`
- Server error: `"Internal server error"`

All error messages are displayed to users via toast notifications, so keep them user-friendly and concise.

---

## Testing the API

Use the Settings page (`/app/settings`) to test your API:

1. Check system health: Click "Check Health" to call `GET /api/v1/health`
2. Validate token: Paste your token and click "Validate Token"
3. View health response: The response is displayed on the page

---

## Development vs. Production

### Development
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### Production
```env
NEXT_PUBLIC_API_BASE_URL=https://api.example.com
```

Set the environment variable during deployment (Vercel Settings → Environment Variables).

---

## Troubleshooting

### CORS Errors
- Backend not accepting requests from frontend domain
- Solution: Add proper CORS headers to backend responses

### 401 Unauthorized
- Token is missing, invalid, or expired
- Solution: Check that token is in localStorage and valid
- Check Settings page to validate token

### 400 Bad Request
- Request body format doesn't match API expectations
- Solution: Check the request format in this guide

### 500 Server Error
- Backend encountered an error
- Solution: Check backend logs and ensure database is accessible

### API Base URL Issues
- Frontend can't reach backend
- Solution: Verify `NEXT_PUBLIC_API_BASE_URL` is correct
- Ensure backend is running and network accessible
- Check browser DevTools Network tab for failed requests

---

For more details on the frontend implementation, see [README.md](./README.md).
