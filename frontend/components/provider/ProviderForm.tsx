"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useCreateProvider } from "@/hooks/useProviders";
import type { CreateProviderInput } from "@/types";

const schema = z.object({
  name: z.string().min(2),
  baseUrl: z.string().url(),
  apiKey: z.string().min(8),
  modelName: z.string().min(1),
  temperature: z.coerce.number().min(0).max(2),
  maxTokens: z.coerce.number().int().positive(),
  isDefault: z.boolean().default(false),
});
type FormValues = z.infer<typeof schema>;

const PRESETS: { name: string; baseUrl: string; modelName: string }[] = [
  { name: "OpenAI", baseUrl: "https://api.openai.com/v1", modelName: "gpt-4o-mini" },
  { name: "OpenRouter", baseUrl: "https://openrouter.ai/api/v1", modelName: "openai/gpt-4o-mini" },
  { name: "Groq", baseUrl: "https://api.groq.com/openai/v1", modelName: "llama-3.1-70b-versatile" },
  { name: "DeepSeek", baseUrl: "https://api.deepseek.com/v1", modelName: "deepseek-chat" },
  { name: "Ollama (local)", baseUrl: "http://localhost:11434/v1", modelName: "llama3.1" },
  { name: "LM Studio (local)", baseUrl: "http://localhost:1234/v1", modelName: "local-model" },
];

export function ProviderForm({ onDone }: { onDone?: () => void }) {
  const create = useCreateProvider();
  const [preset, setPreset] = useState("");

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
      baseUrl: "",
      apiKey: "",
      modelName: "",
      temperature: 0.2,
      maxTokens: 2048,
      isDefault: false,
    },
  });

  const applyPreset = (key: string) => {
    const p = PRESETS.find((x) => x.name === key);
    if (!p) return;
    setPreset(key);
    form.setValue("name", p.name);
    form.setValue("baseUrl", p.baseUrl);
    form.setValue("modelName", p.modelName);
  };

  const onSubmit = async (v: FormValues) => {
    try {
      const payload: CreateProviderInput = { ...v };
      await create.mutateAsync(payload);
      toast.success("Provider created");
      form.reset();
      onDone?.();
    } catch {
      toast.error("Failed to create provider");
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Add provider</CardTitle>
        <CardDescription>
          Configure any OpenAI-compatible endpoint. API keys are encrypted at rest and never returned to the browser.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <Label>Preset</Label>
            <select
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              value={preset}
              onChange={(e) => applyPreset(e.target.value)}
            >
              <option value="">— select preset —</option>
              {PRESETS.map((p) => (
                <option key={p.name} value={p.name}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" {...form.register("name")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="modelName">Model</Label>
              <Input id="modelName" {...form.register("modelName")} />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="baseUrl">Base URL</Label>
            <Input id="baseUrl" placeholder="https://api.openai.com/v1" {...form.register("baseUrl")} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="apiKey">API key</Label>
            <Input id="apiKey" type="password" {...form.register("apiKey")} />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="temperature">Temperature</Label>
              <Input
                id="temperature"
                type="number"
                step="0.05"
                min={0}
                max={2}
                {...form.register("temperature", { valueAsNumber: true })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="maxTokens">Max tokens</Label>
              <Input
                id="maxTokens"
                type="number"
                min={1}
                {...form.register("maxTokens", { valueAsNumber: true })}
              />
            </div>
          </div>
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" {...form.register("isDefault")} className="h-4 w-4" />
            Set as default provider
          </label>
          <Textarea className="hidden" aria-hidden />
          <Button type="submit" disabled={create.isPending}>
            {create.isPending ? "Creating..." : "Create provider"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
