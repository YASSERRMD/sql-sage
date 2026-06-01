import { api } from "@/lib/api";
import type {
  CreateProviderInput,
  Provider,
  TestConnectionResult,
} from "@/types";

export const providerService = {
  list: async (): Promise<Provider[]> => (await api.get<Provider[]>("/providers")).data,
  get: async (id: string): Promise<Provider> => (await api.get<Provider>(`/providers/${id}`)).data,
  create: async (input: CreateProviderInput): Promise<Provider> =>
    (await api.post<Provider>("/providers", input)).data,
  update: async (id: string, input: Partial<CreateProviderInput>): Promise<Provider> =>
    (await api.put<Provider>(`/providers/${id}`, input)).data,
  remove: async (id: string): Promise<void> => {
    await api.delete(`/providers/${id}`);
  },
  test: async (id: string): Promise<TestConnectionResult> =>
    (await api.post<TestConnectionResult>(`/providers/${id}/test`)).data,
  setDefault: async (id: string): Promise<Provider> =>
    (await api.post<Provider>(`/providers/${id}/default`)).data,
};
