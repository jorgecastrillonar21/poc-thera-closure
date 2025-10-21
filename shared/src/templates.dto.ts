export interface TemplateDTO {
  id: string;
  title: string;
  content: string;
  category: 'closure' | 'referral' | 'admin' | 'general';
  isActive: boolean;
  createdBy: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface CreateTemplateDTO {
  title: string;
  content: string;
  category: 'closure' | 'referral' | 'admin' | 'general';
}

export interface UpdateTemplateDTO {
  title?: string;
  content?: string;
  category?: 'closure' | 'referral' | 'admin' | 'general';
  isActive?: boolean;
}

export interface SupportTicketDTO {
  id: string;
  userId: string;
  subject: string;
  message: string;
  priority: 'low' | 'medium' | 'high';
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  createdAt: Date;
  updatedAt: Date;
}

export interface CreateSupportTicketDTO {
  subject: string;
  message: string;
  priority: 'low' | 'medium' | 'high';
}