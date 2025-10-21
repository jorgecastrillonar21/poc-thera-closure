import React from 'react'
import { Container, Typography } from '@mui/material'

const Coverage: React.FC = () => {
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Coverage & Services
      </Typography>
      <Typography variant="body1">
        Information about our coverage and service areas for therapy practice closure.
      </Typography>
    </Container>
  )
}

export default Coverage