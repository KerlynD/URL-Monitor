"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Plus, Lock } from "lucide-react";

interface AddMonitorFormProps {
  onAddMonitor: (url: string, password: string) => Promise<{ success: boolean; error?: string }>;
}

export function AddMonitorForm({ onAddMonitor }: AddMonitorFormProps) {
  const [url, setUrl] = useState("");
  const [password, setPassword] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const handleInitialSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!url.trim()) {
      setMessage({ type: "error", text: "Please enter a URL" });
      return;
    }

    // Open password dialog
    setMessage(null);
    setIsDialogOpen(true);
  };

  const handleFinalSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!password.trim()) {
      setMessage({ type: "error", text: "Please enter the password" });
      return;
    }

    setIsSubmitting(true);
    setMessage(null);

    const result = await onAddMonitor(url.trim(), password.trim());

    if (result.success) {
      setUrl("");
      setPassword("");
      setIsDialogOpen(false);
      setMessage({ type: "success", text: "Monitor added and checked successfully!" });
      setTimeout(() => setMessage(null), 3000);
    } else {
      setMessage({ type: "error", text: result.error || "Failed to add monitor" });
    }

    setIsSubmitting(false);
  };

  const handleDialogClose = () => {
    setIsDialogOpen(false);
    setPassword("");
    setMessage(null);
  };

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>Add New Monitor</CardTitle>
          <CardDescription>
            Enter a URL to start monitoring its uptime and performance
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleInitialSubmit} className="space-y-4">
            <div className="flex gap-3">
              <Input
                type="url"
                placeholder="https://example.com"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                disabled={isSubmitting}
                className="flex-1"
                required
              />
              <Button type="submit" disabled={isSubmitting} className="shrink-0">
                <Plus className="h-4 w-4 mr-2" />
                Add Monitor
              </Button>
            </div>

            {message && (
              <div
                className={`p-3 rounded-md text-sm font-medium ${
                  message.type === "success"
                    ? "bg-green-50 text-green-700 border border-green-200 dark:bg-green-950/30 dark:text-green-400 dark:border-green-900"
                    : "bg-red-50 text-red-700 border border-red-200 dark:bg-red-950/30 dark:text-red-400 dark:border-red-900"
                }`}
              >
                {message.text}
              </div>
            )}
          </form>
        </CardContent>
      </Card>

      <Dialog open={isDialogOpen} onOpenChange={handleDialogClose}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Lock className="h-5 w-5" />
              Password Required
            </DialogTitle>
            <DialogDescription>
              Enter the admin password to add a new monitor for <strong className="text-foreground">{url}</strong>
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleFinalSubmit}>
            <div className="space-y-4 py-4">
              <Input
                type="password"
                placeholder="Enter password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isSubmitting}
                autoFocus
                required
              />
              
              {message && message.type === "error" && (
                <div className="bg-red-50 text-red-700 border border-red-200 dark:bg-red-950/30 dark:text-red-400 dark:border-red-900 p-3 rounded-md text-sm font-medium">
                  {message.text}
                </div>
              )}
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={handleDialogClose}
                disabled={isSubmitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Adding..." : "Add Monitor"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}
