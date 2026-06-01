"use client";

import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  useDeleteProvider,
  useProviders,
  useSetDefaultProvider,
  useTestProvider,
} from "@/hooks/useProviders";
import type { Provider } from "@/types";

export function ProviderList() {
  const { data, isLoading } = useProviders();
  const del = useDeleteProvider();
  const test = useTestProvider();
  const setDef = useSetDefaultProvider();

  if (isLoading) {
    return <p className="text-sm text-muted-foreground">Loading providers...</p>;
  }
  if (!data || data.length === 0) {
    return (
      <p className="text-sm text-muted-foreground">
        No providers yet. Add your first provider to start analyzing.
      </p>
    );
  }
  return (
    <div className="grid gap-4 md:grid-cols-2">
      {data.map((p) => (
        <ProviderCard
          key={p.id}
          provider={p}
          onDelete={() => del.mutate(p.id)}
          onTest={() =>
            test.mutate(p.id, {
              onSuccess: (r) => {
                if (r.ok) toast.success(`Connected in ${r.latencyMs}ms`);
                else toast.error(r.message || "Connection failed");
              },
              onError: () => toast.error("Test failed"),
            })
          }
          onDefault={() => setDef.mutate(p.id, {
            onSuccess: () => toast.success("Default updated"),
            onError: () => toast.error("Failed to set default"),
          })}
        />
      ))}
    </div>
  );
}

function ProviderCard({
  provider,
  onDelete,
  onTest,
  onDefault,
}: {
  provider: Provider;
  onDelete: () => void;
  onTest: () => void;
  onDefault: () => void;
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div>
          <CardTitle className="text-base">{provider.name}</CardTitle>
          <CardDescription className="break-all">{provider.baseUrl}</CardDescription>
        </div>
        {provider.isDefault && <Badge>Default</Badge>}
      </CardHeader>
      <CardContent className="space-y-3 text-sm">
        <dl className="grid grid-cols-2 gap-y-1">
          <dt className="text-muted-foreground">Model</dt>
          <dd>{provider.modelName}</dd>
          <dt className="text-muted-foreground">Temperature</dt>
          <dd>{provider.temperature}</dd>
          <dt className="text-muted-foreground">Max tokens</dt>
          <dd>{provider.maxTokens}</dd>
          <dt className="text-muted-foreground">Key</dt>
          <dd className="font-mono text-xs">{provider.apiKeyPreview || "—"}</dd>
        </dl>
        <div className="flex flex-wrap gap-2 pt-2">
          <Button size="sm" variant="outline" onClick={onTest}>
            Test
          </Button>
          {!provider.isDefault && (
            <Button size="sm" variant="secondary" onClick={onDefault}>
              Make default
            </Button>
          )}
          <Button size="sm" variant="destructive" onClick={onDelete}>
            Delete
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
