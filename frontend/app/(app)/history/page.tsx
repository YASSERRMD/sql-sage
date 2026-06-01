"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAnalyses, useDeleteAnalysis } from "@/hooks/useAnalyses";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { formatDate } from "@/lib/utils";
import type { ObjectType, RiskLevel } from "@/types";

const RISK_VARIANT: Record<RiskLevel, "success" | "secondary" | "warning" | "destructive"> = {
  low: "success",
  medium: "secondary",
  high: "warning",
  critical: "destructive",
};

export default function HistoryPage() {
  const router = useRouter();
  const [q, setQ] = useState("");
  const [objectType, setObjectType] = useState<ObjectType | "">("");
  const [risk, setRisk] = useState<RiskLevel | "">("");
  const { data, isLoading } = useAnalyses({ q, objectType: objectType || undefined, risk: risk || undefined });
  const del = useDeleteAnalysis();

  return (
    <div className="space-y-4">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Analysis history</h1>
          <p className="text-muted-foreground">Browse, search, and manage your analyses.</p>
        </div>
        <Button onClick={() => router.push("/workspace")}>New analysis</Button>
      </header>

      <Card>
        <CardContent className="flex flex-wrap gap-2 p-4">
          <Input
            placeholder="Search by object or summary"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            className="w-64"
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
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-0">
          {isLoading ? (
            <p className="p-6 text-sm text-muted-foreground">Loading...</p>
          ) : !data?.items.length ? (
            <p className="p-6 text-sm text-muted-foreground">No analyses.</p>
          ) : (
            <table className="w-full text-sm">
              <thead className="border-b text-muted-foreground">
                <tr>
                  <th className="p-3 text-left">Object</th>
                  <th className="p-3 text-left">Type</th>
                  <th className="p-3 text-left">Risk</th>
                  <th className="p-3 text-left">Date</th>
                  <th className="p-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {data.items.map((a) => (
                  <tr key={a.id} className="border-b">
                    <td className="p-3">
                      <button
                        className="font-medium text-primary hover:underline"
                        onClick={() => router.push(`/history/${a.id}`)}
                      >
                        {a.objectName}
                      </button>
                    </td>
                    <td className="p-3">
                      <Badge variant="outline">{a.objectType}</Badge>
                    </td>
                    <td className="p-3">
                      <Badge variant={RISK_VARIANT[(a.riskScore as RiskLevel) ?? "low"] ?? "secondary"}>
                        {a.riskScore || "—"}
                      </Badge>
                    </td>
                    <td className="p-3 text-muted-foreground">{formatDate(a.createdAt)}</td>
                    <td className="p-3 text-right">
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => router.push(`/history/${a.id}`)}
                      >
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
          )}
        </CardContent>
      </Card>
    </div>
  );
}
