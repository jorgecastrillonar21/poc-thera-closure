import React from 'react'
import { Container, Typography, Box, Grid, TextField, Button } from '@mui/material'
import { Link } from 'react-router-dom'

const Retirement: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Retirement Planning</Typography>
    <Typography variant="body1">Retirement closure planning and transition services.</Typography>
  </Container>
)

const Testimonials: React.FC = () => (
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
          What Our Clients Say
        </Typography>
        <Typography 
          variant="h5" 
          sx={{ 
            opacity: 0.9,
            maxWidth: '600px',
            mx: 'auto'
          }}
        >
          Real experiences from mental health professionals
        </Typography>
      </Container>
    </Box>

    <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
      <Grid container spacing={6}>
        {/* Featured Testimonial */}
        <Grid item xs={12}>
          <Box sx={{ 
            backgroundColor: '#F5F1E8', 
            p: { xs: 4, md: 6 }, 
            borderRadius: 3, 
            textAlign: 'center',
            mb: 4
          }}>
            <Typography 
              variant="h5" 
              sx={{ 
                fontStyle: 'italic', 
                mb: 4, 
                lineHeight: 1.6,
                maxWidth: '800px',
                mx: 'auto'
              }}
            >
              "I cannot speak more highly of my consultation with Dr. Robyn Miller. I contacted Robyn after putting off drafting and even thinking about my professional will for the umpteenth time. Robyn was extremely knowledgeable about the process and gave me a fantastic framework and advice to think through how to approach the whole process."
            </Typography>
            <Typography 
              variant="h6" 
              color="primary" 
              fontWeight={600}
            >
              Dr. Emma Basch, Ph.D.
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.8 }}>
              Licensed Clinical Psychologist
            </Typography>
          </Box>
        </Grid>

        {/* Additional Testimonials */}
        <Grid item xs={12} md={6}>
          <Box sx={{ p: 4, border: '1px solid #E0E0E0', borderRadius: 2, height: '100%' }}>
            <Typography 
              variant="body1" 
              sx={{ 
                fontStyle: 'italic', 
                mb: 3, 
                lineHeight: 1.6
              }}
            >
              "I appreciated her thoughtfulness around all the clinical nuances of this process. Most of all, she was extremely compassionate and supportive in helping me think through my emotional barriers to identifying executors and moving forward with drafting the contract."
            </Typography>
            <Typography variant="h6" color="primary" fontWeight={600}>
              Dr. Sarah Johnson
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.8 }}>
              Private Practice Therapist
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={12} md={6}>
          <Box sx={{ p: 4, border: '1px solid #E0E0E0', borderRadius: 2, height: '100%' }}>
            <Typography 
              variant="body1" 
              sx={{ 
                fontStyle: 'italic', 
                mb: 3, 
                lineHeight: 1.6
              }}
            >
              "Meeting with Robyn eliminated so much stress and uncertainty for me and would highly recommend her services. TheraClosure provides exactly what our profession needs - ethical, compassionate practice closure support."
            </Typography>
            <Typography variant="h6" color="primary" fontWeight={600}>
              Dr. Michael Chen
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.8 }}>
              Group Practice Owner
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={12} md={6}>
          <Box sx={{ p: 4, border: '1px solid #E0E0E0', borderRadius: 2, height: '100%' }}>
            <Typography 
              variant="body1" 
              sx={{ 
                fontStyle: 'italic', 
                mb: 3, 
                lineHeight: 1.6
              }}
            >
              "As someone approaching retirement, TheraClosure gave me the peace of mind I needed. Knowing my clients will be cared for with the same compassion I've shown throughout my career means everything to me."
            </Typography>
            <Typography variant="h6" color="primary" fontWeight={600}>
              Dr. Patricia Williams
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.8 }}>
              Retiring Psychologist
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={12} md={6}>
          <Box sx={{ p: 4, border: '1px solid #E0E0E0', borderRadius: 2, height: '100%' }}>
            <Typography 
              variant="body1" 
              sx={{ 
                fontStyle: 'italic', 
                mb: 3, 
                lineHeight: 1.6
              }}
            >
              "The team at TheraClosure understands our world as clinicians. They've thought through every detail and made the process of creating a professional will straightforward and thorough."
            </Typography>
            <Typography variant="h6" color="primary" fontWeight={600}>
              Dr. James Rodriguez
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.8 }}>
              Licensed Marriage & Family Therapist
            </Typography>
          </Box>
        </Grid>
      </Grid>

      {/* Call to Action */}
      <Box sx={{ textAlign: 'center', mt: 10 }}>
        <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
          Ready to Experience TheraClosure?
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
          Join the growing number of mental health professionals who trust TheraClosure to protect their clients and practices.
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

const FAQ: React.FC = () => (
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
          Frequently Asked Questions
        </Typography>
        <Typography 
          variant="h5" 
          sx={{ 
            opacity: 0.9,
            maxWidth: '600px',
            mx: 'auto'
          }}
        >
          Everything you need to know about TheraClosure
        </Typography>
      </Container>
    </Box>

    <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Typography variant="h4" gutterBottom sx={{ mb: 6, textAlign: 'center' }}>
            Common Questions About Our Services
          </Typography>
        </Grid>

        <Grid item xs={12} md={6}>
          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
              What is a Professional Will?
            </Typography>
            <Typography sx={{ lineHeight: 1.6, mb: 2 }}>
              A Professional Will is a legally reviewed document that outlines your wishes for client care and practice closure in the event of your incapacitation, death, or retirement. It names TheraClosure as your Practice Executor to carry out your detailed directives.
            </Typography>
          </Box>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
              How quickly can TheraClosure respond if needed?
            </Typography>
            <Typography sx={{ lineHeight: 1.6, mb: 2 }}>
              We maintain 24/7 availability and can typically respond within hours of notification. Our team is always on standby with the resources and expertise needed to immediately support your practice and clients.
            </Typography>
          </Box>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
              Is my information secure with TheraClosure?
            </Typography>
            <Typography sx={{ lineHeight: 1.6, mb: 2 }}>
              Absolutely. We maintain HIPAA-compliant, encrypted storage systems and follow the highest security standards. All client information is protected according to legal requirements and ethical obligations.
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={12} md={6}>
          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
              What happens to my client records?
            </Typography>
            <Typography sx={{ lineHeight: 1.6, mb: 2 }}>
              We secure the confidentiality of psychotherapy records in compliance with ethical standards and retain them according to legal requirements. We facilitate notification and transfer of records as needed for patient care.
            </Typography>
          </Box>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
              How are my clients notified and cared for?
            </Typography>
            <Typography sx={{ lineHeight: 1.6, mb: 2 }}>
              Our team of grief-trained therapists provides compassionate, clinically-informed notification to patients. We facilitate personalized referrals for each client and ensure no patient slips through the cracks.
            </Typography>
          </Box>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary" fontWeight={600}>
              Can I update my Professional Will?
            </Typography>
            <Typography sx={{ lineHeight: 1.6, mb: 2 }}>
              Yes, you can update your Professional Will at any time. We recommend reviewing it annually or whenever there are significant changes to your practice or personal circumstances.
            </Typography>
          </Box>
        </Grid>
      </Grid>

      <Box sx={{ textAlign: 'center', mt: 8 }}>
        <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
          Still Have Questions?
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
          We're here to help. Schedule a free consultation to discuss your specific needs and learn how TheraClosure can provide peace of mind for your practice.
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
            CONTACT US
          </Button>
        </Box>
      </Box>
    </Container>
  </>
)

const Contact: React.FC = () => (
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
          Contact Us
        </Typography>
        <Typography 
          variant="h5" 
          sx={{ 
            opacity: 0.9,
            maxWidth: '600px',
            mx: 'auto'
          }}
        >
          Get in touch with our compassionate team
        </Typography>
      </Container>
    </Box>

    <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
      <Grid container spacing={6}>
        <Grid item xs={12} md={6}>
          <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
            Send us a message
          </Typography>
          
          <Box component="form" sx={{ mb: 4 }}>
            <TextField
              fullWidth
              label="Full Name"
              margin="normal"
              required
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Email Address"
              type="email"
              margin="normal"
              required
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Phone Number"
              margin="normal"
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Message"
              multiline
              rows={4}
              margin="normal"
              required
              sx={{ mb: 3 }}
            />
            <Button
              variant="contained"
              size="large"
              sx={{ 
                fontWeight: 600,
                px: 4,
                py: 1.5,
                fontSize: '1.1rem'
              }}
            >
              SEND MESSAGE
            </Button>
          </Box>
        </Grid>

        <Grid item xs={12} md={6}>
          <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
            Get in Touch
          </Typography>
          
          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary">
              Schedule a Free Consultation
            </Typography>
            <Typography sx={{ mb: 3, lineHeight: 1.6 }}>
              Ready to learn more about TheraClosure? Schedule a free 20-minute informational session with one of our clinicians.
            </Typography>
            <Button
              component={Link}
              to="/contact"
              variant="contained"
              sx={{ 
                backgroundColor: '#D4A574',
                color: 'white',
                fontWeight: 600,
                px: 3,
                py: 1,
                '&:hover': { backgroundColor: '#B8935F' }
              }}
            >
              SCHEDULE NOW
            </Button>
          </Box>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" gutterBottom color="primary">
              Questions About Our Services?
            </Typography>
            <Typography sx={{ mb: 2, lineHeight: 1.6 }}>
              Have questions about our coverage options and whether TheraClosure is right for you? Contact our friendly team today.
            </Typography>
          </Box>

          <Box>
            <Typography variant="h6" gutterBottom color="primary">
              Ready to Sign Up?
            </Typography>
            <Typography sx={{ mb: 3, lineHeight: 1.6 }}>
              Let us make it easy. Begin today with our simple enrollment process.
            </Typography>
            <Button
              component={Link}
              to="/coverage"
              variant="outlined"
              sx={{ 
                fontWeight: 600,
                px: 3,
                py: 1
              }}
            >
              ENROLL NOW
            </Button>
          </Box>
        </Grid>
      </Grid>
    </Container>
  </>
)

const Enrollment: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Enrollment</Typography>
    <Typography variant="body1">Begin your practice closure journey.</Typography>
  </Container>
)

const Templates: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Document Templates</Typography>
    <Typography variant="body1">Professional templates for client communication.</Typography>
  </Container>
)

const Billing: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Billing & Payments</Typography>
    <Typography variant="body1">Manage your subscription and billing.</Typography>
  </Container>
)

const Support: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Support Center</Typography>
    <Typography variant="body1">Get help with your account and services.</Typography>
  </Container>
)

export { Retirement, Testimonials, FAQ, Contact, Enrollment, Templates, Billing, Support }