"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { providerService } from "@/services/provider.service";
import type { CreateProviderInput } from "@/types";

export const providerKeys = {
  all: ["providers"] as const,
  list: () => [...providerKeys.all, "list"] as const,
  detail: (id: string) => [...providerKeys.all, "detail", id] as const,
};

export function useProviders() {
  return useQuery({ queryKey: providerKeys.list(), queryFn: providerService.list });
}

export function useProvider(id: string) {
  return useQuery({
    queryKey: providerKeys.detail(id),
    queryFn: () => providerService.get(id),
    enabled: !!id,
  });
}

export function useCreateProvider() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateProviderInput) => providerService.create(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: providerKeys.all }),
  });
}

export function useUpdateProvider() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: Partial<CreateProviderInput> }) =>
      providerService.update(id, input),
    onSuccess: (_d, vars) => {
      qc.invalidateQueries({ queryKey: providerKeys.all });
      qc.invalidateQueries({ queryKey: providerKeys.detail(vars.id) });
    },
  });
}

export function useDeleteProvider() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => providerService.remove(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: providerKeys.all }),
  });
}

export function useTestProvider() {
  return useMutation({ mutationFn: (id: string) => providerService.test(id) });
}

export function useSetDefaultProvider() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => providerService.setDefault(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: providerKeys.all }),
  });
}
