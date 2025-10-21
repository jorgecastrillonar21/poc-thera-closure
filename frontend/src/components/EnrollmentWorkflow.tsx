import React, { useState, useEffect } from 'react'
import {
  Box,
  Paper,
  Typography,
  Stepper,
  Step,
  StepLabel,
  Button,
  CircularProgress,
  Alert,
  LinearProgress,
  Card,
  CardContent,
  Divider,
  Chip,
} from '@mui/material'
import {
  PersonOutlined,
  VerifiedUserOutlined,
  BusinessOutlined,
  AdminPanelSettingsOutlined,
  ScheduleOutlined,
  CheckCircleOutlined,
} from '@mui/icons-material'
import { EnrollmentDataDTO, usersService } from '../services/usersService'

interface EnrollmentWorkflowProps {
  userId: string
  onComplete?: (enrollment: EnrollmentDataDTO) => void
}

const ENROLLMENT_STEPS = [
  {
    label: 'Personal Information',
    description: 'Complete your personal and contact details',
    icon: <PersonOutlined />,
  },
  {
    label: 'License Verification',
    description: 'Verify your professional license and credentials',
    icon: <VerifiedUserOutlined />,
  },
  {
    label: 'Practice Information',
    description: 'Provide details about your practice',
    icon: <BusinessOutlined />,
  },
  {
    label: 'Admin Setup',
    description: 'Configure administrative settings and preferences',
    icon: <AdminPanelSettingsOutlined />,
  },
  {
    label: 'Schedule Configuration',
    description: 'Set up your availability and scheduling preferences',
    icon: <ScheduleOutlined />,
  },
]

const EnrollmentWorkflow: React.FC<EnrollmentWorkflowProps> = ({ userId, onComplete }) => {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [enrollment, setEnrollment] = useState<EnrollmentDataDTO | null>(null)
  const [completingStep, setCompletingStep] = useState<number | null>(null)

  useEffect(() => {
    loadEnrollment()
  }, [userId])

  const loadEnrollment = async () => {
    try {
      setLoading(true)
      setError(null)
      const enrollmentData = await usersService.getEnrollment(userId)
      setEnrollment(enrollmentData)
    } catch (err: any) {
      if (err.response?.status === 404) {
        // No enrollment exists, we'll start one when user clicks start
        setEnrollment(null)
      } else {
        setError(`Failed to load enrollment: ${err.response?.data?.error || err.message}`)
      }
    } finally {
      setLoading(false)
    }
  }

  const startEnrollment = async () => {
    try {
      setLoading(true)
      setError(null)
      await usersService.startEnrollment(userId, 'professional') // Default plan
      await loadEnrollment() // Reload to get the created enrollment
    } catch (err: any) {
      setError(`Failed to start enrollment: ${err.response?.data?.error || err.message}`)
    } finally {
      setLoading(false)
    }
  }

  const completeStep = async (step: number) => {
    try {
      setCompletingStep(step)
      setError(null)
      await usersService.completeStep(userId, step)
      await loadEnrollment() // Reload to get updated status
      
      // Call onComplete callback when all steps are finished
      if (step === 5 && onComplete && enrollment) {
        onComplete(enrollment)
      }
    } catch (err: any) {
      setError(`Failed to complete step ${step}: ${err.response?.data?.error || err.message}`)
    } finally {
      setCompletingStep(null)
    }
  }

  const getStepStatus = (stepNumber: number): 'completed' | 'active' | 'disabled' => {
    if (!enrollment) return 'disabled'
    
    const stepCompleted = getStepCompletedStatus(stepNumber)
    if (stepCompleted) return 'completed'
    if (stepNumber === enrollment.current_step) return 'active'
    if (stepNumber < enrollment.current_step) return 'completed'
    return 'disabled'
  }

  const getStepCompletedStatus = (stepNumber: number): boolean => {
    if (!enrollment) return false
    
    switch (stepNumber) {
      case 1: return enrollment.personal_info_complete
      case 2: return enrollment.licensure_details_complete
      case 3: return enrollment.practice_info_complete
      case 4: return enrollment.admin_setup_complete
      case 5: return enrollment.schedule_config_complete
      default: return false
    }
  }

  const canCompleteStep = (stepNumber: number): boolean => {
    if (!enrollment) return false
    return stepNumber === enrollment.current_step && !getStepCompletedStatus(stepNumber)
  }

  const getProgressPercentage = (): number => {
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

  const isEnrollmentComplete = (): boolean => {
    return enrollment?.enrollment_status === 'completed'
  }

  if (loading && !enrollment) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto' }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Typography variant="h4" gutterBottom>
          Enrollment Workflow
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {!enrollment ? (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="h6" gutterBottom>
              Ready to start your enrollment?
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              Complete the 5-step enrollment process to activate your account
            </Typography>
            <Button
              variant="contained"
              size="large"
              onClick={startEnrollment}
              disabled={loading}
              startIcon={loading ? <CircularProgress size={20} /> : undefined}
            >
              {loading ? 'Starting...' : 'Start Enrollment'}
            </Button>
          </Box>
        ) : (
          <>
            {/* Progress Overview */}
            <Card sx={{ mb: 4, bgcolor: 'primary.50' }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">
                    Enrollment Progress
                  </Typography>
                  <Chip
                    label={enrollment.enrollment_status.toUpperCase()}
                    color={enrollment.enrollment_status === 'completed' ? 'success' : 'primary'}
                    icon={enrollment.enrollment_status === 'completed' ? <CheckCircleOutlined /> : undefined}
                  />
                </Box>
                <Box sx={{ mb: 1 }}>
                  <LinearProgress
                    variant="determinate"
                    value={getProgressPercentage()}
                    sx={{ height: 8, borderRadius: 4 }}
                  />
                </Box>
                <Typography variant="body2" color="text.secondary">
                  {Math.round(getProgressPercentage())}% Complete - Step {enrollment.current_step} of {enrollment.total_steps}
                </Typography>
                {enrollment.completion_date && (
                  <Typography variant="body2" color="success.main" sx={{ mt: 1 }}>
                    Completed on {new Date(enrollment.completion_date).toLocaleDateString()}
                  </Typography>
                )}
              </CardContent>
            </Card>

            {/* Stepper */}
            <Stepper orientation="vertical" sx={{ mb: 4 }}>
              {ENROLLMENT_STEPS.map((step, index) => {
                const stepNumber = index + 1
                const stepStatus = getStepStatus(stepNumber)
                const canComplete = canCompleteStep(stepNumber)
                const isCompletingThisStep = completingStep === stepNumber

                return (
                  <Step key={step.label} active={stepStatus === 'active'} completed={stepStatus === 'completed'}>
                    <StepLabel
                      icon={stepStatus === 'completed' ? <CheckCircleOutlined color="success" /> : step.icon}
                    >
                      <Box>
                        <Typography variant="h6">{step.label}</Typography>
                        <Typography variant="body2" color="text.secondary">
                          {step.description}
                        </Typography>
                        {canComplete && (
                          <Box sx={{ mt: 2 }}>
                            <Button
                              variant="contained"
                              size="small"
                              onClick={() => completeStep(stepNumber)}
                              disabled={isCompletingThisStep}
                              startIcon={isCompletingThisStep ? <CircularProgress size={16} /> : undefined}
                            >
                              {isCompletingThisStep ? 'Completing...' : 'Complete Step'}
                            </Button>
                          </Box>
                        )}
                        {stepStatus === 'completed' && (
                          <Chip
                            label="Completed"
                            size="small"
                            color="success"
                            sx={{ mt: 1 }}
                          />
                        )}
                      </Box>
                    </StepLabel>
                  </Step>
                )
              })}
            </Stepper>

            {/* Completion Message */}
            {isEnrollmentComplete() && (
              <Alert severity="success" sx={{ mt: 3 }}>
                <Typography variant="h6" gutterBottom>
                  🎉 Enrollment Complete!
                </Typography>
                <Typography variant="body2">
                  Congratulations! You have successfully completed all enrollment steps. 
                  Your account is now fully activated and ready to use.
                </Typography>
              </Alert>
            )}

            {/* Enrollment Details */}
            <Divider sx={{ my: 3 }} />
            <Typography variant="h6" gutterBottom>
              Enrollment Details
            </Typography>
            <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', mb: 2 }}>
              <Chip label={`Plan: ${enrollment.selected_plan.toUpperCase()}`} />
              <Chip 
                label={`Payment: ${enrollment.payment_status.toUpperCase()}`}
                color={enrollment.payment_status === 'completed' ? 'success' : 'warning'}
              />
            </Box>
            <Typography variant="body2" color="text.secondary">
              Started on {new Date(enrollment.created_at).toLocaleDateString()}
            </Typography>
          </>
        )}
      </Paper>
    </Box>
  )
}

export default EnrollmentWorkflow