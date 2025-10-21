import React from 'react'
import { Container, Typography, Box, Grid, Card } from '@mui/material'
import { Psychology, Group, Security } from '@mui/icons-material'

const About: React.FC = () => {
  return (
    <>
      {/* Hero Section */}
      <Box
        sx={{
          background: 'linear-gradient(135deg, #2C5F5D 0%, #4A8B8A 100%)',
          color: 'white',
          py: { xs: 6, md: 8 },
          textAlign: 'center',
        }}
      >
        <Container maxWidth="lg">
          <Typography 
            variant="h2" 
            sx={{ 
              fontWeight: 600, 
              mb: 3,
              fontSize: { xs: '2rem', md: '3rem' }
            }}
          >
            Our Origin Story
          </Typography>
          <Typography 
            variant="h5" 
            sx={{ 
              opacity: 0.9,
              maxWidth: '600px',
              mx: 'auto'
            }}
          >
            Founded by clinicians, for clinicians
          </Typography>
        </Container>
      </Box>

      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
        
        {/* Our Story */}
        <Box sx={{ mb: 8 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 4 }}>
            We Are Clinicians, Too
          </Typography>
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              mb: 4,
              textAlign: 'center',
              maxWidth: '800px',
              mx: 'auto'
            }}
          >
            TheraClosure was founded by licensed psychotherapists who understand the unique challenges and ethical obligations that come with our profession. We recognized a critical gap in the field - the need for specialized, compassionate support when therapists face incapacitation, death, or retirement.
          </Typography>
          
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              mb: 4,
              textAlign: 'center',
              maxWidth: '800px',
              mx: 'auto'
            }}
          >
            Having witnessed firsthand the challenges that arise when practices close unexpectedly, we developed TheraClosure to provide the expertise, resources, and clinical understanding necessary to protect both therapists and their clients during these difficult transitions.
          </Typography>
        </Box>

        {/* Our Values */}
        <Box sx={{ backgroundColor: '#F5F1E8', py: { xs: 6, md: 8 }, px: 4, borderRadius: 2, mb: 8 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 6 }}>
            Our Core Values
          </Typography>
          
          <Grid container spacing={4}>
            <Grid item xs={12} md={4}>
              <Card sx={{ height: '100%', p: 4, textAlign: 'center' }}>
                <Psychology sx={{ fontSize: '3rem', color: 'primary.main', mb: 2 }} />
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Clinical Excellence
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  We bring deep clinical understanding to every aspect of practice closure, ensuring that therapeutic relationships are honored and clients receive appropriate care transitions.
                </Typography>
              </Card>
            </Grid>

            <Grid item xs={12} md={4}>
              <Card sx={{ height: '100%', p: 4, textAlign: 'center' }}>
                <Group sx={{ fontSize: '3rem', color: 'primary.main', mb: 2 }} />
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Compassionate Care
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  We understand that practice closure often occurs during times of grief and loss. Our team provides supportive, empathetic service to both therapists and their clients.
                </Typography>
              </Card>
            </Grid>

            <Grid item xs={12} md={4}>
              <Card sx={{ height: '100%', p: 4, textAlign: 'center' }}>
                <Security sx={{ fontSize: '3rem', color: 'primary.main', mb: 2 }} />
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Ethical Integrity
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  We maintain the highest ethical standards in all our work, ensuring compliance with professional obligations and protecting client confidentiality throughout the process.
                </Typography>
              </Card>
            </Grid>
          </Grid>
        </Box>

        {/* Our Expertise */}
        <Box sx={{ mb: 8 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 4 }}>
            Our Expertise
          </Typography>
          
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              mb: 4,
              textAlign: 'center',
              maxWidth: '800px',
              mx: 'auto'
            }}
          >
            Our team consists of licensed mental health professionals who have experienced the complexities of practice management and closure firsthand. We combine clinical training with specialized knowledge in:
          </Typography>

          <Grid container spacing={3} sx={{ mt: 2 }}>
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  • Professional Will Development
                </Typography>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  • Client Care Coordination
                </Typography>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  • Grief Counseling
                </Typography>
              </Box>
            </Grid>
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  • HIPAA Compliance & Record Management
                </Typography>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  • Business Closure Procedures
                </Typography>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  • Referral Network Development
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Box>

        {/* Our Commitment */}
        <Box sx={{ textAlign: 'center' }}>
          <Typography variant="h3" gutterBottom sx={{ mb: 4 }}>
            Our Commitment to You
          </Typography>
          
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              mb: 4,
              maxWidth: '800px',
              mx: 'auto'
            }}
          >
            We understand that planning for practice closure can feel overwhelming. That's why we've dedicated ourselves to making this process as smooth and compassionate as possible. When you work with TheraClosure, you're not just getting a service - you're partnering with fellow clinicians who genuinely understand your world and are committed to protecting what you've built.
          </Typography>

          <Typography 
            variant="h5" 
            sx={{ 
              fontWeight: 500,
              fontStyle: 'italic',
              color: 'primary.main',
              mt: 4
            }}
          >
            "You can count on us to be there."
          </Typography>
        </Box>
      </Container>
    </>
  )
}

export default About