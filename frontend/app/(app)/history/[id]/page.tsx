"use client";

import { useState } from "react";
import { toast } from "sonner";
import { useAnalysis } from "@/hooks/useAnalyses";
import { api, loadAccessToken } from "@/lib/api";
import { AnalysisResultView } from "@/components/analysis/AnalysisResultView";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function HistoryDetailPage({ params }: { params: { id: string } }) {
  const { id } = params;
  const { data, isLoading, isError } = useAnalysis(id);
  const [downloading, setDownloading] = useState<"md" | "html" | null>(null);

  const download = async (format: "md" | "html") => {
    setDownloading(format);
    try {
      const token = loadAccessToken();
      const base = api.defaults.baseURL || "";
      const res = await fetch(`${base}/analyses/${id}/report?format=${format}`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) throw new Error("download failed");
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `${data?.objectName ?? id}.${format}`;
      a.click();
      URL.revokeObjectURL(url);
    } catch {
      toast.error("Download failed");
    } finally {
      setDownloading(null);
    }
  };

  if (isLoading) return <p className="p-6 text-sm text-muted-foreground">Loading...</p>;
  if (isError || !data)
    return <p className="p-6 text-sm text-destructive">Could not load analysis.</p>;

  return (
    <div className="space-y-4">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{data.objectName}</h1>
          <p className="text-muted-foreground">
            {data.objectType} • {new Date(data.createdAt).toLocaleString()}
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => download("md")} disabled={downloading === "md"}>
            {downloading === "md" ? "..." : "Markdown"}
          </Button>
          <Button variant="outline" onClick={() => download("html")} disabled={downloading === "html"}>
            {downloading === "html" ? "..." : "HTML"}
          </Button>
        </div>
      </header>
      <Card>
        <CardHeader>
          <CardTitle>Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">{data.summary}</p>
        </CardContent>
      </Card>
      {data.result && <AnalysisResultView result={data.result} />}
    </div>
  );
}
