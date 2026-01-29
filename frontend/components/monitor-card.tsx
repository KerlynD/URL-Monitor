"use client";

import { Monitor } from "@/types/monitor";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useState } from "react";

interface MonitorCardProps {
  monitor: Monitor;
  onCheckNow: (monitorId: string) => Promise<{ success: boolean; error?: string }>;
}

function simplifyError(error: string): string {
  if (error.includes("no such host") || error.includes("lookup")) {
    return "Website does not exist or cannot be reached";
  }
  if (error.includes("connection refused")) {
    return "Connection refused by server";
  }
  if (error.includes("timeout")) {
    return "Request timed out";
  }
  if (error.includes("certificate")) {
    return "SSL certificate error";
  }
  // Return simplified version if it's a long technical error
  if (error.length > 100) {
    return "Unable to reach website";
  }
  return error;
}

export function MonitorCard({ monitor, onCheckNow }: MonitorCardProps) {
  const [isChecking, setIsChecking] = useState(false);

  const status = monitor.last_result?.is_up ? "up" : monitor.last_result ? "down" : "unknown";
  const statusColor = status === "up" ? "success" : status === "down" ? "destructive" : "secondary";

  const formatResponseTime = (nanoseconds?: number) => {
    if (!nanoseconds) return "N/A";
    const ms = nanoseconds / 1000000;
    if (ms < 1000) {
      return `${Math.round(ms)}ms`;
    } else {
      return `${(ms / 1000).toFixed(2)}s`;
    }
  };

  const formatLastChecked = (timestamp?: string) => {
    if (!timestamp) return "Never";
    
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    
    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;
    
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
  };

  const handleCheckNow = async () => {
    setIsChecking(true);
    await onCheckNow(monitor.id);
    setIsChecking(false);
  };

  const formatUrl = (url: string) => {
    return url.replace(/^https?:\/\//, "").replace(/\/$/, "");
  };

  return (
    <Card className={`group relative overflow-hidden transition-all hover:shadow-lg hover:scale-[1.02] ${
      status === "up" ? "border-l-4 border-l-green-500" : 
      status === "down" ? "border-l-4 border-l-red-500" : 
      "border-l-4 border-l-gray-400"
    }`}>
      {/* Status indicator dot */}
      <div className="absolute top-4 right-4">
        <div className={`w-3 h-3 rounded-full ${
          status === "up" ? "bg-green-500" : 
          status === "down" ? "bg-red-500" : 
          "bg-gray-400"
        } ${status === "up" ? "animate-pulse" : ""}`} />
      </div>

      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-2 pr-6">
          <CardTitle className="text-base font-semibold line-clamp-2 break-all group-hover:text-primary transition-colors">
            {formatUrl(monitor.url)}
          </CardTitle>
        </div>
        <Badge variant={statusColor} className="w-fit uppercase text-xs mt-2">
          {status}
        </Badge>
      </CardHeader>
      
      <CardContent className="space-y-4">
        <div className="grid grid-cols-3 gap-3 text-sm">
          <div className="space-y-1">
            <p className="text-muted-foreground text-[10px] uppercase tracking-wider font-medium">Response</p>
            <p className="font-bold text-base">
              {formatResponseTime(monitor.last_result?.response_time)}
            </p>
          </div>
          <div className="space-y-1">
            <p className="text-muted-foreground text-[10px] uppercase tracking-wider font-medium">Status</p>
            <p className="font-bold text-base">
              {monitor.last_result?.status_code || "—"}
            </p>
          </div>
          <div className="space-y-1">
            <p className="text-muted-foreground text-[10px] uppercase tracking-wider font-medium">Checked</p>
            <p className="font-bold text-base">
              {formatLastChecked(monitor.last_result?.timestamp)}
            </p>
          </div>
        </div>

        {monitor.last_result?.error && (
          <div className="bg-destructive/10 border border-destructive/20 rounded-md p-2.5">
            <p className="text-[11px] text-destructive font-medium break-words leading-tight">
              {simplifyError(monitor.last_result.error)}
            </p>
          </div>
        )}
      </CardContent>

      <CardFooter className="pt-0 pb-4">
        <Button 
          variant="outline" 
          size="sm"
          className="w-full hover:bg-primary hover:text-primary-foreground transition-colors" 
          onClick={handleCheckNow}
          disabled={isChecking}
        >
          {isChecking ? "Checking..." : "Check Now"}
        </Button>
      </CardFooter>
    </Card>
  );
}
