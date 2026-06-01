"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { analysisService, type AnalysisListParams } from "@/services/analysis.service";

export const analysisKeys = {
  all: ["analyses"] as const,
  list: (p: AnalysisListParams) => [...analysisKeys.all, "list", p] as const,
  detail: (id: string) => [...analysisKeys.all, "detail", id] as const,
};

export function useAnalyses(params: AnalysisListParams = {}) {
  return useQuery({
    queryKey: analysisKeys.list(params),
    queryFn: () => analysisService.list(params),
  });
}

export function useAnalysis(id: string) {
  return useQuery({
    queryKey: analysisKeys.detail(id),
    queryFn: () => analysisService.get(id),
    enabled: !!id,
  });
}

export function useCreateAnalysis() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: Parameters<typeof analysisService.create>[0]) => analysisService.create(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: analysisKeys.all }),
  });
}

export function useDeleteAnalysis() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => analysisService.remove(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: analysisKeys.all }),
  });
}
