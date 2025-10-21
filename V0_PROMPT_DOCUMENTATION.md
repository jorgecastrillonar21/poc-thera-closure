# TheraClosure Website - v0.ly.ai Development Prompt

## Project Overview
Create a comprehensive, professional healthcare website for TheraClosure - a specialized service providing ethical practice closure solutions for psychotherapists. This is a full-stack React application using ConvexDB for backend services.

## 🎯 Core Requirements

### Technology Stack
- **Frontend**: React 18+ with TypeScript
- **Backend**: ConvexDB for real-time database and authentication
- **Styling**: Material-UI (MUI) v5 with custom theme
- **Routing**: React Router v6
- **State Management**: React Context API + ConvexDB hooks
- **Authentication**: ConvexDB Auth
- **Payments**: Stripe integration
- **Deployment**: Vercel (frontend) + ConvexDB hosting

### Brand Identity & Design
- **Primary Color**: Teal (#2C5F5D) - Professional healthcare tone
- **Secondary Color**: Cream (#F5F1E8) - Warm, approachable background
- **Accent Color**: Gold (#D4A574) - Premium service indicator
- **Typography**: Professional, readable fonts with proper contrast
- **Style**: Clean, modern, healthcare-focused with trust-building elements

## 📋 Required Pages & Features

### 1. Landing Page (Home)
**Hero Section:**
- Main headline: "Protecting your clients, your practice, and your peace of mind"
- Subheadline about ethical practice closure for therapists
- Two CTA buttons: "Schedule Free Call" & "View Coverage Plans"
- Professional hero image or healthcare-themed graphics

**Key Sections:**
- **What We Do**: 6 service cards with icons
  1. Practice Closure Planning
  2. Client Transition Services  
  3. Record Management
  4. Legal Compliance
  5. Emergency Response
  6. Retirement Services

- **Why Enroll**: 4 benefit highlights
  1. Peace of Mind (Shield icon)
  2. Professional Legacy (Heart icon) 
  3. Client Care Continuity (Group icon)
  4. Ethical Compliance (Security icon)

- **How It Works**: 3-step process
  1. Enroll & Create Plan
  2. We Monitor & Maintain
  3. Execute When Needed

### 2. About Page
**Content Structure:**
- Hero: "Our Origin Story"
- Mission statement about licensed psychotherapists understanding unique challenges
- Core values section with professional imagery
- Team credentials and expertise
- Trust indicators and professional affiliations

### 3. Coverage & Pricing Page
**Three-Tier Pricing Structure:**

**Essential Plan - $297/year**
- Basic practice closure planning
- Emergency contact protocol
- Essential documentation
- Email support

**Professional Plan - $497/year** (Most Popular)
- Everything in Essential
- Comprehensive closure planning
- Client transition services
- Dedicated support
- Legal document review

**Enterprise Plan - $897/year**
- Everything in Professional
- Multiple practice locations
- Advanced compliance features
- Priority support
- Custom legal documentation

### 4. How We Work Page
**Service Portfolio:**
- Practice Assessment & Planning
- Client Notification Services
- Record Transfer & Management
- Legal & Ethical Compliance
- Emergency Response Protocol
- Retirement Planning Services

### 5. Contact Page
**Contact Form Fields:**
- Full Name (required)
- Email (required) 
- Phone Number
- Practice Type (dropdown)
- State/Location (dropdown)
- Message (textarea)
- Preferred Contact Method (radio buttons)

**Additional Elements:**
- Business hours
- Phone number with professional formatting
- Email address
- Physical address (if applicable)
- Response time expectations

### 6. FAQ Page
**Key Questions Categories:**
- Service Overview
- Pricing & Plans
- Legal & Compliance
- Emergency Procedures
- Account Management

### 7. Testimonials Page
**Client Testimonials:**
- Professional headshots (stock photos)
- Authentic-sounding quotes from therapists
- Variety of practice types and situations
- Location diversity (states/regions)

### 8. Authentication Pages
- **Login**: Email/password with "Remember Me" option
- **Registration**: Multi-step form with plan selection
- **Password Reset**: Email-based recovery
- **Dashboard**: User account management

### 9. User Dashboard
**Protected Features:**
- Account overview
- Current plan details
- Document upload area
- Emergency contacts management
- Billing & payment history
- Support ticket system

## 🛠 Technical Implementation Guide

### ConvexDB Schema Design

```typescript
// Users table
users: {
  _id: Id<"users">,
  email: string,
  firstName: string,
  lastName: string,
  phone?: string,
  practiceType?: string,
  location?: string,
  plan: "essential" | "professional" | "enterprise",
  status: "active" | "inactive" | "pending",
  createdAt: number,
  updatedAt: number
}

// Plans table
plans: {
  _id: Id<"plans">,
  name: string,
  price: number,
  features: string[],
  isPopular: boolean,
  description: string
}

// Contact form submissions
contacts: {
  _id: Id<"contacts">,
  name: string,
  email: string,
  phone?: string,
  practiceType?: string,
  location?: string,
  message: string,
  preferredContact: "email" | "phone",
  status: "new" | "contacted" | "closed",
  createdAt: number
}

// User documents
documents: {
  _id: Id<"documents">,
  userId: Id<"users">,
  fileName: string,
  fileUrl: string,
  fileType: string,
  category: "closure_plan" | "emergency_contacts" | "legal_docs",
  uploadedAt: number
}
```

### Material-UI Theme Configuration

```typescript
import { createTheme } from '@mui/material/styles'

export const theme = createTheme({
  palette: {
    primary: {
      main: '#2C5F5D',
      light: '#4A8B8A',
      dark: '#1A3A38'
    },
    secondary: {
      main: '#F5F1E8',
      dark: '#E8E0D5'
    },
    accent: {
      main: '#D4A574'
    },
    text: {
      primary: '#2C2C2C',
      secondary: '#666666'
    }
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 600,
      fontSize: '2.5rem'
    },
    h2: {
      fontWeight: 600,
      fontSize: '2rem'
    }
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 600
        }
      }
    }
  }
})
```

### Key React Components Structure

```
src/
├── components/
│   ├── layout/
│   │   ├── Navbar.tsx
│   │   ├── Footer.tsx
│   │   └── Layout.tsx
│   ├── ui/
│   │   ├── PricingCard.tsx
│   │   ├── ServiceCard.tsx
│   │   ├── TestimonialCard.tsx
│   │   └── ContactForm.tsx
│   └── auth/
│       ├── ProtectedRoute.tsx
│       └── AuthProvider.tsx
├── pages/
│   ├── Home.tsx
│   ├── About.tsx
│   ├── Coverage.tsx
│   ├── HowWeWork.tsx
│   ├── Contact.tsx
│   ├── FAQ.tsx
│   ├── Testimonials.tsx
│   ├── Login.tsx
│   ├── Register.tsx
│   └── Dashboard.tsx
├── hooks/
│   ├── useAuth.ts
│   └── useConvex.ts
├── utils/
│   ├── constants.ts
│   └── formatters.ts
└── convex/
    ├── auth.ts
    ├── users.ts
    ├── plans.ts
    └── contacts.ts
```

## 🎨 Design Requirements

### Visual Elements
- **Icons**: Use Material-UI icons for consistency
- **Images**: Professional healthcare stock photos
- **Cards**: Subtle shadows with rounded corners
- **Buttons**: Primary (teal), secondary (cream), accent (gold)
- **Forms**: Clean, spacious with proper validation

### Responsive Design
- **Mobile-first**: Ensure all layouts work on mobile devices
- **Tablet**: Optimize for iPad and similar devices  
- **Desktop**: Full-width layouts with proper content containers

### Accessibility
- **Color Contrast**: WCAG AA compliance
- **Keyboard Navigation**: Full keyboard accessibility
- **Screen Readers**: Proper ARIA labels and semantic HTML
- **Focus States**: Clear focus indicators

## 🔐 Authentication & Security

### ConvexDB Auth Setup
- Email/password authentication
- Password reset functionality
- Session management
- Protected routes implementation
- User role management

### Data Security
- Input validation on all forms
- Sanitize user uploads
- Secure file storage
- HIPAA-compliant considerations for healthcare data

## 💰 Stripe Integration

### Payment Features
- Plan subscription management
- Secure checkout flow
- Customer portal integration
- Webhook handling for subscription events
- Invoice generation and management

### Required Stripe Products
- Essential Plan: $297/year
- Professional Plan: $497/year  
- Enterprise Plan: $897/year

## 📱 Key User Flows

### 1. New User Registration
1. Landing page → "Schedule Free Call" or "View Plans"
2. Plan selection → Registration form
3. Account creation → Payment setup
4. Dashboard access → Onboarding

### 2. Contact Form Submission
1. Contact page → Form completion
2. Form validation → Submission
3. Confirmation message → Admin notification
4. Follow-up communication

### 3. Authenticated User Experience
1. Login → Dashboard
2. Plan management → Document upload
3. Emergency contacts → Support access
4. Billing management → Account settings

## 📋 Content Guidelines

### Tone & Voice
- **Professional**: Healthcare industry appropriate
- **Compassionate**: Understanding therapist challenges
- **Trustworthy**: Building confidence in service
- **Clear**: Avoiding jargon, explaining complex concepts

### Key Messaging
- Ethical compliance and professional responsibility
- Peace of mind for practicing therapists
- Seamless client care transition
- Professional legacy protection

## 🚀 Deployment Instructions

### ConvexDB Setup
1. Create ConvexDB project
2. Configure authentication providers
3. Set up database schema
4. Deploy functions and queries

### Vercel Deployment
1. Connect GitHub repository
2. Configure environment variables
3. Set up custom domain
4. Enable automatic deployments

### Environment Variables
```
CONVEX_DEPLOYMENT=your_convex_deployment_url
NEXT_PUBLIC_CONVEX_URL=your_public_convex_url
STRIPE_PUBLIC_KEY=pk_live_or_test_key
STRIPE_SECRET_KEY=sk_live_or_test_key
```

## ✅ Success Criteria

### Functional Requirements
- [ ] All 9+ pages fully functional
- [ ] Complete user authentication system
- [ ] Working contact form with admin notifications
- [ ] Stripe payment integration
- [ ] Responsive design across all devices
- [ ] Fast loading times (<3 seconds)

### Quality Standards
- [ ] No console errors in production
- [ ] Accessibility compliance (WCAG AA)
- [ ] SEO optimization with meta tags
- [ ] Cross-browser compatibility
- [ ] Mobile-responsive design

### Business Goals
- [ ] Clear value proposition communication
- [ ] Professional healthcare industry appearance
- [ ] Trust-building elements throughout
- [ ] Clear call-to-action placement
- [ ] Conversion-optimized user flows

## 📞 Final Prompt for v0.ly.ai

"Create a complete React TypeScript application for TheraClosure, a professional healthcare service helping psychotherapists with ethical practice closure. Use ConvexDB for backend, Material-UI for styling with teal (#2C5F5D) primary color, and include: 

1. **9 main pages**: Home (hero + services), About, Coverage/Pricing (3 tiers: $297, $497, $897), How We Work, Contact form, FAQ, Testimonials, Login, Dashboard
2. **ConvexDB integration**: User auth, contact forms, plan management, file uploads
3. **Stripe integration**: Three subscription tiers with secure checkout
4. **Professional design**: Healthcare-focused, trustworthy, responsive
5. **Key features**: Contact form, user registration, protected dashboard, plan comparison

Focus on creating a conversion-optimized, accessible, and professional website that builds trust with healthcare professionals while maintaining clean, modern design principles."

---

*This documentation provides a complete blueprint for recreating the TheraClosure website using modern React and ConvexDB technologies. Follow this guide to build a professional, scalable healthcare service platform.*