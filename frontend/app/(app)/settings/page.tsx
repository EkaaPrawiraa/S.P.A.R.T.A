"use client";

import { useRouter } from "next/navigation";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Loader2, LogOut } from "lucide-react";
import { clearAuth, getAuthState } from "@/lib/auth";
import { toast } from "sonner";

interface Health {
  status: string;
  message?: string;
}

export default function SettingsPage() {
  const router = useRouter();
  const { theme, setTheme } = useTheme();
  const [token, setToken] = useState("");
  const [userId, setUserId] = useState("");
  const [isValidating, setIsValidating] = useState(false);
  const [health, setHealth] = useState<Health | null>(null);
  const [isCheckingHealth, setIsCheckingHealth] = useState(false);

  useEffect(() => {
    const auth = getAuthState();
    if (auth) {
      setToken(auth.token);
      setUserId(auth.userId);
    }
  }, []);

  const baseUrl =
    process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

  const fetchHealth = async (): Promise<Health> => {
    const res = await fetch(`${baseUrl}/api/v1/health`, { method: "GET" });
    const json = (await res.json()) as Health;
    if (!res.ok) {
      return {
        status: "error",
        message: json?.message || `HTTP ${res.status}`,
      };
    }
    return json;
  };

  const handleValidateToken = async () => {
    if (!token || !userId) {
      toast.error("Please enter both token and user ID");
      return;
    }

    setIsValidating(true);
    try {
      const res = await fetch(`${baseUrl}/api/v1/workouts/user/${userId}`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        let message = `API error: ${res.status}`;
        try {
          const body = (await res.json()) as { message?: string };
          if (body?.message) message = body.message;
        } catch {
          // ignore
        }
        setHealth({ status: "error", message });
        throw new Error(message);
      }

      setHealth({ status: "ok" });
      toast.success("Token validated successfully");
    } catch (err) {
      const message = err instanceof Error ? err.message : "Invalid token";
      toast.error(message);
    } finally {
      setIsValidating(false);
    }
  };

  const handleCheckHealth = async () => {
    setIsCheckingHealth(true);
    try {
      const health = await fetchHealth();
      setHealth(health);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to check health";
      toast.error(message);
    } finally {
      setIsCheckingHealth(false);
    }
  };

  const handleLogout = () => {
    clearAuth();
    toast.success("Logged out successfully");
    router.push("/login");
  };

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Settings</h1>
          <p className="text-muted-foreground mt-2">
            Manage your account and preferences
          </p>
        </div>

        {/* Appearance Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Appearance</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="theme">Theme</Label>
              <Select value={theme || "system"} onValueChange={setTheme}>
                <SelectTrigger id="theme">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="light">Light</SelectItem>
                  <SelectItem value="dark">Dark</SelectItem>
                  <SelectItem value="system">System</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* Session Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Session</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="token">Auth Token</Label>
              <Input
                id="token"
                type="password"
                value={token}
                onChange={(e) => setToken(e.target.value)}
                placeholder="Your auth token"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="userId">User ID</Label>
              <Input
                id="userId"
                type="text"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder="Your user ID"
              />
            </div>

            <Button
              onClick={handleValidateToken}
              disabled={isValidating}
              className="w-full"
            >
              {isValidating && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Validate Token
            </Button>
          </CardContent>
        </Card>

        {/* System Status */}
        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button
              onClick={handleCheckHealth}
              disabled={isCheckingHealth}
              className="w-full"
            >
              {isCheckingHealth && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Check Health
            </Button>

            {health && (
              <div className="p-4 bg-primary/10 rounded-lg border border-primary/20">
                <div className="space-y-2 text-sm">
                  <div>
                    <span className="text-muted-foreground">Status: </span>
                    <span className="font-medium text-foreground">
                      {health.status}
                    </span>
                  </div>
                  {health.message && (
                    <div>
                      <span className="text-muted-foreground">Message: </span>
                      <span className="font-medium text-foreground">
                        {health.message}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Logout */}
        <Card className="border-destructive/40">
          <CardHeader>
            <CardTitle>Logout</CardTitle>
          </CardHeader>
          <CardContent>
            <Button
              onClick={handleLogout}
              variant="destructive"
              className="w-full"
            >
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
