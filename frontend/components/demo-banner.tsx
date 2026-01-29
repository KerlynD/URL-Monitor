"use client";

import { Info } from "lucide-react";

export function DemoBanner() {
  return (
    <div className="bg-blue-50 dark:bg-blue-950/30 border-l-4 border-l-blue-500 rounded-lg p-4 mb-6">
      <div className="flex items-start gap-3">
        <Info className="h-5 w-5 text-blue-600 dark:text-blue-400 mt-0.5 shrink-0" />
        <div className="space-y-1">
          <p className="text-sm font-medium text-blue-900 dark:text-blue-100">
            Demo Portfolio Project
          </p>
          <p className="text-xs text-blue-700 dark:text-blue-300">
            This is a public demo. All monitors are shared across visitors. You can view and check any monitor, but you need a password to add new ones. Contact me for access!
          </p>
        </div>
      </div>
    </div>
  );
}
