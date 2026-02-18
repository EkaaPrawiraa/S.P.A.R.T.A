"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { ArrowLeft } from "lucide-react";
import { toast } from "sonner";
import type { SplitTemplateResponseDTO } from "@/lib/backend-dto";

export default function SplitDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [split, setSplit] = useState<SplitTemplateResponseDTO | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const splitId = Array.isArray(params.id) ? params.id[0] : params.id;

  useEffect(() => {
    loadSplit();
  }, [splitId]);

  const loadSplit = async () => {
    if (!splitId || typeof splitId !== "string") {
      toast.error("Invalid split ID");
      router.push("/app/splits");
      return;
    }

    try {
      const data = await api.get<SplitTemplateResponseDTO>(
        `/api/v1/splits/${splitId}`,
      );
      setSplit(data);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load split";
      toast.error(message);
      router.push("/app/splits");
    } finally {
      setIsLoading(false);
    }
  };

  const handleActivate = async () => {
    if (!splitId || typeof splitId !== "string") {
      toast.error("Invalid split ID");
      return;
    }

    try {
      await api.post(`/api/v1/splits/${splitId}/activate`, {});
      toast.success("Split activated successfully");
      loadSplit();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to activate split";
      toast.error(message);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background p-6 md:p-8">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (!split) {
    return (
      <div className="min-h-screen bg-background p-6 md:p-8">
        <p className="text-muted-foreground">Split not found</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => router.back()}>
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-foreground">{split.name}</h1>
            <p className="text-muted-foreground">
              {split.days.length} days per week {split.is_active && "• Active"}
            </p>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Split Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Training Days</p>
                <p className="text-2xl font-bold text-foreground">
                  {split.days.length}
                </p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Status</p>
                <p className="text-2xl font-bold text-foreground">
                  {split.is_active ? "Active" : "Inactive"}
                </p>
              </div>
            </div>

            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Focus Muscle</p>
              <p className="text-base text-foreground">{split.focus_muscle}</p>
            </div>

            {split.description && split.description.trim().length > 0 && (
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Description</p>
                <p className="text-base text-foreground whitespace-pre-wrap">
                  {split.description}
                </p>
              </div>
            )}

            {!split.is_active && (
              <Button onClick={handleActivate} className="w-full">
                Activate This Split
              </Button>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
