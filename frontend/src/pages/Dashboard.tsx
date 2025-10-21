import React, { useState, useEffect } from 'react'
import { 
  Container, 
  Typography, 
  Card, 
  CardContent, 
  Grid, 
  Button,
  Box,
  Alert,
  Chip,
  LinearProgress,
} from '@mui/material'
import {
  PersonOutlined,
  NotificationsOutlined,
  AssignmentOutlined,
  AccountCircleOutlined,
  SchoolOutlined,
} from '@mui/icons-material'
import { useAuth } from '../contexts/AuthContext'
import { EnrollmentDataDTO, usersService } from '../services/usersService'

const Dashboard: React.FC = () => {
  const { user: currentUser } = useAuth()
  const [enrollment, setEnrollment] = useState<EnrollmentDataDTO | null>(null)
  const [loadingEnrollment, setLoadingEnrollment] = useState(false)
  
  // Demo user ID for enrollment status
  const demoUserId = "demo-user-123"

  useEffect(() => {
    loadEnrollmentStatus()
  }, [])

  const loadEnrollmentStatus = async () => {
    try {
      setLoadingEnrollment(true)
      const enrollmentData = await usersService.getEnrollment(currentUser?.id || demoUserId)
      setEnrollment(enrollmentData)
    } catch (err: any) {
      // Enrollment might not exist yet, which is fine
      if (err.response?.status !== 404) {
        console.error('Failed to load enrollment:', err)
      }
    } finally {
      setLoadingEnrollment(false)
    }
  }

  const getEnrollmentProgress = (): number => {
    if (!enrollment) return 0
    
    const completedSteps = [
      enrollment.personal_info_complete,
      enrollment.licensure_details_complete,
      enrollment.practice_info_complete,
      enrollment.admin_setup_complete,
      enrollment.schedule_config_complete,
    ].filter(Boolean).length
    
    return (completedSteps / 5) * 100
  }

  const handleStartEnrollment = () => {
    window.location.href = '/enrollment'
  }

  const handleManageUsers = () => {
    window.location.href = '/user-management'
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Dashboard
      </Typography>

      {currentUser && (
        <Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
          Welcome back, {currentUser.firstName} {currentUser.lastName}!
        </Typography>
      )}
      
      {/* Enrollment Status Card */}
      {!loadingEnrollment && (
        <Grid container spacing={3} sx={{ mb: 3 }}>
          <Grid item xs={12}>
            <Card sx={{ bgcolor: enrollment?.enrollment_status === 'completed' ? 'success.50' : 'warning.50' }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6" gutterBottom>
                    <SchoolOutlined sx={{ mr: 1, verticalAlign: 'middle' }} />
                    Account Enrollment
                  </Typography>
                  {enrollment ? (
                    <Chip
                      label={enrollment.enrollment_status.toUpperCase()}
                      color={enrollment.enrollment_status === 'completed' ? 'success' : 'warning'}
                    />
                  ) : (
                    <Chip label="NOT STARTED" color="error" />
                  )}
                </Box>

                {enrollment ? (
                  <>
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        Progress: {Math.round(getEnrollmentProgress())}% Complete
                      </Typography>
                      <LinearProgress
                        variant="determinate"
                        value={getEnrollmentProgress()}
                        sx={{ height: 8, borderRadius: 4 }}
                      />
                    </Box>
                    <Typography variant="body2" sx={{ mb: 2 }}>
                      Step {enrollment.current_step} of {enrollment.total_steps} - {enrollment.selected_plan} plan
                    </Typography>
                    {enrollment.enrollment_status !== 'completed' && (
                      <Button variant="outlined" onClick={handleStartEnrollment}>
                        Continue Enrollment
                      </Button>
                    )}
                  </>
                ) : (
                  <>
                    <Typography variant="body2" sx={{ mb: 2 }}>
                      Complete your account setup to access all features
                    </Typography>
                    <Button variant="contained" onClick={handleStartEnrollment}>
                      Start Enrollment
                    </Button>
                  </>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}

      {/* Main Dashboard Cards */}
      <Grid container spacing={3} sx={{ mt: 2 }}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <PersonOutlined sx={{ mr: 1 }} />
                <Typography variant="h6">
                  User Management
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                Manage therapist profiles and accounts
              </Typography>
              <Button variant="outlined" onClick={handleManageUsers} fullWidth>
                Manage Users
              </Button>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <NotificationsOutlined sx={{ mr: 1 }} />
                <Typography variant="h6">
                  Client Notifications
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                15 pending notifications
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AssignmentOutlined sx={{ mr: 1 }} />
                <Typography variant="h6">
                  Records Transfer
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                8 transfers completed
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AccountCircleOutlined sx={{ mr: 1 }} />
                <Typography variant="h6">
                  Practice Status
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                Active - Closure Planning Phase
              </Typography>
              <Alert severity="info">
                Your practice closure process is in progress. Complete your enrollment to access all closure management tools.
              </Alert>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Quick Actions
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                <Button variant="outlined" onClick={handleStartEnrollment} fullWidth>
                  {enrollment?.enrollment_status === 'completed' ? 'View Enrollment' : 'Complete Enrollment'}
                </Button>
                <Button variant="outlined" onClick={handleManageUsers} fullWidth>
                  User Profiles
                </Button>
                <Button variant="outlined" disabled fullWidth>
                  Payment Settings (Coming Soon)
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Container>
  )
}

export default Dashboard