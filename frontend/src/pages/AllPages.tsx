import React from 'react'
import { Container, Typography } from '@mui/material'

const Retirement: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Retirement Planning</Typography>
    <Typography variant="body1">Retirement closure planning and transition services.</Typography>
  </Container>
)

const Testimonials: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Client Testimonials</Typography>
    <Typography variant="body1">What our clients say about our services.</Typography>
  </Container>
)

const FAQ: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Frequently Asked Questions</Typography>
    <Typography variant="body1">Common questions about practice closure.</Typography>
  </Container>
)

const Contact: React.FC = () => (
  <Container maxWidth="lg" sx={{ py: 4 }}>
    <Typography variant="h4" component="h1" gutterBottom>Contact Us</Typography>
    <Typography variant="body1">Get in touch with our team.</Typography>
  </Container>
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