export type ObjectType =
  | "procedure"
  | "function"
  | "package"
  | "trigger"
  | "view"
  | "sql_script"
  | "unknown";

export type RiskLevel = "low" | "medium" | "high" | "critical";

export interface User {
  id: string;
  email: string;
  name: string;
  role: "admin" | "user";
  isActive: boolean;
  createdAt: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface Provider {
  id: string;
  name: string;
  baseUrl: string;
  modelName: string;
  temperature: number;
  maxTokens: number;
  isDefault: boolean;
  apiKeyPreview: string;
  createdAt: string;
}

export interface CreateProviderInput {
  name: string;
  baseUrl: string;
  apiKey: string;
  modelName: string;
  temperature: number;
  maxTokens: number;
  isDefault: boolean;
}

export interface TestConnectionResult {
  ok: boolean;
  message: string;
  latencyMs: number;
}

export interface TableDependency {
  name: string;
  operation: "SELECT" | "INSERT" | "UPDATE" | "DELETE" | "MERGE";
  columns: string[];
  purpose: string;
}

export interface Parameter {
  name: string;
  direction: "IN" | "OUT" | "IN OUT";
  type: string;
  description?: string;
}

export interface BusinessRule {
  rule: string;
  rationale?: string;
}

export interface Risk {
  level: RiskLevel;
  description: string;
  location?: string;
}

export interface PossibleBug {
  description: string;
  severity: RiskLevel;
  evidence?: string;
}

export interface ModernizationSuggestion {
  title: string;
  description: string;
  impact: "low" | "medium" | "high";
}

export interface AnalysisResult {
  objectName: string;
  objectType: ObjectType;
  summary: string;
  simpleExplanation: string;
  executionFlow: string[];
  mermaidDiagram: string;
  tables: TableDependency[];
  parameters: Parameter[];
  businessRules: BusinessRule[];
  risks: Risk[];
  possibleBugs: PossibleBug[];
  modernizationSuggestions: ModernizationSuggestion[];
  markdownReport: string;
}

export interface Analysis {
  id: string;
  providerId: string;
  objectName: string;
  objectType: ObjectType;
  summary: string;
  riskScore: RiskLevel;
  result: AnalysisResult;
  tokensUsed: number;
  createdAt: string;
}

export interface AnalysisListItem {
  id: string;
  objectName: string;
  objectType: ObjectType;
  summary: string;
  riskScore: RiskLevel;
  createdAt: string;
}

export interface CreateAnalysisInput {
  objectName: string;
  objectType: ObjectType;
  sourceCode: string;
  providerId?: string;
}

export interface DashboardSummary {
  totalAnalyses: number;
  totalProcedures: number;
  totalFunctions: number;
  totalPackages: number;
  highRiskFindings: number;
  providerUsage: { providerId: string; providerName: string; count: number }[];
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
