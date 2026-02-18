"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import { Split } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import type {
  ExerciseResponseDTO,
  SplitTemplateResponseDTO,
} from "@/lib/backend-dto";
import { FOCUS_MUSCLE_OPTIONS } from "@/lib/grammar";

export default function SplitsPage() {
  const [splits, setSplits] = useState<SplitTemplateResponseDTO[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [splitName, setSplitName] = useState("");
  const [focusMuscle, setFocusMuscle] = useState("");
  const [dayName, setDayName] = useState("Day 1");
  const [exerciseId, setExerciseId] = useState("");
  const [targetSets, setTargetSets] = useState("3");
  const [targetReps, setTargetReps] = useState("10");
  const [targetWeight, setTargetWeight] = useState("0");
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

    if (!splitName || !focusMuscle || !exerciseId) {
      toast.error("Please enter a split name");
      return;
    }

    setIsSubmitting(true);
    try {
      await api.post("/api/v1/splits", {
        user_id: userId,
        name: splitName,
        description: "",
        created_by: userId,
        focus_muscle: focusMuscle,
        is_active: false,
        days: [
          {
            day_order: 1,
            name: dayName || "Day 1",
            exercises: [
              {
                exercise_id: exerciseId,
                target_sets: parseInt(targetSets) || 3,
                target_reps: parseInt(targetReps) || 10,
                target_weight: parseFloat(targetWeight) || 0,
                notes: "",
              },
            ],
          },
        ],
      });
      toast.success("Split created successfully");
      setSplitName("");
      setFocusMuscle("");
      setDayName("Day 1");
      setExerciseId("");
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

                <div className="space-y-2">
                  <Label htmlFor="dayName">Day Name</Label>
                  <Input
                    id="dayName"
                    placeholder="Day 1"
                    value={dayName}
                    onChange={(e) => setDayName(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="exercise">Day 1 Exercise</Label>
                  <Select value={exerciseId} onValueChange={setExerciseId}>
                    <SelectTrigger id="exercise">
                      <SelectValue placeholder="Select an exercise" />
                    </SelectTrigger>
                    <SelectContent>
                      {exercises.map((ex) => (
                        <SelectItem key={ex.id} value={ex.id}>
                          {ex.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="grid grid-cols-3 gap-3">
                  <div className="space-y-2">
                    <Label htmlFor="targetSets">Sets</Label>
                    <Input
                      id="targetSets"
                      type="number"
                      value={targetSets}
                      onChange={(e) => setTargetSets(e.target.value)}
                      disabled={isSubmitting}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="targetReps">Reps</Label>
                    <Input
                      id="targetReps"
                      type="number"
                      value={targetReps}
                      onChange={(e) => setTargetReps(e.target.value)}
                      disabled={isSubmitting}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="targetWeight">Weight</Label>
                    <Input
                      id="targetWeight"
                      type="number"
                      value={targetWeight}
                      onChange={(e) => setTargetWeight(e.target.value)}
                      disabled={isSubmitting}
                    />
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
        ) : splits.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Split className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">
                No splits yet. Create your first split!
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {splits.map((split) => (
              <Link key={split.id} href={`/app/splits/${split.id}`}>
                <Card className="h-full hover:shadow-lg transition-shadow cursor-pointer">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <CardTitle className="text-lg">{split.name}</CardTitle>
                      {split.is_active && (
                        <span className="text-xs bg-primary text-primary-foreground px-2 py-1 rounded">
                          Active
                        </span>
                      )}
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-2">
                    <div className="text-sm text-muted-foreground">
                      {split.days.length} days per week
                    </div>
                    {split.focus_muscle && (
                      <div className="text-sm text-muted-foreground">
                        Focus: {split.focus_muscle}
                      </div>
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
