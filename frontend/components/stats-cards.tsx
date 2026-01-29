"use client";

import { Monitor } from "@/types/monitor";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Activity, Globe, TrendingUp, Zap } from "lucide-react";

interface StatsCardsProps {
  monitors: Monitor[];
}

export function StatsCards({ monitors }: StatsCardsProps) {
  const totalMonitors = monitors.length;
  const upMonitors = monitors.filter(m => m.last_result?.is_up).length;
  const downMonitors = monitors.filter(m => m.last_result && !m.last_result.is_up).length;
  const uptime = totalMonitors > 0 ? ((upMonitors / totalMonitors) * 100).toFixed(1) : "0.0";
  
  const avgResponseTime = monitors.length > 0
    ? Math.round(
        monitors
          .filter(m => m.last_result?.response_time)
          .reduce((acc, m) => acc + (m.last_result?.response_time || 0), 0) /
          (monitors.filter(m => m.last_result?.response_time).length || 1) /
          1000000
      )
    : 0;

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card className="border-l-4 border-l-primary">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Monitors</CardTitle>
          <Globe className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{totalMonitors}</div>
          <p className="text-xs text-muted-foreground">
            Active monitoring endpoints
          </p>
        </CardContent>
      </Card>

      <Card className="border-l-4 border-l-green-500">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Services Up</CardTitle>
          <Activity className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">{upMonitors}</div>
          <p className="text-xs text-muted-foreground">
            {downMonitors > 0 ? `${downMonitors} down` : "All systems operational"}
          </p>
        </CardContent>
      </Card>

      <Card className="border-l-4 border-l-blue-500">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Uptime</CardTitle>
          <TrendingUp className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{uptime}%</div>
          <p className="text-xs text-muted-foreground">
            Overall availability
          </p>
        </CardContent>
      </Card>

      <Card className="border-l-4 border-l-orange-500">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Avg Response</CardTitle>
          <Zap className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{avgResponseTime}ms</div>
          <p className="text-xs text-muted-foreground">
            Average response time
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
