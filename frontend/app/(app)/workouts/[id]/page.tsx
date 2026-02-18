"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { ArrowLeft, Dumbbell } from "lucide-react";
import { toast } from "sonner";
import type {
  ExerciseResponseDTO,
  WorkoutSessionResponseDTO,
} from "@/lib/backend-dto";

export default function WorkoutDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [workout, setWorkout] = useState<WorkoutSessionResponseDTO | null>(
    null,
  );
  const [exerciseNameById, setExerciseNameById] = useState<
    Record<string, string>
  >({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadWorkout();
  }, [params.id]);

  const loadWorkout = async () => {
    try {
      const [workoutData, exercises] = await Promise.all([
        api.get<WorkoutSessionResponseDTO>(`/api/v1/workouts/${params.id}`),
        api.get<ExerciseResponseDTO[]>("/api/v1/exercises"),
      ]);

      setWorkout(workoutData);
      const map = (exercises || []).reduce(
        (acc, ex) => {
          acc[ex.id] = ex.name;
          return acc;
        },
        {} as Record<string, string>,
      );
      setExerciseNameById(map);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load workout";
      toast.error(message);
      router.push("/app/workouts");
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background p-6 md:p-8">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (!workout) {
    return (
      <div className="min-h-screen bg-background p-6 md:p-8">
        <p className="text-muted-foreground">Workout not found</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => router.back()}>
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-foreground">
              {new Date(workout.session_date).toLocaleDateString()}
            </h1>
            <p className="text-muted-foreground">
              {workout.duration_minutes} minutes • {workout.exercises.length}{" "}
              exercises
            </p>
          </div>
        </div>

        {workout.notes && workout.notes.trim().length > 0 && (
          <Card>
            <CardContent className="pt-6">
              <p className="text-muted-foreground">{workout.notes}</p>
            </CardContent>
          </Card>
        )}

        <div>
          <h2 className="text-2xl font-bold text-foreground mb-4">Exercises</h2>
          <div className="space-y-4">
            {workout.exercises.length === 0 ? (
              <Card>
                <CardContent className="pt-6 text-center py-12">
                  <Dumbbell className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
                  <p className="text-muted-foreground">No exercises recorded</p>
                </CardContent>
              </Card>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/40">
                      <th className="text-left py-3 px-4 font-semibold text-foreground">
                        Exercise
                      </th>
                      <th className="text-center py-3 px-4 font-semibold text-foreground">
                        Set #
                      </th>
                      <th className="text-center py-3 px-4 font-semibold text-foreground">
                        Reps
                      </th>
                      <th className="text-center py-3 px-4 font-semibold text-foreground">
                        Weight
                      </th>
                      <th className="text-center py-3 px-4 font-semibold text-foreground">
                        RPE
                      </th>
                      <th className="text-center py-3 px-4 font-semibold text-foreground">
                        Type
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {workout.exercises.flatMap((exercise) =>
                      (exercise.sets || []).map((set) => (
                        <tr key={set.id} className="border-b border-border/40">
                          <td className="py-3 px-4 text-foreground">
                            {exerciseNameById[exercise.exercise_id] ||
                              exercise.exercise_id}
                          </td>
                          <td className="text-center py-3 px-4 text-muted-foreground">
                            {set.set_order}
                          </td>
                          <td className="text-center py-3 px-4 text-muted-foreground">
                            {set.reps}
                          </td>
                          <td className="text-center py-3 px-4 text-muted-foreground">
                            {set.weight}
                          </td>
                          <td className="text-center py-3 px-4 text-muted-foreground">
                            {set.rpe}
                          </td>
                          <td className="text-center py-3 px-4 text-muted-foreground">
                            {set.set_type}
                          </td>
                        </tr>
                      )),
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
