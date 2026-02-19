// TypeScript mirrors of backend/internal/delivery/http/dto (snake_case)

export interface ExerciseMediaResponseDTO {
  id: string;
  media_type: string;
  media_url: string;
  thumbnail_url?: string | null;
  created_at: string;
}

export interface ExerciseResponseDTO {
  id: string;
  name: string;
  primary_muscle: string;
  secondary_muscles: string[];
  equipment: string;
  created_at: string;
  media: ExerciseMediaResponseDTO[];
}

export interface WorkoutSetResponseDTO {
  id: string;
  set_order: number;
  reps: number;
  weight: number;
  rpe: number;
  set_type: string;
  created_at: string;
}

export interface WorkoutExerciseResponseDTO {
  id: string;
  exercise_id: string;
  sets: WorkoutSetResponseDTO[];
}

export interface WorkoutSessionResponseDTO {
  id: string;
  user_id: string;
  split_day_id?: string | null;
  session_date: string;
  duration_minutes: number;
  notes: string;
  exercises: WorkoutExerciseResponseDTO[];
  created_at: string;
}

export interface SplitExerciseResponseDTO {
  exercise_id: string;
  exercise_name: string;
  target_sets: number;
  target_reps: number;
  target_weight: number;
  notes: string;
}

export interface SplitDayResponseDTO {
  id: string;
  day_order: number;
  name: string;
  exercises: SplitExerciseResponseDTO[];
}

export interface SplitTemplateResponseDTO {
  id: string;
  user_id: string;
  name: string;
  description: string;
  created_by: string;
  focus_muscle: string;
  is_active: boolean;
  days: SplitDayResponseDTO[];
  created_at: string;
}

export interface DailyNutritionResponseDTO {
  id: string;
  user_id: string;
  date: string;
  protein_grams: number;
  calories: number;
  notes: string;
}

export interface PlannerRecommendationResponseDTO {
  id: string;
  user_id: string;
  workout_session_id?: string | null;
  recommendation: string;
  recommendation_type: string;
  created_at: string;
}

export interface CoachingSuggestionsResponseDTO {
  date: string;
  suggestions: string[];
}

export interface WorkoutPlanExerciseResponseDTO {
  name: string;
  sets: number;
  rep_range: string;
  weight: number;
}

export interface WorkoutPlanResponseDTO {
  user_id: string;
  split_day_id: string;
  date: string;
  exercises: WorkoutPlanExerciseResponseDTO[];
}

export interface WorkoutExplanationExerciseNoteDTO {
  name: string;
  note: string;
}

export interface WorkoutExplanationResponseDTO {
  summary: string;
  exercise_notes: WorkoutExplanationExerciseNoteDTO[];
}

export interface DailyMotivationResponseDTO {
  date: string;
  message: string;
}

export interface HealthDTO {
  status: string;
}
