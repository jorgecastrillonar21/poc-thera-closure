import axios, { AxiosResponse } from 'axios'
import { UserDTO, LoginDTO, RegisterDTO, AuthResponseDTO } from '@theraclosure/shared'
import Cookies from 'js-cookie'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

// Create axios instance
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
})

// Request interceptor to add auth token
apiClient.interceptors.request.use((config) => {
  const token = Cookies.get('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor to handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Try to refresh token
      try {
        const refreshToken = Cookies.get('refreshToken')
        if (refreshToken) {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refreshToken,
          })
          const { accessToken } = response.data
          Cookies.set('accessToken', accessToken, { expires: 1 })
          // Retry original request
          error.config.headers.Authorization = `Bearer ${accessToken}`
          return axios.request(error.config)
        }
      } catch (refreshError) {
        // Refresh failed, redirect to login
        Cookies.remove('accessToken')
        Cookies.remove('refreshToken')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

class AuthService {
  async login(credentials: LoginDTO): Promise<AuthResponseDTO> {
    const response: AxiosResponse<AuthResponseDTO> = await apiClient.post('/auth/login', credentials)
    const { accessToken, refreshToken, user } = response.data
    
    // Store tokens in cookies
    Cookies.set('accessToken', accessToken, { expires: 1 }) // 1 day
    Cookies.set('refreshToken', refreshToken, { expires: 7 }) // 7 days
    
    return response.data
  }

  async register(userData: RegisterDTO): Promise<AuthResponseDTO> {
    const response: AxiosResponse<AuthResponseDTO> = await apiClient.post('/auth/register', userData)
    const { accessToken, refreshToken } = response.data
    
    // Store tokens in cookies
    Cookies.set('accessToken', accessToken, { expires: 1 })
    Cookies.set('refreshToken', refreshToken, { expires: 7 })
    
    return response.data
  }

  async getCurrentUser(): Promise<UserDTO> {
    const response: AxiosResponse<UserDTO> = await apiClient.get('/auth/me')
    return response.data
  }

  async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      Cookies.remove('accessToken')
      Cookies.remove('refreshToken')
    }
  }

  async initiateCognitoLogin(): Promise<string> {
    const response: AxiosResponse<{ url: string }> = await apiClient.get('/auth/cognito/login')
    return response.data.url
  }
}

export const authService = new AuthService()