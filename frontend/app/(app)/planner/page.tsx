"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import {
  Calendar,
  ChevronLeft,
  ChevronRight,
  FileText,
  Loader2,
  Sparkles,
} from "lucide-react";
import { toast } from "sonner";
import type { PlannerRecommendationResponseDTO } from "@/lib/backend-dto";
import {
  addDays,
  addMonths,
  endOfMonth,
  endOfWeek,
  format,
  isSameDay,
  isSameMonth,
  startOfMonth,
  startOfWeek,
  subMonths,
} from "date-fns";
import { cn } from "@/lib/utils";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

export default function PlannerPage() {
  const [recommendations, setRecommendations] = useState<
    PlannerRecommendationResponseDTO[]
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isGenerating, setIsGenerating] = useState(false);
  const [monthCursor, setMonthCursor] = useState(() =>
    startOfMonth(new Date()),
  );
  const [selectedDate, setSelectedDate] = useState(() => new Date());
  const [isNoteOpen, setIsNoteOpen] = useState(false);

  useEffect(() => {
    loadRecommendations();
  }, []);

  const loadRecommendations = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      setIsLoading(true);
      const data = await api.get<PlannerRecommendationResponseDTO[]>(
        `/api/v1/planner/user/${userId}`,
      );
      setRecommendations(data || []);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load recommendations";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleGenerate = async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      setIsGenerating(true);
      await api.post(`/api/v1/planner/generate/${userId}`, {});
      toast.success("Recommendation generated successfully");
      loadRecommendations();
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Failed to generate recommendation";
      toast.error(message);
    } finally {
      setIsGenerating(false);
    }
  };

  const recommendationsByDate = (recommendations || []).reduce(
    (acc, rec) => {
      const key = (rec.created_at || "").slice(0, 10);
      if (!key) return acc;
      if (!acc[key]) acc[key] = [];
      acc[key].push(rec);
      return acc;
    },
    {} as Record<string, PlannerRecommendationResponseDTO[]>,
  );

  const selectedISO = format(selectedDate, "yyyy-MM-dd");
  const selectedRecs = recommendationsByDate[selectedISO] || [];

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
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Planner</h1>
            <p className="text-muted-foreground mt-2">
              Personalized training recommendations
            </p>
          </div>
          <Button onClick={handleGenerate} disabled={isGenerating}>
            {isGenerating ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Sparkles className="mr-2 h-4 w-4" />
            )}
            Generate
          </Button>
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading recommendations...</p>
          </div>
        ) : (
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg flex items-center gap-2">
                  <Calendar className="h-5 w-5" />
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
              {recommendations.length === 0 ? (
                <div className="pt-6 text-center py-12">
                  <Calendar className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
                  <p className="text-muted-foreground">
                    No recommendations yet. Generate your first recommendation!
                  </p>
                </div>
              ) : (
                <>
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
                      const dayRecs = recommendationsByDate[iso] || [];
                      const inMonth = isSameMonth(d, monthCursor);
                      const isSelected = isSameDay(d, selectedDate);
                      const hasRec = dayRecs.length > 0;

                      return (
                        <button
                          key={iso}
                          type="button"
                          onClick={() => {
                            setSelectedDate(d);
                            if (hasRec) setIsNoteOpen(true);
                          }}
                          className={cn(
                            "rounded-md border border-border text-left p-2 min-h-20 overflow-hidden",
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
                            {hasRec && (
                              <div className="text-xs text-muted-foreground">
                                {dayRecs.length}×
                              </div>
                            )}
                          </div>

                          {hasRec ? (
                            <div className="mt-2 space-y-1 overflow-hidden">
                              <div className="flex items-center gap-2 text-xs text-muted-foreground min-w-0">
                                <FileText className="h-3 w-3" />
                                <span className="truncate">Recommendation</span>
                              </div>
                              <div className="text-xs text-muted-foreground truncate whitespace-nowrap">
                                Tap to view
                              </div>
                            </div>
                          ) : (
                            <div className="mt-2 text-xs text-muted-foreground truncate whitespace-nowrap">
                              No recommendation
                            </div>
                          )}
                        </button>
                      );
                    })}
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        )}

        <Dialog open={isNoteOpen} onOpenChange={setIsNoteOpen}>
          <DialogContent className="sm:max-w-2xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5" />
                AI Recommendation
              </DialogTitle>
              <DialogDescription>
                {format(selectedDate, "PPP")}
              </DialogDescription>
            </DialogHeader>

            <div className="rounded-md border border-dashed border-border bg-muted/30 p-4">
              {selectedRecs.length === 0 ? (
                <p className="text-sm text-muted-foreground">
                  No recommendation for this day.
                </p>
              ) : (
                <div className="space-y-4">
                  {selectedRecs.map((rec) => (
                    <div key={rec.id} className="space-y-2">
                      <div className="text-xs text-muted-foreground">
                        {new Date(rec.created_at).toLocaleString()}
                      </div>
                      <div className="text-sm text-foreground whitespace-pre-wrap leading-relaxed">
                        {rec.recommendation}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
