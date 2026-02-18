-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(150) UNIQUE,
    password_hash TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- EXERCISES
CREATE TABLE exercises (
    id UUID PRIMARY KEY,
    name VARCHAR(150),
    primary_muscle VARCHAR(100),
    secondary_muscles TEXT[],
    equipment VARCHAR(100),
    created_at TIMESTAMP
);

-- FAVORITE EXERCISES
CREATE TABLE user_favorite_exercises (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, exercise_id)
);

-- SPLIT TEMPLATES
CREATE TABLE split_templates (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100),
    description TEXT,
    created_by VARCHAR(20), -- user | ai
    focus_muscle VARCHAR(100),
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP
);

-- SPLIT DAYS
CREATE TABLE split_days (
    id UUID PRIMARY KEY,
    split_template_id UUID REFERENCES split_templates(id) ON DELETE CASCADE,
    day_order INT,
    name VARCHAR(100)
);

-- SPLIT DAY EXERCISES (PLANNED)
CREATE TABLE split_day_exercises (
    id UUID PRIMARY KEY,
    split_day_id UUID REFERENCES split_days(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id),
    target_sets INT,
    target_reps INT,
    target_weight NUMERIC(6,2),
    notes TEXT
);

-- WORKOUT SESSIONS
CREATE TABLE workout_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    split_day_id UUID NULL,
    session_date DATE,
    duration_minutes INT,
    notes TEXT,
    created_at TIMESTAMP
);

-- WORKOUT EXERCISES
CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY,
    workout_session_id UUID REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id)
);

-- WORKOUT SETS
CREATE TABLE workout_sets (
    id UUID PRIMARY KEY,
    workout_exercise_id UUID REFERENCES workout_exercises(id) ON DELETE CASCADE,
    set_order INT,
    reps INT,
    weight NUMERIC(6,2),
    rpe NUMERIC(3,1),
    set_type VARCHAR(20), -- warmup | working | failure
    created_at TIMESTAMP
);

-- DAILY NUTRITION
CREATE TABLE daily_nutritions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    date DATE,
    protein_grams INT,
    calories INT,
    notes TEXT
);

-- AI RECOMMENDATIONS
CREATE TABLE planner_recommendations (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    workout_session_id UUID NULL,
    recommendation TEXT,
    recommendation_type VARCHAR(50),
    created_at TIMESTAMP
);
CREATE TABLE exercise_media (
    id UUID PRIMARY KEY,
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,

    media_type VARCHAR(20), -- image | video
    media_url TEXT,

    thumbnail_url TEXT NULL,
    created_at TIMESTAMP
);
