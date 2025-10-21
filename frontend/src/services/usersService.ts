import axios, { AxiosResponse } from 'axios'
import Cookies from 'js-cookie'

// User Profile DTOs
export interface UserProfileDTO {
  id: string
  user_id: string
  first_name: string
  last_name: string
  email: string
  phone?: string
  date_of_birth?: string
  address?: string
  city?: string
  state?: string
  zip_code?: string
  country?: string
  license_number?: string
  license_state?: string
  license_expiration?: string
  professional_title?: string
  specializations?: string[]
  years_of_experience?: number
  practice_name?: string
  practice_type?: string
  practice_address?: string
  practice_city?: string
  practice_state?: string
  practice_zip_code?: string
  practice_phone?: string
  emergency_contact_name?: string
  emergency_contact_phone?: string
  emergency_contact_email?: string
  profile_complete?: boolean
  status?: string
  created_at?: string
  updated_at?: string
}

export interface CreateUserProfileDTO {
  first_name: string
  last_name: string
  email: string
  phone?: string
  date_of_birth?: string
  license_number?: string
  license_state?: string
  license_expiration?: string
  years_of_experience?: number
  practice_name?: string
  practice_address?: string
  practice_phone?: string
}

export interface UpdateUserProfileDTO extends Partial<CreateUserProfileDTO> {
  // All fields are optional for updates
}

// Enrollment DTOs
export interface EnrollmentDataDTO {
  id: string
  user_id: string
  personal_info_complete: boolean
  licensure_details_complete: boolean
  practice_info_complete: boolean
  admin_setup_complete: boolean
  schedule_config_complete: boolean
  enrollment_status: 'in_progress' | 'completed' | 'paused'
  current_step: number
  total_steps: number
  completion_date?: string
  selected_plan: 'essential' | 'professional' | 'enterprise'
  payment_status: 'pending' | 'completed' | 'failed'
  preferred_contact_method?: string
  referral_source?: string
  special_requests?: string
  created_at: string
  updated_at: string
}

export interface StartEnrollmentDTO {
  user_id: string
  selected_plan: 'essential' | 'professional' | 'enterprise'
}

export interface EnrollmentProgressDTO {
  user_id: string
  current_step: number
  total_steps: number
  progress: number // percentage
}

// API Response DTOs
export interface UserProfileResponse {
  message: string
  profile: UserProfileDTO
}

export interface UserProfilesListResponse {
  profiles: UserProfileDTO[]
  limit: number
  offset: number
  count: number
  total: number
}

export interface UserProfilesSearchResponse {
  profiles: UserProfileDTO[]
  query: string
  count: number
  limit: number
  offset: number
}

export interface EnrollmentResponse {
  enrollment: EnrollmentDataDTO
}

export interface EnrollmentStartResponse {
  message: string
  user_id: string
  plan: string
}

export interface EnrollmentStepResponse {
  message: string
  user_id: string
  step: number
}

export interface EnrollmentProgressResponse {
  user_id: string
  current_step: number
  total_steps: number
  progress: number
}

const USERS_API_BASE_URL = import.meta.env.VITE_USERS_API_URL || 'http://localhost:3002/api/v1'

// Create axios instance for Users service
export const usersApiClient = axios.create({
  baseURL: USERS_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
usersApiClient.interceptors.request.use((config) => {
  const token = Cookies.get('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor for error handling
usersApiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('Users API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)

class UsersService {
  // User Profile Management
  
  async createProfile(profileData: CreateUserProfileDTO): Promise<UserProfileDTO> {
    const response: AxiosResponse<UserProfileResponse> = await usersApiClient.post('/users/profiles', profileData)
    return response.data.profile
  }

  async getProfile(userId: string): Promise<UserProfileDTO> {
    const response: AxiosResponse<UserProfileResponse> = await usersApiClient.get(`/users/profiles/${userId}`)
    return response.data.profile
  }

  async updateProfile(userId: string, profileData: UpdateUserProfileDTO): Promise<UserProfileDTO> {
    const response: AxiosResponse<UserProfileResponse> = await usersApiClient.put(`/users/profiles/${userId}`, profileData)
    return response.data.profile
  }

  async deleteProfile(userId: string): Promise<void> {
    await usersApiClient.delete(`/users/profiles/${userId}`)
  }

  async listProfiles(limit: number = 20, offset: number = 0): Promise<UserProfilesListResponse> {
    const response: AxiosResponse<UserProfilesListResponse> = await usersApiClient.get(
      `/users/profiles?limit=${limit}&offset=${offset}`
    )
    return response.data
  }

  async searchProfiles(query: string, limit: number = 20, offset: number = 0): Promise<UserProfilesSearchResponse> {
    const response: AxiosResponse<UserProfilesSearchResponse> = await usersApiClient.get(
      `/users/profiles/search?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}`
    )
    return response.data
  }

  // Enrollment Management

  async startEnrollment(userId: string, selectedPlan: 'essential' | 'professional' | 'enterprise'): Promise<EnrollmentStartResponse> {
    const response: AxiosResponse<EnrollmentStartResponse> = await usersApiClient.post('/enrollments/start', {
      user_id: userId,
      selected_plan: selectedPlan,
    })
    return response.data
  }

  async getEnrollment(userId: string): Promise<EnrollmentDataDTO> {
    const response: AxiosResponse<EnrollmentResponse> = await usersApiClient.get(`/enrollments/${userId}`)
    return response.data.enrollment
  }

  async completeStep(userId: string, step: number): Promise<EnrollmentStepResponse> {
    const response: AxiosResponse<EnrollmentStepResponse> = await usersApiClient.post(
      `/enrollments/${userId}/steps/${step}/complete`
    )
    return response.data
  }

  async getEnrollmentProgress(userId: string): Promise<EnrollmentProgressResponse> {
    const response: AxiosResponse<EnrollmentProgressResponse> = await usersApiClient.get(`/enrollments/${userId}/progress`)
    return response.data
  }

  // Utility Methods

  async healthCheck(): Promise<{ service: string; status: string; timestamp: number }> {
    const response = await usersApiClient.get('/health')
    return response.data
  }
}

export const usersService = new UsersService()
export default usersService