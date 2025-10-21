export enum UserRole {
  ADMIN = 'ADMIN',
  THERAPIST = 'THERAPIST',
  STAFF = 'STAFF',
}

export enum SubscriptionStatus {
  ACTIVE = 'active',
  CANCELED = 'canceled',
  TRIALING = 'trialing',
  PAST_DUE = 'past_due',
  INCOMPLETE = 'incomplete',
}

export enum EnrollmentStep {
  PERSONAL_INFO = 'personal_info',
  LICENSURE = 'licensure',
  PRACTICE = 'practice',
  ADMIN = 'admin',
  SCHEDULE = 'schedule',
}