"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { SqlEditor } from "@/components/editor/SqlEditor";
import { AnalysisResultView } from "@/components/analysis/AnalysisResultView";
import { useCreateAnalysis } from "@/hooks/useAnalyses";
import { useProviders } from "@/hooks/useProviders";
import type { Analysis, ObjectType } from "@/types";

const SAMPLE = `CREATE OR REPLACE PROCEDURE close_account(p_account_id IN NUMBER) IS
  v_balance NUMBER;
BEGIN
  SELECT balance INTO v_balance FROM accounts WHERE id = p_account_id FOR UPDATE;
  IF v_balance <> 0 THEN
    UPDATE accounts SET status = 'CLOSED' WHERE id = p_account_id;
    INSERT INTO account_events(account_id, event, amount) VALUES (p_account_id, 'CLOSE', v_balance);
  END IF;
  COMMIT;
EXCEPTION
  WHEN OTHERS THEN
    ROLLBACK;
    RAISE;
END;
`;

export default function WorkspacePage() {
  const [code, setCode] = useState<string>(SAMPLE);
  const [name, setName] = useState("close_account");
  const [objectType, setObjectType] = useState<ObjectType>("procedure");
  const [providerId, setProviderId] = useState<string>("");
  const [result, setResult] = useState<Analysis | null>(null);
  const { data: providers } = useProviders();
  const create = useCreateAnalysis();

  const onAnalyze = async () => {
    if (!code.trim()) {
      toast.error("Source code is empty");
      return;
    }
    try {
      const a = await create.mutateAsync({
        objectName: name,
        objectType,
        sourceCode: code,
        providerId: providerId || undefined,
      });
      setResult(a);
      toast.success("Analysis complete");
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : "Analysis failed";
      toast.error(msg);
    }
  };

  const onSave = () => toast.info("Saved to history");
  const onExport = () => {
    if (!result) return;
    const blob = new Blob([result.result.markdownReport ?? ""], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${result.objectName}.md`;
    a.click();
    URL.revokeObjectURL(url);
  };
  const onClear = () => {
    setCode("");
    setResult(null);
  };

  return (
    <div className="space-y-4">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Analysis workspace</h1>
          <p className="text-muted-foreground">
            Paste or type PL/SQL and click Analyze. SQL is never executed.
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button onClick={onAnalyze} disabled={create.isPending}>
            {create.isPending ? "Analyzing..." : "Analyze"}
          </Button>
          <Button variant="outline" onClick={onSave} disabled={!result}>
            Save
          </Button>
          <Button variant="outline" onClick={onExport} disabled={!result}>
            Export
          </Button>
          <Button variant="ghost" onClick={onClear}>
            Clear
          </Button>
        </div>
      </header>

      <Card>
        <CardHeader>
          <CardTitle>Source</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
            <div className="space-y-2">
              <Label>Object name</Label>
              <input
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              />
            </div>
            <div className="space-y-2">
              <Label>Object type</Label>
              <select
                value={objectType}
                onChange={(e) => setObjectType(e.target.value as ObjectType)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="procedure">Procedure</option>
                <option value="function">Function</option>
                <option value="package">Package</option>
                <option value="trigger">Trigger</option>
                <option value="view">View</option>
                <option value="sql_script">SQL Script</option>
                <option value="unknown">Unknown</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label>Provider (optional)</Label>
              <select
                value={providerId}
                onChange={(e) => setProviderId(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="">Use default</option>
                {providers?.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>
            </div>
          </div>
          <SqlEditor value={code} onChange={setCode} />
        </CardContent>
      </Card>

      {result && <AnalysisResultView result={result.result} />}
    </div>
  );
}
