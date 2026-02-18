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
}

export function WorkoutForm({ onSuccess }: WorkoutFormProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [sessionDate, setSessionDate] = useState(
    new Date().toISOString().split("T")[0],
  );
  const selectedDate = fromISODate(sessionDate);
  const [exerciseId, setExerciseId] = useState("");
  const [duration, setDuration] = useState("");
  const [notes, setNotes] = useState("");
  const [reps, setReps] = useState("8");
  const [weight, setWeight] = useState("0");
  const [rpe, setRpe] = useState("7");
  const [setType, setSetType] = useState("working");
  const [exercises, setExercises] = useState<ExerciseResponseDTO[]>([]);

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

    if (!exerciseId || !duration) {
      toast.error("Please fill in all required fields");
      return;
    }

    setIsLoading(true);

    try {
      await api.post("/api/v1/workouts", {
        user_id: userId,
        split_day_id: null,
        session_date: sessionDate,
        duration_minutes: parseInt(duration),
        notes: notes || "",
        exercises: [
          {
            exercise_id: exerciseId,
            sets: [
              {
                set_order: 1,
                reps: parseInt(reps),
                weight: parseFloat(weight),
                rpe: parseFloat(rpe),
                set_type: setType,
              },
            ],
          },
        ],
      });
      toast.success("Workout created successfully");
      setExerciseId("");
      setDuration("");
      setNotes("");
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
          <Label htmlFor="exercise">Exercise</Label>
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

        <div className="space-y-2">
          <Label htmlFor="reps">Reps</Label>
          <Input
            id="reps"
            type="number"
            value={reps}
            onChange={(e) => setReps(e.target.value)}
            disabled={isLoading}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="weight">Weight</Label>
          <Input
            id="weight"
            type="number"
            value={weight}
            onChange={(e) => setWeight(e.target.value)}
            disabled={isLoading}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="rpe">RPE</Label>
          <Input
            id="rpe"
            type="number"
            value={rpe}
            onChange={(e) => setRpe(e.target.value)}
            disabled={isLoading}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="setType">Set Type</Label>
          <Select value={setType} onValueChange={setSetType}>
            <SelectTrigger id="setType">
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
