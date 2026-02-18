"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import { BarChart3 } from "lucide-react";
import { toast } from "sonner";
import type {
  ExerciseResponseDTO,
  WorkoutSessionResponseDTO,
} from "@/lib/backend-dto";

interface ChartData {
  date: string;
  volume: number;
  weight: number;
  sets: number;
}

export default function AnalyticsPage() {
  const [selectedMuscle, setSelectedMuscle] = useState("All");
  const [chartData, setChartData] = useState<ChartData[]>([]);
  const [muscleGroups, setMuscleGroups] = useState<string[]>(["All"]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadAnalytics();
  }, [selectedMuscle]);

  const loadAnalytics = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      setIsLoading(true);

      // Fetch workouts and exercises
      const [workouts, exercises] = await Promise.all([
        api.get<WorkoutSessionResponseDTO[]>(`/api/v1/workouts/user/${userId}`),
        api.get<ExerciseResponseDTO[]>("/api/v1/exercises"),
      ]);

      // Build a map of exercise_id -> primary_muscle
      const exerciseMap = (exercises || []).reduce(
        (acc, ex) => {
          acc[ex.id] = ex.primary_muscle;
          return acc;
        },
        {} as Record<string, string>,
      );

      const groups = Array.from(
        new Set(
          (exercises || [])
            .map((e) => e.primary_muscle)
            .filter((m): m is string => Boolean(m && m.trim())),
        ),
      ).sort((a, b) => a.localeCompare(b));

      const nextGroups = ["All", ...groups];
      setMuscleGroups(nextGroups);
      if (selectedMuscle !== "All" && !nextGroups.includes(selectedMuscle)) {
        setSelectedMuscle("All");
      }

      // Compute metrics
      const data: ChartData[] = [];

      (workouts || []).forEach((workout) => {
        let volume = 0;
        let totalWeight = 0;
        let totalSets = 0;
        let weightSamples = 0;

        (workout.exercises || []).forEach((ex) => {
          const muscle = exerciseMap[ex.exercise_id];
          if (selectedMuscle !== "All" && muscle !== selectedMuscle) return;
          (ex.sets || []).forEach((set) => {
            const reps = set.reps || 0;
            const weight = set.weight || 0;
            volume += reps * weight;
            totalWeight += weight;
            weightSamples++;
            totalSets += 1;
          });
        });

        if (totalSets > 0) {
          data.push({
            date: workout.session_date,
            volume,
            weight: weightSamples > 0 ? totalWeight / weightSamples : 0,
            sets: totalSets,
          });
        }
      });

      setChartData(
        data.sort(
          (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime(),
        ),
      );
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load analytics";
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
            <h1 className="text-3xl font-bold text-foreground">Analytics</h1>
            <p className="text-muted-foreground mt-2">
              Visualize your training progress
            </p>
          </div>
          <Select value={selectedMuscle} onValueChange={setSelectedMuscle}>
            <SelectTrigger className="w-48">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {muscleGroups.map((muscle) => (
                <SelectItem key={muscle} value={muscle}>
                  {muscle}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading analytics...</p>
          </div>
        ) : chartData.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <BarChart3 className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">
                No workout data available. Start logging workouts to see
                analytics.
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Volume Trend ({selectedMuscle})</CardTitle>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={chartData}>
                    <CartesianGrid
                      strokeDasharray="3 3"
                      stroke="rgba(255,255,255,0.1)"
                    />
                    <XAxis
                      dataKey="date"
                      stroke="var(--foreground)"
                      style={{ fontSize: "12px" }}
                    />
                    <YAxis
                      stroke="var(--foreground)"
                      style={{ fontSize: "12px" }}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--background)",
                        border: "1px solid var(--border)",
                        borderRadius: "8px",
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="volume"
                      stroke="var(--primary)"
                      dot={false}
                      strokeWidth={2}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Average Working Weight ({selectedMuscle})</CardTitle>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={chartData}>
                    <CartesianGrid
                      strokeDasharray="3 3"
                      stroke="rgba(255,255,255,0.1)"
                    />
                    <XAxis
                      dataKey="date"
                      stroke="var(--foreground)"
                      style={{ fontSize: "12px" }}
                    />
                    <YAxis
                      stroke="var(--foreground)"
                      style={{ fontSize: "12px" }}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--background)",
                        border: "1px solid var(--border)",
                        borderRadius: "8px",
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="weight"
                      stroke="var(--primary)"
                      dot={false}
                      strokeWidth={2}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Total Sets ({selectedMuscle})</CardTitle>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={chartData}>
                    <CartesianGrid
                      strokeDasharray="3 3"
                      stroke="rgba(255,255,255,0.1)"
                    />
                    <XAxis
                      dataKey="date"
                      stroke="var(--foreground)"
                      style={{ fontSize: "12px" }}
                    />
                    <YAxis
                      stroke="var(--foreground)"
                      style={{ fontSize: "12px" }}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--background)",
                        border: "1px solid var(--border)",
                        borderRadius: "8px",
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="sets"
                      stroke="var(--primary)"
                      dot={false}
                      strokeWidth={2}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}
