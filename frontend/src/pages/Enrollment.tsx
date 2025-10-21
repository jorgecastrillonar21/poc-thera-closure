import React, { useState } from 'react'
import { Container, Typography, Box, Alert, Button } from '@mui/material'
import { ArrowBackOutlined } from '@mui/icons-material'
import EnrollmentWorkflow from '../components/EnrollmentWorkflow'
import { useAuth } from '../contexts/AuthContext'
import { EnrollmentDataDTO } from '../services/usersService'

const Enrollment: React.FC = () => {
  const { user: currentUser } = useAuth()
  const [enrollmentComplete, setEnrollmentComplete] = useState(false)

  // For demo purposes, we'll use a hardcoded user ID
  // In a real app, this would come from the authenticated user context
  const demoUserId = "demo-user-123"

  const handleEnrollmentComplete = (enrollment: EnrollmentDataDTO) => {
    console.log('Enrollment completed:', enrollment)
    setEnrollmentComplete(true)
  }

  const handleBackToDashboard = () => {
    // Navigate back to dashboard
    window.location.href = '/dashboard'
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Box sx={{ mb: 3 }}>
        <Button
          startIcon={<ArrowBackOutlined />}
          onClick={handleBackToDashboard}
          sx={{ mb: 2 }}
        >
          Back to Dashboard
        </Button>
        
        <Typography variant="h4" component="h1" gutterBottom>
          Account Enrollment
        </Typography>
        
        <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
          Complete your account setup with our 5-step enrollment process
        </Typography>
      </Box>

      {!currentUser && (
        <Alert severity="info" sx={{ mb: 3 }}>
          <Typography variant="body2">
            <strong>Demo Mode:</strong> In a real application, you would need to be logged in to access enrollment. 
            For demonstration purposes, we're using a sample user ID.
          </Typography>
        </Alert>
      )}

      {enrollmentComplete && (
        <Alert severity="success" sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            🎉 Welcome to TheraClosure!
          </Typography>
          <Typography variant="body2">
            Your enrollment is complete and your account is fully activated. 
            You can now access all features of the platform.
          </Typography>
        </Alert>
      )}

      <EnrollmentWorkflow
        userId={currentUser?.id || demoUserId}
        onComplete={handleEnrollmentComplete}
      />

      {enrollmentComplete && (
        <Box sx={{ textAlign: 'center', mt: 4 }}>
          <Button
            variant="contained"
            size="large"
            onClick={handleBackToDashboard}
          >
            Go to Dashboard
          </Button>
        </Box>
      )}
    </Container>
  )
}

export default Enrollment