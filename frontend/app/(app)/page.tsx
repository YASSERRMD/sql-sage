"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useDeleteAnalysis, useAnalyses } from "@/hooks/useAnalyses";
import { useProviders } from "@/hooks/useProviders";
import { useDashboard } from "@/hooks/useDashboard";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { formatDate, maskKey } from "@/lib/utils";
import type { ObjectType, RiskLevel } from "@/types";

const RISK_VARIANT: Record<RiskLevel, "success" | "secondary" | "warning" | "destructive"> = {
  low: "success",
  medium: "secondary",
  high: "warning",
  critical: "destructive",
};

export default function DashboardPage() {
  const router = useRouter();
  const { data: summary } = useDashboard();
  const { data: providers } = useProviders();
  const { data, isLoading } = useAnalyses({ pageSize: 8 });
  const del = useDeleteAnalysis();
  const [q, setQ] = useState("");
  const [objectType, setObjectType] = useState<ObjectType | "">("");
  const [risk, setRisk] = useState<RiskLevel | "">("");

  return (
    <div className="space-y-6">
      <header>
        <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">High-level view of analyses, risk, and provider usage.</p>
      </header>

      <section className="grid gap-4 md:grid-cols-4">
        <Stat label="Total analyses" value={summary?.totalAnalyses ?? 0} />
        <Stat label="Procedures" value={summary?.totalProcedures ?? 0} />
        <Stat label="Functions" value={summary?.totalFunctions ?? 0} />
        <Stat label="High risk findings" value={summary?.highRisk ?? 0} />
      </section>

      <section className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Recent analyses</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <p className="text-sm text-muted-foreground">Loading...</p>
            ) : !data?.items.length ? (
              <p className="text-sm text-muted-foreground">No analyses yet.</p>
            ) : (
              <ul className="divide-y">
                {data.items.slice(0, 6).map((a) => (
                  <li
                    key={a.id}
                    className="flex cursor-pointer items-center justify-between py-2"
                    onClick={() => router.push(`/history/${a.id}`)}
                  >
                    <div>
                      <p className="font-medium">{a.objectName}</p>
                      <p className="text-xs text-muted-foreground">{formatDate(a.createdAt)}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline">{a.objectType}</Badge>
                      <Badge variant={RISK_VARIANT[(a.riskScore as RiskLevel) ?? "low"] ?? "secondary"}>
                        {a.riskScore || "—"}
                      </Badge>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Providers</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="space-y-2">
              {providers?.map((p) => (
                <li key={p.id} className="flex items-center justify-between rounded-md border p-2">
                  <div>
                    <p className="font-medium">{p.name}</p>
                    <p className="text-xs text-muted-foreground">{p.baseUrl}</p>
                  </div>
                  <span className="font-mono text-xs">{maskKey(p.apiKeyPreview)}</span>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </section>

      <section>
        <header className="mb-3 flex flex-wrap items-end justify-between gap-2">
          <h2 className="text-lg font-semibold">Analysis history</h2>
          <div className="flex flex-wrap gap-2">
            <Input
              placeholder="Search..."
              value={q}
              onChange={(e) => setQ(e.target.value)}
              className="w-48"
            />
            <select
              value={objectType}
              onChange={(e) => setObjectType(e.target.value as ObjectType | "")}
              className="h-10 rounded-md border border-input bg-background px-3 text-sm"
            >
              <option value="">All types</option>
              <option value="procedure">Procedure</option>
              <option value="function">Function</option>
              <option value="package">Package</option>
              <option value="trigger">Trigger</option>
              <option value="view">View</option>
              <option value="sql_script">SQL Script</option>
            </select>
            <select
              value={risk}
              onChange={(e) => setRisk(e.target.value as RiskLevel | "")}
              className="h-10 rounded-md border border-input bg-background px-3 text-sm"
            >
              <option value="">All risks</option>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="critical">Critical</option>
            </select>
            <Button onClick={() => router.push("/history?q=" + encodeURIComponent(q))}>Search</Button>
          </div>
        </header>
        <Card>
          <CardContent className="p-0">
            <table className="w-full text-sm">
              <thead className="border-b text-muted-foreground">
                <tr>
                  <th className="p-3 text-left">Object</th>
                  <th className="p-3 text-left">Type</th>
                  <th className="p-3 text-left">Risk</th>
                  <th className="p-3 text-left">Summary</th>
                  <th className="p-3 text-left">Date</th>
                  <th className="p-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {data?.items.map((a) => (
                  <tr key={a.id} className="border-b">
                    <td className="p-3 font-medium">{a.objectName}</td>
                    <td className="p-3">
                      <Badge variant="outline">{a.objectType}</Badge>
                    </td>
                    <td className="p-3">
                      <Badge variant={RISK_VARIANT[(a.riskScore as RiskLevel) ?? "low"] ?? "secondary"}>
                        {a.riskScore || "—"}
                      </Badge>
                    </td>
                    <td className="p-3 text-muted-foreground">{a.summary?.slice(0, 80)}</td>
                    <td className="p-3 text-muted-foreground">{formatDate(a.createdAt)}</td>
                    <td className="p-3 text-right">
                      <Button size="sm" variant="ghost" onClick={() => router.push(`/history/${a.id}`)}>
                        View
                      </Button>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() =>
                          del.mutate(a.id, {
                            onSuccess: () => toast.success("Deleted"),
                            onError: () => toast.error("Delete failed"),
                          })
                        }
                      >
                        Delete
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <Card>
      <CardContent className="p-6">
        <p className="text-sm text-muted-foreground">{label}</p>
        <p className="mt-2 text-3xl font-semibold">{value}</p>
      </CardContent>
    </Card>
  );
}
