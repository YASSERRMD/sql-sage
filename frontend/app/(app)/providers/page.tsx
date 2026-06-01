"use client";

import { useState } from "react";
import { ProviderForm } from "@/components/provider/ProviderForm";
import { ProviderList } from "@/components/provider/ProviderList";

export default function ProvidersPage() {
  const [creating, setCreating] = useState(false);
  return (
    <div className="space-y-6">
      <header className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">AI Providers</h1>
          <p className="text-muted-foreground">
            Manage OpenAI-compatible endpoints used to analyze your SQL.
          </p>
        </div>
        <button
          onClick={() => setCreating((v) => !v)}
          className="text-sm font-medium text-primary hover:underline"
        >
          {creating ? "Close" : "+ Add provider"}
        </button>
      </header>
      {creating && <ProviderForm onDone={() => setCreating(false)} />}
      <ProviderList />
    </div>
  );
}
