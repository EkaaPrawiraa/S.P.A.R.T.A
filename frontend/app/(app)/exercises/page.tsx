"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { Zap } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import type { ExerciseResponseDTO } from "@/lib/backend-dto";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";
import { EQUIPMENT_OPTIONS, MUSCLE_OPTIONS } from "@/lib/grammar";

export default function ExercisesPage() {
  const [exercises, setExercises] = useState<ExerciseResponseDTO[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [exerciseName, setExerciseName] = useState("");
  const [primaryMuscle, setPrimaryMuscle] = useState("");
  const [equipment, setEquipment] = useState("");
  const [secondaryMuscles, setSecondaryMuscles] = useState<string[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    loadExercises();
  }, []);

  const loadExercises = async () => {
    try {
      const data = await api.get<ExerciseResponseDTO[]>("/api/v1/exercises");
      setExercises(data || []);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load exercises";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateExercise = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!exerciseName || !primaryMuscle || !equipment) {
      toast.error("Please fill in all fields");
      return;
    }

    setIsSubmitting(true);
    try {
      await api.post("/api/v1/exercises", {
        name: exerciseName,
        primary_muscle: primaryMuscle,
        equipment,
        secondary_muscles: secondaryMuscles,
      });
      toast.success("Exercise created successfully");
      setExerciseName("");
      setPrimaryMuscle("");
      setEquipment("");
      setSecondaryMuscles([]);
      setShowForm(false);
      loadExercises();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to create exercise";
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
            <h1 className="text-3xl font-bold text-foreground">Exercises</h1>
            <p className="text-muted-foreground mt-2">
              Browse and manage your exercise database
            </p>
          </div>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? "Cancel" : "New Exercise"}
          </Button>
        </div>

        {showForm && (
          <Card>
            <CardHeader>
              <CardTitle>Create New Exercise</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleCreateExercise} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Exercise Name</Label>
                  <Input
                    id="name"
                    placeholder="e.g., Barbell Bench Press"
                    value={exerciseName}
                    onChange={(e) => setExerciseName(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="muscle">Primary Muscle</Label>
                  <Select
                    value={primaryMuscle}
                    onValueChange={setPrimaryMuscle}
                  >
                    <SelectTrigger id="muscle" disabled={isSubmitting}>
                      <SelectValue placeholder="Select a muscle" />
                    </SelectTrigger>
                    <SelectContent>
                      {MUSCLE_OPTIONS.map((m) => (
                        <SelectItem key={m} value={m}>
                          {m}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="equipment">Equipment</Label>
                  <Select value={equipment} onValueChange={setEquipment}>
                    <SelectTrigger id="equipment" disabled={isSubmitting}>
                      <SelectValue placeholder="Select equipment" />
                    </SelectTrigger>
                    <SelectContent>
                      {EQUIPMENT_OPTIONS.map((eq) => (
                        <SelectItem key={eq} value={eq}>
                          {eq}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="secondaryMuscles">
                    Secondary Muscles (choose any)
                  </Label>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 rounded-md border border-border p-3 max-h-56 overflow-auto">
                    {MUSCLE_OPTIONS.map((m) => {
                      const checked = secondaryMuscles.includes(m);
                      return (
                        <label
                          key={m}
                          className="flex items-center gap-2 text-sm text-foreground"
                        >
                          <Checkbox
                            checked={checked}
                            onCheckedChange={(v) => {
                              const isChecked = v === true;
                              setSecondaryMuscles((prev) =>
                                isChecked
                                  ? Array.from(new Set([...prev, m]))
                                  : prev.filter((x) => x !== m),
                              );
                            }}
                            disabled={isSubmitting}
                          />
                          <span className="text-muted-foreground">{m}</span>
                        </label>
                      );
                    })}
                  </div>
                </div>
                <Button type="submit" disabled={isSubmitting}>
                  Create Exercise
                </Button>
              </form>
            </CardContent>
          </Card>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading exercises...</p>
          </div>
        ) : exercises.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Zap className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">
                No exercises yet. Create your first exercise!
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {exercises.map((exercise) => (
              <Link key={exercise.id} href={`/app/exercises/${exercise.id}`}>
                <Card className="h-full hover:shadow-lg transition-shadow cursor-pointer">
                  <CardHeader>
                    <CardTitle className="text-lg">{exercise.name}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-sm text-muted-foreground">
                      Primary: {exercise.primary_muscle}
                    </div>
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
