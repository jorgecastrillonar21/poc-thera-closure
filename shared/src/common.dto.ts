export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
  timestamp: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  meta: PaginationMeta;
}

export interface HealthCheckDTO {
  service: string;
  status: 'healthy' | 'unhealthy';
  version: string;
  timestamp: string;
  uptime: number;
  dependencies?: {
    [key: string]: 'healthy' | 'unhealthy';
  };
}