import React, { createContext, useContext, useReducer, useEffect } from 'react'
import { UserDTO, LoginDTO, RegisterDTO, authService } from '../services/authService'
import Cookies from 'js-cookie'

interface AuthState {
  user: UserDTO | null
  isLoading: boolean
  isAuthenticated: boolean
  error: string | null
}

type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: UserDTO }
  | { type: 'AUTH_ERROR'; payload: string }
  | { type: 'AUTH_LOGOUT' }
  | { type: 'CLEAR_ERROR' }

const initialState: AuthState = {
  user: null,
  isLoading: false,
  isAuthenticated: false,
  error: null,
}

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'AUTH_START':
      return { ...state, isLoading: true, error: null }
    case 'AUTH_SUCCESS':
      return {
        ...state,
        isLoading: false,
        isAuthenticated: true,
        user: action.payload,
        error: null,
      }
    case 'AUTH_ERROR':
      return {
        ...state,
        isLoading: false,
        isAuthenticated: false,
        user: null,
        error: action.payload,
      }
    case 'AUTH_LOGOUT':
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        error: null,
      }
    case 'CLEAR_ERROR':
      return { ...state, error: null }
    default:
      return state
  }
}

interface AuthContextType extends AuthState {
  login: (credentials: LoginDTO) => Promise<void>
  register: (userData: RegisterDTO) => Promise<void>
  logout: () => void
  clearError: () => void
  checkAuth: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: React.ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(authReducer, initialState)

  const login = async (credentials: LoginDTO) => {
    dispatch({ type: 'AUTH_START' })
    try {
      const response = await authService.login(credentials)
      dispatch({ type: 'AUTH_SUCCESS', payload: response.user })
    } catch (error: any) {
      dispatch({ type: 'AUTH_ERROR', payload: error.message || 'Login failed' })
    }
  }

  const register = async (userData: RegisterDTO) => {
    dispatch({ type: 'AUTH_START' })
    try {
      const response = await authService.register(userData)
      dispatch({ type: 'AUTH_SUCCESS', payload: response.user })
    } catch (error: any) {
      dispatch({ type: 'AUTH_ERROR', payload: error.message || 'Registration failed' })
    }
  }

  const logout = () => {
    authService.logout()
    dispatch({ type: 'AUTH_LOGOUT' })
  }

  const clearError = () => {
    dispatch({ type: 'CLEAR_ERROR' })
  }

  const checkAuth = async () => {
    const token = Cookies.get('accessToken')
    if (token) {
      dispatch({ type: 'AUTH_START' })
      try {
        const user = await authService.getCurrentUser()
        dispatch({ type: 'AUTH_SUCCESS', payload: user })
      } catch (error) {
        dispatch({ type: 'AUTH_LOGOUT' })
      }
    }
  }

  useEffect(() => {
    checkAuth()
  }, [])

  const value: AuthContextType = {
    ...state,
    login,
    register,
    logout,
    clearError,
    checkAuth,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}