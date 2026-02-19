"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { HeroSection } from "@/components/dashboard/hero-section";
import { QuickTiles } from "@/components/dashboard/quick-tiles";
import { Card } from "@/components/ui/card";
import { api } from "@/lib/api";
import { getAuthState } from "@/lib/auth";
import type {
  DailyNutritionResponseDTO,
  PlannerRecommendationResponseDTO,
  WorkoutSessionResponseDTO,
} from "@/lib/backend-dto";

function toISODate(date: Date): string {
  return date.toISOString().slice(0, 10);
}

function startOfWeekMonday(date: Date): Date {
  const d = new Date(date);
  const day = d.getDay(); // 0=Sun ... 6=Sat
  const diff = day === 0 ? -6 : 1 - day;
  d.setDate(d.getDate() + diff);
  d.setHours(0, 0, 0, 0);
  return d;
}

export default function DashboardPage() {
  const contentRef = useRef<HTMLDivElement>(null);

  const [workoutsThisWeek, setWorkoutsThisWeek] = useState(0);
  const [proteinToday, setProteinToday] = useState(0);
  const [latestRecommendation, setLatestRecommendation] = useState("");
  const [apiConnected, setApiConnected] = useState<boolean | null>(null);
  const [authActive, setAuthActive] = useState<boolean | null>(null);

  const auth = useMemo(() => getAuthState(), []);

  const handleHeroScroll = () => {
    contentRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    const load = async () => {
      if (!auth?.userId) {
        setAuthActive(false);
        setApiConnected(null);
        setWorkoutsThisWeek(0);
        setProteinToday(0);
        setLatestRecommendation("");
        return;
      }

      setAuthActive(true);

      const today = new Date();
      const todayISO = toISODate(today);
      const weekStart = startOfWeekMonday(today);
      const weekStartMs = weekStart.getTime();

      const [workoutsRes, nutritionRes, plannerRes] = await Promise.allSettled([
        api.get<WorkoutSessionResponseDTO[]>(
          `/api/v1/workouts/user/${auth.userId}`,
        ),
        api.get<DailyNutritionResponseDTO>(
          `/api/v1/nutrition/user/${auth.userId}?date=${todayISO}`,
        ),
        api.get<PlannerRecommendationResponseDTO[]>(
          `/api/v1/planner/user/${auth.userId}`,
        ),
      ]);

      const anyOk =
        workoutsRes.status === "fulfilled" ||
        nutritionRes.status === "fulfilled" ||
        plannerRes.status === "fulfilled";
      setApiConnected(anyOk);

      const maybeAuthError = (err: unknown) => {
        const msg = err instanceof Error ? err.message : "";
        return msg.includes("401") || msg.includes("403");
      };

      if (
        (workoutsRes.status === "rejected" &&
          maybeAuthError(workoutsRes.reason)) ||
        (nutritionRes.status === "rejected" &&
          maybeAuthError(nutritionRes.reason)) ||
        (plannerRes.status === "rejected" && maybeAuthError(plannerRes.reason))
      ) {
        setAuthActive(false);
      }

      if (workoutsRes.status === "fulfilled") {
        const count = (workoutsRes.value || []).filter((w) => {
          const ms = new Date(w.session_date).getTime();
          return !Number.isNaN(ms) && ms >= weekStartMs;
        }).length;
        setWorkoutsThisWeek(count);
      } else {
        setWorkoutsThisWeek(0);
      }

      if (nutritionRes.status === "fulfilled") {
        setProteinToday(nutritionRes.value?.protein_grams ?? 0);
      } else {
        setProteinToday(0);
      }

      if (plannerRes.status === "fulfilled") {
        const recs = plannerRes.value || [];
        const latest = [...recs].sort(
          (a, b) =>
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
        )[0];
        setLatestRecommendation(latest?.recommendation ?? "");
      } else {
        setLatestRecommendation("");
      }
    };

    void load();
  }, [auth]);

  return (
    <div className="w-full">
      <HeroSection onScroll={handleHeroScroll} />

      <div
        ref={contentRef}
        className="min-h-screen bg-background p-6 md:p-8 scroll-mt-0"
      >
        <div className="max-w-7xl mx-auto space-y-12">
          <div>
            <h2 className="text-3xl font-bold text-foreground mb-2">
              Dashboard
            </h2>
            <p className="text-muted-foreground">
              Track your progress and stay motivated
            </p>
          </div>

          <QuickTiles
            workoutsThisWeek={workoutsThisWeek}
            proteinToday={proteinToday}
            latestRecommendation={latestRecommendation}
          />

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-stretch">
            <Card className="p-6 h-full flex flex-col">
              <h3 className="text-lg font-bold text-foreground mb-4">
                Quick Start
              </h3>
              <ul className="space-y-3 text-sm text-muted-foreground">
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold">→</span>
                  <span>Create your first workout in the Workouts section</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold">→</span>
                  <span>Set up your workout split for better organization</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold">→</span>
                  <span>Track your nutrition to optimize performance</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="text-primary font-bold">→</span>
                  <span>Use AI tools for personalized recommendations</span>
                </li>
              </ul>
            </Card>

            <Card className="p-6 h-full flex flex-col">
              <h3 className="text-lg font-bold text-foreground mb-4">
                System Status
              </h3>
              <div className="space-y-3 text-sm">
                <div className="flex items-center justify-between pb-3 border-b border-border/40">
                  <span className="text-muted-foreground">API Connection</span>
                  <span
                    className={`font-medium ${apiConnected === null ? "text-muted-foreground" : apiConnected === false ? "text-destructive" : "text-green-600"}`}
                  >
                    {apiConnected === null
                      ? "Checking…"
                      : apiConnected === false
                        ? "Disconnected"
                        : "Connected"}
                  </span>
                </div>
                <div className="flex items-center justify-between pb-3 border-b border-border/40">
                  <span className="text-muted-foreground">Authentication</span>
                  <span
                    className={`font-medium ${authActive === null ? "text-muted-foreground" : authActive === false ? "text-destructive" : "text-green-600"}`}
                  >
                    {authActive === null
                      ? "Checking…"
                      : authActive === false
                        ? "Inactive"
                        : "Active"}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Data Sync</span>
                  <span
                    className={`font-medium ${apiConnected === null ? "text-muted-foreground" : apiConnected === false ? "text-destructive" : "text-green-600"}`}
                  >
                    {apiConnected === null
                      ? "Checking…"
                      : apiConnected === false
                        ? "Not ready"
                        : "Ready"}
                  </span>
                </div>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
