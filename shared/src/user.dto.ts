import { UserRole, SubscriptionStatus } from './enums';

export interface UserDTO {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: UserRole;
  subscriptionStatus: SubscriptionStatus;
  stripeCustomerId?: string;
  cognitoId?: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface CreateUserDTO {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  role?: UserRole;
}

export interface UpdateUserDTO {
  firstName?: string;
  lastName?: string;
  role?: UserRole;
  subscriptionStatus?: SubscriptionStatus;
}

export interface LoginDTO {
  email: string;
  password: string;
}

export interface RegisterDTO {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface AuthResponseDTO {
  user: UserDTO;
  accessToken: string;
  refreshToken: string;
}

export interface RefreshTokenDTO {
  refreshToken: string;
}

export interface CognitoUserDTO {
  sub: string;
  email: string;
  given_name: string;
  family_name: string;
  email_verified: boolean;
}