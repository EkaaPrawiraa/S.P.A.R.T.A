"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { api } from "@/lib/api";
import {
  Loader2,
  Brain,
  Sparkles,
  FileText,
  Wand2,
  TrendingUp,
  GitBranch,
  BookmarkPlus,
  Info,
} from "lucide-react";
import { toast } from "sonner";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type {
  CoachingSuggestionsResponseDTO,
  ExerciseResponseDTO,
  PlannerRecommendationResponseDTO,
  SplitTemplateResponseDTO,
  WorkoutExplanationResponseDTO,
  WorkoutPlanResponseDTO,
} from "@/lib/backend-dto";
import { FOCUS_MUSCLE_OPTIONS } from "@/lib/grammar";
import { getUserId } from "@/lib/auth";
import { ExerciseCombobox } from "@/components/exercise-combobox";

export default function AIToolsPage() {
  const [coachingTips, setCoachingTips] = useState<string[]>([]);

  const [explainSplitDayName, setExplainSplitDayName] = useState("");
  const [explainFatigue, setExplainFatigue] = useState("5");
  const [explainExerciseName, setExplainExerciseName] = useState("");
  const [explainExerciseSets, setExplainExerciseSets] = useState("3");
  const [explainExerciseRepRange, setExplainExerciseRepRange] =
    useState("8-12");
  const [explainExerciseWeight, setExplainExerciseWeight] = useState("0");
  const [explainResponse, setExplainResponse] =
    useState<WorkoutExplanationResponseDTO | null>(null);

  const [planSplitDayId, setPlanSplitDayId] = useState("");
  const [planFatigue, setPlanFatigue] = useState("5");
  const [planResponse, setPlanResponse] =
    useState<WorkoutPlanResponseDTO | null>(null);

  const [userSplitDays, setUserSplitDays] = useState<
    { id: string; label: string }[]
  >([]);

  const [overloadExerciseId, setOverloadExerciseId] = useState("");
  const [overloadResponse, setOverloadResponse] =
    useState<PlannerRecommendationResponseDTO | null>(null);

  const [splitDaysPerWeek, setSplitDaysPerWeek] = useState("4");
  const [splitFocusMuscle, setSplitFocusMuscle] = useState("");
  const [splitTemplate, setSplitTemplate] =
    useState<SplitTemplateResponseDTO | null>(null);

  const [exercises, setExercises] = useState<ExerciseResponseDTO[]>([]);

  const [coachingLoading, setCoachingLoading] = useState(false);
  const [explainLoading, setExplainLoading] = useState(false);
  const [planLoading, setPlanLoading] = useState(false);
  const [overloadLoading, setOverloadLoading] = useState(false);
  const [splitLoading, setSplitLoading] = useState(false);

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

  useEffect(() => {
    const loadUserSplitDays = async () => {
      try {
        const userId = getUserId();
        if (!userId) {
          setUserSplitDays([]);
          return;
        }

        const splits = await api.get<SplitTemplateResponseDTO[]>(
          `/api/v1/splits/user/${userId}`,
        );

        const days = (splits || []).flatMap((s) =>
          (s.days || []).map((d) => ({
            id: d.id,
            label: `${s.name} — Day ${d.day_order}: ${d.name}`,
          })),
        );
        setUserSplitDays(days);

        if (days.length > 0 && !days.some((d) => d.id === planSplitDayId)) {
          setPlanSplitDayId(days[0].id);
        }
      } catch {
        setUserSplitDays([]);
      }
    };

    loadUserSplitDays();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const exerciseNameById = (exercises || []).reduce(
    (acc, ex) => {
      acc[ex.id] = ex.name;
      return acc;
    },
    {} as Record<string, string>,
  );

  const normalizeExerciseName = (name: string) =>
    name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, " ")
      .trim();

  const exerciseIdByNormalizedName = (exercises || []).reduce(
    (acc, ex) => {
      const key = normalizeExerciseName(ex.name);
      if (key && !acc[key]) acc[key] = ex.id;
      return acc;
    },
    {} as Record<string, string>,
  );

  const resolveSplitExercise = (ex: {
    exercise_id: string;
    notes?: string;
  }) => {
    const rawId = (ex.exercise_id || "").trim();
    if (rawId) {
      return {
        displayName: exerciseNameById[rawId] || rawId,
        linkedExerciseId: rawId,
        linked: true,
      };
    }

    const suggestedName = (ex.notes || "").trim();
    const matchedId = suggestedName
      ? exerciseIdByNormalizedName[normalizeExerciseName(suggestedName)]
      : "";

    return {
      displayName: suggestedName || "Exercise",
      linkedExerciseId: matchedId || "",
      linked: Boolean(matchedId),
    };
  };

  const loadCoaching = async () => {
    setCoachingLoading(true);
    try {
      const data = await api.get<CoachingSuggestionsResponseDTO>(
        "/api/v1/ai/coaching",
      );
      setCoachingTips(data.suggestions || []);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load coaching";
      toast.error(message);
    } finally {
      setCoachingLoading(false);
    }
  };

  const handleExplain = async () => {
    if (!explainExerciseName || !explainExerciseRepRange) {
      toast.error("Please enter exercise details");
      return;
    }

    setExplainLoading(true);
    try {
      const data = await api.post<WorkoutExplanationResponseDTO>(
        "/api/v1/ai/explain-workout",
        {
          split_day_name: explainSplitDayName,
          fatigue: parseInt(explainFatigue) || 0,
          exercises: [
            {
              name: explainExerciseName,
              sets: parseInt(explainExerciseSets) || 1,
              rep_range: explainExerciseRepRange,
              weight: parseFloat(explainExerciseWeight) || 0,
            },
          ],
        },
      );
      setExplainResponse(data);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to explain workout";
      toast.error(message);
    } finally {
      setExplainLoading(false);
    }
  };

  const handleGeneratePlan = async () => {
    if (!planSplitDayId) {
      toast.error("Please select a split day");
      return;
    }

    setPlanLoading(true);
    try {
      const data = await api.post<WorkoutPlanResponseDTO>(
        "/api/v1/ai/workout",
        {
          split_day_id: planSplitDayId,
          fatigue: parseInt(planFatigue) || 0,
        },
      );
      setPlanResponse(data);
    } catch (err) {
      let message =
        err instanceof Error ? err.message : "Failed to generate plan";

      if (message.toLowerCase().includes("openai api key not configured")) {
        message =
          "OpenAI API key not configured. Set OPENAI_API_KEY in backend/.env and restart the backend.";
      }

      toast.error(message);
    } finally {
      setPlanLoading(false);
    }
  };

  const handleOverload = async () => {
    setOverloadLoading(true);
    try {
      const data = await api.post<PlannerRecommendationResponseDTO>(
        "/api/v1/ai/overload",
        {
          exercise_id: overloadExerciseId,
        },
      );
      setOverloadResponse(data);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to get overload strategy";
      toast.error(message);
    } finally {
      setOverloadLoading(false);
    }
  };

  const handleGenerateSplit = async () => {
    if (!splitFocusMuscle) {
      toast.error("Please select a focus muscle");
      return;
    }

    setSplitLoading(true);
    try {
      const data = await api.post<SplitTemplateResponseDTO>(
        "/api/v1/ai/generate-split",
        {
          days_per_week: parseInt(splitDaysPerWeek) || 1,
          focus_muscle: splitFocusMuscle,
        },
      );
      setSplitTemplate(data);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to generate split";
      toast.error(message);
    } finally {
      setSplitLoading(false);
    }
  };

  const canSaveSplitSuggestion = (() => {
    if (!splitTemplate) return false;
    const userId = getUserId();
    if (!userId) return false;

    const days = splitTemplate.days || [];
    if (days.length === 0) return false;

    for (const day of days) {
      const exs = day.exercises || [];
      if (exs.length === 0) return false;
      for (const ex of exs) {
        const resolved = resolveSplitExercise(ex);
        if (!resolved.linkedExerciseId) return false;
      }
    }
    return true;
  })();

  const handleSaveSplitSuggestion = async () => {
    if (!splitTemplate) return;

    const userId = getUserId();
    if (!userId) {
      toast.error("You must be logged in to save splits");
      return;
    }

    const days = splitTemplate.days || [];
    let missingLinks = 0;
    for (const day of days) {
      for (const ex of day.exercises || []) {
        const resolved = resolveSplitExercise(ex);
        if (!resolved.linkedExerciseId) missingLinks += 1;
      }
    }

    if (missingLinks > 0) {
      toast.error(
        `Can't save yet: ${missingLinks} exercises aren't linked to your library. Add/match those exercises first.`,
      );
      return;
    }

    try {
      await api.post<SplitTemplateResponseDTO>("/api/v1/splits", {
        user_id: userId,
        name: splitTemplate.name || "AI Split",
        description: splitTemplate.description || "AI Generated Split",
        created_by: "ai",
        focus_muscle: splitTemplate.focus_muscle || splitFocusMuscle,
        is_active: false,
        days: (splitTemplate.days || []).map((day) => ({
          day_order: day.day_order,
          name: (day.name || "").trim() || `Day ${day.day_order}`,
          exercises: (day.exercises || []).map((ex) => {
            const resolved = resolveSplitExercise(ex);
            return {
              exercise_id: resolved.linkedExerciseId,
              target_sets: ex.target_sets,
              target_reps: ex.target_reps,
              target_weight: ex.target_weight,
              notes: ex.notes,
            };
          }),
        })),
      });

      toast.success("Saved split to your templates");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to save split template";
      toast.error(message);
    }
  };

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-5xl mx-auto space-y-6">
        <Card className="overflow-hidden">
          <CardHeader className="space-y-2">
            <div className="flex items-start justify-between gap-4">
              <div className="min-w-0">
                <CardTitle className="text-2xl sm:text-3xl font-bold text-foreground flex items-center gap-3">
                  <div className="size-10 rounded-lg bg-muted flex items-center justify-center shrink-0">
                    <Brain className="h-5 w-5" />
                  </div>
                  <span className="truncate">AI Tools</span>
                </CardTitle>
                <p className="text-muted-foreground mt-2 text-sm sm:text-base">
                  Get personalized coaching, workout plans, overload strategies,
                  and split templates.
                </p>
              </div>

              <div className="hidden sm:flex flex-wrap items-center justify-end gap-2">
                <Badge variant="secondary">Coaching</Badge>
                <Badge variant="secondary">Explain</Badge>
                <Badge variant="secondary">Plan</Badge>
                <Badge variant="secondary">Overload</Badge>
                <Badge variant="secondary">Split</Badge>
              </div>
            </div>
          </CardHeader>
          <CardContent className="pt-0">
            <Alert className="bg-muted/30">
              <Info className="h-4 w-4" />
              <AlertTitle>Use AI as guidance</AlertTitle>
              <AlertDescription>
                Recommendations can be wrong or too aggressive. Adjust for your
                experience level, fatigue, and any injuries.
              </AlertDescription>
            </Alert>
          </CardContent>
        </Card>

        <Tabs defaultValue="coaching" className="w-full">
          <TabsList className="grid w-full grid-cols-2 sm:grid-cols-5 h-auto gap-2">
            <TabsTrigger value="coaching" className="py-2">
              Coaching
            </TabsTrigger>
            <TabsTrigger value="explain" className="py-2">
              Explain
            </TabsTrigger>
            <TabsTrigger value="plan" className="py-2">
              Plan
            </TabsTrigger>
            <TabsTrigger value="overload" className="py-2">
              Overload
            </TabsTrigger>
            <TabsTrigger value="split" className="py-2">
              Split
            </TabsTrigger>
          </TabsList>

          <TabsContent value="coaching" className="mt-4">
            <div className="grid gap-4 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Sparkles className="h-5 w-5" />
                    Daily Coaching Tips
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <p className="text-sm text-muted-foreground">
                    Quick, actionable reminders you can apply today.
                  </p>
                  <Button
                    onClick={loadCoaching}
                    disabled={coachingLoading}
                    className="w-full"
                  >
                    {!coachingLoading && <Sparkles className="mr-2 h-4 w-4" />}
                    {coachingLoading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Load Tips
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Results</CardTitle>
                </CardHeader>
                <CardContent>
                  {coachingTips.length === 0 ? (
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <Sparkles />
                        </EmptyMedia>
                        <EmptyTitle>No tips yet</EmptyTitle>
                        <EmptyDescription>
                          Tap “Load Tips” to generate your daily coaching.
                        </EmptyDescription>
                      </EmptyHeader>
                      <EmptyContent />
                    </Empty>
                  ) : (
                    <ScrollArea className="max-h-[420px] pr-2">
                      <div className="space-y-3">
                        {coachingTips.map((tip, idx) => (
                          <div
                            key={idx}
                            className="rounded-lg border bg-muted/30 p-4"
                          >
                            <p className="text-sm text-foreground leading-relaxed">
                              {tip}
                            </p>
                          </div>
                        ))}
                      </div>
                    </ScrollArea>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="explain" className="mt-4">
            <div className="grid gap-4 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FileText className="h-5 w-5" />
                    Explain Workout
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="splitDayName">
                      Split Day Name (optional)
                    </Label>
                    <Input
                      id="splitDayName"
                      placeholder="e.g., Push"
                      value={explainSplitDayName}
                      onChange={(e) => setExplainSplitDayName(e.target.value)}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-2">
                      <Label htmlFor="fatigue">Fatigue (0-10)</Label>
                      <Input
                        id="fatigue"
                        type="number"
                        value={explainFatigue}
                        onChange={(e) => setExplainFatigue(e.target.value)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="exSets">Sets</Label>
                      <Input
                        id="exSets"
                        type="number"
                        value={explainExerciseSets}
                        onChange={(e) => setExplainExerciseSets(e.target.value)}
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="exName">Exercise Name</Label>
                    <Input
                      id="exName"
                      placeholder="e.g., Bench Press"
                      value={explainExerciseName}
                      onChange={(e) => setExplainExerciseName(e.target.value)}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-2">
                      <Label htmlFor="repRange">Rep Range</Label>
                      <Input
                        id="repRange"
                        placeholder="8-12"
                        value={explainExerciseRepRange}
                        onChange={(e) =>
                          setExplainExerciseRepRange(e.target.value)
                        }
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="exWeight">Weight</Label>
                      <Input
                        id="exWeight"
                        type="number"
                        value={explainExerciseWeight}
                        onChange={(e) =>
                          setExplainExerciseWeight(e.target.value)
                        }
                      />
                    </div>
                  </div>

                  <Button
                    onClick={handleExplain}
                    disabled={explainLoading}
                    className="w-full"
                  >
                    {!explainLoading && <FileText className="mr-2 h-4 w-4" />}
                    {explainLoading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Explain
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Explanation</CardTitle>
                </CardHeader>
                <CardContent>
                  {!explainResponse ? (
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <FileText />
                        </EmptyMedia>
                        <EmptyTitle>Nothing to show yet</EmptyTitle>
                        <EmptyDescription>
                          Fill the fields and hit “Explain” to get feedback.
                        </EmptyDescription>
                      </EmptyHeader>
                      <EmptyContent />
                    </Empty>
                  ) : (
                    <ScrollArea className="max-h-[420px] pr-2">
                      <div className="rounded-lg border bg-muted/30 p-4">
                        <p className="text-sm text-foreground whitespace-pre-wrap leading-relaxed">
                          {explainResponse.summary}
                        </p>

                        {explainResponse.exercise_notes?.length > 0 && (
                          <>
                            <Separator className="my-4" />
                            <div className="space-y-2">
                              {explainResponse.exercise_notes.map((n) => (
                                <div
                                  key={n.name}
                                  className="text-sm text-muted-foreground"
                                >
                                  <span className="font-medium text-foreground">
                                    {n.name}:
                                  </span>{" "}
                                  {n.note}
                                </div>
                              ))}
                            </div>
                          </>
                        )}
                      </div>
                    </ScrollArea>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="plan" className="mt-4">
            <div className="grid gap-4 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Wand2 className="h-5 w-5" />
                    Generate Workout Plan
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="splitDayId">Split Day</Label>
                    <Select
                      value={planSplitDayId}
                      onValueChange={setPlanSplitDayId}
                    >
                      <SelectTrigger id="splitDayId">
                        <SelectValue placeholder="Select a split day" />
                      </SelectTrigger>
                      <SelectContent>
                        {userSplitDays.map((d) => (
                          <SelectItem key={d.id} value={d.id}>
                            {d.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="planFatigue">Fatigue (0-10)</Label>
                    <Input
                      id="planFatigue"
                      type="number"
                      value={planFatigue}
                      onChange={(e) => setPlanFatigue(e.target.value)}
                    />
                  </div>

                  <Button
                    onClick={handleGeneratePlan}
                    disabled={planLoading || !planSplitDayId}
                    className="w-full"
                  >
                    {!planLoading && <Wand2 className="mr-2 h-4 w-4" />}
                    {planLoading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Generate Plan
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Plan</CardTitle>
                </CardHeader>
                <CardContent>
                  {!planResponse ? (
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <Wand2 />
                        </EmptyMedia>
                        <EmptyTitle>No plan generated</EmptyTitle>
                        <EmptyDescription>
                          Select a split day and generate a plan.
                        </EmptyDescription>
                      </EmptyHeader>
                      <EmptyContent />
                    </Empty>
                  ) : (
                    <ScrollArea className="max-h-[420px] pr-2">
                      <div className="rounded-lg border bg-muted/30 p-4">
                        <div className="flex items-center justify-between gap-3">
                          <p className="text-sm text-foreground">
                            {planResponse.date}
                          </p>
                          <Badge variant="secondary">
                            {planResponse.exercises.length} exercises
                          </Badge>
                        </div>
                        <Separator className="my-4" />
                        <div className="space-y-2">
                          {planResponse.exercises.map((ex) => (
                            <div
                              key={ex.name}
                              className="text-sm text-muted-foreground"
                            >
                              <span className="font-medium text-foreground">
                                {ex.name}
                              </span>{" "}
                              — {ex.sets}x {ex.rep_range} @ {ex.weight}
                            </div>
                          ))}
                        </div>
                      </div>
                    </ScrollArea>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="overload" className="mt-4">
            <div className="grid gap-4 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <TrendingUp className="h-5 w-5" />
                    Progressive Overload
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="overloadExercise">Exercise</Label>
                    <ExerciseCombobox
                      id="overloadExercise"
                      value={overloadExerciseId}
                      onValueChange={setOverloadExerciseId}
                      options={exercises}
                      placeholder="Select an exercise"
                      searchPlaceholder="Search exercise..."
                    />
                  </div>

                  <Button
                    onClick={handleOverload}
                    disabled={overloadLoading}
                    className="w-full"
                  >
                    {!overloadLoading && (
                      <TrendingUp className="mr-2 h-4 w-4" />
                    )}
                    {overloadLoading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Get Strategy
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Strategy</CardTitle>
                </CardHeader>
                <CardContent>
                  {!overloadResponse ? (
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <TrendingUp />
                        </EmptyMedia>
                        <EmptyTitle>No strategy yet</EmptyTitle>
                        <EmptyDescription>
                          Pick an exercise and get progressive overload advice.
                        </EmptyDescription>
                      </EmptyHeader>
                      <EmptyContent />
                    </Empty>
                  ) : (
                    <ScrollArea className="max-h-[420px] pr-2">
                      <div className="rounded-lg border bg-muted/30 p-4">
                        <p className="text-sm text-foreground whitespace-pre-wrap leading-relaxed">
                          {overloadResponse.recommendation}
                        </p>
                      </div>
                    </ScrollArea>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="split" className="mt-4">
            <div className="grid gap-4 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <GitBranch className="h-5 w-5" />
                    Generate Split Template
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-2">
                      <Label htmlFor="daysPerWeek">Days per week</Label>
                      <Input
                        id="daysPerWeek"
                        type="number"
                        value={splitDaysPerWeek}
                        onChange={(e) => setSplitDaysPerWeek(e.target.value)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="focus">Focus Muscle</Label>
                      <Select
                        value={splitFocusMuscle}
                        onValueChange={setSplitFocusMuscle}
                      >
                        <SelectTrigger id="focus">
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
                  </div>

                  <Button
                    onClick={handleGenerateSplit}
                    disabled={splitLoading}
                    className="w-full"
                  >
                    {!splitLoading && <GitBranch className="mr-2 h-4 w-4" />}
                    {splitLoading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Generate Template
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Template</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {!splitTemplate ? (
                    <Empty>
                      <EmptyHeader>
                        <EmptyMedia variant="icon">
                          <GitBranch />
                        </EmptyMedia>
                        <EmptyTitle>No template generated</EmptyTitle>
                        <EmptyDescription>
                          Choose focus + days per week, then generate.
                        </EmptyDescription>
                      </EmptyHeader>
                      <EmptyContent />
                    </Empty>
                  ) : (
                    <>
                      <div className="space-y-1">
                        <p className="text-sm font-medium text-foreground">
                          {splitTemplate.name || "Split Template"}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          Focus: {splitTemplate.focus_muscle} •{" "}
                          {splitTemplate.days?.length || 0} days/week
                        </p>
                      </div>

                      <div>
                        <Button
                          type="button"
                          variant="secondary"
                          onClick={handleSaveSplitSuggestion}
                          disabled={!canSaveSplitSuggestion}
                        >
                          <BookmarkPlus className="mr-2 h-4 w-4" />
                          Save to My Splits
                        </Button>
                        {!canSaveSplitSuggestion && (
                          <p className="text-xs text-muted-foreground mt-2">
                            To save, the AI split must match exercises in your
                            library (so each exercise has a valid ID).
                          </p>
                        )}
                      </div>
                    </>
                  )}
                </CardContent>
              </Card>
            </div>

            {splitTemplate && (
              <div className="mt-4 space-y-4">
                {(splitTemplate.days || []).map((day) => (
                  <Card key={day.id || String(day.day_order)}>
                    <CardHeader>
                      <CardTitle className="text-base">
                        Day {day.day_order}: {day.name?.trim() || "(Unnamed)"}
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      {(day.exercises || []).length === 0 ? (
                        <p className="text-sm text-muted-foreground">
                          No exercises
                        </p>
                      ) : (
                        <div className="space-y-2">
                          {(day.exercises || []).map((ex, idx) => {
                            const resolved = resolveSplitExercise(ex);
                            return (
                              <div
                                key={`${ex.exercise_id || "unlinked"}-${day.day_order}-${idx}`}
                                className="flex items-center justify-between gap-4 text-sm"
                              >
                                <div className="font-medium text-foreground flex items-center gap-2 min-w-0">
                                  <span className="truncate">
                                    {resolved.displayName}
                                  </span>
                                  {resolved.linked ? (
                                    <Badge variant="secondary">Linked</Badge>
                                  ) : (
                                    <Badge variant="outline">Not linked</Badge>
                                  )}
                                </div>
                                <div className="text-muted-foreground whitespace-nowrap">
                                  {ex.target_sets}x{ex.target_reps}
                                  {ex.target_weight > 0
                                    ? ` @ ${ex.target_weight}`
                                    : ""}
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
