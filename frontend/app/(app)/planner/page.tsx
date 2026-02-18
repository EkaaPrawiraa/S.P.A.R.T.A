"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { getUserId } from "@/lib/auth";
import { Calendar, Loader2 } from "lucide-react";
import { toast } from "sonner";
import type { PlannerRecommendationResponseDTO } from "@/lib/backend-dto";

export default function PlannerPage() {
  const [recommendations, setRecommendations] = useState<
    PlannerRecommendationResponseDTO[]
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isGenerating, setIsGenerating] = useState(false);

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
            {isGenerating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Generate
          </Button>
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading recommendations...</p>
          </div>
        ) : recommendations.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Calendar className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">
                No recommendations yet. Generate your first recommendation!
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-4">
            {recommendations.map((rec) => (
              <Card key={rec.id}>
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <CardTitle className="text-base">
                      {new Date(rec.created_at).toLocaleDateString()}
                    </CardTitle>
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-foreground whitespace-pre-wrap">
                    {rec.recommendation}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
