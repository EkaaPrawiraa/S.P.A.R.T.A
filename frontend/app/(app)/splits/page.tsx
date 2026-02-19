"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import {
  CheckCircle2,
  CircleOff,
  MinusCircle,
  Plus,
  Split,
  Trash2,
} from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import type {
  ExerciseResponseDTO,
  SplitTemplateResponseDTO,
} from "@/lib/backend-dto";
import { FOCUS_MUSCLE_OPTIONS } from "@/lib/grammar";
import { ExerciseCombobox } from "@/components/exercise-combobox";

type DayExerciseForm = {
  id: string;
  exerciseId: string;
  targetSets: string;
  targetReps: string;
  targetWeight: string;
  notes: string;
};

type DayForm = {
  id: string;
  name: string;
  exercises: DayExerciseForm[];
};

function newId() {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

function createExercise(): DayExerciseForm {
  return {
    id: newId(),
    exerciseId: "",
    targetSets: "3",
    targetReps: "10",
    targetWeight: "0",
    notes: "",
  };
}

function createDay(dayOrder: number): DayForm {
  return {
    id: newId(),
    name: `Day ${dayOrder}`,
    exercises: [createExercise()],
  };
}

export default function SplitsPage() {
  const [splits, setSplits] = useState<SplitTemplateResponseDTO[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [focusFilter, setFocusFilter] = useState<string>("");
  const [daysPerWeekFilter, setDaysPerWeekFilter] = useState<string>("");
  const [splitName, setSplitName] = useState("");
  const [focusMuscle, setFocusMuscle] = useState("");
  const [days, setDays] = useState<DayForm[]>([createDay(1)]);
  const [exercises, setExercises] = useState<ExerciseResponseDTO[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    loadSplits();
    loadExercises();
  }, []);

  const loadSplits = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      const data = await api.get<SplitTemplateResponseDTO[]>(
        `/api/v1/splits/user/${userId}`,
      );
      setSplits(data || []);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load splits";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  const loadExercises = async () => {
    try {
      const data = await api.get<ExerciseResponseDTO[]>("/api/v1/exercises");
      setExercises(data || []);
    } catch {
      setExercises([]);
    }
  };

  const handleCreateSplit = async (e: React.FormEvent) => {
    e.preventDefault();
    const userId = getUserId();
    if (!userId) {
      toast.error("You must be logged in");
      return;
    }

    const trimmedName = splitName.trim();
    if (!trimmedName || !focusMuscle.trim()) {
      toast.error("Please enter split name and focus muscle");
      return;
    }

    if (!days.length) {
      toast.error("Please add at least 1 day");
      return;
    }

    for (const day of days) {
      if (!day.exercises.length) {
        toast.error("Each day must have at least 1 exercise");
        return;
      }
      for (const ex of day.exercises) {
        if (!ex.exerciseId) {
          toast.error("Please select an exercise for each row");
          return;
        }
      }
    }

    setIsSubmitting(true);
    try {
      await api.post("/api/v1/splits", {
        user_id: userId,
        name: trimmedName,
        description: "",
        created_by: "user",
        focus_muscle: focusMuscle.trim(),
        is_active: false,
        days: days.map((day, idx) => ({
          day_order: idx + 1,
          name: (day.name || `Day ${idx + 1}`).trim(),
          exercises: day.exercises.map((ex) => ({
            exercise_id: ex.exerciseId,
            target_sets: parseInt(ex.targetSets) || 3,
            target_reps: parseInt(ex.targetReps) || 10,
            target_weight: parseFloat(ex.targetWeight) || 0,
            notes: ex.notes || "",
          })),
        })),
      });
      toast.success("Split created successfully");
      setSplitName("");
      setFocusMuscle("");
      setDays([createDay(1)]);
      setShowForm(false);
      loadSplits();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to create split";
      toast.error(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeactivate = async (templateId: string) => {
    try {
      await api.post(`/api/v1/splits/${templateId}/deactivate`);
      toast.success("Split deactivated");
      loadSplits();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to deactivate split";
      toast.error(message);
    }
  };

  const filteredSplits = (splits || []).filter((s) => {
    const norm = (v: string | null | undefined) =>
      (v || "").trim().toLowerCase();
    if (focusFilter && norm(s.focus_muscle) !== norm(focusFilter)) return false;
    if (daysPerWeekFilter) {
      const n = parseInt(daysPerWeekFilter, 10);
      if (!Number.isNaN(n) && (s.days || []).length !== n) return false;
    }
    return true;
  });

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">
              Workout Splits
            </h1>
            <p className="text-muted-foreground mt-2">
              Organize your training days and muscle groups
            </p>
          </div>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? "Cancel" : "New Split"}
          </Button>
        </div>

        <div className="flex flex-col gap-4 md:flex-row md:items-end">
          <div className="space-y-2 md:flex-1">
            <Label htmlFor="focusFilter">Focus Muscle</Label>
            <Select value={focusFilter} onValueChange={setFocusFilter}>
              <SelectTrigger id="focusFilter">
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent>
                {FOCUS_MUSCLE_OPTIONS.map((m) => (
                  <SelectItem key={m} value={m}>
                    {m}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2 md:flex-1">
            <Label htmlFor="daysPerWeekFilter">Days / Week</Label>
            <Select
              value={daysPerWeekFilter}
              onValueChange={setDaysPerWeekFilter}
            >
              <SelectTrigger id="daysPerWeekFilter">
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent>
                {Array.from({ length: 7 }).map((_, i) => {
                  const v = String(i + 1);
                  return (
                    <SelectItem key={v} value={v}>
                      {v}
                    </SelectItem>
                  );
                })}
              </SelectContent>
            </Select>
          </div>

          <div className="md:flex md:justify-end">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => {
                setFocusFilter("");
                setDaysPerWeekFilter("");
              }}
              disabled={!focusFilter && !daysPerWeekFilter}
              className="w-full md:w-auto"
            >
              Reset
            </Button>
          </div>
        </div>

        {showForm && (
          <Card>
            <CardHeader>
              <CardTitle>Create New Split</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleCreateSplit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="splitName">Split Name</Label>
                  <Input
                    id="splitName"
                    placeholder="e.g., Push/Pull/Legs"
                    value={splitName}
                    onChange={(e) => setSplitName(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="focusMuscle">Focus Muscle</Label>
                  <Select value={focusMuscle} onValueChange={setFocusMuscle}>
                    <SelectTrigger id="focusMuscle" disabled={isSubmitting}>
                      <SelectValue placeholder="Select focus" />
                    </SelectTrigger>
                    <SelectContent>
                      {FOCUS_MUSCLE_OPTIONS.map((m) => (
                        <SelectItem key={m} value={m}>
                          {m}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center justify-between gap-3">
                    <Label>Days</Label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        setDays((prev) => [...prev, createDay(prev.length + 1)])
                      }
                      disabled={isSubmitting}
                    >
                      <Plus className="h-4 w-4" />
                      <span>Add Day</span>
                    </Button>
                  </div>

                  <div className="space-y-4">
                    {days.map((day, dayIndex) => (
                      <Card key={day.id} className="border">
                        <CardHeader className="pb-3">
                          <div className="flex items-start justify-between gap-3">
                            <div className="space-y-2 flex-1">
                              <Label
                                htmlFor={`dayName-${day.id}`}
                              >{`Day ${dayIndex + 1} Name`}</Label>
                              <Input
                                id={`dayName-${day.id}`}
                                value={day.name}
                                onChange={(e) => {
                                  const v = e.target.value;
                                  setDays((prev) =>
                                    prev.map((d) =>
                                      d.id === day.id ? { ...d, name: v } : d,
                                    ),
                                  );
                                }}
                                disabled={isSubmitting}
                              />
                            </div>

                            <Button
                              type="button"
                              variant="outline"
                              size="icon"
                              onClick={() =>
                                setDays((prev) =>
                                  prev.filter((d) => d.id !== day.id),
                                )
                              }
                              disabled={isSubmitting || days.length === 1}
                              className={days.length === 1 ? "invisible" : ""}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </CardHeader>
                        <CardContent className="space-y-3">
                          <div className="flex items-center justify-between gap-3">
                            <Label>Exercises</Label>
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={() => {
                                setDays((prev) =>
                                  prev.map((d) =>
                                    d.id === day.id
                                      ? {
                                          ...d,
                                          exercises: [
                                            ...d.exercises,
                                            createExercise(),
                                          ],
                                        }
                                      : d,
                                  ),
                                );
                              }}
                              disabled={isSubmitting}
                            >
                              <Plus className="h-4 w-4" />
                              <span>Add Exercise</span>
                            </Button>
                          </div>

                          <div className="space-y-4">
                            {day.exercises.map((exRow) => {
                              return (
                                <div
                                  key={exRow.id}
                                  className="grid grid-cols-1 gap-3 md:grid-cols-[1fr_auto] md:items-start"
                                >
                                  <div className="space-y-3">
                                    <div className="space-y-2">
                                      <Label>Exercise</Label>
                                      <ExerciseCombobox
                                        value={exRow.exerciseId}
                                        onValueChange={(v) => {
                                          setDays((prev) =>
                                            prev.map((d) => {
                                              if (d.id !== day.id) return d;
                                              return {
                                                ...d,
                                                exercises: d.exercises.map(
                                                  (row) =>
                                                    row.id === exRow.id
                                                      ? {
                                                          ...row,
                                                          exerciseId: v,
                                                        }
                                                      : row,
                                                ),
                                              };
                                            }),
                                          );
                                        }}
                                        options={exercises}
                                        disabled={isSubmitting}
                                        placeholder="Select an exercise"
                                        searchPlaceholder="Search exercise..."
                                      />
                                    </div>

                                    <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                                      <div className="space-y-2">
                                        <Label>Sets</Label>
                                        <Input
                                          type="number"
                                          value={exRow.targetSets}
                                          onChange={(e) => {
                                            const v = e.target.value;
                                            setDays((prev) =>
                                              prev.map((d) => {
                                                if (d.id !== day.id) return d;
                                                return {
                                                  ...d,
                                                  exercises: d.exercises.map(
                                                    (row) =>
                                                      row.id === exRow.id
                                                        ? {
                                                            ...row,
                                                            targetSets: v,
                                                          }
                                                        : row,
                                                  ),
                                                };
                                              }),
                                            );
                                          }}
                                          disabled={isSubmitting}
                                        />
                                      </div>
                                      <div className="space-y-2">
                                        <Label>Reps</Label>
                                        <Input
                                          type="number"
                                          value={exRow.targetReps}
                                          onChange={(e) => {
                                            const v = e.target.value;
                                            setDays((prev) =>
                                              prev.map((d) => {
                                                if (d.id !== day.id) return d;
                                                return {
                                                  ...d,
                                                  exercises: d.exercises.map(
                                                    (row) =>
                                                      row.id === exRow.id
                                                        ? {
                                                            ...row,
                                                            targetReps: v,
                                                          }
                                                        : row,
                                                  ),
                                                };
                                              }),
                                            );
                                          }}
                                          disabled={isSubmitting}
                                        />
                                      </div>
                                      <div className="space-y-2">
                                        <Label>Weight</Label>
                                        <Input
                                          type="number"
                                          value={exRow.targetWeight}
                                          onChange={(e) => {
                                            const v = e.target.value;
                                            setDays((prev) =>
                                              prev.map((d) => {
                                                if (d.id !== day.id) return d;
                                                return {
                                                  ...d,
                                                  exercises: d.exercises.map(
                                                    (row) =>
                                                      row.id === exRow.id
                                                        ? {
                                                            ...row,
                                                            targetWeight: v,
                                                          }
                                                        : row,
                                                  ),
                                                };
                                              }),
                                            );
                                          }}
                                          disabled={isSubmitting}
                                        />
                                      </div>
                                    </div>

                                    <div className="space-y-2">
                                      <Label>Notes</Label>
                                      <Input
                                        value={exRow.notes}
                                        onChange={(e) => {
                                          const v = e.target.value;
                                          setDays((prev) =>
                                            prev.map((d) => {
                                              if (d.id !== day.id) return d;
                                              return {
                                                ...d,
                                                exercises: d.exercises.map(
                                                  (row) =>
                                                    row.id === exRow.id
                                                      ? { ...row, notes: v }
                                                      : row,
                                                ),
                                              };
                                            }),
                                          );
                                        }}
                                        disabled={isSubmitting}
                                      />
                                    </div>
                                  </div>

                                  <div className="flex md:flex-col gap-2 md:pt-7">
                                    <Button
                                      type="button"
                                      variant="outline"
                                      size="icon"
                                      onClick={() => {
                                        setDays((prev) =>
                                          prev.map((d) => {
                                            if (d.id !== day.id) return d;
                                            return {
                                              ...d,
                                              exercises: d.exercises.filter(
                                                (row) => row.id !== exRow.id,
                                              ),
                                            };
                                          }),
                                        );
                                      }}
                                      disabled={
                                        isSubmitting ||
                                        (day.exercises || []).length === 1
                                      }
                                      className={
                                        (day.exercises || []).length === 1
                                          ? "invisible"
                                          : ""
                                      }
                                    >
                                      <Trash2 className="h-4 w-4" />
                                    </Button>
                                  </div>
                                </div>
                              );
                            })}
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </div>

                <Button type="submit" disabled={isSubmitting}>
                  Create Split
                </Button>
              </form>
            </CardContent>
          </Card>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading splits...</p>
          </div>
        ) : filteredSplits.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Split className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">No splits found.</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 items-stretch">
            {filteredSplits.map((split) => (
              <Card key={split.id} className="h-full flex flex-col">
                <CardHeader className="space-y-3">
                  <div className="min-w-0 space-y-3">
                    <Link
                      href={`/app/splits/${split.id}`}
                      className="hover:underline"
                    >
                      <div className="min-h-[2.75rem]">
                        <CardTitle className="text-lg leading-snug line-clamp-2">
                          {split.name}
                        </CardTitle>
                      </div>
                    </Link>

                    <div className="flex flex-wrap items-center gap-2">
                      {split.is_active ? (
                        <Badge>
                          <CheckCircle2 className="h-4 w-4" />
                          Active
                        </Badge>
                      ) : (
                        <Badge variant="outline">
                          <MinusCircle className="h-4 w-4" />
                          Inactive
                        </Badge>
                      )}

                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => handleDeactivate(split.id)}
                        disabled={!split.is_active}
                        className={
                          split.is_active
                            ? "h-8"
                            : "h-8 invisible pointer-events-none"
                        }
                      >
                        <CircleOff className="h-4 w-4" />
                        <span>Deactivate</span>
                      </Button>
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="space-y-2 mt-auto">
                  <div className="text-sm text-muted-foreground">
                    {(split.days || []).length} days per week
                  </div>
                  {split.focus_muscle && (
                    <div className="text-sm text-muted-foreground">
                      Focus: {split.focus_muscle}
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
