"use client";

import { motion } from "framer-motion";
import { ChevronDown } from "lucide-react";

import { useEffect, useRef, useState } from "react";

import { api } from "@/lib/api";
import type { DailyMotivationResponseDTO } from "@/lib/backend-dto";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

interface HeroSectionProps {
  onScroll?: () => void;
}

export function HeroSection({ onScroll }: HeroSectionProps) {
  const today = new Date().toLocaleDateString("en-US", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  const [quote, setQuote] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const cancelledRef = useRef(false);

  const refreshMotivation = async (opts?: { silent?: boolean }) => {
    try {
      if (!opts?.silent && !cancelledRef.current) setIsLoading(true);
      const data = await api.get<DailyMotivationResponseDTO>(
        "/api/v1/ai/motivation",
      );
      if (!cancelledRef.current) setQuote(data.message);
    } catch {
      if (!cancelledRef.current) setQuote("");
    } finally {
      if (!opts?.silent && !cancelledRef.current) setIsLoading(false);
    }
  };

  const handleResetQuote = async () => {
    try {
      if (!cancelledRef.current) setIsLoading(true);
      await api.post("/api/v1/ai/motivation/reset");
      await refreshMotivation({ silent: true });
      toast.success("Motivation quote reset");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to reset quote";
      toast.error(message);
    } finally {
      if (!cancelledRef.current) setIsLoading(false);
    }
  };

  useEffect(() => {
    cancelledRef.current = false;
    refreshMotivation().catch(() => {
      // handled inside refreshMotivation
    });
    return () => {
      cancelledRef.current = true;
    };
  }, []);

  return (
    <div className="relative h-screen flex items-center justify-center bg-gradient-to-br from-background via-background to-muted/20 overflow-hidden">
      {/* Subtle background pattern */}
      <div className="absolute inset-0 opacity-5">
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `linear-gradient(45deg, transparent 25%, rgba(255,255,255,.05) 25%, rgba(255,255,255,.05) 50%, transparent 50%, transparent 75%, rgba(255,255,255,.05) 75%, rgba(255,255,255,.05))`,
            backgroundSize: "60px 60px",
          }}
        />
      </div>

      <motion.div
        className="relative text-center px-4 max-w-2xl"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8 }}
      >
        <div className="mb-6">
          <div className="text-6xl font-bold text-foreground mb-4 leading-tight text-balance">
            φ
          </div>
          <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6 leading-tight text-balance">
            {isLoading
              ? "Loading your daily motivation…"
              : quote || "Motivation unavailable"}
          </h1>
        </div>

        <p className="text-lg text-muted-foreground mb-8">{today}</p>

        <div className="mb-8">
          <Button
            type="button"
            variant="outline"
            onClick={handleResetQuote}
            disabled={isLoading}
          >
            Motivate Me
          </Button>
        </div>

        <p className="text-sm text-muted-foreground uppercase tracking-wider">
          Scroll to Dashboard
        </p>
      </motion.div>

      {/* Scroll indicator */}
      <motion.div
        className="absolute bottom-8 left-1/2 transform -translate-x-1/2"
        animate={{ y: [0, 10, 0] }}
        transition={{ duration: 2, repeat: Infinity }}
      >
        <ChevronDown className="h-6 w-6 text-muted-foreground/60" />
      </motion.div>
    </div>
  );
}
