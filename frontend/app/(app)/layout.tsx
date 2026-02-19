"use client";

import { AppHeader } from "@/components/app-header";
import { AppNav } from "@/components/app-nav";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { isAuthenticated } from "@/lib/auth";
import { SpartanHelmetIcon } from "@/components/spartan-helmet-icon";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [sidebarOpen, setSidebarOpen] = useState(() => {
    try {
      const saved = localStorage.getItem("sidebarOpen");
      if (saved === "true") return true;
      if (saved === "false") return false;
    } catch {
      // ignore
    }
    return true;
  });

  const authed = isAuthenticated();

  useEffect(() => {
    if (!authed) {
      router.push("/login");
    }
  }, [authed, router]);

  if (!authed) {
    return null;
  }

  const toggleSidebar = () => {
    setSidebarOpen((prev) => {
      const next = !prev;
      try {
        localStorage.setItem("sidebarOpen", String(next));
      } catch {
        // ignore
      }
      return next;
    });
  };

  return (
    <div className="flex flex-col h-screen md:flex-row">
      {/* Desktop Sidebar */}
      {sidebarOpen && (
        <aside className="hidden md:flex md:flex-col md:w-64 border-r border-border/40 bg-background">
          <div className="p-6">
            <div className="flex items-center gap-3">
              <div className="size-10 rounded-lg bg-muted flex items-center justify-center">
                <SpartanHelmetIcon className="h-6 w-6" />
              </div>
              <h1 className="text-2xl font-bold text-foreground">
                S.P.A.R.T.A
              </h1>
            </div>
          </div>
          <nav className="flex-1 overflow-y-auto px-4 pb-4">
            <AppNav variant="sidebar" />
          </nav>
        </aside>
      )}

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        <AppHeader sidebarOpen={sidebarOpen} onToggleSidebar={toggleSidebar} />
        <main className="flex-1 overflow-y-auto">{children}</main>

        {/* Mobile Bottom Nav */}
        <nav className="md:hidden border-t border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky bottom-0">
          <div className="px-4 py-3 overflow-x-auto">
            <AppNav variant="mobile" />
          </div>
        </nav>
      </div>
    </div>
  );
}
