import React from 'react'
import { Routes, Route } from 'react-router-dom'
import { Box } from '@mui/material'
import Layout from './components/Layout/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import About from './pages/About'
import Coverage from './pages/Coverage'
import HowWeWork from './pages/HowWeWork'
import Retirement from './pages/Retirement'
import Testimonials from './pages/Testimonials'
import FAQ from './pages/FAQ'
import Contact from './pages/Contact'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Enrollment from './pages/Enrollment'
import Templates from './pages/Templates'
import Billing from './pages/Billing'
import Support from './pages/Support'

const App: React.FC = () => {
  return (
    <Box>
      <Routes>
        <Route path="/" element={<Layout />}>
          {/* Public routes */}
          <Route index element={<Home />} />
          <Route path="about" element={<About />} />
          <Route path="coverage" element={<Coverage />} />
          <Route path="how-we-work" element={<HowWeWork />} />
          <Route path="retirement" element={<Retirement />} />
          <Route path="testimonials" element={<Testimonials />} />
          <Route path="faq" element={<FAQ />} />
          <Route path="contact" element={<Contact />} />
          <Route path="login" element={<Login />} />
          
          {/* Protected routes */}
          <Route path="dashboard" element={<ProtectedRoute />}>
            <Route index element={<Dashboard />} />
            <Route path="enrollment" element={<Enrollment />} />
            <Route path="templates" element={<Templates />} />
            <Route path="billing" element={<Billing />} />
            <Route path="support" element={<Support />} />
          </Route>
        </Route>
      </Routes>
    </Box>
  )
}

export default App