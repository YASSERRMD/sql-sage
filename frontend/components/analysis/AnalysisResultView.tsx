"use client";

import { useMemo } from "react";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import type { AnalysisResult, RiskLevel } from "@/types";

const RISK_VARIANT: Record<RiskLevel, "success" | "secondary" | "warning" | "destructive"> = {
  low: "success",
  medium: "secondary",
  high: "warning",
  critical: "destructive",
};

export function AnalysisResultView({ result }: { result: AnalysisResult }) {
  return (
    <div className="grid gap-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <span>{result.objectName}</span>
            <Badge variant="outline">{result.objectType}</Badge>
          </CardTitle>
          <CardDescription>{result.summary}</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">{result.simpleExplanation}</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Execution flow</CardTitle>
        </CardHeader>
        <CardContent>
          <ol className="list-decimal space-y-1 pl-5 text-sm">
            {result.executionFlow.map((s, i) => (
              <li key={i}>{s}</li>
            ))}
          </ol>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Diagram</CardTitle>
        </CardHeader>
        <CardContent>
          <pre className="overflow-auto rounded-md bg-muted p-3 text-xs">
            <code>{result.mermaidDiagram}</code>
          </pre>
        </CardContent>
      </Card>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Tables</CardTitle>
          </CardHeader>
          <CardContent>
            <TableList items={result.tables} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Parameters</CardTitle>
          </CardHeader>
          <CardContent>
            <ParamList items={result.parameters} />
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Business rules</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="list-disc space-y-1 pl-5 text-sm">
              {result.businessRules.map((r, i) => (
                <li key={i}>
                  <span className="font-medium">{r.rule}</span>
                  {r.rationale && <span className="text-muted-foreground"> — {r.rationale}</span>}
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Modernization</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="list-disc space-y-1 pl-5 text-sm">
              {result.modernizationSuggestions.map((m, i) => (
                <li key={i}>
                  <span className="font-medium">{m.title}</span> — {m.description}
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Risks</CardTitle>
        </CardHeader>
        <CardContent>
          <ul className="space-y-2 text-sm">
            {result.risks.map((r, i) => (
              <li key={i} className="flex items-start gap-2">
                <Badge variant={RISK_VARIANT[r.level] ?? "secondary"}>{r.level}</Badge>
                <div>
                  <p>{r.description}</p>
                  {r.location && <p className="text-xs text-muted-foreground">{r.location}</p>}
                </div>
              </li>
            ))}
          </ul>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Possible bugs</CardTitle>
        </CardHeader>
        <CardContent>
          <ul className="space-y-2 text-sm">
            {result.possibleBugs.map((b, i) => (
              <li key={i} className="flex items-start gap-2">
                <Badge variant={RISK_VARIANT[b.severity] ?? "secondary"}>{b.severity}</Badge>
                <div>
                  <p>{b.description}</p>
                  {b.evidence && <p className="text-xs text-muted-foreground">{b.evidence}</p>}
                </div>
              </li>
            ))}
          </ul>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Report (Markdown)</CardTitle>
        </CardHeader>
        <CardContent>
          <pre className="overflow-auto rounded-md bg-muted p-3 text-xs">
            <code>{result.markdownReport}</code>
          </pre>
        </CardContent>
      </Card>
    </div>
  );
}

function TableList({ items }: { items: AnalysisResult["tables"] }) {
  if (!items.length) return <p className="text-sm text-muted-foreground">No tables detected.</p>;
  return (
    <ul className="space-y-2 text-sm">
      {items.map((t, i) => (
        <li key={i} className="rounded-md border p-2">
          <div className="flex items-center gap-2">
            <span className="font-medium">{t.name}</span>
            <Badge variant="outline">{t.operation}</Badge>
          </div>
          {t.columns.length > 0 && (
            <p className="mt-1 text-xs text-muted-foreground">Columns: {t.columns.join(", ")}</p>
          )}
          <p className="text-xs">{t.purpose}</p>
        </li>
      ))}
    </ul>
  );
}

function ParamList({ items }: { items: AnalysisResult["parameters"] }) {
  if (!items.length) return <p className="text-sm text-muted-foreground">No parameters detected.</p>;
  return (
    <ul className="space-y-1 text-sm">
      {items.map((p, i) => (
        <li key={i} className="flex items-center gap-2">
          <Badge variant="outline">{p.direction}</Badge>
          <span className="font-mono">{p.name}</span>
          <span className="text-xs text-muted-foreground">{p.type}</span>
          {p.description && <span className="text-xs">— {p.description}</span>}
        </li>
      ))}
    </ul>
  );
}
