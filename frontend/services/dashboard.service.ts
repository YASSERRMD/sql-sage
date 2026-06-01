import { api } from "@/lib/api";
import type { DashboardSummary } from "@/types";

export const dashboardService = {
  summary: async (): Promise<DashboardSummary> =>
    (await api.get<DashboardSummary>("/dashboard/summary")).data,
  trend: async (): Promise<{ date: string; count: number }[]> =>
    (await api.get("/dashboard/trend")).data,
  riskDistribution: async (): Promise<{ level: string; count: number }[]> =>
    (await api.get("/dashboard/risk-distribution")).data,
  objectTypes: async (): Promise<{ type: string; count: number }[]> =>
    (await api.get("/dashboard/object-types")).data,
};
