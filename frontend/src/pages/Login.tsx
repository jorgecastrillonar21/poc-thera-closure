import React, { useState } from 'react'
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  Link as MuiLink,
  Grid
} from '@mui/material'
import { useNavigate, Link } from 'react-router-dom'
import { authService } from '../services/authService'
import { Shield } from '@mui/icons-material'

const Login: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await authService.login({ email, password })
      console.log('Login successful:', response)
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed. Please try again.')
      console.error('Login error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      {/* Hero Section */}
      <Box
        sx={{
          background: 'linear-gradient(135deg, #2C5F5D 0%, #4A8B8A 100%)',
          color: 'white',
          py: { xs: 4, md: 6 },
        }}
      >
        <Container maxWidth="lg">
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography 
                variant="h3" 
                sx={{ 
                  fontWeight: 600, 
                  mb: 3,
                  fontSize: { xs: '1.8rem', md: '2.2rem' }
                }}
              >
                Welcome Back to TheraClosure
              </Typography>
              <Typography 
                variant="h6" 
                sx={{ 
                  opacity: 0.9,
                  mb: 3
                }}
              >
                Access your professional will and practice protection services
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Shield sx={{ mr: 2, fontSize: '2rem' }} />
                <Typography variant="body1">
                  Secure, HIPAA-compliant access to your account
                </Typography>
              </Box>
            </Grid>
            <Grid item xs={12} md={6}>
              {/* Login Form */}
              <Paper 
                elevation={8} 
                sx={{ 
                  p: 4, 
                  borderRadius: 3,
                  backgroundColor: 'rgba(255, 255, 255, 0.95)',
                  backdropFilter: 'blur(10px)'
                }}
              >
                <Box sx={{ textAlign: 'center', mb: 3 }}>
                  <Typography variant="h4" gutterBottom color="primary" fontWeight={600}>
                    Sign In
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Access your TheraClosure dashboard
                  </Typography>
                </Box>

                {error && (
                  <Alert severity="error" sx={{ mb: 3 }}>
                    {error}
                  </Alert>
                )}

                <Box component="form" onSubmit={handleSubmit} noValidate>
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="email"
                    label="Email Address"
                    name="email"
                    autoComplete="email"
                    autoFocus
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    sx={{ mb: 2 }}
                  />
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    name="password"
                    label="Password"
                    type="password"
                    id="password"
                    autoComplete="current-password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    sx={{ mb: 3 }}
                  />
                  <Button
                    type="submit"
                    fullWidth
                    variant="contained"
                    size="large"
                    disabled={loading}
                    sx={{ 
                      mb: 2, 
                      py: 1.5,
                      fontSize: '1.1rem',
                      fontWeight: 600
                    }}
                  >
                    {loading ? 'Signing In...' : 'Sign In'}
                  </Button>
                  
                  <Box sx={{ textAlign: 'center', mt: 2 }}>
                    <MuiLink 
                      href="#" 
                      variant="body2" 
                      sx={{ 
                        color: 'primary.main',
                        textDecoration: 'none',
                        '&:hover': { textDecoration: 'underline' }
                      }}
                    >
                      Forgot your password?
                    </MuiLink>
                  </Box>
                </Box>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Information Section */}
      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 8 } }}>
        <Box textAlign="center">
          <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
            New to TheraClosure?
          </Typography>
          
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              mb: 6,
              maxWidth: '600px',
              mx: 'auto'
            }}
          >
            Protect your clients, practice, and peace of mind with our comprehensive professional will and practice closure services designed specifically for mental health professionals.
          </Typography>
          
          <Box sx={{ display: 'flex', gap: 3, justifyContent: 'center', flexWrap: 'wrap' }}>
            <Button
              component={Link}
              to="/contact"
              variant="contained"
              size="large"
              sx={{ 
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem'
              }}
            >
              SCHEDULE FREE CONSULTATION
            </Button>
            <Button
              component={Link}
              to="/about"
              variant="outlined"
              size="large"
              sx={{ 
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem'
              }}
            >
              LEARN MORE
            </Button>
          </Box>
        </Box>
      </Container>
    </>
  )
}

export default Login