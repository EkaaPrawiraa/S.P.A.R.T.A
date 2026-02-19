"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { WorkoutForm } from "@/components/workouts/workout-form";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import {
  Clock,
  Dumbbell,
  ChevronLeft,
  ChevronRight,
  FileText,
} from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import type {
  ExerciseResponseDTO,
  WorkoutSessionResponseDTO,
} from "@/lib/backend-dto";
import { cn } from "@/lib/utils";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  addDays,
  addMonths,
  endOfMonth,
  endOfWeek,
  format,
  isSameMonth,
  isSameDay,
  startOfMonth,
  startOfWeek,
  subMonths,
} from "date-fns";

export default function WorkoutsPage() {
  const [workouts, setWorkouts] = useState<WorkoutSessionResponseDTO[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [monthCursor, setMonthCursor] = useState(() =>
    startOfMonth(new Date()),
  );
  const [selectedDate, setSelectedDate] = useState(() => new Date());
  const [isWorkoutModalOpen, setIsWorkoutModalOpen] = useState(false);
  const [exerciseNameById, setExerciseNameById] = useState<
    Record<string, string>
  >({});

  useEffect(() => {
    loadWorkouts();
  }, []);

  useEffect(() => {
    const loadExercises = async () => {
      try {
        const items = await api.get<ExerciseResponseDTO[]>("/api/v1/exercises");
        const map = (items || []).reduce(
          (acc, ex) => {
            acc[ex.id] = ex.name;
            return acc;
          },
          {} as Record<string, string>,
        );
        setExerciseNameById(map);
      } catch {
        setExerciseNameById({});
      }
    };

    loadExercises();
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

  const workoutsByDate = (workouts || []).reduce(
    (acc, w) => {
      const key = (w.session_date || "").slice(0, 10);
      if (!key) return acc;
      if (!acc[key]) acc[key] = [];
      acc[key].push(w);
      return acc;
    },
    {} as Record<string, WorkoutSessionResponseDTO[]>,
  );

  const selectedISO = format(selectedDate, "yyyy-MM-dd");
  const selectedWorkouts = workoutsByDate[selectedISO] || [];

  const monthStart = startOfMonth(monthCursor);
  const monthEnd = endOfMonth(monthCursor);
  const gridStart = startOfWeek(monthStart, { weekStartsOn: 0 });
  const gridEnd = endOfWeek(monthEnd, { weekStartsOn: 0 });

  const days: Date[] = [];
  for (let d = gridStart; d <= gridEnd; d = addDays(d, 1)) {
    days.push(d);
  }

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
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading workouts...</p>
          </div>
        ) : (
          <>
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle className="text-lg">
                    {format(monthCursor, "MMMM yyyy")}
                  </CardTitle>
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      onClick={() => setMonthCursor((d) => subMonths(d, 1))}
                      aria-label="Previous month"
                    >
                      <ChevronLeft className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      onClick={() => setMonthCursor((d) => addMonths(d, 1))}
                      aria-label="Next month"
                    >
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-7 gap-2 text-xs text-muted-foreground mb-2">
                  {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map(
                    (d) => (
                      <div key={d} className="text-center">
                        {d}
                      </div>
                    ),
                  )}
                </div>

                <div className="grid grid-cols-7 gap-2">
                  {days.map((d) => {
                    const iso = format(d, "yyyy-MM-dd");
                    const dayWorkouts = workoutsByDate[iso] || [];
                    const inMonth = isSameMonth(d, monthCursor);
                    const isSelected = isSameDay(d, selectedDate);
                    const hasWorkout = dayWorkouts.length > 0;

                    return (
                      <button
                        key={iso}
                        type="button"
                        onClick={() => {
                          setSelectedDate(d);
                          setIsWorkoutModalOpen(true);
                        }}
                        className={cn(
                          "rounded-md border border-border text-left p-2 min-h-20",
                          inMonth ? "bg-background" : "bg-muted/30",
                          isSelected && "ring-2 ring-ring",
                        )}
                      >
                        <div className="flex items-start justify-between">
                          <div
                            className={cn(
                              "text-sm font-semibold",
                              inMonth
                                ? "text-foreground"
                                : "text-muted-foreground",
                            )}
                          >
                            {format(d, "d")}
                          </div>
                          {hasWorkout && (
                            <div className="text-xs text-muted-foreground">
                              {dayWorkouts.length}×
                            </div>
                          )}
                        </div>

                        {hasWorkout ? (
                          <div className="mt-2 space-y-1">
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                              <Dumbbell className="h-3 w-3" />
                              <span>Workout</span>
                            </div>
                            <div className="text-xs text-muted-foreground">
                              {dayWorkouts[0]?.exercises?.length || 0} exercises
                            </div>
                          </div>
                        ) : (
                          <div className="mt-2 text-xs text-muted-foreground">
                            No workout
                          </div>
                        )}
                      </button>
                    );
                  })}
                </div>
              </CardContent>
            </Card>

            <Dialog
              open={isWorkoutModalOpen}
              onOpenChange={setIsWorkoutModalOpen}
            >
              <DialogContent className="sm:max-w-3xl max-h-[80vh] overflow-y-auto">
                <DialogHeader>
                  <DialogTitle className="flex items-center gap-2">
                    <FileText className="h-5 w-5" />
                    Workout
                  </DialogTitle>
                  <DialogDescription>
                    {format(selectedDate, "PPP")}
                  </DialogDescription>
                </DialogHeader>

                {selectedWorkouts.length === 0 ? (
                  <div className="space-y-3">
                    <div className="text-sm text-muted-foreground">
                      No workout recorded for this day.
                    </div>
                    <WorkoutForm
                      initialSessionDate={selectedISO}
                      onSuccess={() => {
                        loadWorkouts();
                        setIsWorkoutModalOpen(false);
                      }}
                    />
                  </div>
                ) : (
                  <div className="space-y-4">
                    {selectedWorkouts.map((w) => (
                      <div
                        key={w.id}
                        className="rounded-md border border-border p-3 space-y-3"
                      >
                        <div className="flex items-center justify-between gap-3">
                          <div className="text-sm font-semibold text-foreground">
                            Workout Session
                          </div>
                          <Link
                            href={`/app/workouts/${w.id}`}
                            className="text-sm text-muted-foreground underline"
                          >
                            View details
                          </Link>
                        </div>

                        <div className="text-sm text-muted-foreground flex items-center gap-2">
                          <Clock className="h-4 w-4" />
                          <span>
                            {w.session_date
                              ? format(new Date(w.session_date), "p")
                              : ""}
                          </span>
                        </div>

                        {(w.exercises || []).length === 0 ? (
                          <div className="text-sm text-muted-foreground">
                            No exercises.
                          </div>
                        ) : (
                          <div className="space-y-2">
                            {(w.exercises || []).map((ex, idx) => (
                              <div
                                key={`${w.id}-${ex.exercise_id}-${idx}`}
                                className="rounded-md border border-border p-2"
                              >
                                <div className="text-sm font-medium text-foreground">
                                  {exerciseNameById[ex.exercise_id] ||
                                    ex.exercise_id}
                                </div>

                                {(ex.sets || []).length === 0 ? (
                                  <div className="text-xs text-muted-foreground mt-1">
                                    No sets.
                                  </div>
                                ) : (
                                  <div className="mt-2 space-y-1">
                                    {(ex.sets || []).map((s) => (
                                      <div
                                        key={`${w.id}-${ex.exercise_id}-${s.set_order}`}
                                        className="text-xs text-muted-foreground flex items-center justify-between gap-3"
                                      >
                                        <span className="min-w-0 truncate">
                                          Set {s.set_order} • {s.reps} reps
                                          {s.weight ? ` @ ${s.weight}` : ""}
                                        </span>
                                        <span className="whitespace-nowrap">
                                          {s.set_type}
                                          {s.rpe ? ` • RPE ${s.rpe}` : ""}
                                        </span>
                                      </div>
                                    ))}
                                  </div>
                                )}
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </DialogContent>
            </Dialog>
          </>
        )}
      </div>
    </div>
  );
}
