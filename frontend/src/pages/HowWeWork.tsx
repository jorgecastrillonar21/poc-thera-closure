import React from 'react'
import { Container, Typography, Box, Grid, Card, Button } from '@mui/material'
import { Link } from 'react-router-dom'
import { Assessment, PersonAdd, HealthAndSafety, VpnKey, Payment } from '@mui/icons-material'

const HowWeWork: React.FC = () => {
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
            How We Work
          </Typography>
          <Typography 
            variant="h5" 
            sx={{ 
              opacity: 0.9,
              maxWidth: '700px',
              mx: 'auto'
            }}
          >
            A comprehensive, compassionate approach to practice closure
          </Typography>
        </Container>
      </Box>

      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
        
        {/* Process Overview */}
        <Box sx={{ mb: 10 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 6 }}>
            Our Complete Service Portfolio
          </Typography>
          
          <Grid container spacing={4}>
            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                  <Assessment sx={{ fontSize: '3rem', color: 'primary.main', mr: 3 }} />
                  <Box>
                    <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
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
              <Card sx={{ height: '100%', p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                  <PersonAdd sx={{ fontSize: '3rem', color: 'primary.main', mr: 3 }} />
                  <Box>
                    <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
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
              <Card sx={{ height: '100%', p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                  <HealthAndSafety sx={{ fontSize: '3rem', color: 'primary.main', mr: 3 }} />
                  <Box>
                    <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                      Retirement Services
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      Our goal is to provide you with peace of mind that you have closed your practice with the same degree of conscientiousness, ethical compliance, and compassion with which you provided care for so many years. We assist you in enjoying a full retirement while protecting your legacy.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>

            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                  <HealthAndSafety sx={{ fontSize: '3rem', color: 'primary.main', mr: 3 }} />
                  <Box>
                    <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                      Confidentiality & Record Retention
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      We secure the confidentiality of psychotherapy records in compliance with ethical standards and retain them according to legal requirements. We facilitate notification and transfer of records as needed for patient care.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>

            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                  <VpnKey sx={{ fontSize: '3rem', color: 'primary.main', mr: 3 }} />
                  <Box>
                    <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                      Practice Password Manager
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      We provide you with an enterprise-grade password manager for the online services your practice uses, including your electronic health records system. Your passwords will be kept up-to-date, safely encrypted, and only made available to TheraClosure to administer your will in the event of your incapacitation.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>

            <Grid item xs={12} md={6}>
              <Card sx={{ height: '100%', p: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                  <Payment sx={{ fontSize: '3rem', color: 'primary.main', mr: 3 }} />
                  <Box>
                    <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                      Billing Completion & Practice Closure
                    </Typography>
                    <Typography sx={{ lineHeight: 1.6 }}>
                      We discontinue your business accounts, notifying licensure boards and malpractice insurance, and consulting with your Personal Executor about remaining tasks. For those who bill using EHRs, we issue statements to your active clients and to insurance so that your estate can collect for services you have provided.
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Grid>
          </Grid>
        </Box>

        {/* How We Gain Access */}
        <Box sx={{ backgroundColor: '#F5F1E8', py: { xs: 6, md: 8 }, px: 4, borderRadius: 2, mb: 8 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 4 }}>
            How We Access Your Records
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
            In order for TheraClosure to fulfill duties of the Practice Executor, we would need remote access to your patient schedule, contact information, and your patient records. We can gain this access in the following ways:
          </Typography>
          
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  1. Electronic Health Record (EHR) System
                </Typography>
                <Typography sx={{ mb: 3, lineHeight: 1.6 }}>
                  Direct secure access to your existing EHR platform
                </Typography>
                
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  2. Web-based System
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Access through other secure online platforms you currently use
                </Typography>
              </Box>
            </Grid>
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  3. Digital File Transmission
                </Typography>
                <Typography sx={{ mb: 3, lineHeight: 1.6 }}>
                  Secure digital files transmitted by someone you designate
                </Typography>
                
                <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
                  4. Digitized Paper Records
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Paper records that have been digitized (which we can help with)
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Box>

        {/* Our Approach */}
        <Box sx={{ textAlign: 'center', mb: 8 }}>
          <Typography variant="h3" gutterBottom sx={{ mb: 4 }}>
            Our Approach
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
            Having served as Practice Executors, we know the tremendous amount of time and expertise required to transfer client care, safeguard their privacy, and manage the business aspects of another person's practice. We have learned that when therapists rely on a trusted colleague to be their Practice Executor, the colleague will not have the time or resources to handle all the demands that are suddenly placed on them.
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
            TheraClosure takes the burden off your colleagues who otherwise may scramble to inform clients who show up in the waiting room and to provide coverage and referrals. As this is our sole responsibility, you can count on us to be available with the resources needed to jump into action for your practice.
          </Typography>
        </Box>

        {/* Call to Action */}
        <Box sx={{ textAlign: 'center' }}>
          <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
            Ready to Learn More?
          </Typography>
          
          <Box sx={{ display: 'flex', gap: 3, justifyContent: 'center', flexWrap: 'wrap' }}>
            <Button
              component={Link}
              to="/contact"
              variant="contained"
              size="large"
              sx={{ 
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem'
              }}
            >
              SCHEDULE FREE CONSULTATION
            </Button>
            <Button
              component={Link}
              to="/coverage"
              variant="outlined"
              size="large"
              sx={{ 
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem'
              }}
            >
              VIEW COVERAGE OPTIONS
            </Button>
          </Box>
        </Box>
      </Container>
    </>
  )
}

export default HowWeWork