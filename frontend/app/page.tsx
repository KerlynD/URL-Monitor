"use client";

import { useState, useEffect } from "react";
import { MonitorCard } from "@/components/monitor-card";
import { AddMonitorForm } from "@/components/add-monitor-form";
import { StatsCards } from "@/components/stats-cards";
import { DemoBanner } from "@/components/demo-banner";
import { Monitor } from "@/types/monitor";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const REFRESH_INTERVAL = 60000; // 60 seconds

export default function Home() {
  const [monitors, setMonitors] = useState<Monitor[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadMonitors = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/monitor`);
      
      if (!response.ok) {
        throw new Error("Failed to fetch monitors");
      }
      
      const data = await response.json();
      setMonitors(data || []);
      setError(null);
    } catch (err) {
      console.error("Error loading monitors:", err);
      setError("Failed to load monitors. Make sure the backend is running.");
    } finally {
      setLoading(false);
    }
  };

  const handleAddMonitor = async (url: string, password: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/monitor`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          url: url,
          check_interval: 60,
          password: password,
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to add monitor");
      }

      const newMonitor = await response.json();
      
      // Trigger immediate check
      await fetch(`${API_BASE_URL}/monitor/${newMonitor.id}/check`, {
        method: "POST",
      });

      // Reload monitors
      setTimeout(loadMonitors, 500);
      
      return { success: true };
    } catch (err: any) {
      console.error("Error adding monitor:", err);
      return { success: false, error: err.message };
    }
  };

  const handleCheckNow = async (monitorId: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/monitor/${monitorId}/check`, {
        method: "POST",
      });

      if (!response.ok) {
        throw new Error("Failed to trigger check");
      }

      // Reload monitors after a short delay
      setTimeout(loadMonitors, 1000);
      
      return { success: true };
    } catch (err: any) {
      console.error("Error checking monitor:", err);
      return { success: false, error: err.message };
    }
  };

  useEffect(() => {
    loadMonitors();
    
    const interval = setInterval(loadMonitors, REFRESH_INTERVAL);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-background">
      {/* Navigation Bar */}
      <nav className="border-b bg-card/50 backdrop-blur-sm sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4 max-w-7xl">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="text-2xl">🔍</div>
              <div>
                <h1 className="text-xl font-bold">URL Monitor</h1>
                <p className="text-xs text-muted-foreground">Real-time uptime monitoring</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-full bg-green-500/10 border border-green-500/20">
                <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
                <span className="text-xs font-medium text-green-700 dark:text-green-400">Live</span>
              </div>
            </div>
          </div>
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8 max-w-7xl">
        {/* Demo Banner */}
        <DemoBanner />

        {/* Stats Cards */}
        {!loading && !error && monitors.length > 0 && (
          <div className="mb-8">
            <StatsCards monitors={monitors} />
          </div>
        )}

        {/* Add Monitor Form */}
        <div className="mb-8">
          <AddMonitorForm onAddMonitor={handleAddMonitor} />
        </div>

        {/* Monitors List */}
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-semibold">Monitored Endpoints</h2>
              <p className="text-sm text-muted-foreground mt-1">
                {monitors.length} {monitors.length === 1 ? "monitor" : "monitors"} • Auto-refresh every 60s
              </p>
            </div>
          </div>

          {loading && (
            <div className="text-center py-16">
              <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-primary border-r-transparent" />
              <p className="mt-4 text-sm text-muted-foreground">Loading monitors...</p>
            </div>
          )}

          {error && (
            <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-6 text-center">
              <p className="text-destructive font-medium">{error}</p>
            </div>
          )}

          {!loading && !error && monitors.length === 0 && (
            <div className="text-center py-16 space-y-4 bg-muted/30 rounded-lg border-2 border-dashed">
              <div className="text-5xl opacity-50">📭</div>
              <div className="space-y-2">
                <p className="text-lg font-medium">No monitors yet</p>
                <p className="text-sm text-muted-foreground">Add your first URL above to start monitoring!</p>
              </div>
            </div>
          )}

          {!loading && !error && monitors.length > 0 && (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {monitors.map((monitor) => (
                <MonitorCard
                  key={monitor.id}
                  monitor={monitor}
                  onCheckNow={handleCheckNow}
                />
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <footer className="mt-16 text-center text-sm text-muted-foreground border-t pt-8">
          <p>URL Monitor by Angel Difo</p>
        </footer>
      </div>
    </div>
  );
}
