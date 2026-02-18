"use client";

import { motion } from "framer-motion";
import { ChevronDown } from "lucide-react";

import { useEffect, useState } from "react";

import { api } from "@/lib/api";
import type { DailyMotivationResponseDTO } from "@/lib/backend-dto";

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

  useEffect(() => {
    let cancelled = false;

    async function loadMotivation() {
      try {
        setIsLoading(true);
        const data = await api.get<DailyMotivationResponseDTO>(
          "/api/v1/ai/motivation",
        );
        if (!cancelled) setQuote(data.message);
      } catch {
        if (!cancelled) setQuote("");
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    loadMotivation();
    return () => {
      cancelled = true;
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
