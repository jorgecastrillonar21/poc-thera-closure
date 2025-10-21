import React from 'react'
import { Container, Typography, Card, CardContent, Grid } from '@mui/material'

const Dashboard: React.FC = () => {
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Dashboard
      </Typography>
      
      <Grid container spacing={3} sx={{ mt: 2 }}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Practice Status
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Active - Closure Planning
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Client Notifications
              </Typography>
              <Typography variant="body2" color="text.secondary">
                15 pending notifications
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Records Transfer
              </Typography>
              <Typography variant="body2" color="text.secondary">
                8 transfers completed
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Container>
  )
}

export default Dashboard