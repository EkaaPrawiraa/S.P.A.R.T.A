"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { Zap } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
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
import { AspectRatio } from "@/components/ui/aspect-ratio";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

function isLikelyImageUrl(url: string): boolean {
  const u = url.toLowerCase();
  return (
    u.endsWith(".png") ||
    u.endsWith(".jpg") ||
    u.endsWith(".jpeg") ||
    u.endsWith(".gif") ||
    u.endsWith(".webp") ||
    u.endsWith(".svg")
  );
}

function isLikelyVideoUrl(url: string): boolean {
  const u = url.toLowerCase();
  return u.endsWith(".mp4") || u.endsWith(".webm") || u.endsWith(".ogg");
}

export default function ExercisesPage() {
  const router = useRouter();
  const [exercises, setExercises] = useState<ExerciseResponseDTO[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [exerciseName, setExerciseName] = useState("");
  const [primaryMuscle, setPrimaryMuscle] = useState("");
  const [primaryMuscleFilter, setPrimaryMuscleFilter] = useState<string>("");
  const [illustrationOnly, setIllustrationOnly] = useState(false);
  const [sortMode, setSortMode] = useState<"az" | "date">("az");
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

  const filteredExercises = (exercises || []).filter((ex) => {
    if (illustrationOnly && (ex.media?.length ?? 0) === 0) return false;
    if (!primaryMuscleFilter) return true;

    const norm = (v: string | null | undefined) =>
      (v || "").trim().toLowerCase();
    return norm(ex.primary_muscle) === norm(primaryMuscleFilter);
  });

  const sortedExercises = [...filteredExercises].sort((a, b) => {
    if (sortMode === "date") {
      const ta = Date.parse(a.created_at || "") || 0;
      const tb = Date.parse(b.created_at || "") || 0;
      if (tb !== ta) return tb - ta;
    }
    return (a.name || "").localeCompare(b.name || "", undefined, {
      sensitivity: "base",
    });
  });

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

        <div className="flex flex-col gap-4 md:flex-row md:items-end">
          <div className="space-y-2 md:flex-1">
            <Label htmlFor="primaryMuscleFilter">Primary Muscle</Label>
            <Select
              value={primaryMuscleFilter}
              onValueChange={setPrimaryMuscleFilter}
            >
              <SelectTrigger id="primaryMuscleFilter">
                <SelectValue placeholder="All" />
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

          <div className="space-y-2 md:w-56">
            <Label htmlFor="sortMode">Sort</Label>
            <Select
              value={sortMode}
              onValueChange={(v) => {
                if (v === "az" || v === "date") setSortMode(v);
              }}
            >
              <SelectTrigger id="sortMode">
                <SelectValue placeholder="Sort" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="az">A–Z</SelectItem>
                <SelectItem value="date">Date added</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label></Label>
            <label className="flex items-center gap-2 rounded-md border border-border px-3 py-2 text-sm">
              <Checkbox
                checked={illustrationOnly}
                onCheckedChange={(v) => setIllustrationOnly(v === true)}
              />
              <span className="text-muted-foreground">Media only</span>
            </label>
          </div>

          <div className="md:flex md:justify-end">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setPrimaryMuscleFilter("");
                setIllustrationOnly(false);
              }}
              disabled={!primaryMuscleFilter && !illustrationOnly}
              className="w-full md:w-auto"
            >
              Reset
            </Button>
          </div>
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
        ) : filteredExercises.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Zap className="h-12 w-12 text-muted-foreground/40 mx-auto mb-4" />
              <p className="text-muted-foreground mb-4">No exercises found.</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 items-stretch">
            {sortedExercises.map((exercise) => (
              <Card
                key={exercise.id}
                className="h-full flex flex-col hover:shadow-lg transition-shadow cursor-pointer"
                role="button"
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault();
                    router.push(`/app/exercises/${exercise.id}`);
                  }
                }}
                onClick={(e) => {
                  const target = e.target as HTMLElement | null;
                  if (target?.closest('[data-slot="carousel"]')) return;
                  if (target?.closest("video")) return;
                  router.push(`/app/exercises/${exercise.id}`);
                }}
              >
                <CardHeader className="space-y-3">
                  <div className="overflow-hidden rounded-md border border-border">
                    <AspectRatio ratio={16 / 9}>
                      {exercise.media && exercise.media.length > 0 ? (
                        <Carousel
                          className="h-full w-full"
                          opts={{ loop: false }}
                        >
                          <CarouselContent className="h-full">
                            {exercise.media.map((m) => {
                              const canPreviewImage =
                                m.media_type === "image" ||
                                isLikelyImageUrl(m.media_url);
                              const canPreviewVideo =
                                m.media_type === "video" ||
                                isLikelyVideoUrl(m.media_url);

                              return (
                                <CarouselItem
                                  key={m.id}
                                  className="h-full w-full"
                                >
                                  {canPreviewImage ? (
                                    // eslint-disable-next-line @next/next/no-img-element
                                    <img
                                      src={m.media_url}
                                      alt={`${exercise.name} media`}
                                      className="h-full w-full object-cover"
                                      loading="lazy"
                                    />
                                  ) : canPreviewVideo ? (
                                    <video
                                      className="h-full w-full"
                                      controls
                                      preload="metadata"
                                      poster={m.thumbnail_url ?? undefined}
                                    >
                                      <source src={m.media_url} />
                                    </video>
                                  ) : (
                                    <div className="h-full w-full bg-muted" />
                                  )}
                                </CarouselItem>
                              );
                            })}
                          </CarouselContent>

                          {exercise.media.length > 1 && (
                            <>
                              <CarouselPrevious className="left-2 top-1/2 -translate-y-1/2" />
                              <CarouselNext className="right-2 top-1/2 -translate-y-1/2" />
                            </>
                          )}
                        </Carousel>
                      ) : (
                        <div className="h-full w-full bg-muted flex items-center justify-center">
                          <span className="text-xs text-muted-foreground">
                            No media
                          </span>
                        </div>
                      )}
                    </AspectRatio>
                  </div>

                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0 min-h-[2.75rem]">
                      <Link
                        href={`/app/exercises/${exercise.id}`}
                        className="hover:underline"
                      >
                        <CardTitle className="text-lg leading-snug line-clamp-2">
                          {exercise.name}
                        </CardTitle>
                      </Link>
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="space-y-2 mt-auto">
                  <div className="text-sm text-muted-foreground">
                    Primary: {exercise.primary_muscle}
                  </div>

                  <div className="text-xs text-muted-foreground">
                    Media: {exercise.media?.length ?? 0}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
