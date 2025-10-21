import React from 'react'
import { Container, Typography, Box, Grid, Card, Button, List, ListItem, ListItemIcon, ListItemText } from '@mui/material'
import { Link } from 'react-router-dom'
import { CheckCircle, Star, Shield } from '@mui/icons-material'

const Coverage: React.FC = () => {
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
            Coverage & Costs
          </Typography>
          <Typography 
            variant="h5" 
            sx={{ 
              opacity: 0.9,
              maxWidth: '600px',
              mx: 'auto'
            }}
          >
            Comprehensive protection for your practice and peace of mind for you
          </Typography>
        </Container>
      </Box>

      <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
        
        {/* Pricing Plans */}
        <Box sx={{ mb: 10 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 6 }}>
            Choose Your Coverage Plan
          </Typography>
          
          <Grid container spacing={4} justifyContent="center">
            {/* Essential Plan */}
            <Grid item xs={12} md={4}>
              <Card sx={{ 
                height: '100%', 
                p: 4, 
                textAlign: 'center',
                border: '2px solid transparent',
                '&:hover': {
                  border: '2px solid #2C5F5D',
                  boxShadow: '0 8px 25px rgba(44, 95, 93, 0.15)'
                }
              }}>
                <Shield sx={{ fontSize: '3rem', color: 'primary.main', mb: 2 }} />
                <Typography variant="h4" gutterBottom color="primary" fontWeight={600}>
                  Essential
                </Typography>
                <Typography variant="h3" gutterBottom sx={{ color: 'primary.main', fontWeight: 700 }}>
                  $297<Typography component="span" variant="h6">/year</Typography>
                </Typography>
                <Typography variant="body1" sx={{ mb: 4, opacity: 0.8 }}>
                  Basic protection for solo practitioners
                </Typography>
                
                <List sx={{ mb: 4 }}>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Professional Will Creation" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Practice Executor Services" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Client Notification (up to 50 clients)" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Basic Record Management" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Referral Coordination" />
                  </ListItem>
                </List>

                <Button
                  component={Link}
                  to="/contact"
                  variant="outlined"
                  size="large"
                  fullWidth
                  sx={{ 
                    fontWeight: 600,
                    py: 1.5,
                    fontSize: '1.1rem'
                  }}
                >
                  GET STARTED
                </Button>
              </Card>
            </Grid>

            {/* Professional Plan */}
            <Grid item xs={12} md={4}>
              <Card sx={{ 
                height: '100%', 
                p: 4, 
                textAlign: 'center',
                border: '2px solid #D4A574',
                position: 'relative',
                '&:hover': {
                  boxShadow: '0 8px 25px rgba(212, 165, 116, 0.25)'
                }
              }}>
                <Box sx={{ 
                  position: 'absolute',
                  top: -15,
                  left: '50%',
                  transform: 'translateX(-50%)',
                  backgroundColor: '#D4A574',
                  color: 'white',
                  px: 3,
                  py: 1,
                  borderRadius: 2,
                  fontSize: '0.9rem',
                  fontWeight: 600,
                  zIndex: 10,
                  boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
                }}>
                  MOST POPULAR
                </Box>
                
                <Star sx={{ fontSize: '3rem', color: '#D4A574', mb: 2, mt: 2 }} />
                <Typography variant="h4" gutterBottom color="primary" fontWeight={600}>
                  Professional
                </Typography>
                <Typography variant="h3" gutterBottom sx={{ color: 'primary.main', fontWeight: 700 }}>
                  $497<Typography component="span" variant="h6">/year</Typography>
                </Typography>
                <Typography variant="body1" sx={{ mb: 4, opacity: 0.8 }}>
                  Complete coverage for established practices
                </Typography>
                
                <List sx={{ mb: 4 }}>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Everything in Essential" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Client Notification (up to 150 clients)" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Enterprise Password Manager" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Insurance & Billing Completion" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Priority Support" />
                  </ListItem>
                </List>

                <Button
                  component={Link}
                  to="/contact"
                  variant="contained"
                  size="large"
                  fullWidth
                  sx={{ 
                    fontWeight: 600,
                    py: 1.5,
                    fontSize: '1.1rem',
                    backgroundColor: '#D4A574',
                    '&:hover': { backgroundColor: '#B8935F' }
                  }}
                >
                  GET STARTED
                </Button>
              </Card>
            </Grid>

            {/* Enterprise Plan */}
            <Grid item xs={12} md={4}>
              <Card sx={{ 
                height: '100%', 
                p: 4, 
                textAlign: 'center',
                border: '2px solid transparent',
                '&:hover': {
                  border: '2px solid #2C5F5D',
                  boxShadow: '0 8px 25px rgba(44, 95, 93, 0.15)'
                }
              }}>
                <Shield sx={{ fontSize: '3rem', color: 'primary.main', mb: 2 }} />
                <Typography variant="h4" gutterBottom color="primary" fontWeight={600}>
                  Enterprise
                </Typography>
                <Typography variant="h3" gutterBottom sx={{ color: 'primary.main', fontWeight: 700 }}>
                  $897<Typography component="span" variant="h6">/year</Typography>
                </Typography>
                <Typography variant="body1" sx={{ mb: 4, opacity: 0.8 }}>
                  Full-service for large practices & groups
                </Typography>
                
                <List sx={{ mb: 4 }}>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Everything in Professional" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Unlimited Client Notifications" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Multi-location Support" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Custom Integration" />
                  </ListItem>
                  <ListItem sx={{ px: 0 }}>
                    <ListItemIcon>
                      <CheckCircle sx={{ color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary="Dedicated Account Manager" />
                  </ListItem>
                </List>

                <Button
                  component={Link}
                  to="/contact"
                  variant="outlined"
                  size="large"
                  fullWidth
                  sx={{ 
                    fontWeight: 600,
                    py: 1.5,
                    fontSize: '1.1rem'
                  }}
                >
                  CONTACT US
                </Button>
              </Card>
            </Grid>
          </Grid>
        </Box>

        {/* What's Included */}
        <Box sx={{ backgroundColor: '#F5F1E8', py: { xs: 6, md: 8 }, px: 4, borderRadius: 2, mb: 8 }}>
          <Typography variant="h3" gutterBottom textAlign="center" sx={{ mb: 6 }}>
            What's Included in Every Plan
          </Typography>
          
          <Grid container spacing={4}>
            <Grid item xs={12} md={6}>
              <Box sx={{ mb: 4 }}>
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Professional Will Creation
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Legally reviewed contractual Professional Will naming TheraClosure as your Practice Executor, customized to your specific practice needs.
                </Typography>
              </Box>

              <Box sx={{ mb: 4 }}>
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Client Care Coordination
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Compassionate, clinically-informed notification and personalized referrals by grief-trained therapists.
                </Typography>
              </Box>

              <Box>
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Record Management
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  HIPAA-compliant secure storage, retention according to legal requirements, and transfer facilitation.
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box sx={{ mb: 4 }}>
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Business Closure Support
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Account discontinuation, licensure board notification, malpractice insurance coordination.
                </Typography>
              </Box>

              <Box sx={{ mb: 4 }}>
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  24/7 Availability
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Ready to act immediately when needed, with resources and expertise always on standby.
                </Typography>
              </Box>

              <Box>
                <Typography variant="h5" gutterBottom color="primary" fontWeight={600}>
                  Clinical Expertise
                </Typography>
                <Typography sx={{ lineHeight: 1.6 }}>
                  Licensed mental health professionals who understand your world and ethical obligations.
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Box>

        {/* Call to Action */}
        <Box sx={{ textAlign: 'center' }}>
          <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
            Start Protecting Your Practice Today
          </Typography>
          
          <Typography 
            variant="body1" 
            sx={{ 
              fontSize: '1.1rem', 
              lineHeight: 1.7,
              mb: 6,
              maxWidth: '600px',
              mx: 'auto'
            }}
          >
            Don't wait until it's too late. Schedule a free consultation to learn how TheraClosure can provide peace of mind for you and protection for your clients.
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
              to="/how-we-work"
              variant="outlined"
              size="large"
              sx={{ 
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem'
              }}
            >
              LEARN HOW WE WORK
            </Button>
          </Box>
        </Box>
      </Container>
    </>
  )
}

export default Coverage