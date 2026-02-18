"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { ArrowLeft, Link as LinkIcon } from "lucide-react";
import { toast } from "sonner";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { ExerciseResponseDTO } from "@/lib/backend-dto";

export default function ExerciseDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [exercise, setExercise] = useState<ExerciseResponseDTO | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [mediaType, setMediaType] = useState<"image" | "video">("video");
  const [mediaUrl, setMediaUrl] = useState("");
  const [thumbnailUrl, setThumbnailUrl] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const exerciseId = Array.isArray(params.id) ? params.id[0] : params.id;

  useEffect(() => {
    loadExercise();
  }, [exerciseId]);

  const loadExercise = async () => {
    if (!exerciseId || typeof exerciseId !== "string") {
      toast.error("Invalid exercise ID");
      router.push("/app/exercises");
      return;
    }

    try {
      const data = await api.get<ExerciseResponseDTO>(
        `/api/v1/exercises/${exerciseId}`,
      );
      setExercise(data);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to load exercise";
      toast.error(message);
      router.push("/app/exercises");
    } finally {
      setIsLoading(false);
    }
  };

  const handleAttachMedia = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!mediaUrl) {
      toast.error("Please enter a media URL");
      return;
    }

    if (!exerciseId || typeof exerciseId !== "string") {
      toast.error("Invalid exercise ID");
      return;
    }

    setIsSubmitting(true);
    try {
      await api.post(`/api/v1/exercises/${exerciseId}/media`, {
        media_type: mediaType,
        media_url: mediaUrl,
        thumbnail_url: thumbnailUrl ? thumbnailUrl : undefined,
      });
      toast.success("Media attached successfully");
      setMediaUrl("");
      setThumbnailUrl("");
      loadExercise();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to attach media";
      toast.error(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background p-6 md:p-8">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (!exercise) {
    return (
      <div className="min-h-screen bg-background p-6 md:p-8">
        <p className="text-muted-foreground">Exercise not found</p>
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
            <h1 className="text-3xl font-bold text-foreground">
              {exercise.name}
            </h1>
            <p className="text-muted-foreground">
              Primary: {exercise.primary_muscle}
            </p>
          </div>
        </div>

        {exercise.media && exercise.media.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle>Media</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-2">
                {exercise.media.map((m) => (
                  <li
                    key={m.id}
                    className="flex items-center gap-2 text-sm text-muted-foreground"
                  >
                    <LinkIcon className="h-4 w-4" />
                    <a
                      href={m.media_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:text-primary"
                    >
                      {m.media_url}
                    </a>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        )}

        <Card>
          <CardHeader>
            <CardTitle>Attach Media</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleAttachMedia} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="mediaType">Media Type</Label>
                <Select
                  value={mediaType}
                  onValueChange={(v) => setMediaType(v as "image" | "video")}
                >
                  <SelectTrigger id="mediaType">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="image">Image</SelectItem>
                    <SelectItem value="video">Video</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="media">Media URL</Label>
                <Input
                  id="media"
                  type="url"
                  placeholder="https://example.com/video.mp4"
                  value={mediaUrl}
                  onChange={(e) => setMediaUrl(e.target.value)}
                  disabled={isSubmitting}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="thumb">Thumbnail URL (optional)</Label>
                <Input
                  id="thumb"
                  type="url"
                  placeholder="https://example.com/thumb.jpg"
                  value={thumbnailUrl}
                  onChange={(e) => setThumbnailUrl(e.target.value)}
                  disabled={isSubmitting}
                />
              </div>
              <Button type="submit" disabled={isSubmitting}>
                Attach Media
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
