import React from 'react'
import { Container, Typography, Box } from '@mui/material'

const About: React.FC = () => {
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box>
        <Typography variant="h3" component="h1" gutterBottom>
          About TheraClosure
        </Typography>
        <Typography variant="h6" color="text.secondary" paragraph>
          Professional therapy practice closure and client transition services.
        </Typography>
        <Typography variant="body1" paragraph>
          TheraClosure specializes in helping therapists manage practice closures 
          with dignity and care. We understand that closing a practice is a significant 
          life transition, and we're here to support both practitioners and their clients 
          through this important process.
        </Typography>
        <Typography variant="body1" paragraph>
          Our comprehensive services include client notification management, record 
          transfer coordination, referral assistance, and compliance support to ensure 
          all ethical and legal requirements are met during the closure process.
        </Typography>
      </Box>
    </Container>
  )
}

export default About