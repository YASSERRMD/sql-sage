"use client";

import dynamic from "next/dynamic";
import { useEffect, useState } from "react";

const Monaco = dynamic(() => import("@monaco-editor/react"), { ssr: false });

interface Props {
  value: string;
  onChange?: (v: string) => void;
  height?: number | string;
  readOnly?: boolean;
  language?: string;
}

export function SqlEditor({ value, onChange, height = 480, readOnly, language = "sql" }: Props) {
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) {
    return (
      <div
        style={{ height: typeof height === "number" ? `${height}px` : height }}
        className="rounded-md border bg-muted/30"
      />
    );
  }
  return (
    <Monaco
      height={height}
      defaultLanguage={language}
      theme="vs-dark"
      value={value}
      onChange={(v) => onChange?.(v ?? "")}
      options={{
        readOnly: !!readOnly,
        minimap: { enabled: false },
        fontSize: 13,
        wordWrap: "on",
        scrollBeyondLastLine: false,
        automaticLayout: true,
        renderLineHighlight: "gutter",
        tabSize: 2,
      }}
    />
  );
}
