"use client";

import { useRouter } from "next/navigation";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { Moon, Sun, LogOut, Smartphone, Menu } from "lucide-react";
import { clearAuth } from "@/lib/auth";
import { toast } from "sonner";
import { SpartanHelmetIcon } from "@/components/spartan-helmet-icon";

interface AppHeaderProps {
  sidebarOpen?: boolean;
  onToggleSidebar?: () => void;
}

export function AppHeader({
  sidebarOpen = true,
  onToggleSidebar,
}: AppHeaderProps) {
  const router = useRouter();
  const { theme, setTheme } = useTheme();

  const handleLogout = () => {
    clearAuth();
    toast.success("Logged out successfully");
    router.push("/login");
  };

  return (
    <header className="border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-16 items-center justify-between px-4 md:px-6">
        <div className="flex items-center gap-2">
          {onToggleSidebar && (
            <Button
              type="button"
              variant="ghost"
              size="icon"
              onClick={onToggleSidebar}
              aria-label={sidebarOpen ? "Hide sidebar" : "Show sidebar"}
              className="hidden md:inline-flex"
            >
              <Menu className="h-5 w-5" />
            </Button>
          )}
          <div className="flex items-center gap-2 text-lg font-bold text-foreground">
            <div className="size-8 rounded-md bg-muted flex items-center justify-center">
              <SpartanHelmetIcon className="h-5 w-5" />
            </div>
            <span>S.P.A.R.T.A</span>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                {theme === "dark" ? (
                  <Moon className="h-5 w-5" />
                ) : (
                  <Sun className="h-5 w-5" />
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => setTheme("light")}>
                Light
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setTheme("dark")}>
                Dark
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setTheme("system")}>
                System
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <Smartphone className="h-5 w-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <a href="/app/settings">Settings</a>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={handleLogout}>
                <LogOut className="mr-2 h-4 w-4" />
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
}
