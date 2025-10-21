import React from 'react'
import { Container, Typography, Box, Button, Grid, Card } from '@mui/material'
import { Link } from 'react-router-dom'
import { Shield, FavoriteRounded, Group, Security, EventNote, Business } from '@mui/icons-material'

const Home: React.FC = () => {
  return (
    <>
      {/* Hero Section */}
      <Box
        sx={{
          background: 'linear-gradient(135deg, #2C5F5D 0%, #4A8B8A 100%)',
          color: 'white',
          py: { xs: 6, md: 10 },
          textAlign: 'center',
        }}
      >
        <Container maxWidth="lg">
          <Typography 
            variant="h2" 
            sx={{ 
              fontWeight: 600, 
              mb: 3,
              fontSize: { xs: '1.8rem', md: '2.5rem' },
              lineHeight: 1.2
            }}
          >
            Protecting your clients, your practice, and your peace of mind
          </Typography>
          
          <Box sx={{ display: 'flex', gap: 3, justifyContent: 'center', flexWrap: 'wrap', mt: 4 }}>
            <Button
              component={Link}
              to="/contact"
              variant="contained"
              size="large"
              sx={{ 
                backgroundColor: 'secondary.main',
                color: 'primary.main',
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem',
                '&:hover': { backgroundColor: 'secondary.dark' }
              }}
            >
              SCHEDULE A FREE INFORMATIONAL CALL
            </Button>
            <Button
              component={Link}
              to="/coverage"
              variant="contained"
              size="large"
              sx={{ 
                backgroundColor: '#D4A574',
                color: 'white',
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem',
                '&:hover': { backgroundColor: '#B8935F' }
              }}
            >
              PURCHASE COVERAGE
            </Button>
          </Box>
        </Container>
      </Box>

      {/* Main Message Section */}
      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
        <Typography 
          variant="h3" 
          textAlign="center" 
          gutterBottom
          sx={{ 
            fontSize: { xs: '1.5rem', md: '2rem' },
            fontWeight: 500,
            mb: 4
          }}
        >
          You've devoted your professional life to your practice and your clients. But what happens to them if something happens to you?
        </Typography>
        
        <Typography 
          variant="body1" 
          textAlign="center" 
          sx={{ 
            fontSize: '1.1rem', 
            lineHeight: 1.7,
            maxWidth: '900px',
            mx: 'auto',
            mb: 6 
          }}
        >
          We know planning for the worst is troubling to think about. However, failing to plan can place an overwhelming logistical and emotional burden on family, friends and colleagues, and can be quite costly. TheraClosure is the compassionate, efficient, cost-effective, and secure way to fulfill your ethical obligations and to ensure that your patients will be well cared for. We make it simple for psychotherapists to create a thoughtful Professional Will, and we learn about your practice and then serve, if needed, as your Practice Executor, administering your detailed directive with our caring expertise.
        </Typography>
      </Container>

      {/* What We Do Section */}
      <Box sx={{ backgroundColor: '#F5F1E8', py: { xs: 6, md: 10 } }}>
        <Container maxWidth="lg">
          <Typography 
            variant="h3" 
            textAlign="center" 
            gutterBottom
            sx={{ mb: 6 }}
          >
            What We Do
          </Typography>
          
          <Grid container spacing={4}>
            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                  <Shield sx={{ color: 'primary.main', mr: 2, fontSize: '2rem', mt: 0.5 }} />
                  <Box>
                    <Typography variant="h6" color="primary" fontWeight={600} sx={{ mb: 1 }}>
                      Professional Will Consultation & Practice Executor Provided
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      We take time to learn about you and your practice. Together, we consider best practices and devise a plan for your clients' well-being and for the closing of your business in the event of a sudden incapacitation. A legally reviewed contractual Professional Will is created, naming TheraClosure as your Practice Executor.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>

            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                  <FavoriteRounded sx={{ color: 'primary.main', mr: 2, fontSize: '2rem', mt: 0.5 }} />
                  <Box>
                    <Typography variant="h6" color="primary" fontWeight={600} sx={{ mb: 1 }}>
                      Notification & Continuing Care Coordination
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      If called upon, our team specializes in compassionate, clinically-informed notification to patients. Trained in grief counseling, our therapists provide support and facilitate personalized referrals for each client. We ensure that no patient slips through the cracks.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>

            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                  <EventNote sx={{ color: 'primary.main', mr: 2, fontSize: '2rem', mt: 0.5 }} />
                  <Box>
                    <Typography variant="h6" color="primary" fontWeight={600} sx={{ mb: 1 }}>
                      Retirement Services
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      Our goal is to provide you with peace of mind that you have closed your practice with the same degree of conscientiousness, ethical compliance, and compassion with which you provided care for so many years.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>

            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                  <Security sx={{ color: 'primary.main', mr: 2, fontSize: '2rem', mt: 0.5 }} />
                  <Box>
                    <Typography variant="h6" color="primary" fontWeight={600} sx={{ mb: 1 }}>
                      Confidentiality & Record Retention
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      We secure the confidentiality of psychotherapy records in compliance with ethical standards and retain them according to legal requirements. We facilitate notification and transfer of records as needed for patient care.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Mission Statement */}
      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
        <Box textAlign="center">
          <Typography 
            variant="h4" 
            gutterBottom
            sx={{ 
              fontWeight: 500,
              mb: 4,
              fontSize: { xs: '1.3rem', md: '1.75rem' }
            }}
          >
            At TheraClosure, we are dedicated to supporting therapists in safeguarding the well-being and privacy of their clients.
          </Typography>
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              maxWidth: '800px',
              mx: 'auto' 
            }}
          >
            Our focus is to provide seamless and reliable ethically-compliant solutions for therapists to ensure confidentiality and continuity of care for clients in the event of the therapist's incapacitation, death, or retirement.
          </Typography>
        </Box>
      </Container>

      {/* Why Enroll Section */}
      <Box sx={{ backgroundColor: 'primary.main', color: 'white', py: { xs: 6, md: 10 } }}>
        <Container maxWidth="lg">
          <Typography 
            variant="h3" 
            textAlign="center" 
            gutterBottom
            sx={{ mb: 6, color: 'white' }}
          >
            Why Enroll?
          </Typography>
          
          <Grid container spacing={4}>
            <Grid item xs={12} md={6}>
              <Box textAlign="center" sx={{ p: 3 }}>
                <Shield sx={{ fontSize: '3rem', mb: 2, color: 'secondary.main' }} />
                <Typography variant="h5" gutterBottom sx={{ color: 'white' }}>
                  Fulfill your ethical obligations with an easy, fail-safe solution
                </Typography>
                <Typography sx={{ lineHeight: 1.6, opacity: 0.9 }}>
                  You practice psychotherapy with integrity and the highest ethical standards which must include planning for client care in the event of your sudden incapacitation or death.
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box textAlign="center" sx={{ p: 3 }}>
                <FavoriteRounded sx={{ fontSize: '3rem', mb: 2, color: 'secondary.main' }} />
                <Typography variant="h5" gutterBottom sx={{ color: 'white' }}>
                  Protect your clients from added trauma
                </Typography>
                <Typography sx={{ lineHeight: 1.6, opacity: 0.9 }}>
                  Ensure your clients receive compassionate, professional care during their most vulnerable time with grief-trained clinicians who understand their needs.
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box textAlign="center" sx={{ p: 3 }}>
                <Group sx={{ fontSize: '3rem', mb: 2, color: 'secondary.main' }} />
                <Typography variant="h5" gutterBottom sx={{ color: 'white' }}>
                  Protect your loved ones from stress and financial risk
                </Typography>
                <Typography sx={{ lineHeight: 1.6, opacity: 0.9 }}>
                  Take the burden off family and colleagues who otherwise would scramble to manage your practice during an already difficult time.
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box textAlign="center" sx={{ p: 3 }}>
                <Business sx={{ fontSize: '3rem', mb: 2, color: 'secondary.main' }} />
                <Typography variant="h5" gutterBottom sx={{ color: 'white' }}>
                  Retiring? Step away from your practice while protecting your legacy
                </Typography>
                <Typography sx={{ lineHeight: 1.6, opacity: 0.9 }}>
                  We assist you in enjoying a full retirement while protecting your legacy with the same conscientiousness you've shown throughout your career.
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Final CTA Section */}
      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 8 }, textAlign: 'center' }}>
        <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
          We are clinicians, too. You can count on us to be there.
        </Typography>
        
        <Typography variant="h3" gutterBottom sx={{ mb: 4 }}>
          Ready to get started?
        </Typography>
        
        <Box sx={{ display: 'flex', gap: 3, justifyContent: 'center', flexWrap: 'wrap' }}>
          <Button
            component={Link}
            to="/coverage"
            variant="contained"
            size="large"
            sx={{ 
              fontWeight: 600,
              px: 4,
              py: 1.5,
              fontSize: '1.1rem'
            }}
          >
            SIGN UP
          </Button>
          <Button
            component={Link}
            to="/contact"
            variant="outlined"
            size="large"
            sx={{ 
              fontWeight: 600,
              px: 4,
              py: 1.5,
              fontSize: '1.1rem'
            }}
          >
            BOOK FREE INFORMATIONAL CALL
          </Button>
        </Box>
      </Container>
    </>
  )
}

export default Home