"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CalendarIcon, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import type { ExerciseResponseDTO } from "@/lib/backend-dto";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import { ExerciseCombobox } from "@/components/exercise-combobox";

function toISODate(date: Date): string {
  return date.toISOString().slice(0, 10);
}

function fromISODate(value: string): Date {
  // Interpret as local date to avoid timezone shifts.
  const [y, m, d] = value.split("-").map((x) => parseInt(x, 10));
  return new Date(y, (m || 1) - 1, d || 1);
}

interface WorkoutFormProps {
  onSuccess?: () => void;
  initialSessionDate?: string;
}

type WorkoutSetDraft = {
  set_order: number;
  reps: string;
  weight: string;
  rpe: string;
  set_type: string;
};

type WorkoutExerciseDraft = {
  exercise_id: string;
  sets: WorkoutSetDraft[];
};

function newSetDraft(setOrder: number): WorkoutSetDraft {
  return {
    set_order: setOrder,
    reps: "8",
    weight: "0",
    rpe: "7",
    set_type: "working",
  };
}

function newExerciseDraft(): WorkoutExerciseDraft {
  return {
    exercise_id: "",
    sets: [newSetDraft(1)],
  };
}

export function WorkoutForm({
  onSuccess,
  initialSessionDate,
}: WorkoutFormProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [sessionDate, setSessionDate] = useState(() =>
    initialSessionDate && initialSessionDate.trim().length > 0
      ? initialSessionDate
      : new Date().toISOString().split("T")[0],
  );
  const selectedDate = fromISODate(sessionDate);
  const [duration, setDuration] = useState("");
  const [notes, setNotes] = useState("");
  const [workoutExercises, setWorkoutExercises] = useState<
    WorkoutExerciseDraft[]
  >(() => [newExerciseDraft()]);
  const [exercises, setExercises] = useState<ExerciseResponseDTO[]>([]);

  useEffect(() => {
    if (!initialSessionDate) return;
    const d = initialSessionDate.trim();
    if (d.length === 0) return;
    setSessionDate(d);
  }, [initialSessionDate]);

  useEffect(() => {
    const loadExercises = async () => {
      try {
        const data = await api.get<ExerciseResponseDTO[]>("/api/v1/exercises");
        setExercises(data || []);
      } catch {
        setExercises([]);
      }
    };
    loadExercises();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const userId = getUserId();
    if (!userId) {
      toast.error("You must be logged in");
      return;
    }

    if (!duration) {
      toast.error("Please fill in all required fields");
      return;
    }

    const cleanedExercises = (workoutExercises || []).filter(
      (ex) => (ex.exercise_id || "").trim().length > 0,
    );
    if (cleanedExercises.length === 0) {
      toast.error("Please add at least one exercise");
      return;
    }

    for (const ex of cleanedExercises) {
      if (!ex.sets || ex.sets.length === 0) {
        toast.error("Each exercise must have at least one set");
        return;
      }
    }

    setIsLoading(true);

    try {
      await api.post("/api/v1/workouts", {
        user_id: userId,
        split_day_id: null,
        session_date: sessionDate,
        duration_minutes: parseInt(duration),
        notes: notes || "",
        exercises: cleanedExercises.map((ex) => ({
          exercise_id: ex.exercise_id,
          sets: (ex.sets || []).map((s) => ({
            set_order: s.set_order,
            reps: parseInt(s.reps) || 0,
            weight: parseFloat(s.weight) || 0,
            rpe: parseFloat(s.rpe) || 0,
            set_type: s.set_type,
          })),
        })),
      });
      toast.success("Workout created successfully");
      setDuration("");
      setNotes("");
      setWorkoutExercises([newExerciseDraft()]);
      onSuccess?.();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to create workout";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="date">Session Date</Label>
          <Popover>
            <PopoverTrigger asChild>
              <Button
                id="date"
                type="button"
                variant="outline"
                disabled={isLoading}
                className={cn(
                  "w-full justify-start text-left font-normal",
                  !sessionDate && "text-muted-foreground",
                )}
              >
                <CalendarIcon className="mr-2 h-4 w-4" />
                {sessionDate ? format(selectedDate, "PPP") : "Pick a date"}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="start">
              <Calendar
                mode="single"
                selected={selectedDate}
                onSelect={(d) => {
                  if (!d) return;
                  setSessionDate(toISODate(d));
                }}
                initialFocus
              />
            </PopoverContent>
          </Popover>
        </div>

        <div className="space-y-2">
          <Label htmlFor="duration">Duration (minutes)</Label>
          <Input
            id="duration"
            type="number"
            placeholder="60"
            value={duration}
            onChange={(e) => setDuration(e.target.value)}
            disabled={isLoading}
          />
        </div>
      </div>

      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label>Exercises</Label>
          <Button
            type="button"
            variant="outline"
            disabled={isLoading}
            onClick={() =>
              setWorkoutExercises((prev) => [
                ...(prev || []),
                newExerciseDraft(),
              ])
            }
          >
            Add Exercise
          </Button>
        </div>

        <div className="space-y-4">
          {(workoutExercises || []).map((ex, exIdx) => (
            <div
              key={exIdx}
              className="rounded-md border border-border p-3 space-y-3"
            >
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3 items-end">
                <div className="space-y-2">
                  <Label>Exercise</Label>
                  <ExerciseCombobox
                    value={ex.exercise_id}
                    onValueChange={(v) =>
                      setWorkoutExercises((prev) =>
                        (prev || []).map((item, idx) =>
                          idx === exIdx ? { ...item, exercise_id: v } : item,
                        ),
                      )
                    }
                    options={exercises}
                    disabled={isLoading}
                    placeholder="Select an exercise"
                    searchPlaceholder="Search exercise..."
                  />
                </div>

                <div className="flex gap-2 justify-end">
                  <Button
                    type="button"
                    variant="outline"
                    disabled={isLoading}
                    onClick={() =>
                      setWorkoutExercises((prev) =>
                        (prev || []).map((item, idx) => {
                          if (idx !== exIdx) return item;
                          const nextOrder = (item.sets?.length || 0) + 1;
                          return {
                            ...item,
                            sets: [
                              ...(item.sets || []),
                              newSetDraft(nextOrder),
                            ],
                          };
                        }),
                      )
                    }
                  >
                    Add Set
                  </Button>

                  <Button
                    type="button"
                    variant="outline"
                    disabled={isLoading || (workoutExercises || []).length <= 1}
                    onClick={() =>
                      setWorkoutExercises((prev) =>
                        (prev || []).filter((_, idx) => idx !== exIdx),
                      )
                    }
                  >
                    Delete Exercise
                  </Button>
                </div>
              </div>

              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/40">
                      <th className="text-left py-2 pr-2 font-semibold text-foreground">
                        Set
                      </th>
                      <th className="text-left py-2 pr-2 font-semibold text-foreground">
                        Reps
                      </th>
                      <th className="text-left py-2 pr-2 font-semibold text-foreground">
                        Weight
                      </th>
                      <th className="text-left py-2 pr-2 font-semibold text-foreground">
                        RPE
                      </th>
                      <th className="text-left py-2 pr-2 font-semibold text-foreground">
                        Type
                      </th>
                      <th className="text-left py-2 font-semibold text-foreground">
                        &nbsp;
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {(ex.sets || []).map((s, setIdx) => (
                      <tr key={setIdx} className="border-b border-border/40">
                        <td className="py-2 pr-2 text-muted-foreground">
                          {s.set_order}
                        </td>
                        <td className="py-2 pr-2">
                          <Input
                            type="number"
                            value={s.reps}
                            disabled={isLoading}
                            onChange={(e) =>
                              setWorkoutExercises((prev) =>
                                (prev || []).map((item, idx) => {
                                  if (idx !== exIdx) return item;
                                  const sets = (item.sets || []).map((x, j) =>
                                    j === setIdx
                                      ? { ...x, reps: e.target.value }
                                      : x,
                                  );
                                  return { ...item, sets };
                                }),
                              )
                            }
                          />
                        </td>
                        <td className="py-2 pr-2">
                          <Input
                            type="number"
                            value={s.weight}
                            disabled={isLoading}
                            onChange={(e) =>
                              setWorkoutExercises((prev) =>
                                (prev || []).map((item, idx) => {
                                  if (idx !== exIdx) return item;
                                  const sets = (item.sets || []).map((x, j) =>
                                    j === setIdx
                                      ? { ...x, weight: e.target.value }
                                      : x,
                                  );
                                  return { ...item, sets };
                                }),
                              )
                            }
                          />
                        </td>
                        <td className="py-2 pr-2">
                          <Input
                            type="number"
                            value={s.rpe}
                            disabled={isLoading}
                            onChange={(e) =>
                              setWorkoutExercises((prev) =>
                                (prev || []).map((item, idx) => {
                                  if (idx !== exIdx) return item;
                                  const sets = (item.sets || []).map((x, j) =>
                                    j === setIdx
                                      ? { ...x, rpe: e.target.value }
                                      : x,
                                  );
                                  return { ...item, sets };
                                }),
                              )
                            }
                          />
                        </td>
                        <td className="py-2 pr-2">
                          <Select
                            value={s.set_type}
                            onValueChange={(v) =>
                              setWorkoutExercises((prev) =>
                                (prev || []).map((item, idx) => {
                                  if (idx !== exIdx) return item;
                                  const sets = (item.sets || []).map((x, j) =>
                                    j === setIdx ? { ...x, set_type: v } : x,
                                  );
                                  return { ...item, sets };
                                }),
                              )
                            }
                          >
                            <SelectTrigger disabled={isLoading}>
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="working">Working</SelectItem>
                              <SelectItem value="warmup">Warmup</SelectItem>
                              <SelectItem value="failure">Failure</SelectItem>
                              <SelectItem value="dropset">Dropset</SelectItem>
                              <SelectItem value="superset">Superset</SelectItem>
                            </SelectContent>
                          </Select>
                        </td>
                        <td className="py-2">
                          <Button
                            type="button"
                            variant="outline"
                            disabled={isLoading || (ex.sets || []).length <= 1}
                            onClick={() =>
                              setWorkoutExercises((prev) =>
                                (prev || []).map((item, idx) => {
                                  if (idx !== exIdx) return item;
                                  const nextSets = (item.sets || [])
                                    .filter((_, j) => j !== setIdx)
                                    .map((x, k) => ({
                                      ...x,
                                      set_order: k + 1,
                                    }));
                                  return { ...item, sets: nextSets };
                                }),
                              )
                            }
                          >
                            Delete
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="notes">Notes</Label>
        <Input
          id="notes"
          placeholder="How did the workout feel?"
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          disabled={isLoading}
        />
      </div>

      <Button type="submit" disabled={isLoading} className="w-full md:w-auto">
        {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Create Workout
      </Button>
    </form>
  );
}
