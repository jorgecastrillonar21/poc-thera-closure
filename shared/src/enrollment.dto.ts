import { EnrollmentStep } from './enums';

export interface EnrollmentDataDTO {
  userId: string;
  step: EnrollmentStep;
  data: Record<string, any>;
  isCompleted: boolean;
  completedAt?: Date;
}

export interface PersonalInfoDTO {
  firstName: string;
  lastName: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
}

export interface LicensureDTO {
  licenseNumber: string;
  licenseState: string;
  licenseExpiration: Date;
  specializations: string[];
  certifications: string[];
}

export interface PracticeDTO {
  practiceName: string;
  practiceAddress: string;
  practicePhone: string;
  practiceType: 'individual' | 'group' | 'hospital' | 'clinic';
  yearsExperience: number;
}

export interface AdminDTO {
  preferredContactMethod: 'email' | 'phone' | 'text';
  notificationPreferences: {
    email: boolean;
    sms: boolean;
    push: boolean;
  };
}

export interface ScheduleDTO {
  availability: {
    monday: { start: string; end: string; available: boolean };
    tuesday: { start: string; end: string; available: boolean };
    wednesday: { start: string; end: string; available: boolean };
    thursday: { start: string; end: string; available: boolean };
    friday: { start: string; end: string; available: boolean };
    saturday: { start: string; end: string; available: boolean };
    sunday: { start: string; end: string; available: boolean };
  };
  timezone: string;
}