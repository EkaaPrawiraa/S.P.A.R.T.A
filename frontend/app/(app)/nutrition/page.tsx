"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
import { Apple, CalendarIcon } from "lucide-react";
import { toast } from "sonner";
import type { DailyNutritionResponseDTO } from "@/lib/backend-dto";
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
  const [y, m, d] = value.split("-").map((x) => parseInt(x, 10));
  return new Date(y, (m || 1) - 1, d || 1);
}

interface ChartData {
  date: string;
  protein: number;
}

export default function NutritionPage() {
  const [selectedDate, setSelectedDate] = useState(
    new Date().toISOString().split("T")[0],
  );
  const selectedDateObj = fromISODate(selectedDate);
  const [nutrition, setNutrition] = useState<DailyNutritionResponseDTO | null>(
    null,
  );
  const [chartData, setChartData] = useState<ChartData[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [protein, setProtein] = useState("");
  const [calories, setCalories] = useState("");
  const [notes, setNotes] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    loadNutrition();
    loadWeeklyData();
  }, [selectedDate]);

  const loadNutrition = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      setIsLoading(true);
      const data = await api.get<DailyNutritionResponseDTO>(
        `/api/v1/nutrition/user/${userId}?date=${selectedDate}`,
      );
      setNutrition(data);
      setProtein(data?.protein_grams.toString() || "");
      setCalories(data?.calories?.toString() || "");
      setNotes(data?.notes || "");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load nutrition";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  const loadWeeklyData = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      const data: ChartData[] = [];
      const today = fromISODate(selectedDate);

      for (let i = 6; i >= 0; i--) {
        const date = new Date(today);
        date.setDate(date.getDate() - i);
        const dateStr = toISODate(date);

        try {
          const nutrition = await api.get<DailyNutritionResponseDTO>(
            `/api/v1/nutrition/user/${userId}?date=${dateStr}`,
          );
          data.push({
            date: dateStr,
            protein: nutrition?.protein_grams || 0,
          });
        } catch {
          data.push({
            date: dateStr,
            protein: 0,
          });
        }
      }

      setChartData(data);
    } catch (err) {
      console.error("Failed to load weekly data");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!protein) {
      toast.error("Please enter protein amount");
      return;
    }

    setIsSubmitting(true);

    try {
      const userId = getUserId();
      if (!userId) return;

      await api.post("/api/v1/nutrition", {
        user_id: userId,
        date: selectedDate,
        protein_grams: parseInt(protein),
        calories: calories ? parseInt(calories) : undefined,
        notes: notes || "",
      });

      toast.success("Nutrition updated successfully");
      loadNutrition();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to update nutrition";
      toast.error(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Nutrition</h1>
          <p className="text-muted-foreground mt-2">
            Track your daily nutrition intake
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Weekly Protein Trend</CardTitle>
              </CardHeader>
              <CardContent>
                {chartData.length > 0 ? (
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
                        dataKey="protein"
                        stroke="var(--primary)"
                        dot={false}
                        strokeWidth={2}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                ) : (
                  <div className="h-60 flex items-center justify-center text-muted-foreground">
                    No data available
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Apple className="h-5 w-5" />
                Add Nutrition
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="date">Date</Label>
                  <Popover>
                    <PopoverTrigger asChild>
                      <Button
                        id="date"
                        type="button"
                        variant="outline"
                        disabled={isSubmitting}
                        className={cn(
                          "w-full justify-start text-left font-normal",
                          !selectedDate && "text-muted-foreground",
                        )}
                      >
                        <CalendarIcon className="mr-2 h-4 w-4" />
                        {selectedDate
                          ? format(selectedDateObj, "PPP")
                          : "Pick a date"}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={selectedDateObj}
                        onSelect={(d) => {
                          if (!d) return;
                          setSelectedDate(toISODate(d));
                        }}
                        initialFocus
                      />
                    </PopoverContent>
                  </Popover>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="protein">Protein (g)</Label>
                  <Input
                    id="protein"
                    type="number"
                    placeholder="150"
                    value={protein}
                    onChange={(e) => setProtein(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="calories">Calories</Label>
                  <Input
                    id="calories"
                    type="number"
                    placeholder="2500"
                    value={calories}
                    onChange={(e) => setCalories(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="notes">Notes</Label>
                  <Input
                    id="notes"
                    placeholder="Optional notes"
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    disabled={isSubmitting}
                  />
                </div>
                <Button
                  type="submit"
                  disabled={isSubmitting}
                  className="w-full"
                >
                  Save Nutrition
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
