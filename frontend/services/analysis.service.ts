import { api } from "@/lib/api";
import type {
  Analysis,
  AnalysisListItem,
  CreateAnalysisInput,
  ObjectType,
  RiskLevel,
} from "@/types";

export interface AnalysisListParams {
  q?: string;
  objectType?: ObjectType;
  risk?: RiskLevel;
  page?: number;
  pageSize?: number;
}

export const analysisService = {
  create: async (input: CreateAnalysisInput): Promise<Analysis> =>
    (await api.post<Analysis>("/analyses", input)).data,
  list: async (params: AnalysisListParams = {}): Promise<{
    items: AnalysisListItem[];
    total: number;
  }> => (await api.get("/analyses", { params })).data,
  get: async (id: string): Promise<Analysis> => (await api.get<Analysis>(`/analyses/${id}`)).data,
  remove: async (id: string): Promise<void> => {
    await api.delete(`/analyses/${id}`);
  },
  reportUrl: (id: string, format: "md" | "html" | "pdf" = "md"): string =>
    `${api.defaults.baseURL}/analyses/${id}/report?format=${format}`,
};
