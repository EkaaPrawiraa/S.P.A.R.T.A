"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { WorkoutForm } from "@/components/workouts/workout-form";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import { Clock, Dumbbell } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import type { WorkoutSessionResponseDTO } from "@/lib/backend-dto";

export default function WorkoutsPage() {
  const [workouts, setWorkouts] = useState<WorkoutSessionResponseDTO[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  useEffect(() => {
    loadWorkouts();
  }, []);

  const loadWorkouts = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      const data = await api.get<WorkoutSessionResponseDTO[]>(
        `/api/v1/workouts/user/${userId}`,
      );
      setWorkouts(data || []);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load workouts";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Workouts</h1>
            <p className="text-muted-foreground mt-2">
              Track and manage your training sessions
            </p>
          </div>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? "Cancel" : "New Workout"}
          </Button>
        </div>

        {showForm && (
          <Card>
            <CardHeader>
              <CardTitle>Create New Workout</CardTitle>
            </CardHeader>
            <CardContent>
              <WorkoutForm
                onSuccess={() => {
                  setShowForm(false);
                  loadWorkouts();
                }}
              />
            </CardContent>
          </Card>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading workouts...</p>
          </div>
        ) : workouts.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Dumbbell className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">
                No workouts yet. Start by creating your first workout!
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {workouts.map((workout) => (
              <Link key={workout.id} href={`/app/workouts/${workout.id}`}>
                <Card className="h-full hover:shadow-lg transition-shadow cursor-pointer">
                  <CardHeader>
                    <CardTitle className="text-lg">
                      {new Date(workout.session_date).toLocaleDateString()}
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Clock className="h-4 w-4" />
                      <span>{workout.duration_minutes} minutes</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Dumbbell className="h-4 w-4" />
                      <span>{workout.exercises.length} exercises</span>
                    </div>
                    {workout.notes && workout.notes.trim().length > 0 && (
                      <p className="text-sm text-muted-foreground italic">
                        {workout.notes}
                      </p>
                    )}
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
