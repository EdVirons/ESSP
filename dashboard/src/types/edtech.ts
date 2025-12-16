// EdTech Profile Status
export type EdTechProfileStatus = 'draft' | 'completed';

// Device Types breakdown
export interface DeviceTypes {
  laptops: number;
  chromebooks: number;
  tablets: number;
  desktops: number;
  other: number;
}

// AI Recommendation
export interface AIRecommendation {
  category: 'infrastructure' | 'training' | 'software' | 'security' | 'support';
  title: string;
  description: string;
  priority: 'high' | 'medium' | 'low';
}

// Follow-up Question
export interface FollowUpQuestion {
  id: string;
  question: string;
  context: string;
}

// EdTech Profile
export interface EdTechProfile {
  id: string;
  tenantId: string;
  schoolId: string;

  // Infrastructure Section
  totalDevices: number;
  deviceTypes: DeviceTypes;
  networkQuality: string;
  internetSpeed: string;
  lmsPlatform: string;
  existingSoftware: string[];
  itStaffCount: number;
  deviceAge: string;

  // Pain Points Section
  painPoints: string[];
  supportSatisfaction: number;
  biggestChallenges: string[];
  supportFrequency: string;
  avgResolutionTime: string;
  biggestFrustration: string;
  wishList: string;

  // Goals Section
  strategicGoals: string[];
  budgetRange: string;
  timeline: string;
  expansionPlans: string;
  priorityRanking: string[];
  decisionMakers: string[];

  // AI Section
  aiSummary: string;
  aiRecommendations: AIRecommendation[];
  followUpQuestions: FollowUpQuestion[];
  followUpResponses: Record<string, string>;

  // Metadata
  status: EdTechProfileStatus;
  completedAt?: string;
  completedBy?: string;
  version: number;
  createdAt: string;
  updatedAt: string;
}

// Profile History Entry
export interface EdTechProfileHistory {
  id: string;
  profileId: string;
  snapshot: EdTechProfile;
  changedBy: string;
  changeReason: string;
  changedAt: string;
}

// Form Options returned by the API
export interface EdTechFormOptions {
  networkQuality: string[];
  internetSpeed: string[];
  deviceAge: string[];
  supportFrequency: string[];
  resolutionTime: string[];
  budgetRange: string[];
  timeline: string[];
  lmsPlatforms: string[];
  existingSoftware: string[];
  painPoints: string[];
  strategicGoals: string[];
}

// Request types
export interface SaveProfileRequest {
  schoolId: string;
  totalDevices?: number;
  deviceTypes?: DeviceTypes;
  networkQuality?: string;
  internetSpeed?: string;
  lmsPlatform?: string;
  existingSoftware?: string[];
  itStaffCount?: number;
  deviceAge?: string;
  painPoints?: string[];
  supportSatisfaction?: number;
  biggestChallenges?: string[];
  supportFrequency?: string;
  avgResolutionTime?: string;
  biggestFrustration?: string;
  wishList?: string;
  strategicGoals?: string[];
  budgetRange?: string;
  timeline?: string;
  expansionPlans?: string;
  priorityRanking?: string[];
  decisionMakers?: string[];
}

export interface SubmitFollowUpRequest {
  responses: Record<string, string>;
}

// Response types
export interface EdTechProfileResponse {
  profile: EdTechProfile | null;
}

export interface EdTechHistoryResponse {
  history: EdTechProfileHistory[];
}

// Assessment step data types
export interface InfrastructureStepData {
  totalDevices: number;
  deviceTypes: DeviceTypes;
  networkQuality: string;
  internetSpeed: string;
  lmsPlatform: string;
  existingSoftware: string[];
  itStaffCount: number;
  deviceAge: string;
}

export interface PainPointsStepData {
  painPoints: string[];
  supportSatisfaction: number;
  biggestChallenges: string[];
  supportFrequency: string;
  avgResolutionTime: string;
  biggestFrustration: string;
  wishList: string;
}

export interface GoalsStepData {
  strategicGoals: string[];
  budgetRange: string;
  timeline: string;
  expansionPlans: string;
  priorityRanking: string[];
  decisionMakers: string[];
}

// Combined assessment form data
export interface EdTechAssessmentData extends
  InfrastructureStepData,
  PainPointsStepData,
  GoalsStepData {}
