import React from 'react'
import { Container, Typography, Box, Button, Grid, Card, CardContent } from '@mui/material'
import { Link } from 'react-router-dom'

const Home: React.FC = () => {
  return (
    <>
      {/* Hero Section */}
      <Box
        sx={{
          background: 'linear-gradient(135deg, #2C5F5D 0%, #4A8B8A 100%)',
          color: 'white',
          py: 8,
        }}
      >
        <Container maxWidth="lg">
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="h1" gutterBottom>
                Professional Practice Closure Services
              </Typography>
              <Typography variant="h5" sx={{ mb: 4, opacity: 0.9 }}>
                Helping therapists close their practices with dignity, care, and professional integrity
              </Typography>
              <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                <Button
                  component={Link}
                  to="/contact"
                  variant="contained"
                  size="large"
                  sx={{ 
                    backgroundColor: 'secondary.main',
                    color: 'primary.main',
                    '&:hover': { backgroundColor: 'secondary.dark' }
                  }}
                >
                  Get Started
                </Button>
                <Button
                  component={Link}
                  to="/how-we-work"
                  variant="outlined"
                  size="large"
                  sx={{ 
                    borderColor: 'white',
                    color: 'white',
                    '&:hover': { borderColor: 'secondary.main', backgroundColor: 'rgba(255,255,255,0.1)' }
                  }}
                >
                  Learn More
                </Button>
              </Box>
            </Grid>
            <Grid item xs={12} md={6}>
              {/* Placeholder for hero image */}
              <Box
                sx={{
                  backgroundColor: 'rgba(255,255,255,0.1)',
                  height: 300,
                  borderRadius: 2,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                <Typography variant="h6" sx={{ opacity: 0.7 }}>
                  Hero Image Placeholder
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Services Section */}
      <Container maxWidth="lg" sx={{ py: 8 }}>
        <Typography variant="h2" textAlign="center" gutterBottom>
          Why Choose TheraClosure?
        </Typography>
        <Typography variant="h6" textAlign="center" sx={{ mb: 6, opacity: 0.8 }}>
          We understand the challenges of closing a therapy practice
        </Typography>
        
        <Grid container spacing={4}>
          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
              <CardContent>
                <Typography variant="h5" gutterBottom color="primary">
                  Professional Guidance
                </Typography>
                <Typography>
                  Expert support through every step of the closure process with industry best practices
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
              <CardContent>
                <Typography variant="h5" gutterBottom color="primary">
                  Client Care
                </Typography>
                <Typography>
                  Compassionate client transition services ensuring continuity of care
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
              <CardContent>
                <Typography variant="h5" gutterBottom color="primary">
                  Complete Support
                </Typography>
                <Typography>
                  Full-service assistance with legal, administrative, and emotional aspects
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        <Box textAlign="center" sx={{ mt: 6 }}>
          <Button
            component={Link}
            to="/about"
            variant="contained"
            size="large"
          >
            Learn About Our Process
          </Button>
        </Box>
      </Container>
    </>
  )
}

export default Home