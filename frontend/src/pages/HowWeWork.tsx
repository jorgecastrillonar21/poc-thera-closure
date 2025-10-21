import React from 'react'
import { Container, Typography, Box, Grid, Card, CardContent } from '@mui/material'

const HowWeWork: React.FC = () => {
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        How We Work
      </Typography>
      
      <Grid container spacing={4} sx={{ mt: 2 }}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom color="primary">
                1. Assessment
              </Typography>
              <Typography variant="body2">
                We assess your practice closure needs and create a customized plan.
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom color="primary">
                2. Client Notification
              </Typography>
              <Typography variant="body2">
                Professional client notification and referral coordination.
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom color="primary">
                3. Record Transfer
              </Typography>
              <Typography variant="body2">
                Secure record transfer and compliance management.
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Container>
  )
}

export default HowWeWork